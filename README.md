# BotWeb

BotWeb is the shared web control plane for PBMP-compatible bots, initially LuCa and Engo, with AmBot planned. Bots remain fully functional without BotWeb.

## Standalone-first rule

BotWeb is strictly optional. It MUST NOT become a runtime dependency of LuCa, Engo, AmBot, or any other managed IRC bot. Bots MUST continue normal IRC operation when BotWeb is absent, disabled, unreachable, or removed. BotWeb manages capabilities exposed by bots; it does not own their IRC core, configuration authority required for startup, or essential runtime lifecycle.

M1 provides the Go backend, multi-bot registry and PBMP/1 client. M1.1 adds the fleet dashboard. M1.2 adds per-bot detail routes and capability-driven management navigation without implementation-specific LuCa/Engo/AmBot frontend branches.

## Run

    cp bots.example.json bots.json
    go run ./cmd/botweb -registry bots.json

Open `http://127.0.0.1:8080/`. The dashboard shows registered bots, reachability/online state, implementation/version and networks. Selecting a bot opens its detail view with all advertised capabilities and navigation for optional management surfaces only when the bot advertises them.

The backend binds to `127.0.0.1:8080` by default. Put an authenticated TLS reverse proxy in front before exposing it beyond localhost. Security headers and a restrictive same-origin CSP are emitted by BotWeb itself.

Bot PBMP endpoints remain backend-only and are deliberately omitted from the browser-facing registry response.

License: MIT.


M1.3 wires the first optional PBMP capability end-to-end: bots advertising `channels.list` get a real Channels view. Bots without the capability do not expose that view.

M1.6 adds a generic, read-only BotAI status view for bots advertising the optional PBMP `botai.status` capability. BotWeb never connects to BotAI directly; see `docs/M1.6-botai.md`.


## M2.1 authentication and roles

BotWeb supports optional bearer-token authentication at the web control-plane boundary.

- `BOTWEB_VIEWER_TOKEN` grants read-only access.
- `BOTWEB_OPERATOR_TOKEN` grants read and write access.
- `/healthz` remains unauthenticated for local/container health checks.
- With neither token configured, BotWeb retains the simple development mode but refuses a non-loopback `-listen` address.

Tokens protect BotWeb itself; they are not PBMP credentials and are never sent to managed bots. For remote deployments, configure at least one token and terminate TLS at an authenticated reverse proxy. BotWeb remains optional: managed bots continue operating independently when it is absent.


## M2.2 browser sessions

When authentication is configured, browser requests without a session are redirected to `/login`. A configured viewer or operator token can be exchanged for a short-lived server-side session. The browser cookie contains only a cryptographically random session identifier, is `HttpOnly`, uses `SameSite=Strict`, and is marked `Secure` when BotWeb itself receives HTTPS. Sessions expire after 12 hours and can be invalidated through `POST /logout`.

Bearer authentication remains supported for API clients and automation. Browser sessions are a BotWeb concern only and do not alter PBMP or managed-bot authentication.
