---
title: February 2026 releases // Pop-out widgets in gomuks
summary: A summary of the mautrix releases in February 2026
slug: 2026-02-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
- Meowlnir
---
This month's releases include pop-out widgets in gomuks, policy list subscribing
commands in Meowlnir and more bugfixes in bridges. Also, [Go 1.26] was released,
so the minimum version for compiling any new releases is now 1.25. As always,
precompiled binaries have no dependencies and don't require installing Go.

[Go 1.26]: https://go.dev/blog/go1.26

| Software          | Version                                                           |
|-------------------|-------------------------------------------------------------------|
| mautrix-gmessages | [v26.02](https://github.com/mautrix/gmessages/releases/v0.2602.0) |
| mautrix-whatsapp  | [v26.02](https://github.com/mautrix/whatsapp/releases/v0.2602.0)  |
| mautrix-linkedin  | [v26.02](https://github.com/mautrix/linkedin/releases/v0.2602.0)  |
| mautrix-discord   | [v0.7.6](https://github.com/mautrix/discord/releases/v0.7.6)      |
| mautrix-signal    | [v26.02](https://github.com/mautrix/signal/releases/v0.2602.0)    |
| mautrix-slack     | [v26.02](https://github.com/mautrix/slack/releases/v0.2602.0)     |
| mautrix-meta      | [v26.02](https://github.com/mautrix/meta/releases/v0.2602.0)      |
| meowlnir          | [v26.02](https://github.com/maunium/meowlnir/releases/v0.2602.0)  |
| gomuks            | [v26.02](https://github.com/gomuks/gomuks/releases/v0.2602.0)     |
| mautrix-go        | [v0.26.3](https://github.com/mautrix/go/releases/v0.26.3)         |
| go-util           | [v0.9.6](https://github.com/mautrix/go-util/releases/v0.9.6)      |

## gomuks
The widget support in gomuks web was upgraded to include always on screen
widgets, which should be the last piece to actually usable Element Call in
gomuks: you can now safely switch between rooms without the call being
interrupted. It'll just pop out into a floating box that you can move or resize
however you want.

A bunch of smaller improvements were made as well, such as hiding events in JS
to enable accurate date change headers and redirecting read receipts to the
previous visible event. It also improves performance by not adding hidden events
to the DOM.

Terminal didn't get much changes other than a fix for the /reply command, which
was supposed to work in the previous release, but accidentally didn't.

The backend now has C FFI bindings, which can be used to embed gomuks into
different kinds of clients. For example, QuadRadical is building a Flutter
frontend which uses the C FFI: <https://git.federated.nexus/Henry-Hiles/nexus>.
Future gomuks web mobile or desktop clients with embedded backend support will
probably also use the C bindings.

## Bridges
Most bridge changes were bugfixes again, although Messenger got a new login flow
which uses username/password input rather than requiring a browser. It's still
experimental though, so for now it may be safer to stick to the old flows.

Discord had accumulated a bunch of small bugfixes and while the rewrite is
progressing, it's not ready yet, so I decided to make another hopefully final
release of the old bridge.

## Meowlnir
I didn't get around to reviewing the protections PR, but the policy list
subscribing commands were merged, so you can now manage a Meowlnir instance
without sending raw state events.
