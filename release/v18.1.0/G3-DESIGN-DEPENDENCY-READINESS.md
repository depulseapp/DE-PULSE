# v18.1.0 G3 — Design / Dependency Readiness

v18.1 adds **zero external dependencies**. It reuses authenticated Principal/UserID, current persistence migrations, current SSE Hub and the shared engine. New-account administration/presence/session lifecycle remains v18.2. The TEST runtime is isolated as `PersonalMarketTerminal-v18.1.0-TEST` and first-run clones v18.0.6 Stable without modifying Stable.
