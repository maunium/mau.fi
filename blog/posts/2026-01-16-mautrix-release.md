---
title: January 2026 releases // New command systems and gomuks notification center
summary: A summary of the mautrix releases in January 2026
slug: 2026-01-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
- Meowlnir
---
This month's releases include [MSC4391] support in gomuks and Meowlnir, as well
as various bugfixes.

| Software          | Version                                                           |
|-------------------|-------------------------------------------------------------------|
| mautrix-gmessages | [v26.01](https://github.com/mautrix/gmessages/releases/v0.2601.0) |
| mautrix-whatsapp  | [v26.01](https://github.com/mautrix/whatsapp/releases/v0.2601.0)  |
| mautrix-signal    | [v26.01](https://github.com/mautrix/signal/releases/v0.2601.0)    |
| meowlnir          | [v26.01](https://github.com/maunium/meowlnir/releases/v0.2601.0)  |
| gomuks            | [v26.01](https://github.com/gomuks/gomuks/releases/v0.2601.0)     |
| mautrix-go        | [v0.26.2](https://github.com/mautrix/go/releases/v0.26.2)         |
| go-util           | [v0.9.5](https://github.com/mautrix/go-util/releases/v0.9.5)      |

## gomuks
The command system was switched from [MSC4332] to [MSC4391], which is both
simpler and more versatile. As a part of the switch, commands were also
implemented in new gomuks terminal, which means replying/reacting is possible
again (among other things). Note that there's no autocompletion for commands
in the terminal interface yet, but they should work other than that.

[MSC4332]: https://github.com/matrix-org/matrix-spec-proposals/pull/4332
[MSC4391]: https://github.com/matrix-org/matrix-spec-proposals/pull/4391

On the web side, a notification center was added to easily view messages that
triggered notifications. It has the caveat that only messages that gomuks
received will be listed, i.e. if gomuks was offline and there were lots of
messages, those messages won't show up until you scroll up to them.

## Bridges
Bridge releases this month were very bugfix-heavy: Signal had various fixes for
sender keys added last month and WhatsApp got more LID handling fixes.

Telegram is mostly in the same spot as last month, but more bugs were fixed, so
now it's stable for single-account bridges. Discord is progressing on schedule,
so it'll take at least another month, but might take longer too.

The Twitter bridge is also being rewritten again to support the new "encrypted"
XChat protocol. You can already use the `xchat` branch if you depend on Twitter
DMs for some reason. After updating the bridge, run `relogin <your login ID>`
to be prompted for your encryption PIN. The remaining bugs will hopefully be
fixed within the next month so that the rewrite can be included in the next
release cycle.

Switching mautrix-go to use [MSC4391] for bridge commands is planned next month
and maubot will probably also eventually get support.

## Meowlnir
Meowlnir has also switched to [MSC4391] for commands, but other than that there
weren't any major changes. I'm hoping to review nexy's protections PR in the
next month to get it into the next release, but not sure if that will actually
happen.
