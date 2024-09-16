---
title: September 2024 releases & progress
summary: A summary of the mautrix releases in September 2024
slug: 2024-09-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
---
As mentioned in last month's release post, the Google Messages and Meta bridges
have been upgraded to megabridges. As a bonus surprise, there's now a Google
Voice bridge as well.

| Bridge/library    | Version                                                          |
|-------------------|------------------------------------------------------------------|
| mautrix-gmessages | [v0.5.0](https://github.com/mautrix/gmessages/releases/tag/v0.5.0)  |
| mautrix-gvoice    | [v0.1.0](https://github.com/mautrix/gvoice/releases/tag/v0.1.0)  |
| mautrix-meta      | [v0.4.0](https://github.com/mautrix/meta/releases/tag/v0.4.0)    |
| mautrix-slack     | [v0.1.1](https://github.com/mautrix/slack/releases/tag/v0.1.1)   |
| mautrix-signal    | [v0.7.1](https://github.com/mautrix/signal/releases/tag/v0.7.1)   |
| mautrix-go        | [v0.21.0](https://github.com/mautrix/go/releases/tag/v0.21.0)    |
| go-util           | [v0.8.0](https://github.com/mautrix/go-util/releases/tag/v0.8.0) |

## Megabridge progress
The Telegram bridge rewrite has progressed and is being tested at Beeper.
However, there are a number of features that Beeper doesn't use, which means
it'll still take some time before the rewrite is fully released.

The Discord rewrite got delayed due to bugfixes in other bridges, and WhatsApp
also ended up being prioritized over Discord. The WhatsApp rewrite will likely
be released next month.

## Authenticated media
Authenticated media support was added [back in July](https://mau.fi/blog/2024-07-mautrix-release/)
and by now most servers have been forced to update as matrix.org and others
have started freezing unauthenticated media. These releases drop support for
unauthenticated media downloads entirely, which means the bridges will only
work with an up-to-date homeserver. However, the bridges will not enforce this
via `/versions`, as some servers don't advertise v1.11 support yet.

## New bridgev2 features
The bridgev2 interface now includes streaming file download/upload methods:
`DownloadMediaToFile` and `UploadMediaStream`. Some of the network connectors
already use those, and the rest will switch over time. The methods will
roundtrip big attachments through the disk to avoid storing the entire file in
memory. It should be especially useful when backfilling to avoid several files
being stored in memory at the same time.

Another new feature is "split portals", which means segregating portal rooms on
Matrix by login instead of sharing them. With the default shared rooms, if
multiple users are in the same chat on the remote network, they'll all be in
the same Matrix room. With split portals, every login gets a separate room.
In the future, the option may be extended to split ghost users as well, to
allow using contact list names without risk of conflicts.

## Meowlnir
Because I *clearly* don't have enough projects already, I started writing a new
Matrix moderation bot: Meowlnir. It is now somewhat functional (it can watch a
policy list, ban users and redact messages) and is deployed on maunium.net.

The primary goal is to be as fast and efficient as possible, as well as have
native encryption support without the need for Pantalaimon or other hacks.

It's initially only compatible with Synapse as it depends on [MSC2409] &
[MSC3202] for encryption, plus it reads the database directly to efficiently
find events to redact (although that requirement will hopefully go away in the
future with [MSC4914]).

You can find Meowlnir at <https://github.com/maunium/meowlnir>.

[MSC2409]: https://github.com/matrix-org/matrix-spec-proposals/pull/2409
[MSC3202]: https://github.com/matrix-org/matrix-spec-proposals/pull/3202
[MSC4914]: https://github.com/matrix-org/matrix-spec-proposals/pull/4194
