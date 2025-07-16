---
title: July 2025 releases // bandwidthmuks
summary: A summary of the mautrix releases in July 2025
slug: 2025-07-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
- Meowlnir
---
This month's bridge releases include future-proofing against room v12 and many
smaller bug fixes. On the gomuks side, bandwidth usage has been improved.

| Bridge/library    | Version                                                          |
|-------------------|------------------------------------------------------------------|
| mautrix-gmessages | [v0.6.4](https://github.com/mautrix/gmessages/releases/v0.6.4)   |
| mautrix-whatsapp  | [v0.12.3](https://github.com/mautrix/whatsapp/releases/v0.12.3)  |
| mautrix-telegram  | [v0.15.3](https://github.com/mautrix/telegram/releases/v0.15.3)  |
| mautrix-bluesky   | [v0.1.2](https://github.com/mautrix/bluesky/releases/v0.1.2)     |
| mautrix-discord   | [v0.7.5](https://github.com/mautrix/discord/releases/v0.7.5)     |
| mautrix-twitter   | [v0.4.3](https://github.com/mautrix/twitter/releases/v0.4.3)     |
| mautrix-gvoice    | [v0.1.1](https://github.com/mautrix/gvoice/releases/v0.1.1)      |
| mautrix-signal    | [v0.8.5](https://github.com/mautrix/signal/releases/v0.8.5)      |
| mautrix-slack     | [v0.2.2](https://github.com/mautrix/slack/releases/v0.2.2)       |
| mautrix-meta      | [v0.5.2](https://github.com/mautrix/meta/releases/v0.5.2)        |
| mautrix-go        | [v0.24.2](https://github.com/mautrix/go/releases/v0.24.2)        |

## Bridges
Nearly all bridges got a release this month, primarily to future-proof them
against changes in future room versions. [Room v12] is scheduled to be released
next week and it contains breaking changes to power levels and the `/createRoom`
API. In order to ensure bridges can still create rooms when v12 becomes default,
all bridges have been updated to hardcode v11 in create room requests. The next
month's releases are expected to have proper support for v12.

[Room v12]: https://matrix.org/blog/2025/07/security-predisclosure/

WhatsApp's `@lid` support has received further improvements. LID DMs should now
be merged into the phone number DM more consistently, and even if merging doesn't
work, sending to LID DMs should actually work now. The bridge also stores LID
mappings received from the phone.

## Meowlnir
The Meowlnir release this month was skipped, as it is expected that there will
be changes to support room v12 within a week. There won't be a release with the
changes right away, but admins using Meowlnir who want to upgrade their rooms
to v12 are strongly encouraged to update to main/latest as soon as the v12
compatibility lands there.

Using old versions of Meowlnir in a v12 room may cause unexpected behavior, as
the bot won't know that creators have special power. However, if you don't
upgrade any Meowlnir-moderated rooms to v12, then there's no urgency in updating
Meowlnir either.

## gomuks
The bandwidthmuks project has finally started with 2 main changes:

There's now an option in settings to compress server->client traffic. It uses a
connection-wide deflate session similar to Discord. I chose a custom implementation
over the standard permessage-deflate, as the standard isn't supported properly on
Safari, plus I assume Discord did extensive testing to find that a connection-wide
deflate is more efficient than per-message. Compression is not enabled by default,
as it prevents using browser devtools effectively, and it's not really needed if
gomuks is on `localhost`.

If the client is slow to read incoming data, the server will bundle queued events
in one websocket message rather than waiting for each one. This should lead to
a few big packets instead of many small ones, which will hopefully go through
faster. I didn't actually bother benchmarking it though, so who knows.

In addition to actual bandwidth improvements, the media upload indicator now has
an actual progress bar instead of just a spinner. The indicator shows both gomuks
frontend -> backend upload and backend -> homeserver upload as the first and
second half of the progress bar respectively. Between the two halves, there may
be some processing on the backend, which isn't currently measured and depends on
options chosen in the upload dialog (e.g. re-encoding videos may take a while).
The media uploads can even be cancelled now at any stage.

Non-bandwidth changes this month include:

* Fixes for [MSC4293] support, as it turned out the initial implementation was
  broken.
* Displaying why a message was redacted (either "via ban event" or the user who
  sent the redaction event).
* Workaround for servers that return weird mime types for avatars.
* The markdown parser was changed to only allow two tildes (`~~`) for
  strikethrough, as single-tilde strikethrough has a risk of false positives.
* The view source dialog has a "Copy /raw command" button, which can be used
  as a way to forward messages (especially media messages to avoid reuploading).
  Simply paste the command into another chat and send.
* The Element Call widget is now bundled with gomuks web instead of being loaded
  from an external URL.
* [MSC3417] video rooms with Element Call now have special UI that opens straight
  into the Element Call widget. Other room types like spaces and policy list
  rooms may also receive special UI in the future.

[MSC4293]: https://github.com/matrix-org/matrix-spec-proposals/pull/4293
[MSC3417]: https://github.com/matrix-org/matrix-spec-proposals/pull/3417
