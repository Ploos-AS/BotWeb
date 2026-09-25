# BotWeb architecture

BotWeb is a control plane, not part of an IRC bot runtime.

## Boundaries

1. Browser access terminates at BotWeb over HTTPS.
2. BotWeb backend owns credentials for registered PBMP endpoints.
3. Browsers never receive direct bot-management endpoints or credentials.
4. BotWeb discovers capabilities with PBMP instead of branching on LuCa/Engo implementation names.
5. A bot outage must not compromise or block management of other registered bots.

## Instance model

Each registered instance has a stable local ID, display name, PBMP endpoint, transport type and credential reference. Secrets are stored separately from ordinary configuration and are never returned through the frontend API.

## M0 decision

No frontend/backend framework is selected in M0. PBMP must stabilize enough to exercise two independent bot implementations before BotWeb commits to implementation-specific client abstractions.
