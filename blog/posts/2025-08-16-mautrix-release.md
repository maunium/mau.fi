---
title: August 2025 releases // Room v12 support, threadmuks and wasmuks
summary: A summary of the mautrix releases in August 2025
slug: 2025-08-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
- Meowlnir
---
Meowlnir and bridges received room v12 support as expected. The legacy
provisioning APIs in bridges were also deprecated. Bigger changes in gomuks
this month include jumping to events, a thread view and the start of a wasm
build.

Also, with the release of [Go 1.25], the minimum version for compiling any new
releases is now 1.24. As always, precompiled binaries have no dependencies and
don't require installing Go.

[Go 1.25]: https://go.dev/blog/go1.25

| Bridge/library    | Version                                                          |
|-------------------|------------------------------------------------------------------|
| mautrix-gmessages | [v0.6.5](https://github.com/mautrix/gmessages/releases/v0.6.5)   |
| mautrix-whatsapp  | [v0.12.4](https://github.com/mautrix/whatsapp/releases/v0.12.4)  |
| mautrix-twitter   | [v0.5.0](https://github.com/mautrix/twitter/releases/v0.5.0)     |
| mautrix-gvoice    | [v0.1.2](https://github.com/mautrix/gvoice/releases/v0.1.2)      |
| mautrix-signal    | [v0.8.6](https://github.com/mautrix/signal/releases/v0.8.6)      |
| mautrix-slack     | [v0.2.3](https://github.com/mautrix/slack/releases/v0.2.3)       |
| mautrix-meta      | [v0.5.3](https://github.com/mautrix/meta/releases/v0.5.3)        |
| mautrix-go        | [v0.25.0](https://github.com/mautrix/go/releases/v0.25.0)        |
| meowlnir          | [v0.7.0](https://github.com/maunium/meowlnir/releases/v0.7.0)    |
| go-util           | [v0.9.0](https://github.com/mautrix/go-util/releases/v0.9.0)     |

## Bridges
The legacy provisioning APIs (`/_matrix/provision/v1` and `../v2`) are now
formally deprecated in all megabridges and will be released in next month's
releases. Anyone using the legacy APIs should migrate to the `/v3`
[Megabridge provisioning API] immediately.

[Megabridge provisioning API]: https://spec.mau.fi/megabridge/

Room v12 support was added to megabridges and mautrix-discord, although not all
bridges got a release this month. The bridges will also follow tombstones
automatically to make room upgrades easier. If you use mautrix-discord for
public rooms, you may want to update to the main branch.

Internally, mautrix-go got rid of the gorilla toolkit which has been mostly
unmaintained for years (except for some broken releases). HTTP routers were
replaced with the standard library and websockets now use [coder/websocket].
The HTTP router change is a breaking change for some library users who assume
the router is gorilla/mux.

[coder/websocket]: https://github.com/coder/websocket

Other than that, there were minor bridging improvements in multiple bridges,
such as "delete for me" bridging on Signal and better placeholders for various
message types on WhatsApp. Signal also bumped their backup/transfer protocol
version, so some users may have noticed the lack of a history transfer option
when linking the bridge to Signal iOS. The latest release adds support for the
latest protocol version to restore history transfer.

## Meowlnir
Meowlnir received room v12 support as expected. Note that tombstones are not
followed automatically, so you have to re-protect/watch rooms if a protected
room or watched list is upgraded. There will likely be automated tombstone
handling in the future.

## gomuks

### Jumping to events and thread view
Jumping to arbitrary events is now possible using a separate modal. This
includes jumping to reply when someone replies to an old message as well as
clicking matrix.to links with an event ID. If the event is loaded in the visible
timeline, gomuks will still jump to it in place, but when it's not loaded, a
separate modal will open to show the event and context around it. Internally,
the modal calls the event `/context` API on the homeserver rather than using
the local database.

After I made the event context modal, I realized a thread view could be
implemented just as easily using the `/relations` endpoint, so I did that. When
you click on the reply box of a threaded reply, it'll open the thread view in
the right panel instead of jumping to the event.

### Wasmuks
As mentioned in my [very first post](https://mau.fi/blog/2024-h1-mautrix-updates/#high-level-client-framework),
before starting gomuks web I had looked into running it in wasm, but discarded
the idea due to performance concerns. However, while I want maximum performance
for my own account, that doesn't mean having a wasm build as an option would be
a bad idea, so I decided to see if it works at all.

After some updates to the proof-of-concept Go-JS-SQLite bridge I made last year,
I was able to hook it up to gomuks and compile the whole thing to wasm. The
performance is certainly worse than natively, but it appears usable on my test
account with a few dozen rooms.

Matrix e2ee was trivial thanks to goolm. Getting libolm would've been much more
complicated, as there's no cgo for wasm. That's also why the SQLite connector
has to bridge via JS. In the future, the wasm component model may allow more
direct C<->Go interop, but I'm not sure whether that stuff is going to be
supported in browsers any time soon.

The main missing bit is media, as the native version proxies media downloads
via the Go backend for authentication and decryption. I'll probably go for a
service worker, but still need to figure out how to have it talk to the Go
side for encryption keys and authentication.

**Update:** media has been implemented and a demo is now available at
<https://demo.gomuks.app>. It's still very fragile and will explode if you
hold it wrong.

### Other changes
* Room v12+ power level support in various places where it's relevant, like
  displaying moderation buttons and sorting the member list.
* Megolm sessions are now shared when you start typing rather than when you
  send the message.
* A mutex was added around sending to try to preserve outgoing message order
  better.
* Events flagged as soft-failed by Synapse will be rendered with a special
  background color. Server admins can opt into receiving soft-failed events
  using the `io.element.synapse.admin_client_config` account data event.
