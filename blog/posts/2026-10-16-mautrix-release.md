---
title: October 2026 releases // Encrypted Messenger backfill and blocking users
summary: A summary of the mautrix releases in October 2026
slug: 2026-10-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
---
This month's releases include a ton of smaller (and some bigger) fixes and new
features. They're kind of all over the place though, there's no particular big
thread this time.

| Software          | Version                                                           |
|-------------------|-------------------------------------------------------------------|
| mautrix-gmessages | [v26.10](https://github.com/mautrix/gmessages/releases/v0.2610.0) |
| mautrix-whatsapp  | [v26.10](https://github.com/mautrix/whatsapp/releases/v0.2610.0)  |
| mautrix-telegram  | [v26.10](https://github.com/mautrix/telegram/releases/v0.2610.0)  |
| mautrix-signal    | [v26.10](https://github.com/mautrix/signal/releases/v0.2610.0)    |
| mautrix-slack     | [v26.10](https://github.com/mautrix/slack/releases/v0.2610.0)     |
| mautrix-meta      | [v26.10](https://github.com/mautrix/meta/releases/v0.2610.0)      |
| meowlnir          | [v26.10](https://github.com/maunium/meowlnir/releases/v0.2610.0)  |
| gomuks            | [v26.10](https://github.com/gomuks/gomuks/releases/v0.2610.0)     |
| mautrix-go        | [v0.32.0](https://github.com/mautrix/go/releases/v0.32.0)         |
| go-util           | [v0.11.0](https://github.com/mautrix/go-util/releases/v0.11.0)    |

## Bridges
The Messenger bridge can backfill encrypted history now. It'll prompt for your
backup PIN after login, though you can still skip it if you don't want history.
On the Instagram side, logging in with username/password is now possible using
the `android` login flow.

Blocking users was implemented across most of the bridges. It's supported on
Telegram, Signal, WhatsApp, Messenger, Instagram and Google Messages. Reporting
spam is also possible when blocking on Telegram, Signal and Google Messages.

The Signal got support for linking to phonenumberless accounts, which Signal
has started rolling out recently. You can also start chats with usernames or
username links now.

Google Messages changed the ID used for identifying logins, so hopefully future
relogins won't need to delete all existing rooms unnecessarily. Deleting rooms
will still be necessary if you clear app data, reinstall your phone or otherwise
shuffle the internal IDs that Google Messages uses to identify chats, but just
re-pairing isn't supposed to do that. Also, if you get unpaired, but your
cookies are still valid, you can use the `relogin` command to re-pair without
having to redo the Google login.

The Slack bridge got full support for direct media, following the partial
emoji-only support added last month.

There were a bunch of other changes too, so see the changelogs linked above for
the full list of new features.

## gomuks
Various bugs with aggregations and inefficient database indexes were fixed in
the backend, which should show up as certain views like notifications loading
faster and aggregated info like edits and reactions being more accurate.
However, the index improvements did require a few heavy migrations, so if you're
updating from v26.09 or earlier and have a big database, it may take a few
minutes to start up.

Some highlight features are support for removing reactions by clicking on them
(enabled by the improved aggregations) and support for rendering galleries.
See the release notes for the full list of changes.

Sending galleries is implemented too, but hidden in the UI. While I prefer the
single-event approach in [MSC4274] over relation-based options, I consider the
way Element is shipping the feature to be completely irresponsible. The correct
way to ship backwards-incompatible message types is to not allow sending them at
all until the approach and exact format is widely accepted. Right now, Element
just has a single toggle in the settings UI, which of course means every curious
user will enable it and then complain to other clients that don't render the
message.

[MSC4274]: https://github.com/matrix-org/matrix-spec-proposals/pull/4274

The Android wrapper also got a bunch of improvements, such as:

* User and room avatars in notifications
* Reply and mark as read actions on notifications
* Native conversation shortcuts, including a share target for them
* Support for sharing text and more reliable sharing in general
* Proper support for full screen videos
