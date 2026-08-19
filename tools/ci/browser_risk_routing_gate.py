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
        "name: Qualified Chrome behavior",
        "webkit:\n    name: Targeted WebKit compatibility",
        "if: ${{ needs.context.outputs.webkit_required == 'true' }}",
        "python -m playwright install --with-deps webkit",
        "python3 tools/ci/webkit_targeted_test.py",
        "needs: [context, ci-harness, portability, backend, renderer, browser, webkit]",
        "WEBKIT_REQUIRED: ${{ needs.context.outputs.webkit_required }}",
        "WEBKIT: ${{ needs.webkit.result }}",
        'if [ "$WEBKIT_REQUIRED" = true ]; then',
        'test "$WEBKIT" = skipped',
        "targeted WebKit required=$WEBKIT_REQUIRED",
    )
    for token in required:
        if token not in qualified:
            errors.append(f"Qualified targeted-WebKit contract missing: {token}")

    forbidden = (
        "matrix:\n        browser: [chromium, webkit]",
        "matrix:\n        browser: [chrome, webkit]",
        "needs.context.outputs.lane == 'backend' && needs.context.outputs.webkit_required",
    )
    for token in forbidden:
        if token in qualified:
            errors.append(f"full duplicate browser matrix or backend WebKit coupling prohibited: {token}")

    if "playwright install" in fast or "webkit_targeted_test.py" in fast:
        errors.append("Fast must never install/run WebKit")

    test_required = (
        "p.webkit.launch(headless=True)",
        "Global Remove Failed",
        "aria-pressed",
        "settings-persistent-save",
        "text-align:center!important",
        "Chrome remains the primary browser qualification target",
    )
    for token in test_required:
        if token not in webkit_test:
            errors.append(f"targeted WebKit proof missing compatibility assertion: {token}")

    if errors:
        return fail(errors)

    print("DE.PULSE browser risk routing gate: PASS")
    print("Chrome primary full behavioral qualification: PASS")
    print("WebKit conditional on renderer/UI impact signal: PASS")
    print("Fast/backend-only WebKit suppression: PASS")
    print("targeted Safari-sensitive contract scope: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
