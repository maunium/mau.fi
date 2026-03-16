---
title: March 2026 releases // Meowlnir protections and XChat support
summary: A summary of the mautrix releases in March 2026
slug: 2026-03-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
- Meowlnir
---
The initial revision of protections in Meowlnir is finally here. The Telegram
Go rewrite is also nearing completion.

| Software          | Version                                                           |
|-------------------|-------------------------------------------------------------------|
| mautrix-whatsapp  | [v26.03](https://github.com/mautrix/whatsapp/releases/v0.2603.0)  |
| mautrix-twitter   | [v26.03](https://github.com/mautrix/meta/releases/v0.2603.0)      |
| mautrix-signal    | [v26.03](https://github.com/mautrix/signal/releases/v0.2603.0)    |
| mautrix-slack     | [v26.03](https://github.com/mautrix/slack/releases/v0.2603.0)     |
| meowlnir          | [v26.03](https://github.com/maunium/meowlnir/releases/v0.2603.0)  |
| gomuks            | [v26.03](https://github.com/gomuks/gomuks/releases/v0.2603.0)     |
| mautrix-go        | [v0.26.4](https://github.com/mautrix/go/releases/v0.26.4)         |
| go-util           | [v0.9.7](https://github.com/mautrix/go-util/releases/v0.9.7)      |

## Bridges
The Go rewrite of Telegram received support for avoiding contact names, which
was the main blocker for multi-user bridges. Other than the old-style relaybot,
it should now be ready to replace the old bridge. I'm not yet certain whether
I'll implement the actual old-style relaybot, or just add enough config options
to allow emulating the behavior with a dedicated relay user. Either way, I'm
hoping to have it ready for the next release.

Signal had some breaking changes which resulted in extremely rare unscheduled
releases. This latest release also starts sending the new fields to match how
native Signal 8.0 clients work.

The Twitter bridge has been updated to support the new encrypted XChat chats.
The quality of the encryption protocol is obviously extremely dubious and you
shouldn't rely on Twitter for anything sensitive, but supporting the new chat
type is necessary to keep receiving messages. The bridge will prompt you to
provide the encryption PIN after updating, which you can do using the
`relogin <id>` command (where `<id>` can be found with `list-logins`).

## Meowlnir
I finally got around to reviewing and slightly refactoring the protections PR,
which is now included in this release. The list of available protections is
still fairly small, plus most of them aren't designed to work well with the
policy server and are therefore marked as unsafe. However, the most important
one (`no_media`) works fine both with and without a policy server.

Future development plans include more advanced media protections using local
image classification models and making the other protections safer to use in
the policy server mode.

## gomuks
Direct chat classification was switched to use the `m.direct` account data
event, so you may have to manually mark rooms as DMs if they weren't already
there. The `/converttodm` command exists to do that.

Klipy support was also added to the GIF picker, as Tenor is being shut down
in June.
