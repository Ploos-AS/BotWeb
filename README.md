# BotWeb

BotWeb is the shared web control plane for PBMP-compatible bots, initially LuCa and Engo.

BotWeb is intentionally separate from bot processes. Bots remain fully functional without it. The browser talks to the BotWeb backend; the backend talks to bots over authenticated PBMP transports.

## M0 architecture

    Browser --HTTPS--> BotWeb backend --PBMP--> LuCa / Engo / other bots

The backend maintains a registry of bot instances and discovers each instance's PBMP capabilities. The UI renders only supported functions.

## Planned M1

* bot registry
* PBMP/1 client
* capability-driven dashboard
* network/channel status
* authentication and roles
* live events/log view
* Alpine OCI image

See [docs/architecture.md](docs/architecture.md).

License: MIT.
