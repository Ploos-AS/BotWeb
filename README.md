# BotWeb

BotWeb is the shared web control plane for PBMP-compatible bots, initially LuCa and Engo. It is separate from bot processes: bots remain fully functional without BotWeb.

M1 provides a small Go backend, a registry for multiple bot instances, a PBMP/1 Unix-socket client, capability/status HTTP API, tests, CI and an Alpine OCI image.

## Run

    cp bots.example.json bots.json
    go run ./cmd/botweb -registry bots.json

The backend binds to `127.0.0.1:8080` by default. Put an authenticated TLS reverse proxy in front before exposing it beyond localhost.

Initial API: `/healthz`, `/api/v1/bots`, and per-bot `info`, `capabilities`, `networks`, `protocol` endpoints.

Bot endpoint paths are backend-only and are deliberately omitted from the browser-facing registry response.

License: MIT.
