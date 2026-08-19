#!/usr/bin/env python3
from __future__ import annotations

from pathlib import Path
import sys

ROOT = Path(__file__).resolve().parents[2]
QUALIFIED = ROOT / ".github" / "workflows" / "ci-qualified.yml"
FAST = ROOT / ".github" / "workflows" / "ci-fast.yml"
WEBKIT_TEST = ROOT / "tools" / "ci" / "webkit_targeted_test.py"


def fail(errors: list[str]) -> int:
    print("DE.PULSE browser risk routing gate: FAIL", file=sys.stderr)
    for error in errors:
        print(f" - {error}", file=sys.stderr)
    return 1


def main() -> int:
    errors: list[str] = []
    qualified = QUALIFIED.read_text(encoding="utf-8")
    fast = FAST.read_text(encoding="utf-8")
    webkit_test = WEBKIT_TEST.read_text(encoding="utf-8") if WEBKIT_TEST.is_file() else ""

    required = (
        "webkit_required: ${{ steps.resolve.outputs.webkit_required }}",
        "name: Primary Chrome behavior",
        "webkit:\n    name: Primary WebKit compatibility",
        "needs.context.outputs.webkit_required == 'true' || needs.context.outputs.lane == 'full' || needs.context.outputs.lane == 'browser'",
        "runs-on: macos-15",
        "python -m playwright install webkit",
        "python3 tools/ci/webkit_targeted_test.py",
        "needs: [context, ci-harness, portability, backend, renderer, browser, webkit]",
        "WEBKIT_REQUIRED: ${{ needs.context.outputs.webkit_required }}",
        "WEBKIT: ${{ needs.webkit.result }}",
        'if [ "$WEBKIT_REQUIRED" = true ] || [ "$LANE" = full ] || [ "$LANE" = browser ]; then',
        'test "$WEBKIT" = skipped',
        "Chrome+WebKit primary browser policy",
    )
    for token in required:
        if token not in qualified:
            errors.append(f"Qualified Chrome+WebKit primary contract missing: {token}")

    forbidden = (
        "matrix:\n        browser: [chromium, webkit]",
        "matrix:\n        browser: [chrome, webkit]",
        "firefox.launch",
        "playwright install --with-deps firefox",
        "playwright install --with-deps webkit",
    )
    for token in forbidden:
        if token in qualified:
            errors.append(f"secondary-engine/default-matrix or dependency amplification prohibited: {token}")

    if "playwright install" in fast or "webkit_targeted_test.py" in fast:
        errors.append("Fast must never install/run browser engines")

    test_required = (
        "p.webkit.launch(headless=True)",
        "Global Remove Failed",
        "aria-pressed",
        "settings-persistent-save",
        "text-align:center!important",
        "Chrome and WebKit are the primary browser qualification engines",
    )
    for token in test_required:
        if token not in webkit_test:
            errors.append(f"primary WebKit proof missing compatibility assertion/contract: {token}")

    if errors:
        return fail(errors)

    print("DE.PULSE browser risk routing gate: PASS")
    print("Chrome primary broad behavioral qualification: PASS")
    print("WebKit co-primary macOS compatibility qualification: PASS")
    print("full/browser/UI-risk WebKit requirement: PASS")
    print("WebKit apt dependency amplification prevention: PASS")
    print("Fast/backend-only unnecessary browser suppression: PASS")
    print("other browser engines remain secondary by default: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
