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


## M2.3 CSRF and reverse-proxy hardening

Cookie-authenticated state-changing requests require the CSRF token associated with the server-side session. Bearer-authenticated API clients are not subject to browser CSRF validation.

BotWeb ignores `X-Forwarded-Proto` by default. Set `BOTWEB_TRUST_PROXY=1` only when BotWeb is directly behind a trusted reverse proxy that overwrites forwarded headers. In that mode, `X-Forwarded-Proto: https` allows BotWeb to mark session cookies `Secure` when TLS terminates at the proxy. Never enable proxy trust when clients can reach BotWeb directly and supply their own forwarded headers.


## M2.4 session lifecycle and audit

`GET /api/v1/session` exposes the current role and session lifetime to the authenticated UI. It never exposes the session identifier. Cookie sessions also receive their CSRF token so same-origin browser code can authorize state-changing requests.

Operators can rotate a cookie session through `POST /api/v1/session/rotate`; rotation invalidates the old session identifier, generates a new identifier and CSRF token, and renews the 12-hour expiry. Logout revokes the server-side session.

Security audit events record session rotation/logout event types and roles only. BotWeb MUST NOT log bearer tokens, session identifiers, cookies, or CSRF tokens.


## M2.5 authenticated management UI

The browser frontend bootstraps its authenticated role and CSRF token from `GET /api/v1/session`. Cookie-authenticated management POSTs automatically include `X-CSRF-Token`.

Viewer sessions retain read-only fleet, bot, channel, module, metrics, logs, BotAI and safe configuration views, but write controls such as channel join/part and module lifecycle actions are not rendered. Operator sessions receive those controls only when the managed bot also advertises the corresponding PBMP capability.

The unauthenticated loopback development mode retains the pre-auth local management behavior. Authentication and authorization remain BotWeb concerns and do not change the standalone-first PBMP bot contract.


## M2.6 session UX hardening

Authenticated browser sessions now surface their effective role in the header. Viewer sessions are explicitly marked `read-only`.

API responses with HTTP 401 are treated as an expired/revoked browser session and redirect the browser to `/login`. Operator sessions can rotate themselves from the header; the UI then reloads session metadata so subsequent writes use the newly generated CSRF token.

The rotation control is not shown to viewers or in unauthenticated loopback development mode.


## M2.7 management audit observability

State-changing PBMP management operations emit minimal audit events after authorization. Events contain only the operation name, BotWeb bot ID, result, verified authentication method and verified role.

The role and authentication method are attached to the request context by the authentication middleware after credentials have been validated; audit code does not re-parse or re-authenticate credentials.

Audit events intentionally exclude channel names, network names, module identifiers, reasons, request payloads, IRC content, bearer tokens, cookies, session identifiers and CSRF tokens.


## M2.8 structured audit events

BotWeb emits audit records as single-line JSON using the versioned schema `botweb.audit.v1`. The stable fields are `schema`, `event`, optional `bot`, `result`, optional `auth`, and optional `role`.

Management operations and session lifecycle events use the same format. This makes the stream suitable for structured log collectors such as Vector, Fluent Bit and Loki pipelines without parsing human-oriented log text.

The M2.7 privacy boundary remains part of the schema contract: audit records do not add tokens, cookies, session IDs, CSRF values, request payloads, channel/network names, module IDs or IRC content.


## M2.9 operational metrics

`GET /metrics` exposes a small Prometheus-compatible operational surface for BotWeb itself. It currently reports the number of configured bots, PBMP request totals split only by bounded `result=ok|error`, and aggregate PBMP request duration.

Metrics intentionally avoid bot IDs, network/channel names, module IDs, users, authentication material and other high-cardinality or sensitive labels. Per-bot operational detail remains available through the authenticated BotWeb/PBMP views rather than Prometheus labels.


## M2.10 liveness and readiness

`GET /healthz` is the liveness endpoint: it reports that the BotWeb HTTP process is alive and does not depend on bot availability.

`GET /readyz` is the readiness endpoint. For configured Unix/PBMP bots, BotWeb probes `pbmp.info` and reports the number configured, checked and reachable. Readiness returns HTTP 200 when at least one configured Unix/PBMP bot is reachable, and HTTP 503 when bots are configured but none of the supported PBMP endpoints respond. An intentionally empty registry is ready.

This split is suitable for container and orchestrator health checks without turning a transient bot outage into a BotWeb process restart.
