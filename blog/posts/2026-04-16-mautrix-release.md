---
title: April 2026 releases // Telegram Go release and more advanced relays
summary: A summary of the mautrix releases in April 2026
slug: 2026-04-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
- Meowlnir
---
The Go rewrite of Telegram has finally been released. Related to that, there
have been multiple improvements to relay mode in bridgev2.

| Software          | Version                                                           |
|-------------------|-------------------------------------------------------------------|
| mautrix-gmessages | [v26.04](https://github.com/mautrix/gmessages/releases/v0.2604.0) |
| mautrix-whatsapp  | [v26.04](https://github.com/mautrix/whatsapp/releases/v0.2604.0)  |
| mautrix-linkedin  | [v26.04](https://github.com/mautrix/linkedin/releases/v0.2604.0)  |
| mautrix-telegram  | [v26.04](https://github.com/mautrix/telegram/releases/v0.2604.0)  |
| mautrix-twitter   | [v26.04](https://github.com/mautrix/twitter/releases/v0.2604.0)   |
| mautrix-gvoice    | [v26.04](https://github.com/mautrix/gvoice/releases/v0.2604.0)    |
| mautrix-signal    | [v26.04](https://github.com/mautrix/signal/releases/v0.2604.0)    |
| mautrix-slack     | [v26.04](https://github.com/mautrix/slack/releases/v0.2604.0)     |
| mautrix-meta      | [v26.04](https://github.com/mautrix/meta/releases/v0.2604.0)      |
| meowlnir          | [v26.04](https://github.com/maunium/meowlnir/releases/v0.2604.0)  |
| gomuks            | [v26.04](https://github.com/gomuks/gomuks/releases/v0.2604.0)     |
| mautrix-go        | [v0.27.0](https://github.com/mautrix/go/releases/v0.27.0)         |
| go-util           | [v0.9.8](https://github.com/mautrix/go-util/releases/v0.9.8)      |

## Bridges

### Telegram Go rewrite
The Go rewrite of Telegram has been released. Thanks to all the early testers
and people who forgot that the `:latest` docker tag always follows the main
branch.

In terms of new features, the highlights are support for topic groups as well
as sticker/emoji pack import. Since it's a rewrite, there are a lot of other
smaller changes as well, but I'm not going to bother finding all of them. If
something is missing, you can bring it up in the Matrix room.

The relay mode in bridgev2 has been extended to allow emulating the old-style
relaybots more closely. A `bridge` command was added to bridge existing remote
chats to existing Matrix rooms, which also allows bridging through a default
relay user without being logged into the bridge. When doing that, `set-relay`
is called implicitly. Both `bridge` and `set-relay` will now default to using
default relays if any are configured, even for admins who are logged into the
bridge.

Note that there is no migration from the old relaybot to the new relay mode, as
the migration doesn't know who you'd assign the login to. I'm not yet certain
whether a migration will ever be added to the bridge itself, but it's technically
possible to migrate the relevant data manually.

### Matrix relays
I've made a proof of concept that allows running any bridgev2 bridge with just
a single bot instead of an appservice, using [MSC4144](https://github.com/matrix-org/matrix-spec-proposals/pull/4144)
per-message profiles for representing users on the remote network:
<https://github.com/mautrix/go/pull/484>. It still has bugs and it's not
intended to replace appservice mode, but once it's cleaned up and merged,
it might be useful for running light-weight bridges without a homeserver,
e.g. for relaying rooms to many networks.

### Other changes
Most other bridge changes were bug fixes and small features. The Messenger
bridge has improved support for marketplace chats, including moving them to
a separate subspace. See the release notes linked above for more details.

Various crypto bugs were fixed in mautrix-go, bringing goolm closer to being
production-ready. Recent testers haven't reported any unexpected undecryptable
messages caused by goolm, so switching to it by default later this year is
plausible.

The Discord bridgev2 rewrite is progressing, but the current estimate is at
least 2 months, likely more.

## gomuks
gomuks got a whole bunch of small fixes, as well as a few bigger new features.

There's now an arbitrary profile editor in the devtools/state explorer, which
allows you to easily modify any global profile field. Additionally, there's a
dropdown to set your pronouns in your own profile view in the right sidebar.
However, the pronoun dropdown will only work for simple cases; if you want to
use non-English pronouns, use the new arbitrary profile field editor in
devtools.

There's a new message context menu item to view the full list of reactions on
a given message and delete any of them (including reactions from other users if
you have redaction permissions).

The backend and widget API were updated to support sticky events, which should
future-proof the Element Call widget for when they switch to sticky events. It
hasn't actually been tested yet though, so it's possible it'll end up exploding.

Finally, setting up cross-signing and key backup is now supported. When logging
in, gomuks web will detect whether the account already has key backup set up and
offer appropriate buttons based on that. New accounts can set up keys, while old
ones can reset keys if the recovery key was lost. This means that on servers
using SSO or next-gen auth, you should be able to register and fully set up an
account with only gomuks web. The legacy registration API is not supported
though, so for such servers the account still has to be created separately.

## Meowlnir
Meowlnir only had some minor improvements to the policy server, plus support
for obfuscating profiles when banning users to prevent malicious user IDs from
sticking around in the timeline. If you do use policy servers, the new cache
for signatures is somewhat important, as the database can get overloaded
otherwise.
