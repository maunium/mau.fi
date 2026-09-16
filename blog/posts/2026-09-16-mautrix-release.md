---
title: September 2026 releases // gomuks UI polish and policy server mutes
summary: A summary of the mautrix releases in September 2026
slug: 2026-09-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
- Meowlnir
---
Multiple UI/UX improvements were made to gomuks, courtesy of a UI review by
[lveneris](matrix:u/lveneris:kludgecs.com). Also, [Go 1.27] was released, so
the minimum version for compiling any new releases is now 1.26. As always,
precompiled binaries have no dependencies and don't require installing Go.

[Go 1.27]: https://go.dev/blog/go1.27

| Software          | Version                                                           |
|-------------------|-------------------------------------------------------------------|
| mautrix-gmessages | [v26.09](https://github.com/mautrix/gmessages/releases/v0.2609.0) |
| mautrix-whatsapp  | [v26.09](https://github.com/mautrix/whatsapp/releases/v0.2609.0)  |
| mautrix-telegram  | [v26.09](https://github.com/mautrix/telegram/releases/v0.2609.0)  |
| mautrix-linkedin  | [v26.09](https://github.com/mautrix/linkedin/releases/v0.2609.0)  |
| mautrix-twitter   | [v26.09](https://github.com/mautrix/twitter/releases/v0.2609.0)   |
| mautrix-signal    | [v26.09](https://github.com/mautrix/signal/releases/v0.2609.0)    |
| mautrix-slack     | [v26.09](https://github.com/mautrix/slack/releases/v0.2609.0)     |
| mautrix-meta      | [v26.09](https://github.com/mautrix/meta/releases/v0.2609.0)      |
| meowlnir          | [v26.09](https://github.com/maunium/meowlnir/releases/v0.2609.0)  |
| gomuks            | [v26.09](https://github.com/gomuks/gomuks/releases/v0.2609.0)     |
| mautrix-go        | [v0.31.0](https://github.com/mautrix/go/releases/v0.31.0)         |
| go-util           | [v0.10.1](https://github.com/mautrix/go-util/releases/v0.10.1)    |

## gomuks
gomuks got some UX improvements, such as swipe to reply (including haptic
feedback when using the Android wrapper) and a new message context menu style
that opens when tapping a message on mobile. The UI also received polish,
particularly around the room list and timeline message avatars. To ensure
consistency across platforms, the default font was switched to Inter.
Previously, the OS/browser default was used, which made it impossible to
correctly align room list previews everywhere.

The room list has a compact mode now, which hides previews and makes entries
smaller. A future update will likely add a spacious mode, which will expands
the entries to add space for timestamps.

Other changes this month include:

* Tenor support being re-added to the gif picker.
* A `/version` command to easily check the running version.
* Prominent display of message send errors instead hiding them in a small icon.

## Meowlnir
Policy servers can now be used to mute users using the `!ps mute` command.
Power levels could already technically do mutes before, but that's not a very
scalable solution and can only target individual users. Policy-based mutes can
also target globs and entire servers, and they're not limited to 64 KiB like
power levels are.

The policy server also has an option to send a notice to a specific room when
any event is blocked, which can be used to ensure the configuration isn't
blocking things that should be allowed.

## Bridges
The stable IDs migration announced last month for Google Messages was cancelled,
as it turned out the new IDs weren't stable after all. A bunch of other changes
were made though, which should make the bridge more reliable in general.

Twitter and Google Messages got support for streaming files to avoid memory
usage spikes when receiving large files. Other than that, bridge changes were
mostly bugfixes.
