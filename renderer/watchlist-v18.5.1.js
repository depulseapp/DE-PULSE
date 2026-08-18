(() => {
  'use strict';

  const baseBindDynamic = bindDynamic;

  function normalizeTrackedSymbol(value) {
    return String(value || '').trim().toUpperCase();
  }

  function validDeskSymbol(symbol) {
    return /^[A-Z][A-Z0-9.\-]{0,7}$/.test(normalizeTrackedSymbol(symbol));
  }

  async function addSymbolToDesk(symbol, desk) {
    const sym = normalizeTrackedSymbol(symbol);
    if (!validDeskSymbol(sym)) {
      throw new Error('Enter a valid stock/ETF ticker.');
    }
    if (!DESKS.includes(desk)) {
      throw new Error('Invalid trading desk.');
    }

    const membershipResult = await api('/api/desk/membership', {
      symbol: sym,
      desk,
      active: true
    });
    const boot = await api('/api/bootstrap');
    state = boot.state;
    runtime = boot.runtime;

    if (!(deskWL(desk).symbols || []).includes(sym)) {
      throw new Error(`${sym} was not saved to the ${deskCfg[desk].title} watchlist.`);
    }
    return membershipResult;
  }

  deskMembershipStrip = function deskMembershipStripV186(sym) {
    const target = String(sym || '').toUpperCase();
    const items = [['day', 'DAY'], ['swing', 'SWING'], ['long', 'LONG']];
    return `<span class="desk-membership-strip" aria-label="Desk membership for ${esc(target)}">${items.map(([kind, label]) => {
      const active = (deskWL(kind).symbols || []).includes(target);
      const deskName = kind === 'long' ? 'Long-Term' : kind === 'day' ? 'Day Trade' : 'Swing';
      const title = active
        ? `${target} is in ${deskName} · click to remove when another desk remains`
        : `Add ${target} to ${deskName}`;
      return `<button type="button" class="desk-membership-pill ${active ? 'active' : ''}" aria-pressed="${active ? 'true' : 'false'}" data-desk-membership="${kind}:${esc(target)}" title="${esc(title)}">${label}</button>`;
    }).join('')}</span>`;
  };

  async function removeTrackedSymbolEverywhere(symbol, ctx) {
    const sym = normalizeTrackedSymbol(symbol);
    if (!sym) return null;

    const selectedBefore = Object.fromEntries(DESKS.map(kind => [kind, selected[kind] || '']));
    const res = await api('/api/master-symbol/remove', { symbol: sym });
    const removed = res.removed || {};
    const boot = await api('/api/bootstrap');
    state = boot.state;
    runtime = boot.runtime;

    for (const kind of DESKS) {
      if (selected[kind] === sym) selected[kind] = deskWL(kind).symbols[0] || '';
    }

    render();
    restoreSaveContext(ctx);
    toast(
      `${sym} Removed from Tracked Symbols`,
      'Removed from every desk. Undo restores the exact previous desk memberships.',
      'warning'
    );

    window.__masterUndo = {
      sym,
      membership: removed,
      selected: selectedBefore,
      expires: Date.now() + 9000
    };

    const host = $('#header-notification');
    if (host) {
      host.innerHTML += ` <button class="toast-undo" data-master-undo="${esc(sym)}">UNDO</button>`;
      host.querySelector('[data-master-undo]')?.addEventListener('click', async ev => {
        ev.stopPropagation();
        const undo = window.__masterUndo;
        if (!undo || undo.sym !== sym || Date.now() > undo.expires) return;

        const undoCtx = captureSaveContext();
        await api('/api/master-symbol/restore', { symbol: sym, membership: undo.membership });
        const boot2 = await api('/api/bootstrap');
        state = boot2.state;
        runtime = boot2.runtime;
        for (const kind of DESKS) {
          const prior = undo.selected?.[kind] || '';
          selected[kind] = prior || selected[kind] || deskWL(kind).symbols[0] || '';
        }
        window.__masterUndo = null;
        render();
        restoreSaveContext(undoCtx);
        toast(`${sym} Restored`, 'Previous desk memberships and desk selection restored.', 'success');
      });
    }
    return res;
  }

  function bindCanonicalDeskAdds() {
    $$('[data-add-desk]').forEach(button => {
      button.onclick = async () => {
        const desk = button.dataset.addDesk;
        const input = $(`[data-add-input="${desk}"]`);
        const symbol = normalizeTrackedSymbol(input?.value || watchlistDraft[desk] || '');
        if (!symbol) {
          toast('Enter a Ticker', 'Enter a ticker such as NVDA.');
          input?.focus();
          return;
        }

        const ctx = captureSaveContext();
        const original = button.textContent;
        button.disabled = true;
        button.textContent = 'Adding…';
        watchlistDraft[desk] = symbol;
        try {
          await addSymbolToDesk(symbol, desk);
          selected[desk] = symbol;
          watchlistDraft[desk] = '';
          if (input) input.value = '';
          const hasPrice = num(runtime?.quotes?.[symbol]?.price) > 0;
          toast(
            `${symbol} Added`,
            hasPrice ? `${deskCfg[desk].title} · Market data ready.` : `${deskCfg[desk].title} · Saved. Waiting for market data.`,
            'success'
          );
          updateChrome();
          render();
          restoreSaveContext(ctx);
        } catch (err) {
          watchlistDraft[desk] = symbol;
          toast('Could Not Add Symbol', err.message, 'error');
          render();
          restoreSaveContext(ctx);
          setTimeout(() => {
            const next = $(`[data-add-input="${desk}"]`);
            if (next) {
              next.focus({ preventScroll: true });
              next.setSelectionRange(next.value.length, next.value.length);
            }
          }, 0);
        } finally {
          if (button.isConnected) {
            button.disabled = false;
            button.textContent = original;
          }
        }
      };
    });

    $$('[data-add-desk-table]').forEach(button => {
      button.onclick = async () => {
        const desk = button.dataset.addDeskTable;
        const input = $(`[data-add-table-input="${desk}"]`);
        const symbol = normalizeTrackedSymbol(input?.value || watchlistDraft[desk] || '');
        if (!symbol) {
          toast('Enter a Ticker', 'Enter a ticker such as NVDA.');
          input?.focus();
          return;
        }

        const ctx = captureSaveContext();
        const original = button.textContent;
        button.disabled = true;
        button.textContent = 'Adding…';
        watchlistDraft[desk] = symbol;
        try {
          await addSymbolToDesk(symbol, desk);
          selected[desk] = symbol;
          watchlistDraft[desk] = '';
          if (input) input.value = '';
          toast(`${symbol} Added`, `${deskCfg[desk].title} · historical data is hydrating if needed.`, 'success');
          updateChrome();
          render();
          restoreSaveContext(ctx);
        } catch (err) {
          watchlistDraft[desk] = symbol;
          toast('Could Not Add Symbol', err.message, 'error');
          render();
          restoreSaveContext(ctx);
        } finally {
          if (button.isConnected) {
            button.disabled = false;
            button.textContent = original;
          }
        }
      };
    });

    $$('[data-ai-add-desk]').forEach(button => {
      button.onclick = async () => {
        const [desk, rawSymbol] = String(button.dataset.aiAddDesk || '').split(':');
        const symbol = normalizeTrackedSymbol(rawSymbol);
        const ctx = captureSaveContext();
        button.disabled = true;
        try {
          await addSymbolToDesk(symbol, desk);
          selected[desk] = symbol;
          toast(`${symbol} Added`, `${deskCfg[desk].title} watchlist.`, 'success');
          render();
          restoreSaveContext(ctx);
        } catch (err) {
          if (button.isConnected) button.disabled = false;
          toast('Could Not Add Symbol', err.message, 'error');
          restoreSaveContext(ctx);
        }
      };
    });
  }

  function bindGlobalTrackedSymbolRemoval() {
    $$('[data-master-remove]').forEach(button => {
      button.onclick = async event => {
        event.preventDefault();
        event.stopPropagation();
        const ctx = captureSaveContext();
        try {
          await removeTrackedSymbolEverywhere(button.dataset.masterRemove, ctx);
        } catch (err) {
          toast('Global Remove Failed', err.message, 'error');
          restoreSaveContext(ctx);
        }
      };
    });

    $$('[data-desk-remove]').forEach(button => {
      const [, , rawSymbol] = String(button.dataset.deskRemove || '').split(':');
      const symbol = normalizeTrackedSymbol(rawSymbol);
      button.setAttribute('aria-label', `Remove ${symbol} from Tracked Symbols and all desks`);
      button.title = `Remove ${symbol} from all desks`;
      button.onclick = async event => {
        event.preventDefault();
        event.stopPropagation();
        if (button.dataset.busy === '1') return;
        button.dataset.busy = '1';
        button.disabled = true;
        const ctx = captureSaveContext();
        try {
          await removeTrackedSymbolEverywhere(symbol, ctx);
        } catch (err) {
          if (button.isConnected) {
            button.disabled = false;
            button.dataset.busy = '0';
          }
          toast('Global Remove Failed', err.message, 'error');
          restoreSaveContext(ctx);
        }
      };
    });
  }

  bindDynamic = function bindDynamicV186() {
    baseBindDynamic();
    bindCanonicalDeskAdds();
    bindGlobalTrackedSymbolRemoval();
  };

  window.__v186WatchlistContracts = Object.freeze({
    addSymbolToDesk,
    bindCanonicalDeskAdds,
    removeTrackedSymbolEverywhere,
    bindGlobalTrackedSymbolRemoval
  });
})();
