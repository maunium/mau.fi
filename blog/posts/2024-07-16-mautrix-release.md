---
title: July 2024 releases
summary: A summary of the mautrix releases in July 2024
slug: 2024-07-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
---
As always, this month's releases happened on the 16th. Highly unusually, this
month's releases included most of the bridges, including the Python ones that
don't follow release schedules. This is because all bridges needed to be
updated to gain support for authenticated media added in [Matrix v1.11]
(more specifically [MSC3916]).

[Matrix v1.11]: https://matrix.org/blog/2024/06/20/matrix-v1.11-release/
[MSC3916]: https://github.com/matrix-org/matrix-spec-proposals/pull/3916

| Bridge             | Version                                                             |
|--------------------|---------------------------------------------------------------------|
| mautrix-whatsapp   | [v0.10.9](https://github.com/mautrix/whatsapp/releases/tag/v0.10.9) |
| mautrix-meta       | [v0.3.2](https://github.com/mautrix/meta/releases/tag/v0.3.2)       |
| mautrix-discord    | [v0.7.0](https://github.com/mautrix/discord/releases/tag/v0.7.0)    |
| mautrix-signal     | [v0.6.3](https://github.com/mautrix/signal/releases/tag/v0.6.3)     |
| mautrix-gmessages  | [v0.4.3](https://github.com/mautrix/gmessages/releases/tag/v0.4.3)  |
| mautrix-telegram   | [v0.15.2](https://github.com/mautrix/telegram/releases/tag/v0.15.2) |
| mautrix-twitter    | [v0.1.8](https://github.com/mautrix/twitter/releases/tag/v0.1.8)    |
| mautrix-googlechat | [v0.5.2](https://github.com/mautrix/googlechat/releases/tag/v0.5.2) |

## Discord relays & MSC3916
For most of the bridges, MSC3916 will automatically be supported after upgrading
and there's no need to take any other action. However, Discord is a bit of a
special case: when using webhooks for relaying, mautrix-discord sends media repo
links to Discord for Matrix user avatars.

To solve this, the latest version of the bridge will use links to the bridge
itself instead of the media repo. The links will then proxy avatar downloads
with authentication.

For server administrators, this means some config changes:
The `homeserver` → `public_address` config option is unnecessary and has been
removed. Instead, there's a new `bridge` → `public_address` option, which must
be set to an address that is routed to the bridge's appservice port.

There's also a new `avatar_proxy_key` option, which will be automatically
generated on startup if it's not set. That key is used to sign media download
URLs to ensure the bridge can't be used as an unauthenticated proxy.

When the fields are set, the bridge will send URLs of the form
`<public address>/mautrix-discord/avatar/{server}/{id}/{checksum}` to Discord.

This also happens to be the first release of mautrix-discord with the new direct
media system, but I assume anyone who cares about direct media was already using
main branch builds. mautrix-discord is still on an old mautrix-go version, so
further releases are unlikely before the megabridge rewrite.

## Megabridge progress

### Signal
The v2 Signal bridge is nearing feature parity with the old bridge. It already
includes a mostly-seamless migration path: you can simply drop the new bridge
in place of the old one. It will automatically migrate the config file and
database while preserving all existing logins, rooms and messages.

There is one notable exception: in the old bridge, unencrypted rooms did not
include the bridge bot, while in the new one, all rooms include the bot. The
migration does not yet add the bot to all rooms, so it's only safe to use if
you have end-to-bridge encryption enabled.

You may want to double-check the config after migration. If you prevented the
bridge from writing to the config file, you'll have to temporarily allow it or
manually migrate the config, as the config migration system will likely not be
there forever.

### Slack
The Slack bridge was the most behind in terms of general bridge quality, and I
didn't get to finish my original refactor, so I picked it as the next target
for megabridge rewriting. The rewrite was started last week and the bridge is
already more or less functional. The migration system even adds the bridge bot
to DM rooms (which the Signal bridge migration doesn't support yet).

A release of the new bridge will likely happen next month.

### Other bridges
The Telegram and Meta bridges also have rewrites in progress, but are expected
to take a bit longer.

## Fixed replies
Last month's releases included an unfortunate late change to `m.relates_to`
handling in encrypted messages, which broke sending replies in encrypted rooms
from clients that don't use matrix-rust-sdk or mautrix-go. The issue was fixed
quickly and is included in these releases.

## gomuks
Unfortunately, hicli hasn't progressed since my [first post](https://mau.fi/blog/2024-h1-mautrix-updates/).
However, gomuks did get a [release](https://github.com/tulir/gomuks/releases/tag/v0.3.1)
for authenticated media support. It includes a local-only web server to generate
unauthenticated links that are proxied with authentication.
