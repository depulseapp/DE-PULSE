#!/usr/bin/env python3
"""Browser proof for the v18.5.2 configurable display-name contract."""

import os
from pathlib import Path

from playwright.sync_api import sync_playwright

ROOT = Path(__file__).resolve().parents[2]
RENDERER = (ROOT / "renderer" / "renderer.js").read_text(encoding="utf-8")


def between(text: str, start: str, end: str, label: str) -> str:
    s = text.find(start)
    if s < 0:
        raise AssertionError(f"{label}: start anchor missing")
    e = text.find(end, s)
    if e < 0:
        raise AssertionError(f"{label}: end anchor missing")
    return text[s:e]


SYNC_IDENTITY = between(
    RENDERER,
    "function syncIdentityChrome()",
    "function userFacingHealth()",
    "identity chrome sync",
)
PROFILE_BINDING = between(
    RENDERER,
    "$('[data-auth-save-profile]')?.addEventListener",
    "$('[data-auth-set-password]')?.addEventListener",
    "display-name binding",
)

HARNESS = r"""
const $=(s,r=document)=>r.querySelector(s);
let authPrincipal={userId:'bootstrap-owner',username:'owner',displayName:'Local Owner',role:'OWNER',sessionId:'sid-1'};
window.__calls=[];window.__toasts=[];window.__renderCount=0;
async function api(path,payload){
  window.__calls.push({path,payload});
  if(path!=='/api/auth/profile')throw new Error('unexpected API '+path);
  return {ok:true,principal:{...authPrincipal,displayName:payload.displayName}};
}
function toast(title,msg='',tone=''){window.__toasts.push({title,msg,tone})}
function render(){window.__renderCount+=1}
"""


def main() -> None:
    assert "data-auth-display-name" in RENDERER
    assert "data-auth-save-profile" in RENDERER
    assert "/api/auth/profile" in PROFILE_BINDING
    assert "username <b>" in RENDERER
    assert "permissions remain unchanged" in RENDERER

    with sync_playwright() as p:
        kwargs = {"headless": True}
        chrome = os.environ.get("CHROME_BIN", "").strip()
        if chrome:
            assert Path(chrome).is_file(), f"CHROME_BIN missing: {chrome}"
            kwargs["executable_path"] = chrome
        browser = p.chromium.launch(**kwargs)
        page = browser.new_page(viewport={"width": 900, "height": 600})
        page.set_content(
            '<span id="identity-principal" title="Authenticated user"></span>'
            '<button id="identity-signout" hidden>SIGN OUT</button>'
            '<input data-auth-display-name maxlength="64">'
            '<button data-auth-save-profile>Save Display Name</button>'
        )
        page.add_script_tag(content=HARNESS)
        page.add_script_tag(content=SYNC_IDENTITY)
        page.add_script_tag(content="function bindProfile(){"+PROFILE_BINDING+"}")
        page.evaluate("syncIdentityChrome();bindProfile()")
        assert page.locator("#identity-principal").inner_text() == "Local Owner"

        page.locator("[data-auth-display-name]").fill("Deivaram Venkatachalapathy")
        page.locator("[data-auth-save-profile]").click()
        page.wait_for_function("window.__calls.length===1")
        page.wait_for_function("document.querySelector('#identity-principal').textContent==='Deivaram Venkatachalapathy'")

        result = page.evaluate(
            """() => ({
              calls:window.__calls,
              principal:authPrincipal,
              header:document.querySelector('#identity-principal').textContent,
              title:document.querySelector('#identity-principal').title,
              aria:document.querySelector('#identity-principal').getAttribute('aria-label'),
              signoutHidden:document.querySelector('#identity-signout').hidden,
              toasts:window.__toasts,
              renderCount:window.__renderCount
            })"""
        )
        assert result["calls"] == [{
            "path": "/api/auth/profile",
            "payload": {"displayName": "Deivaram Venkatachalapathy"},
        }], result
        assert result["principal"]["username"] == "owner", result
        assert result["principal"]["role"] == "OWNER", result
        assert result["header"] == "Deivaram Venkatachalapathy", result
        assert "Username: owner" in result["title"] and "Role: OWNER" in result["title"], result
        assert "OWNER" in result["aria"], result
        assert not result["signoutHidden"], result
        assert result["renderCount"] == 1, result
        assert any(x["title"] == "Display Name Updated" for x in result["toasts"]), result
        browser.close()

    print("PASS: Settings updates the visible display name immediately while username, OWNER role and session identity remain separate.")


if __name__ == "__main__":
    main()
