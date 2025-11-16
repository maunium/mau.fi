---
title: November 2025 releases // gomuks terminal and mautrix-irc
summary: A summary of the mautrix releases in November 2025
slug: 2025-11-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
- Meowlnir
---

| Software          | Version                                                           |
|-------------------|-------------------------------------------------------------------|
| mautrix-gmessages | [v25.11](https://github.com/mautrix/gmessages/releases/v0.2511.0) |
| mautrix-whatsapp  | [v25.11](https://github.com/mautrix/whatsapp/releases/v0.2511.0)  |
| mautrix-linkedin  | [v25.11](https://github.com/mautrix/linkedin/releases/v0.2511.0)  |
| mautrix-twitter   | [v25.11](https://github.com/mautrix/twitter/releases/v0.2511.0)   |
| mautrix-gvoice    | [v25.11](https://github.com/mautrix/gvoice/releases/v0.2511.0)    |
| mautrix-signal    | [v25.11](https://github.com/mautrix/signal/releases/v0.2511.0)    |
| mautrix-zulip     | [v25.11](https://github.com/mautrix/zulip/releases/v0.2511.0)     |
| mautrix-slack     | [v25.11](https://github.com/mautrix/slack/releases/v0.2511.0)     |
| mautrix-meta      | [v25.11](https://github.com/mautrix/meta/releases/v0.2511.0)      |
| meowlnir          | [v25.11](https://github.com/maunium/meowlnir/releases/v0.2511.0)  |
| gomuks            | [v25.11](https://github.com/gomuks/gomuks/releases/v0.2511.0)     |
| mautrix-go        | [v0.26.0](https://github.com/mautrix/go/releases/v0.26.0)         |
| go-util           | [v0.9.3](https://github.com/mautrix/go-util/releases/v0.9.3)      |

## gomuks
The big news is that gomuks terminal is finally back. As neither of my previous
plans ended up working out, I ended up doing it myself by hacking the legacy
gomuks code to use the new RPC API to talk with the gomuks backend.

It's not quite at feature parity with legacy gomuks yet, but it works for basic
chatting. Some use of the web interface is required, as the terminal frontend
doesn't have Matrix login support yet. For now, gomuks terminal is a separate
binary that only includes the frontend, so the backend has to run separately.
In the future, they'll probably be available in the same binary and you can
choose between embedded and remote backend.

As with bridges, gomuks will use calendar versioning, so the first release of
gomuks web and new gomuks terminal is v25.11.

Other than gomuks terminal, the web frontend also had some changes:

* a space viewer with buttons for adding and removing rooms to a space
* syncing space child changes to the frontend in real time
  (previously you had to reload to apply changes)
* room autocomplete (which was required for adding rooms to a space)
* displaying invited users in the member list and autocomplete
* linkifying event IDs in message pin notices
* commands to change room name/avatar
* support for sending @room mentions

## Bridges

### mautrix-irc
I've started writing the successor to [Heisenbridge](https://github.com/hifi/heisenbridge)
in Go using the bridgev2 framework. It already works for basic chatting and has
some fancy IRCv3 features as well, but it's not very robust, so I wouldn't
recommend switching from Heisenbridge just yet. In particular, support for
legacy (non-v3) IRC networks is not yet good, as fallbacks for various features
are still missing.

### Other bridges
mautrix-zulip and mautrix-linkedin had their first releases. They more or less
work and there haven't been any major changes recently.

mautrix-whatsapp had a lot of fixes. The LID migration on WhatsApp's side is
still in progress and brings up more bugs occasionally. WhatsApp Android also
had a breaking change that broke pairing with the bridge, which is fixed now.

mautrix-signal got initial support for polls, although vote bridging isn't
working yet. As with WhatsApp, bridging polls as Element-compatible extensible
events must be enabled in the config separately, as those polls are not in the
Matrix spec yet.

## Meowlnir
Meowlnir's policy server has been updated to the latest version of [MSC4284],
which has the event sender call `/sign` once instead of every receiver calling
`/check`.
