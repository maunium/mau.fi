---
title: November 2024 releases // More gomuks web progress
summary: A summary of the mautrix releases in November 2024
slug: 2024-11-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
---
This month's bridge releases mostly include bugfixes. On the non-bridge side,
gomuks web is progressing quickly.

| Bridge/library    | Version                                                          |
|-------------------|------------------------------------------------------------------|
| mautrix-gmessages | [v0.5.2](https://github.com/mautrix/gmessages/releases/v0.5.2)   |
| mautrix-whatsapp  | [v0.11.1](https://github.com/mautrix/whatsapp/releases/v0.11.1)  |
| mautrix-discord   | [v0.7.1](https://github.com/mautrix/discord/releases/v0.7.1)     |
| mautrix-signal    | [v0.7.3](https://github.com/mautrix/signal/releases/v0.7.3)      |
| mautrix-slack     | [v0.1.3](https://github.com/mautrix/slack/releases/v0.1.3)       |
| mautrix-meta      | [v0.4.2](https://github.com/mautrix/meta/releases/v0.4.2)        |
| mautrix-go        | [v0.22.0](https://github.com/mautrix/go/releases/v0.22.0)        |
| go-util           | [v0.8.2](https://github.com/mautrix/go-util/releases/v0.8.2)     |

## Megabridge progress
Most of the bridge work has been focused on bugfixes and Beeper's local bridges
(particularly push notifications), which is why nothing particularly exciting
happened.

Beeper has deployed the Telegram bridge to everyone and it's more or less
working fine. There's still no ETA for the full release though. In particular,
migrating the relaybot feature will need some more consideration, and contact
names in multi-user contexts need better handling.

The Twitter rewrite is mostly complete, but testing has stalled, so it's not
ready for release yet. The LinkedIn rewrite is also underway.

The old Discord bridge got a minor bugfix release, but the rewrites for Discord
and Google Chat will likely take until next year.

## gomuks web
Since the last post, gomuks web has received support for

* sending media (including captions), edits and reactions
* custom emoji packs (MSC2545)
* replying in threads (Element X style)
* unread message/notification/highlight counts
* a right panel (pinned messages, member list, user info)
* html sanitization on the backend rather than at render time
* notification sounds
* SSO and Beeper login
* fallback avatars
* basic keybindings (e.g. ctrl+k to focus room search)
* LaTeX rendering

and probably some other things I forgot about.

It has also been moved to the main branch of <https://github.com/tulir/gomuks>,
but there probably won't be any releases before I re-add the terminal frontend.

Some important features like spaces are still missing, but I've already switched
to using gomuks web as my primary desktop client. Using a laptop is much nicer
when the battery lasts 6+ hours instead of 2-3 hours back when I had Element open.
