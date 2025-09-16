---
title: September 2025 releases // Verification in maubot, bridge group creation and commandmuks
summary: A summary of the mautrix releases in September 2025
slug: 2025-09-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
- Meowlnir
- maubot
---

| Bridge/library    | Version                                                          |
|-------------------|------------------------------------------------------------------|
| mautrix-gmessages | [v0.7.0](https://github.com/mautrix/gmessages/releases/v0.7.0)   |
| mautrix-whatsapp  | [v0.12.5](https://github.com/mautrix/whatsapp/releases/v0.12.5)  |
| mautrix-signal    | [v0.8.7](https://github.com/mautrix/signal/releases/v0.8.7)      |
| mautrix-go        | [v0.25.1](https://github.com/mautrix/go/releases/v0.25.1)        |
| meowlnir          | [v0.8.0](https://github.com/maunium/meowlnir/releases/v0.8.0)    |
| go-util           | [v0.9.1](https://github.com/mautrix/go-util/releases/v0.9.1)     |

## maubot
maubot 0.6.0 beta 1 was released, which includes v12 support (e.g. when checking
power levels after following a tombstone) as well as two new buttons for
verifying the maubot device, either using an existing recovery key or generating
a new recovery key.

With the verification support, all maubots can be future-proofed for [MSC4153],
which is expected to land in the near-ish future (within months).

Bridges will likely get basic cross-signing support in the next release.
However, only the bridge bot will be verified initially, and clients will keep
showing a warning for messages from ghost users. MSC4153 includes a caveat to
allow such use until [MSC4350] is implemented, which is likely to happen some
months after the base MSC4153 changes come into effect.

[MSC4153]: https://github.com/matrix-org/matrix-spec-proposals/pull/4153
[MSC4350]: https://github.com/matrix-org/matrix-spec-proposals/pull/4350

## Bridges
There were less bridge releases this month, but all three bridges with releases
now have support for creating new group chats for existing Matrix rooms.
Bridging existing groups to existing rooms will probably happen in the next few
months too.

All three released bridges also removed support for the legacy provisioning APIs
and legacy migration, which means upgrading from the pre-megabridge versions is
no longer supported. If you were using the bridge with config writing disabled,
you must also make sure that the config is updated to the megabridge format
before upgrading to these releases.

While the WhatsApp bridge didn't get any LID-related changes, WhatsApp's rollout
of the feature to existing groups is continuing, so bridges set up before v0.5.0
must ensure their `registration.yaml` file has the appropriate regex. There's a
[new entry on the troubleshooting page](https://docs.mau.fi/bridges/general/troubleshooting.html#mautrix-whatsapp-messages-in-some-groups-are-no-longer-bridged)
with more details.

Google Messages got support for typing notifications in both directions, as well
as read receipts from RCS group chats. The read receipt feature is particularly
hacky, as the messages for web API doesn't actually send them in a sensible
format like `"read_by": ["user_id", ...]`. Instead, it sends
`"status_text": "Read by Name1, Name2, ..."`, which the bridge has to parse.
Fortunately, it wasn't too difficult to cache the raw group metadata in memory
and use those names to map to user IDs. There will be incorrect read receipts if
you have multiple users with the same name, but those are unavoidable.

The Meta bridge has support for typing notifications from/to Instagram, but it
wasn't released yet, as there were some bugs found fairly late and I didn't want
to risk a release in case there are more bugs.

## Meowlnir
Meowlnir's policy server will now actually validate PDUs it receives, which
means it might be somewhat safe to use. However, the PDU validation code in
mautrix-go depends on `GOEXPERIMENT=jsonv2`, so the policy server will be
hard-disabled if Meowlnir is built without that experiment.

## gomuks
gomuks got support for web push, which means you can now use it as a PWA on iOS
with proper push notifications. It of course works on Android and desktop
browsers as well, but those are less useful, as there's gomuks android for
native FCM push and desktop generally doesn't need push.

Thread support has been improved by making thread messages in timeline more
compact by default: instead of rendering entire messages as replies, the main
timeline only shows the first 2 lines and makes it very obvious that they're
thread messages. Clicking on any thread message will open the thread sidebar.

If you prefer the old way of rendering thread messages, you can still disable
"Compact thread messages" in settings.

Experimental support for [MSC4332]: In-room bot commands was added, along with
a bunch of built-in commands built on the same system. The UI is still fairly
rough, but you can still just type parameters instead of using the input fields.

[MSC4332]: https://github.com/matrix-org/matrix-spec-proposals/pull/4332
