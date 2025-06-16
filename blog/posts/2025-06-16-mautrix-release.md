---
title: June 2025 releases // Meowlnir policy server and gomuks updates
summary: A summary of the mautrix releases in June 2025
slug: 2025-06-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
- Meowlnir
---
This month's bridge releases are again fairly boring, but Meowlnir got a
built-in policy server, gomuks terminal development is starting again, and
gomuks android now has prebuilt APKs.

| Bridge/library    | Version                                                          |
|-------------------|------------------------------------------------------------------|
| mautrix-gmessages | [v0.6.3](https://github.com/mautrix/gmessages/releases/v0.6.3)   |
| mautrix-whatsapp  | [v0.12.2](https://github.com/mautrix/whatsapp/releases/v0.12.2)  |
| mautrix-twitter   | [v0.4.2](https://github.com/mautrix/twitter/releases/v0.4.2)     |
| mautrix-discord   | [v0.7.4](https://github.com/mautrix/discord/releases/v0.7.4)     |
| mautrix-signal    | [v0.8.4](https://github.com/mautrix/signal/releases/v0.8.4)      |
| mautrix-meta      | [v0.5.0](https://github.com/mautrix/meta/releases/v0.5.0)        |
| mautrix-go        | [v0.24.1](https://github.com/mautrix/go/releases/v0.24.1)        |
| meowlnir          | [v0.6.0](https://github.com/maunium/meowlnir/releases/v0.6.0)    |
| go-util           | [v0.8.8](https://github.com/mautrix/go-util/releases/v0.8.8)     |

## Bridges
Bridge updates were mostly bug fixes and some small features.

WhatsApp's `@lid` support mentioned in the April update has been improved,
although it's still not quite perfect and the migration is still ongoing on
WhatsApp's side. Further improvements are expected in the next few months.

Admins of old WhatsApp bridge instances may need to regenerate or manually
update the registration files, as old versions generated the namespace with
a `[0-9]+` regex, which will not cover the new `lid-<number>` user IDs.

## Meowlnir
Meowlnir got a bunch of changes, most notably an experimental built-in policy
server by [@nexy7574]. Currently the policy server is just an extra layer for
preventing banned users sending messages, but when Meowlnir gets support for
protections like disallowing images, the policy server will allow blocking such
messages pre-emptively rather than redacting after-the-fact.

Room ban policies are also now fully supported. When receiving a room ban or
discovering a room which was previously banned, Meowlnir will prompt whether
the room should be deleted. To enable the prompts, set the `room_ban_room` in
the config. Rooms can also be deleted manually with the `!rooms delete` and
`!rooms block` commands.

## gomuks
[@nexy7574] has started work on the gomuks terminal frontend. My previous plans
for gomuks terminal didn't end up working out, so the new plan is for nexy to
write it over the summer. It's not usable yet, but you can follow the progress
in the [terminal branch](https://github.com/gomuks/gomuks/tree/terminal).

The main gomuks repo was moved to the [gomuks org](https://github.com/gomuks)
on GitHub and mau.dev. If you were using the Docker image for gomuks web, you'll
have to update the image path.

There are now CI-built APKs for gomuks android: <https://mau.dev/gomuks/android/-/pipelines>.
Additionally, the app now has fancy edge-to-edge rendering so that the room
list background color extends below the OS top bar.

Other changes this month include:

* Added support for [MSC4293] to immediately hide messages from users who are
  banned with the redact_events flag.
* Added support for fallbacks in [MSC4144] per-message profiles (the support
  is that gomuks removes the fallback immediately, as gomuks already has proper
  rendering support).
* Fixed URL-encoded `matrix:` URIs not being handled correctly.

[MSC4293]: https://github.com/matrix-org/matrix-spec-proposals/pull/4293
[MSC4144]: https://github.com/matrix-org/matrix-spec-proposals/pull/4144
[@nexy7574]: https://github.com/nexy7574
