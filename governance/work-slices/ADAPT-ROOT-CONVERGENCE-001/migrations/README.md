# #73 migration waves

Each JSON file in this directory is an exact old-to-new repository migration ledger for one reorganization wave. The canonical repository migration gate composes these ledgers with the permanent `governance/repository-migrations.json` registry before validating stale references, file modes, executable evidence, test identities, and root ownership.

Wave ledgers are evidence, not aliases. Old paths must be absent after a completed move, temporary aliases remain prohibited unless explicitly governed, and historical Stable evidence must not be rewritten.
