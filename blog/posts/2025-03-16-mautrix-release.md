---
title: March 2025 releases // Meowlnir antispam & gomuks widgets
summary: A summary of the mautrix releases in March 2025
slug: 2025-03-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
- Meowlnir
---
Meowlnir has a bunch of new features and gomuks received support for widgets &
Element Call. The Meta bridge received an experimental fix for missing user
profiles in encrypted chats, but most other bridge changes were minor.

| Bridge/library    | Version                                                          |
|-------------------|------------------------------------------------------------------|
| mautrix-gmessages | [v0.6.1](https://github.com/mautrix/gmessages/releases/v0.6.1)   |
| mautrix-whatsapp  | [v0.11.4](https://github.com/mautrix/whatsapp/releases/v0.11.4)  |
| mautrix-bluesky   | [v0.1.1](https://github.com/mautrix/bluesky/releases/v0.1.1)     |
| mautrix-twitter   | [v0.3.0](https://github.com/mautrix/twitter/releases/v0.3.0)     |
| mautrix-signal    | [v0.8.1](https://github.com/mautrix/signal/releases/v0.8.1)      |
| mautrix-slack     | [v0.2.0](https://github.com/mautrix/slack/releases/v0.2.0)       |
| mautrix-meta      | [v0.4.5](https://github.com/mautrix/meta/releases/v0.4.5)        |
| mautrix-go        | [v0.23.2](https://github.com/mautrix/go/releases/v0.23.2)        |
| meowlnir          | [v0.3.0](https://github.com/maunium/meowlnir/releases/v0.3.0)    |
| go-util           | [v0.8.6](https://github.com/mautrix/go-util/releases/v0.8.6)     |

## Meowlnir
Last week I wrote [synapse-http-antispam], which is a Synapse module that
exposes the [spam checker callbacks] over HTTP, allowing other services like
Meowlnir to handle them. With the help of the module, Meowlnir can now be used
to block invites from banned users. It can even automatically reject pending
invites in case the invite happened before the ban policy if you provide a
double puppeting access token in the Meowlnir config.

[synapse-http-antispam]: https://github.com/maunium/synapse-http-antispam
[spam checker callbacks]: https://element-hq.github.io/synapse/latest/modules/spam_checker_callbacks.html

Meowlnir also received support for managing server ACLs, as well as [MSC4204]
and [MSC4205], which allow sending policy events without directly exposing the
banned entity or why they were banned.

[MSC4204]: https://github.com/matrix-org/matrix-spec-proposals/pull/4204
[MSC4205]: https://github.com/matrix-org/matrix-spec-proposals/pull/4205

## gomuks web
gomuks web received support for widgets, which ended up being surprisingly easy
thanks to [matrix-widget-api](https://github.com/matrix-org/matrix-widget-api).
Element Call is also supported, although the UX isn't quite optimal yet (e.g.
the call won't stay open when switching to another room).

Other new features include:

* Viewing edit history (by clicking the `(edited)` text)
* [MSC2815] support to view redacted event content
* Room state explorer and editor
* Context menu in room list for muting rooms and marking as unread
* Starting new DMs in the user info right panel
* Starting new threads via the reply button
* Options to disable inline images and avatars in invites
  * Disabling normal media previews was already an option previously

[MSC2815]: https://github.com/matrix-org/matrix-spec-proposals/pull/2815
