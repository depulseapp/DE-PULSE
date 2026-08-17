#!/usr/bin/env python3
"""Browser/content regression for v18.5.1 first-run Stable profile wording."""
import os
import re
from pathlib import Path

from playwright.sync_api import sync_playwright

ROOT = Path(__file__).resolve().parents[2]
SETUP = ROOT / "renderer" / "setup.html"
AUTH_SURFACES = [ROOT / "renderer" / "setup.html", ROOT / "renderer" / "login.html"]
EXPECTED = "Existing DE.PULSE Stable profile and data are preserved"
STALE = "Existing v17 profile is preserved"


def main() -> None:
    setup = SETUP.read_text(encoding="utf-8")
    assert EXPECTED in setup, "durable Stable profile continuity copy is missing"
    assert STALE not in setup, "stale v17 profile wording remains on setup surface"

    # Current interactive authentication surfaces may not claim a historical app
    # version as the user's current profile/security identity.
    historical_claim = re.compile(r"(?:Existing\s+v\d+(?:\.\d+)*\s+profile|v\d+(?:\.\d+)*\s+security\s+migration)", re.I)
    for path in AUTH_SURFACES:
        text = path.read_text(encoding="utf-8")
        match = historical_claim.search(text)
        assert not match, f"historical-version auth copy in {path.name}: {match.group(0)!r}"

    with sync_playwright() as p:
        kwargs = {"headless": True}
        chrome = os.environ.get("CHROME_BIN", "").strip()
        if chrome:
            assert Path(chrome).is_file(), f"CHROME_BIN missing: {chrome}"
            kwargs["executable_path"] = chrome
        browser = p.chromium.launch(**kwargs)
        pg = browser.new_page(viewport={"width": 1100, "height": 760})
        # set_content avoids invoking auth network behavior; this test is solely
        # the actual user-visible first-run copy contract.
        pg.set_content(setup)
        assert pg.get_by_text(EXPECTED, exact=True).is_visible()
        assert pg.get_by_text("Owner security setup", exact=True).is_visible()
        assert pg.get_by_text(STALE, exact=True).count() == 0
        assert pg.get_by_text(re.compile(r"v\d+(?:\.\d+)* security migration", re.I)).count() == 0
        browser.close()

    print("PASS: first-run auth surface uses durable DE.PULSE Stable profile/data continuity wording with no stale historical-version claim.")


if __name__ == "__main__":
    main()
