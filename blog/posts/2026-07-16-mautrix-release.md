---
title: July 2026 releases // New Instagram bridge and OAuth in gomuks
summary: A summary of the mautrix releases in July 2026
slug: 2026-07-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
---
This month's releases include a new Instagram bridge and various gomuks
improvements, like support for next-gen auth and a better settings screen.

| Software          | Version                                                           |
|-------------------|-------------------------------------------------------------------|
| mautrix-whatsapp  | [v26.07](https://github.com/mautrix/whatsapp/releases/v0.2607.0)  |
| mautrix-telegram  | [v26.07](https://github.com/mautrix/telegram/releases/v0.2607.0)  |
| mautrix-signal    | [v26.07](https://github.com/mautrix/signal/releases/v0.2607.0)    |
| mautrix-slack     | [v26.07](https://github.com/mautrix/slack/releases/v0.2607.0)     |
| mautrix-meta      | [v26.07](https://github.com/mautrix/meta/releases/v0.2607.0)      |
| gomuks            | [v26.07](https://github.com/gomuks/gomuks/releases/v0.2607.0)     |
| mautrix-go        | [v0.29.0](https://github.com/mautrix/go/releases/v0.29.0)         |
| go-util           | [v0.9.11](https://github.com/mautrix/go-util/releases/v0.9.11)    |

## gomuks
Support for next-gen auth was added. To avoid constant problems with redirects
that clients like Element have, gomuks will default to using device code login.
For servers that don't support device code login, redirects are supported to
some extent, but only on web, not the desktop or android wrappers.

Support for legacy SSO has been removed, as most servers that use an external
IdP have already migrated to next-gen auth. Legacy password auth will remain
supported for the forseeable future.

The settings view now has tabs for better organization and more options,
including actual room settings (you don't have to use devtools to change the
room name/topic anymore). You can also change user power levels in the right
panel now.

Finally, polls are now rendered and can be voted on. There's also a /poll
command to create new polls, but no other UI for that yet.

## Bridges

### New Instagram bridge
As mentioned [last month](https://mau.fi/blog/2026-06-mautrix-release/),
Instagram switched to a completely new protocol for DMs, which means a new
Instagram bridge. The old bridge can still be used for Instagram for now, but
Messenger's protocol was also changed slightly in a way that doesn't allow
connecting to Instagram, which means this is the last release where the old
bridge will be able to connect to both.

The new bridge is in the same repo, but is a separate binary and docker image.
Docker image tags have the `ig-` prefix, so this release for example can be
found at `dock.mau.dev/mautrix/meta:ig-v26.07`. Binaries are built in the same
CI job, the old bridge is `mautrix-meta` and the new one is `mautrix-instagram`.
If building manually, use the `./build-ig.sh` script.

If you've been using the bridge with `mode: instagram` in the config, or if all
existing logins are Instagram logins, the new bridge is a drop-in replacement
and even rolling back should more or less work. However, for bridges that have
both Instagram and Messenger logins active, there's no way to split them up, so
you'll have to set up a new bridge and relogin manually.

### WhatsApp passkeys
Last week WhatsApp temporarily enforced passkey login for a fraction of users.
It appears the enforcement was a temporary trial run, or maybe they realized the
UX is horrible and even their official macOS app broke. In either case, the
bridge now supports login using passkeys, so it'll keep working in the event
that they roll out forced passkeys again.

In the bot command login flow, passkeys are implemented by manually running JS
in a browser. If you want a more user-friendly experience, [mautrix-manager](https://github.com/mautrix/manager)
has built-in support for the QR passkey flow (hybrid transport/caBLE).

In other WhatsApp news, their migration from phone numbers to LIDs appears to
be more or less finished, so the bridge will likely switch DMs over to LIDs on
the Matrix side in the next month or two. This will show up as the other user's
ghost in every DM leaving and a new LID ghost joining in its place.
