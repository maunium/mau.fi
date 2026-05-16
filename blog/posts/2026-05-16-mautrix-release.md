---
title: May 2026 releases // Stickers everywhere
summary: A summary of the mautrix releases in May 2026
slug: 2026-05-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
---

| Software          | Version                                                           |
|-------------------|-------------------------------------------------------------------|
| mautrix-gmessages | [v26.05](https://github.com/mautrix/gmessages/releases/v0.2605.0) |
| mautrix-whatsapp  | [v26.05](https://github.com/mautrix/whatsapp/releases/v0.2605.0)  |
| mautrix-telegram  | [v26.05](https://github.com/mautrix/telegram/releases/v0.2605.0)  |
| mautrix-gvoice    | [v26.05](https://github.com/mautrix/gvoice/releases/v0.2605.0)    |
| mautrix-signal    | [v26.05](https://github.com/mautrix/signal/releases/v0.2605.0)    |
| mautrix-slack     | [v26.05](https://github.com/mautrix/slack/releases/v0.2605.0)     |
| mautrix-meta      | [v26.05](https://github.com/mautrix/meta/releases/v0.2605.0)      |
| meowlnir          | [v26.05](https://github.com/maunium/meowlnir/releases/v0.2605.0)  |
| gomuks            | [v26.05](https://github.com/gomuks/gomuks/releases/v0.2605.0)     |
| mautrix-go        | [v0.28.0](https://github.com/mautrix/go/releases/v0.28.0)         |
| go-util           | [v0.9.9](https://github.com/mautrix/go-util/releases/v0.9.9)      |

## Bridges
This month's focus was better support for stickers and emojis. Telegram already
got emoji pack import with the rewrite released last month, and this month it
was extended to work properly with stickers (including sending animated stickers
back to Telegram). Signal and WhatsApp support importing sticker packs (both
static and animated) and Slack supports importing the workspace emojis.

To import packs, simply use `import-image-pack <url>` in your management room
and it'll put it in your personal filtering space. For Slack, the pack URL is
your workspace URL, and if you have multiple logins, you can use
`import-image-pack <login ID> <url>` with the IDs from `list-logins`.

Coincidentally, [MSC2545] was also accepted this month. To encourage adoption
of stable identifiers, the bridges will only send the stable `m.room.image_pack`
event type. Complain to your local client developer if you can't see imported
packs.

[MSC2545]: https://github.com/matrix-org/matrix-spec-proposals/pull/2545

Unrelated to stickers, mautrix-manager got a fix for Google logins, which will
hopefully help with logging into Google Messages now that QR login is no longer
possible. The latest macOS release binary is also signed, which should allow
installing it without workarounds to allow unsigned apps. Windows is probably
too niche and too expensive to bother with, so those binaries remain unsigned.

## gomuks
Support for reading the stable [MSC2545] event types was naturally already
added. Image packs will be deduplicated by state key with the stable event type
taking priority if it exists.

Using packs imported with bridges is still a bit tricky as space image packs
aren't propagated to child rooms automatically. However, you can subscribe to
packs in the space by using the override view -> timeline button in settings,
then using the sticker picker in that room to subscribe to packs.

In non-sticker news, notification sounds can now be customized in settings by
providing the `mxc://` URI to any sound file. You can also disable notification
sounds by clearing the sound field.
