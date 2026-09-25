# BotWeb

BotWeb is the shared web control plane for PBMP-compatible bots, initially LuCa and Engo, with AmBot planned. Bots remain fully functional without BotWeb.

M1 provides the Go backend, multi-bot registry and PBMP/1 client. M1.1 adds the fleet dashboard. M1.2 adds per-bot detail routes and capability-driven management navigation without implementation-specific LuCa/Engo/AmBot frontend branches.

## Run

    cp bots.example.json bots.json
    go run ./cmd/botweb -registry bots.json

Open `http://127.0.0.1:8080/`. The dashboard shows registered bots, reachability/online state, implementation/version and networks. Selecting a bot opens its detail view with all advertised capabilities and navigation for optional management surfaces only when the bot advertises them.

The backend binds to `127.0.0.1:8080` by default. Put an authenticated TLS reverse proxy in front before exposing it beyond localhost. Security headers and a restrictive same-origin CSP are emitted by BotWeb itself.

Bot PBMP endpoints remain backend-only and are deliberately omitted from the browser-facing registry response.

License: MIT.
