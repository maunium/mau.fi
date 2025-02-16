---
title: February 2025 releases // Signal backfill
summary: A summary of the mautrix releases in February 2025
slug: 2025-02-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
---
This month's highlight is the new history transfer feature on Signal, which is
now supported in the bridge.

Signal backfill is similar to WhatsApp in that it only happens on login.
Additionally, the app will ask you if you want to transfer history when pairing
(if you don't get the prompt, it probably means you're using an old version of
the app that doesn't support history transfer).

Other bridges didn't have any major changes, but WhatsApp and Meta had a couple
of bugfixes.

Also, with the release of [Go 1.24](https://go.dev/blog/go1.24), the minimum
version for compiling any new releases is now 1.23. As always, precompiled
binaries have no dependencies and don't require installing Go.

| Bridge/library    | Version                                                          |
|-------------------|------------------------------------------------------------------|
| mautrix-whatsapp  | [v0.11.3](https://github.com/mautrix/whatsapp/releases/v0.11.3)  |
| mautrix-signal    | [v0.8.0](https://github.com/mautrix/signal/releases/v0.8.0)      |
| mautrix-meta      | [v0.4.4](https://github.com/mautrix/meta/releases/v0.4.4)        |
| mautrix-go        | [v0.23.1](https://github.com/mautrix/go/releases/v0.23.1)        |
| go-util           | [v0.8.5](https://github.com/mautrix/go-util/releases/v0.8.5)     |

## gomuks web
There weren't as many changes as previous months, but some features were added
and multiple bugs were fixed. [@nexy7574](https://github.com/nexy7574) contributed
user moderation actions in the right panel (kick, ban, invite, ignore) and
a share button to get links to event. Other than that, I added:

* Touch panning/zooming for mobile devices in the image viewer
* Avatar thumbnails to avoid wasting ram on rendering multi-megabyte avatars
* Manual key export/import (both all keys and keys for a specific room)
