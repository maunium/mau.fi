---
title: August 2026 releases // Bandwidthmuks 2.0 and mascots
summary: A summary of the mautrix releases in August 2026
slug: 2026-08-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
---
This month's releases include the bridge changes announced last month, more
bandwidth-saving features for gomuks, and new mascots!

| Software          | Version                                                           |
|-------------------|-------------------------------------------------------------------|
| mautrix-gmessages | [v26.08](https://github.com/mautrix/gmessages/releases/v0.2608.0) |
| mautrix-whatsapp  | [v26.08](https://github.com/mautrix/whatsapp/releases/v0.2608.0)  |
| mautrix-telegram  | [v26.08](https://github.com/mautrix/telegram/releases/v0.2608.0)  |
| mautrix-linkedin  | [v26.08](https://github.com/mautrix/linkedin/releases/v0.2608.0)  |
| mautrix-twitter   | [v26.08](https://github.com/mautrix/twitter/releases/v0.2608.0)   |
| mautrix-discord   | [v0.7.7](https://github.com/mautrix/discord/releases/v0.7.7)      |
| mautrix-signal    | [v26.08](https://github.com/mautrix/signal/releases/v0.2608.0)    |
| mautrix-slack     | [v26.08](https://github.com/mautrix/slack/releases/v0.2608.0)     |
| mautrix-meta      | [v26.08](https://github.com/mautrix/meta/releases/v0.2608.0)      |
| meowlnir          | [v26.08](https://github.com/maunium/meowlnir/releases/v0.2608.0)  |
| gomuks            | [v26.08](https://github.com/gomuks/gomuks/releases/v0.2608.0)     |
| mautrix-go        | [v0.30.0](https://github.com/mautrix/go/releases/v0.30.0)         |
| go-util           | [v0.10.0](https://github.com/mautrix/go-util/releases/v0.10.0)    |

## gomuks
gomuks got a ton of small changes and a few bigger ones. The biggest one is the
new bandwidth-saving features, which the section below is about. Other notable
changes include:

* Support for [MSC4440] profile biographies.
* Options to pin low priority rooms to bottom and/or stop counting unreads in them.
* Support for clicking checkboxes in own messages.

[MSC4440]: https://github.com/matrix-org/matrix-spec-proposals/pull/4440

See the [changelog](https://github.com/gomuks/gomuks/releases/tag/v0.2608.0) for
the full list of changes.

There's also a [draft PR](https://github.com/gomuks/gomuks/pull/747) for
supporting push notifications from the homeserver to better enable running the
backend on mobile devices. Bundled backends on the Android and iOS wrappers are
probably still several months away at least, but [Nexus](https://git.federated.nexus/Nexus/nexus)
already has a bundled backend and will likely implement push notifications soon.

### Bandwidthmuks 2.0
A year after [the first wave of bandwidthmuks](https://mau.fi/blog/2025-07-mautrix-release/),
the second wave has landed. Server-sent events have been added as an alternative
to websockets. Compression is enabled by default in SSE and zstd is supported in
addition to deflate, which is usually more efficient. Outgoing requests are sent
as separate HTTP requests, which means they won't fail immediately if the socket
disconnects. Requests can also be auto-retried on network errors more easily,
though that part isn't implemented on the client yet.

Using separate HTTP requests does mean that you have to use HTTP v2 or higher
to avoid frequent TCP/TLS handshakes, so make sure your reverse proxy is
configured to support that. In the future, gomuks may automatically enable SSE
when it detects that HTTP/2 is available.

When server-sent events are enabled, the low-bandwidth mode switch will enable
room list caching. Instead of reloading all data from the backend, the frontend
will store the room list in IndexedDB and ask the server to send just the rooms
that changed since it was last online. On my (big) account, the full init sync
is 9mb uncompressed or 2mb with zstd. With room list caching, the init sync size
varies, but it's usually some tens or hundreds of kilobytes uncompressed.

## Bridges
As announced [last month](https://mau.fi/blog/2026-07-mautrix-release/), the old
mautrix-meta bridge can no longer be used for Instagram DMs. Meta also disabled
the old API a couple weeks ago, so the old bridge version won't work at all
anymore. The same repo contains a new mautrix-instagram bridge which replaces
that functionality. mautrix-meta is still used for Facebook Messenger.

As also mentioned last month, the WhatsApp bridge has switched to prefer LIDs
for all chats, including DMs. Phone number ghosts will automatically be replaced
with LID ghosts in all DM portal rooms after upgrading.

In the future, the Google Messages bridge will switch to stable chat identifiers
in the next release. The switch should remove the need to re-backfill chats when
relogining. To prepare for the migration, you must run the bridge on v26.08
before upgrading further.

## Mascots
In non-code news, the bridges and gomuks have mascots now.

![image](/blog/res/mautrix-and-gomuks.png)

The initial designs were commissioned from [@LouLubally](https://x.com/LouLubally).
I'll probably commission a sticker pack next and maybe some room avatars for
gomuks and Meowlnir. Let me know if you have other ideas for them.
