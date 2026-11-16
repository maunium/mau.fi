---
title: November 2026 releases // Bridge parity and preparing for time
summary: A summary of the mautrix releases in November 2026
slug: 2026-11-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
- Meowlnir
---

| Software          | Version                                                           |
|-------------------|-------------------------------------------------------------------|
| mautrix-gmessages | [v26.11](https://github.com/mautrix/gmessages/releases/v0.2611.0) |
| mautrix-whatsapp  | [v26.11](https://github.com/mautrix/whatsapp/releases/v0.2611.0)  |
| mautrix-linkedin  | [v26.11](https://github.com/mautrix/linkedin/releases/v0.2611.0)  |
| mautrix-telegram  | [v26.11](https://github.com/mautrix/telegram/releases/v0.2611.0)  |
| mautrix-gvoice    | [v26.11](https://github.com/mautrix/gvoice/releases/v0.2611.0)    |
| mautrix-signal    | [v26.11](https://github.com/mautrix/signal/releases/v0.2611.0)    |
| mautrix-slack     | [v26.11](https://github.com/mautrix/slack/releases/v0.2611.0)     |
| mautrix-meta      | [v26.11](https://github.com/mautrix/meta/releases/v0.2611.0)      |
| gomuks            | [v26.11](https://github.com/gomuks/gomuks/releases/v0.2611.0)     |
| mautrix-go        | [v0.33.0](https://github.com/mautrix/go/releases/v0.33.0)         |
| go-util           | [v0.11.1](https://github.com/mautrix/go-util/releases/v0.11.1)    |

## Bridges
Bridges had a wide pass over bugs and feature parity with the official apps.
There's still a lot of work to do, so the next month's releases will likely
have a similar number of fixes.

Most of the changes were small, but some worth highlighting are:

* Sending messages over 2048 characters to Signal is now supported.
* Messenger polls are now supported.
* Telegram polls are partially supported (votes can be sent, but not received).
* Google Voice can merge calls and voicemails into the SMS chat.

The parity work is largely driven by AI, which actually produces decent results
as of a few months ago. I don't understand how LLMs, which still essentially
just repeatedly guess the next word, can be so effective at reverse-engineering
apps and debugging code. I'm also not a fan of the various non-technical issues
surrounding AI, but the fact is that it can be used to produce better software
faster than by doing everything by hand. All changes to all the bridges and
other associated software will continue to be reviewed manually, as the LLMs
do still make mistakes occasionally.

## Meowlnir
Meowlnir got a command to pre-sign all events in the Synapse database for the
policy server. This will be used as the basis of a new concept I call "time".

In the CAP theorem, Matrix chooses availability over consistency by default:
if your server is up, you can use any room you're in regardless of whether
other servers in the room are reachable. For private rooms, the lack of a single
point of failure is great. However, it also opens up potential attacks: if
consistency isn't guaranteed, the servers the moderators are on may not see
spam.

Solving that was one of the reasons for introducing policy servers. The spec
defines policy serves to fail closed, which means messages can't be sent if the
policy server isn't reachable. There's just one gap: servers may randomly request
signatures for old events from before the policy server was enabled, so it can't
assume that any old event had its timestamp faked.

If the policy server has pre-signed all events, it becomes reasonably safe to
reject any new events that arrive with a timestamp too far in the past. There
are no protections taking advantage of this properly yet, but the plan is to
add things like slow mode, rate limits and backdating prevention.

## gomuks
gomuks wasn't as active this month, but QuadRadical did add some backend and FFI
features for Nexus.

There's a CSS proxy built into the backend now, which can be used to ensure that
themes imported from [css.gomuks.app](https://css.gomuks.app) keep working even
if you temporarily can't reach the server.

Also, the release and CI builds now use goolm instead of libolm. The next
release of mautrix-go will make it default for everything.
