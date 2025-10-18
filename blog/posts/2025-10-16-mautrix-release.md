---
title: October 2025 releases // Verification in bridges, hosted Meowlnirs & Zulip bridge
summary: A summary of the mautrix releases in October 2025
slug: 2025-10-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
- Meowlnir
---
This month's releases include self-signing support in bridges, better
multi-tenant support in Meowlnir for Asgard.Chat, a Zulip bridge, and a new
versioning scheme for bridges and Meowlnir.

| Bridge/library    | Version                                                           |
|-------------------|-------------------------------------------------------------------|
| mautrix-gmessages | [v25.10](https://github.com/mautrix/gmessages/releases/v0.2510.0) |
| mautrix-whatsapp  | [v25.10](https://github.com/mautrix/whatsapp/releases/v0.2510.0)  |
| mautrix-bluesky   | [v25.10](https://github.com/mautrix/bluesky/releases/v0.2510.0)   |
| mautrix-twitter   | [v25.10](https://github.com/mautrix/twitter/releases/v0.2510.0)   |
| mautrix-gvoice    | [v25.10](https://github.com/mautrix/gvoice/releases/v0.2510.0)    |
| mautrix-signal    | [v25.10](https://github.com/mautrix/signal/releases/v0.2510.0)    |
| mautrix-slack     | [v25.10](https://github.com/mautrix/slack/releases/v0.2510.0)     |
| mautrix-meta      | [v25.10](https://github.com/mautrix/meta/releases/v0.2510.0)      |
| meowlnir          | [v25.10](https://github.com/maunium/meowlnir/releases/v0.2510.0)  |
| mautrix-go        | [v0.25.2](https://github.com/mautrix/go/releases/v0.25.2)         |
| go-util           | [v0.9.2](https://github.com/mautrix/go-util/releases/v0.9.2)      |

## Versioning
All Go bridges have been following a monthly release schedule ever since the
WhatsApp bridge was accidentally released on the same day for 3 months in a row.
To simplify picking version numbers and telling how old a given version is, all
megabridges and Meowlnir have switched to [calendar versioning](https://calver.org/)
using the `vYY.0M.patch` format.

This months releases are `v25.10` with `v0.2510.0` as the git tags. In the
unlikely event of a patch release, it'd be `v25.10.1` and `v0.2510.1`.

Due to restrictions from Go modules, the actual git tags will still follow
[0ver](https://0ver.org/), so they'll be formatted as `v0.YY0M.patch` instead.
Docker tags will be available in both formats.

Other projects like maubot and gomuks are also expected to switch to calendar
versioning in the future. Legacy bridges like Telegram and Discord won't switch
over before being rewritten though.

## Meowlnir
If you've read [This Week In Matrix](https://matrix.org/blog/2025/09/26/this-week-in-matrix-2025-09-26/#asgard-chat)
or watched this morning's keynote at the Matrix conference, you may have already
heard about [Asgard.Chat](https://asgard.chat/), a hosted moderation bot service
I recently launched with some other people. Check it out if you need a
moderation bot (which you do if you run any public Matrix rooms), but can't host
one yourself.

To work with Asgard, Meowlnir has some improved multi-tenancy features. In
particular, the `untrusted` flag in the config can be set to indicate the
instance is meant for multiple tenants to enable some additional checks to the
policy list cache. Without the checks, separate bots may be able to access
policy list info without actually being in the list room.

Additionally, there's a new `meowlnir4all` section in the config which can be
set to enable the `!provision` command. The command can be used from an admin
room to fully provision a new bot and management room for a given user.

## mautrix-zulip
There's a new puppeting bridge for [Zulip](https://zulip.com/). It doesn't have
a lot of features yet, but should work for basic bridging. As with all bridgev2
bridges, it also supports relay mode, but I haven't tested logging in as a bot.

Zulip topics are bridged as Matrix threads. It works best with a client like
gomuks, which shows all messages in the main timeline, but allows focusing on
a thread.

<https://github.com/mautrix/zulip>

## E2EE self-signing in bridges
As mentioned last month, mautrix-go and mautrix-python now have basic support
for cross-signing the bridge bot for future proofing against [MSC4153]. You can
enable it using the `self_sign` option under `encryption` in the bridge config.
When enabled, the bridge will generate a new recovery key and cross-signing keys,
then sign its own device.

Megabridges will save the recovery key in the `kv_store` table in the database,
while legacy bridges (i.e. Telegram) will just log it. If the bridge is reset
such that the bridge has to sign itself again, Go bridges will ask you to put
the old recovery key in the database, or a special `overwrite` value to tell
the bridge you want it to overwrite existing cross-signing keys. However, Python
will just always try to overwrite keys.

Note that the bridge can't do user-interactive auth, which means by default it
can only upload cross-signing keys if the bot doesn't have those keys yet. This
means that overwriting as mentioned above will fail if the bridge bot already
has cross-signing keys uploaded. However, recent changes to [MSC4190] allow
bridges to replace keys without UIA, so the problem can be avoided by enabling
MSC4190 and ensuring you have the latest version of Synapse (v1.139 or higher).
Also note that until Synapse v1.141, enabling MSC4190 requires updating the
bridge registration. See [the docs](https://docs.mau.fi/bridges/general/end-to-bridge-encryption.html#use-with-next-gen-auth-mas-msc4190)
for more info.

Self-signing the bridge bot is only the first step: it is not expected to remove
warnings, but it's enough to be compatible under the [exception that MSC4153
gives for bridges](https://github.com/matrix-org/matrix-spec-proposals/blob/main/proposals/4153-invisible-crypto.md#clients-may-make-provisions-for-encrypted-bridges).
The proper solution to remove warnings will be [MSC4350], which defines a way
to explicitly allow the bridge bot to impersonate ghosts and double puppeted
users.

In addition to verification, the support for [MSC3202] (the `appservice` mode
in the encryption config) was updated to work with Synapse, so it is now safe
to enable as long as you regenerate the registration and enable the relevant
`experimental_features` in the Synapse config as well.

[MSC3202]: https://github.com/matrix-org/matrix-spec-proposals/pull/3202
[MSC4153]: https://github.com/matrix-org/matrix-spec-proposals/pull/4153
[MSC4190]: https://github.com/matrix-org/matrix-spec-proposals/pull/4190
[MSC4350]: https://github.com/matrix-org/matrix-spec-proposals/pull/4350

## gomuks
gomuks was less active this month, but there is an account data explorer now in
addition to the old room state explorer. Also, a dedicated share button was
added to rooms and spaces to avoid having to go through the event share menu to
get a room link.
