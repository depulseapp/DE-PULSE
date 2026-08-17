(() => {
  'use strict';

  const baseBindDynamic = bindDynamic;

  deskMembershipStrip = function deskMembershipStripV1851(sym, currentDesk = '') {
    const target = String(sym || '').toUpperCase();
    const current = String(currentDesk || '').toLowerCase();
    const items = [['day', 'DAY'], ['swing', 'SWING'], ['long', 'LONG']];
    return `<span class="desk-membership-strip" aria-label="Desk membership for ${esc(target)}">${items.map(([kind, label]) => {
      const active = (deskWL(kind).symbols || []).includes(target);
      const isCurrent = active && kind === current;
      const deskName = kind === 'long' ? 'Long-Term' : kind === 'day' ? 'Day Trade' : 'Swing';
      const title = isCurrent
        ? `${label} is the current desk for ${target}`
        : active
          ? `${target} is in ${deskName} · click to remove when another desk remains`
          : `Add ${target} to ${deskName}`;
      return `<button type="button" class="desk-membership-pill ${active ? 'active' : ''}${isCurrent ? ' current-desk' : ''}" aria-pressed="${active ? 'true' : 'false'}"${isCurrent ? ' aria-current="true"' : ''} data-desk-membership="${kind}:${esc(target)}" title="${esc(title)}">${label}${isCurrent ? '<span class="desk-membership-current" aria-hidden="true">CURRENT</span><span class="sr-only"> current desk</span>' : ''}</button>`;
    }).join('')}</span>`;
  };

  async function removeTrackedSymbolEverywhere(symbol, ctx) {
    const sym = normalizeSymbol(symbol);
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
      const symbol = normalizeSymbol(rawSymbol);
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

  bindDynamic = function bindDynamicV1851() {
    baseBindDynamic();
    bindGlobalTrackedSymbolRemoval();
  };

  window.__v1851WatchlistContracts = Object.freeze({
    removeTrackedSymbolEverywhere,
    bindGlobalTrackedSymbolRemoval
  });
})();
