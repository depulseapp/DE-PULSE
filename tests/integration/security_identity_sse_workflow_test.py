#!/usr/bin/env python3
"""Real authenticated identity/workspace/SSE workflow for v18 T3 closure."""
from __future__ import annotations

import http.client
import http.cookiejar
import json
import os
from pathlib import Path
import re
import subprocess
import tempfile
import time
import urllib.error
import urllib.request

ROOT = Path(__file__).resolve().parents[2]
BIN = Path(os.environ.get("DEPULSE_SECURITY_TEST_BINARY", str(Path(tempfile.gettempdir()) / "depulse-v18-security-workflow")))
OWNER_PASSWORD = "DE.PULSE v18 security workflow owner passphrase"
TEMP_PASSWORD = "DE.PULSE v18 temporary analyst passphrase"
USER_PASSWORD = "DE.PULSE v18 permanent analyst passphrase"


def make_opener():
    jar = http.cookiejar.CookieJar()
    return jar, urllib.request.build_opener(urllib.request.HTTPCookieProcessor(jar))


def csrf(jar) -> str:
    for cookie in jar:
        if cookie.name == "depulse_csrf":
            return cookie.value
    return ""


def request(opener, jar, base: str, path: str, payload=None, method=None, expect=200):
    data = None if payload is None else json.dumps(payload).encode()
    req = urllib.request.Request(base + path, data=data, method=method or ("POST" if data is not None else "GET"))
    req.add_header("Connection", "close")
    if data is not None:
        req.add_header("Content-Type", "application/json")
        token = csrf(jar)
        if token:
            req.add_header("X-DE-PULSE-CSRF", token)
    for attempt in range(2):
        try:
            with opener.open(req, timeout=10) as resp:
                body = resp.read()
                code = resp.status
                ctype = resp.headers.get("content-type", "")
            break
        except urllib.error.HTTPError as exc:
            body = exc.read()
            code = exc.code
            ctype = exc.headers.get("content-type", "")
            break
        except (ConnectionResetError, BrokenPipeError, http.client.RemoteDisconnected):
            if attempt:
                raise
            time.sleep(0.05)
    if code != expect:
        raise AssertionError(f"{path}: expected {expect}, got {code}: {body[:300]!r}")
    if "json" in ctype or path.startswith("/api/"):
        try:
            return json.loads(body)
        except Exception:
            return body.decode(errors="replace")
    return body.decode(errors="replace")


def main() -> None:
    subprocess.run(["go", "build", "-o", str(BIN), "."], cwd=ROOT, check=True)
    profile = Path(tempfile.mkdtemp(prefix="depulse-v18-security-"))
    log = profile / "app.log"
    env = os.environ.copy()
    env["XDG_CONFIG_HOME"] = str(profile / "cfg")
    env["DEPULSE_HEADLESS"] = "1"
    with log.open("wb") as out:
        proc = subprocess.Popen([str(BIN)], cwd=ROOT, env=env, stdout=out, stderr=subprocess.STDOUT)
    try:
        base = None
        for _ in range(120):
            if log.exists():
                match = re.search(r"Local terminal: (http://127\.0\.0\.1:\d+/)", log.read_text(errors="ignore"))
                if match:
                    base = match.group(1).rstrip("/")
                    break
            if proc.poll() is not None:
                raise RuntimeError(log.read_text(errors="ignore"))
            time.sleep(0.05)
        if not base:
            raise RuntimeError("app did not publish local URL")

        owner_jar, owner = make_opener()
        request(owner, owner_jar, base, "/")
        status = request(owner, owner_jar, base, "/api/auth/status")
        assert status.get("authenticated") is True and status.get("bootstrapRequired") is True
        secured = request(owner, owner_jar, base, "/api/auth/set-password", {"password": OWNER_PASSWORD})
        assert secured.get("ok") is True

        reauth = request(owner, owner_jar, base, "/api/auth/reauth", {"password": OWNER_PASSWORD})
        assert reauth.get("ok") is True and reauth.get("recentAuthentication") is True
        rotated = request(owner, owner_jar, base, "/api/auth/rotate", {})
        assert rotated.get("ok") is True and rotated.get("principal", {}).get("role") in ("OWNER", "SUPER_OWNER")

        admin = request(owner, owner_jar, base, "/api/admin/identity")
        assert isinstance(admin.get("users"), list) and isinstance(admin.get("sessions"), list)
        created = request(owner, owner_jar, base, "/api/admin/users/create", {
            "username": "qa.analyst",
            "displayName": "QA Analyst",
            "role": "USER",
            "temporaryPassword": TEMP_PASSWORD,
        }, expect=201)
        assert created.get("ok") is True and created.get("user", {}).get("username") == "qa.analyst"

        user_jar, user = make_opener()
        login = request(user, user_jar, base, "/api/auth/login", {"username": "qa.analyst", "password": TEMP_PASSWORD})
        assert login.get("ok") is True and login.get("principal", {}).get("role") == "USER"
        changed = request(user, user_jar, base, "/api/auth/set-password", {"password": USER_PASSWORD})
        assert changed.get("ok") is True

        owner_before = request(owner, owner_jar, base, "/api/bootstrap")
        user_before = request(user, user_jar, base, "/api/bootstrap")
        assert owner_before.get("principal", {}).get("userId") != user_before.get("principal", {}).get("userId")
        request(user, user_jar, base, "/api/ui/ticker", {"symbol": "AMD"})
        user_after = request(user, user_jar, base, "/api/bootstrap")
        owner_after = request(owner, owner_jar, base, "/api/bootstrap")
        assert user_after["state"]["ui"]["selectedTicker"] == "AMD"
        assert owner_after["state"]["ui"]["selectedTicker"] != "AMD" or owner_before["state"]["ui"]["selectedTicker"] == "AMD"

        request(user, user_jar, base, "/api/admin/identity", expect=403)
        request(user, user_jar, base, "/api/cache/clear", {}, expect=403)

        stream_req = urllib.request.Request(base + "/api/events", method="GET")
        stream_req.add_header("Connection", "close")
        with user.open(stream_req, timeout=5) as stream:
            line = stream.readline().decode(errors="replace")
        assert line.startswith("data: ") and '"type":"bootstrap"' in line
        assert '"principal"' not in line, "SSE bootstrap must remain scoped state/runtime data rather than leak credential/session records"

        logged_out = request(user, user_jar, base, "/api/auth/logout", {})
        assert logged_out.get("ok") is True
        post_logout = request(user, user_jar, base, "/api/auth/status")
        assert post_logout.get("authenticated") is False
        relogin = request(user, user_jar, base, "/api/auth/login", {"username": "qa.analyst", "password": USER_PASSWORD})
        assert relogin.get("ok") is True and relogin.get("principal", {}).get("role") == "USER"

        print("SECURITY / IDENTITY / WORKSPACE / SSE WORKFLOW: PASS")
    finally:
        try:
            if proc.poll() is None:
                proc.terminate()
                proc.wait(timeout=3)
        except Exception:
            try:
                proc.kill()
            except Exception:
                pass


if __name__ == "__main__":
    main()
