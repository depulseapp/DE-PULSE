'use strict';

// v18.5.1 live-market reconciliation layer.
// Quote ticks update the active page in place so hover/focus/selection/scroll
// survive streaming refreshes. Structural/non-quote events continue through the
// existing full-render path in renderer.js.
(function installDepulseLiveDomReconciler(){
  if (typeof scheduleLiveRender !== 'function') return;

  const originalScheduleLiveRender = scheduleLiveRender;
  const LIVE_DELAY_MS = 250;
  const DIRECT_KEY_ATTRS = [
    'data-live-key',
    'data-ticker-symbol',
    'data-master-symbol',
    'data-row-symbol',
    'data-watch-symbol',
    'data-candidate-symbol',
    'data-symbol',
    'data-research'
  ];

  function directSemanticKey(el){
    if (!el || el.nodeType !== 1) return '';
    if (el.id) return `id:${el.id}`;

    const symbol = el.getAttribute('data-symbol');
    const desk = el.getAttribute('data-desk') || el.getAttribute('data-desk-membership');
    if (symbol && desk) return `${el.tagName}:symbol:${symbol}:desk:${desk}`;

    for (const attr of DIRECT_KEY_ATTRS){
      const value = el.getAttribute(attr);
      if (value) return `${el.tagName}:${attr}:${value}`;
    }
    return '';
  }

  function semanticKey(node){
    if (!node || node.nodeType !== 1) return '';
    const direct = directSemanticKey(node);
    if (direct) return direct;

    // Rows/cards often carry their symbol on a child action rather than the
    // container. Use that identity for sibling reconciliation without changing
    // source markup or reordering a row under the pointer.
    if (['TR','ARTICLE','LI'].includes(node.tagName)){
      const marker = node.querySelector('[data-live-key],[data-ticker-symbol],[data-master-symbol],[data-row-symbol],[data-watch-symbol],[data-candidate-symbol],[data-symbol],[data-research]');
      const childKey = directSemanticKey(marker);
      if (childKey) return `${node.tagName}:child:${childKey}`;
    }
    return '';
  }

  function containsPinnedInteraction(el){
    if (!el || el.nodeType !== 1) return false;
    const active = document.activeElement;
    if (active && active !== document.body && el.contains(active)) return true;
    try {
      if (el.matches(':hover') || el.querySelector(':hover')) return true;
    } catch (_) {
      // :hover is universally supported in our Chromium/WebView targets, but a
      // selector failure must never force a destructive fallback.
    }
    return false;
  }

  function syncAttributes(current, desired){
    const keepValue = current === document.activeElement && /^(INPUT|TEXTAREA|SELECT)$/.test(current.tagName);
    const value = keepValue ? current.value : undefined;
    const selectionStart = keepValue && typeof current.selectionStart === 'number' ? current.selectionStart : null;
    const selectionEnd = keepValue && typeof current.selectionEnd === 'number' ? current.selectionEnd : null;

    for (const attr of [...current.attributes]){
      if (!desired.hasAttribute(attr.name)) current.removeAttribute(attr.name);
    }
    for (const attr of [...desired.attributes]){
      if (current.getAttribute(attr.name) !== attr.value) current.setAttribute(attr.name, attr.value);
    }

    if (keepValue){
      current.value = value;
      if (selectionStart !== null && typeof current.setSelectionRange === 'function'){
        current.setSelectionRange(selectionStart, selectionEnd);
      }
    }
  }

  function patchNode(current, desired){
    if (!current || !desired) return;
    if (current.nodeType !== desired.nodeType){
      if (!containsPinnedInteraction(current)) current.replaceWith(desired.cloneNode(true));
      return;
    }
    if (current.nodeType === Node.TEXT_NODE || current.nodeType === Node.COMMENT_NODE){
      if (current.nodeValue !== desired.nodeValue) current.nodeValue = desired.nodeValue;
      return;
    }
    if (current.nodeType !== Node.ELEMENT_NODE) return;
    if (current.tagName !== desired.tagName){
      if (!containsPinnedInteraction(current)) current.replaceWith(desired.cloneNode(true));
      return;
    }

    syncAttributes(current, desired);
    reconcileChildren(current, desired);
  }

  function uniqueKeyMap(nodes){
    const map = new Map();
    const duplicates = new Set();
    for (const node of nodes){
      const key = semanticKey(node);
      if (!key) continue;
      if (map.has(key)) duplicates.add(key);
      else map.set(key, node);
    }
    for (const key of duplicates) map.delete(key);
    return map;
  }

  function reconcileChildren(currentParent, desiredParent){
    const current = [...currentParent.childNodes];
    const desired = [...desiredParent.childNodes];
    const currentKeys = uniqueKeyMap(current);
    const desiredKeys = uniqueKeyMap(desired);
    const keyed = currentKeys.size > 0 || desiredKeys.size > 0;

    if (keyed){
      const desiredUnkeyed = desired.filter(node => !semanticKey(node));
      let unkeyedIndex = 0;

      // Never reorder keyed live rows/cards. Patch the matching semantic entity
      // where it already sits so a hovered row cannot jump under the pointer.
      for (const node of current){
        const key = semanticKey(node);
        if (key && desiredKeys.has(key)){
          patchNode(node, desiredKeys.get(key));
        } else if (!key && unkeyedIndex < desiredUnkeyed.length){
          patchNode(node, desiredUnkeyed[unkeyedIndex++]);
        }
      }

      // Quote ticks normally do not change structure. If an entity genuinely
      // appears/disappears, reconcile it only when this container is not under
      // active user interaction; structural events still use the full renderer.
      if (!containsPinnedInteraction(currentParent)){
        const desiredKeySet = new Set(desiredKeys.keys());
        for (const [key, node] of currentKeys){
          if (!desiredKeySet.has(key) && node.isConnected) node.remove();
        }
        const currentKeySet = new Set([...currentParent.childNodes].map(semanticKey).filter(Boolean));
        for (const node of desired){
          const key = semanticKey(node);
          if (key && !currentKeySet.has(key)) currentParent.appendChild(node.cloneNode(true));
        }
      }
      return;
    }

    const common = Math.min(current.length, desired.length);
    for (let i = 0; i < common; i++) patchNode(current[i], desired[i]);

    if (desired.length > current.length && !containsPinnedInteraction(currentParent)){
      for (let i = current.length; i < desired.length; i++) currentParent.appendChild(desired[i].cloneNode(true));
    } else if (current.length > desired.length){
      for (let i = current.length - 1; i >= desired.length; i--){
        if (!containsPinnedInteraction(current[i])) current[i].remove();
      }
    }
  }

  function renderLivePage(changedSymbol){
    if (typeof state === 'undefined' || typeof runtime === 'undefined' || !state || !runtime) return false;
    const main = document.getElementById('main');
    if (!main || typeof pageRenderer !== 'function') return false;

    const renderer = pageRenderer(page);
    if (typeof renderer !== 'function') return false;
    const html = renderer();
    if (typeof html !== 'string' || !html.trim()) return false;

    const scratch = document.createElement('main');
    scratch.innerHTML = html;

    const active = document.activeElement;
    const selectionStart = active && typeof active.selectionStart === 'number' ? active.selectionStart : null;
    const selectionEnd = active && typeof active.selectionEnd === 'number' ? active.selectionEnd : null;
    const windowX = window.scrollX;
    const windowY = window.scrollY;
    const mainLeft = main.scrollLeft;
    const mainTop = main.scrollTop;

    reconcileChildren(main, scratch);

    // DOM identity is preserved by reconciliation. Explicitly restore positions
    // as a guard against layout-side scroll anchoring in Chromium/WebView.
    if (main.scrollLeft !== mainLeft) main.scrollLeft = mainLeft;
    if (main.scrollTop !== mainTop) main.scrollTop = mainTop;
    if (window.scrollX !== windowX || window.scrollY !== windowY) window.scrollTo(windowX, windowY);
    if (active && active.isConnected && document.activeElement !== active) active.focus({preventScroll:true});
    if (active && active.isConnected && selectionStart !== null && typeof active.setSelectionRange === 'function'){
      active.setSelectionRange(selectionStart, selectionEnd);
    }

    if (typeof updateChrome === 'function') updateChrome(changedSymbol);
    if (typeof syncLivePriorityHints === 'function') syncLivePriorityHints();
    return true;
  }

  function incrementalScheduleLiveRender(changedSymbol=''){
    const symbol = String(changedSymbol || '').trim().toUpperCase();

    // Empty-symbol events can represent news, scanner, configuration or other
    // structural changes. Preserve the authoritative existing full-render path.
    if (!symbol) return originalScheduleLiveRender(changedSymbol);

    if (typeof perfMetrics !== 'undefined') perfMetrics.renderRequests++;
    if (typeof quoteAffectsPage === 'function' && !quoteAffectsPage(symbol)){
      if (typeof perfMetrics !== 'undefined') perfMetrics.renderSkipped++;
      return;
    }

    clearTimeout(window.__renderTimer);
    window.__renderTimer = setTimeout(() => {
      if (typeof watchlistEditInProgress === 'function' && watchlistEditInProgress()){
        incrementalScheduleLiveRender(symbol);
        return;
      }
      if (typeof perfMetrics !== 'undefined') perfMetrics.renderExecuted++;
      try {
        if (!renderLivePage(symbol)) originalScheduleLiveRender('');
      } catch (error){
        console.warn('DE.PULSE incremental live render fell back to full render:', error);
        originalScheduleLiveRender('');
      }
    }, LIVE_DELAY_MS);
  }

  scheduleLiveRender = incrementalScheduleLiveRender;
  globalThis.__DEPULSE_LIVE_DOM__ = Object.freeze({
    version: '18.5.1',
    renderLivePage,
    semanticKey
  });
})();
