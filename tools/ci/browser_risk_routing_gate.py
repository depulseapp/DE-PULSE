#!/usr/bin/env python3
from __future__ import annotations

from pathlib import Path
import sys

ROOT = Path(__file__).resolve().parents[2]
QUALIFIED = ROOT / ".github" / "workflows" / "ci-qualified.yml"
FAST = ROOT / ".github" / "workflows" / "ci-fast.yml"
WEBKIT_BROWSER_TEST = ROOT / "tools" / "ci" / "webkit_browser_test.py"
WEBKIT_LEGACY_TEST = ROOT / "tools" / "ci" / "webkit_targeted_test.py"


def fail(errors: list[str]) -> int:
    print("DE.PULSE browser risk routing gate: FAIL", file=sys.stderr)
    for error in errors:
        print(f" - {error}", file=sys.stderr)
    return 1


def require_tokens(errors: list[str], label: str, text: str, tokens: tuple[str, ...]) -> None:
    for token in tokens:
        if token not in text:
            errors.append(f"{label} missing: {token}")


def main() -> int:
    errors: list[str] = []
    qualified = QUALIFIED.read_text(encoding="utf-8")
    fast = FAST.read_text(encoding="utf-8")
    browser_test = WEBKIT_BROWSER_TEST.read_text(encoding="utf-8") if WEBKIT_BROWSER_TEST.is_file() else ""
    legacy_test = WEBKIT_LEGACY_TEST.read_text(encoding="utf-8") if WEBKIT_LEGACY_TEST.is_file() else ""

    require_tokens(
        errors,
        "Planner v3 browser routing contract",
        qualified,
        (
            "chrome_required: ${{ steps.resolve.outputs.chrome_required }}",
            "webkit_required: ${{ steps.resolve.outputs.webkit_required }}",
            "native_macos_required: ${{ steps.resolve.outputs.native_macos_required }}",
            "native_windows_required: ${{ steps.resolve.outputs.native_windows_required }}",
            "name: Primary Chrome behavior",
            "if: ${{ needs.context.outputs.chrome_required == 'true' }}",
            "webkit:\n    name: Primary WebKit browser compatibility",
            "if: ${{ needs.context.outputs.webkit_required == 'true' }}",
            "python -m playwright install webkit",
            "python3 tools/ci/webkit_browser_test.py",
            "name: Qualified macOS native lifecycle rehearsal",
            "if: ${{ needs.context.outputs.native_macos_required == 'true' }}",
            "bash tools/release/native_macos.sh",
            "name: Qualified Windows native runtime rehearsal",
            "if: ${{ needs.context.outputs.native_windows_required == 'true' }}",
            "tools/release/native_windows.ps1",
            "needs: [context, ci-harness, portability, backend, db-integration, security-rights, renderer, browser, webkit, native-macos, native-windows]",
            "CHROME_REQUIRED: ${{ needs.context.outputs.chrome_required }}",
            "WEBKIT_REQUIRED: ${{ needs.context.outputs.webkit_required }}",
            "NATIVE_MACOS_REQUIRED: ${{ needs.context.outputs.native_macos_required }}",
            "NATIVE_WINDOWS_REQUIRED: ${{ needs.context.outputs.native_windows_required }}",
            "require_if \"$CHROME_REQUIRED\" \"$CHROME\" chrome",
            "require_if \"$WEBKIT_REQUIRED\" \"$WEBKIT\" webkit",
            "require_if \"$NATIVE_MACOS_REQUIRED\" \"$NATIVE_MACOS\" native-macos",
            "require_if \"$NATIVE_WINDOWS_REQUIRED\" \"$NATIVE_WINDOWS\" native-windows",
            "Planner v3 selected jobs=",
        ),
    )

    forbidden = (
        "matrix:\n        browser: [chromium, webkit]",
        "matrix:\n        browser: [chrome, webkit]",
        "firefox.launch",
        "playwright install --with-deps firefox",
        "playwright install --with-deps webkit",
        "run: python3 tools/ci/webkit_targeted_test.py",
    )
    for token in forbidden:
        if token in qualified:
            errors.append(f"secondary-engine/default-matrix or browser/native ownership violation: {token}")

    if "playwright install" in fast or "webkit_browser_test.py" in fast or "webkit_targeted_test.py" in fast:
        errors.append("Fast must never install/run browser engines")

    require_tokens(
        errors,
        "primary WebKit browser-only proof",
        browser_test,
        (
            "p.webkit.launch(headless=True)",
            "watchlist_contract(page)",
            "settings_layout_contract(page)",
            "browser-only compatibility evidence",
            "does not package or launch the native macOS app",
        ),
    )
    if "native_packaged_window_contract()" in browser_test or "native_macos.sh" in browser_test:
        errors.append("WebKit browser owner must not invoke native lifecycle packaging")

    # The legacy combined harness may remain temporarily for historical/current
    # release compatibility, but Qualified must no longer call it. Its native
    # assertion implementation remains evidence that ownership was separated,
    # not silently removed.
    require_tokens(
        errors,
        "legacy combined WebKit/native compatibility source",
        legacy_test,
        (
            "def native_packaged_window_contract()",
            "tools/release/native_macos.sh",
            "p.webkit.launch(headless=True)",
        ),
    )

    if errors:
        return fail(errors)

    print("DE.PULSE browser risk routing gate: PASS")
    print("Chrome primary broad behavioral qualification: PASS")
    print("WebKit co-primary browser-only compatibility qualification: PASS")
    print("macOS native lifecycle rehearsal has independent Qualified ownership: PASS")
    print("Windows native runtime rehearsal has independent Qualified ownership: PASS")
    print("Planner v3 browser/native evidence selection: PASS")
    print("WebKit dependency amplification prevention: PASS")
    print("Fast/backend-only unnecessary browser suppression: PASS")
    print("other browser engines remain secondary by default: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
