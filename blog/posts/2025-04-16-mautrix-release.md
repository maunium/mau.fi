---
title: April 2025 releases // WhatsApp store migration and more Meowlnir commands
summary: A summary of the mautrix releases in April 2025
slug: 2025-04-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
- Meowlnir
---
This was a quieter month, so there aren't that many exciting changes. Various
bridges got bugfixes, Meowlnir got a bunch of new commands and gomuks can now
create rooms.

| Bridge/library    | Version                                                          |
|-------------------|------------------------------------------------------------------|
| mautrix-whatsapp  | [v0.12.0](https://github.com/mautrix/whatsapp/releases/v0.12.0)  |
| mautrix-twitter   | [v0.4.0](https://github.com/mautrix/twitter/releases/v0.4.0)     |
| mautrix-discord   | [v0.7.3](https://github.com/mautrix/discord/releases/v0.7.3)     |
| mautrix-signal    | [v0.8.2](https://github.com/mautrix/signal/releases/v0.8.2)      |
| mautrix-slack     | [v0.2.1](https://github.com/mautrix/slack/releases/v0.2.1)       |
| mautrix-meta      | [v0.4.6](https://github.com/mautrix/meta/releases/v0.4.6)        |
| mautrix-go        | [v0.23.3](https://github.com/mautrix/go/releases/v0.23.3)        |
| meowlnir          | [v0.4.0](https://github.com/maunium/meowlnir/releases/v0.4.0)    |

## Bridges
The WhatsApp bridge includes a somewhat dangerous change to the way Signal
sessions are stored, which is necessary as WhatsApp is starting to move away
from using phone numbers as internal user identifiers. Users on the old bridge
version may find certain groups stop working until they update.

Discord had accumulated a bunch of bugfixes, so I decided to make another
release of the old version. Unfortunately there hasn't been any progress on
the megabridge rewrite.

On a non-mautrix note, Heisenbridge [v1.15.3](https://github.com/hifi/heisenbridge/releases/tag/v1.15.3)
was released recently with a crash fix. Updating is recommended ASAP.

## Meowlnir
Meowlnir got a bunch of new commands like `!send-as-bot` for talking as the
moderation bot, `!redact-recent` for redacting all recent messages,
`!powerlevel` for modifying power levels, and more.

[@nexy7574] also contributed support for automatic unbans when a ban policy is
removed.

Meowlnir, my Synapse fork and synapse-http-antispam now have support for an
experimental join rule which can be used to block joins based on a policy list.
The rule is fully compatible with all servers as it (ab)uses the existing
`restricted` join rule mechanism. The only downside is that it centralizes all
room joins to go through a single server.

## gomuks web
Since last month, gomuks web has received a bunch of bugfixes, as well as a few
new features. In particular, there's now a room creation dialog that allows
customizing everything, support for room upgrades, and support for knocking on
rooms. The last 2 features were contributed by [@nexy7574].

[@nexy7574]: https://github.com/nexy7574
