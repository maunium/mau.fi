---
title: December 2025 releases // Image pack editing and megabridge rewrite progress
summary: A summary of the mautrix releases in December 2025
slug: 2025-12-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
- Meowlnir
---
This month's releases don't have too many new features, but gomuks did get an
image pack editor and the remaining megabridge rewrites are progressing again.

| Software          | Version                                                           |
|-------------------|-------------------------------------------------------------------|
| mautrix-whatsapp  | [v25.12](https://github.com/mautrix/whatsapp/releases/v0.2512.0)  |
| mautrix-linkedin  | [v25.12](https://github.com/mautrix/linkedin/releases/v0.2512.0)  |
| mautrix-signal    | [v25.12](https://github.com/mautrix/signal/releases/v0.2512.0)    |
| mautrix-meta      | [v25.12](https://github.com/mautrix/meta/releases/v0.2512.0)      |
| meowlnir          | [v25.12](https://github.com/maunium/meowlnir/releases/v0.2512.0)  |
| gomuks            | [v25.12](https://github.com/gomuks/gomuks/releases/v0.2512.0)     |
| mautrix-go        | [v0.26.1](https://github.com/mautrix/go/releases/v0.26.1)         |
| go-util           | [v0.9.4](https://github.com/mautrix/go-util/releases/v0.9.4)      |

## gomuks
gomuks terminal didn't get any updates this month, but gomuks web got a voice
message recorder and an emoji/sticker pack editor. The pack editor can also be
used to reorder images using [MSC4389](https://github.com/matrix-org/matrix-spec-proposals/pull/4389).

The GeckoView wrapper for Android got a bunch of updates and now supports
requesting microphone permissions for the voice message recorder and Element
Call, as well as downloading files to disk.

Unfortunately asking for camera permission for Element Call is not yet supported
because Android is dumb: apps that never ask for camera permission can use the
system camera to take photos without permission, but if the camera permission is
in the manifest, then it is required even for using the system camera. This
means that using the camera to take a photo also needs a permission prompt
before video calling can be supported.

## Bridges
The Signal bridge got support for sender keys, which makes sending messages to
groups significantly faster. Polls are also fully supported now. Other than
that, most bridge changes were just bugfixes.

On the megabridge rewrites, Telegram is now progressing again. Lots of bugs have
been fixed and support for bot logins was added, so while there's no old-style
relaybot, you can use the generic relay mode with a bot account. It is not yet
safe for multi-user use, as it uses contact names by default (Telegram makes it
very difficult to avoid contact names, the old bridge had complicated logic for
that). If you only have one login though, you can try migrating to the new
bridge.

The Discord rewrite has also finally started, but is expected to take a couple
of months at least. Discord also has a novel relay system (webhooks), which
doesn't quite fit into the generic relay mode.

The Google Chat rewrite is technically ready, but hasn't been tested properly
and unfortunately isn't scheduled to be finished in the near future.

## Meowlnir
Meowlnir didn't have many changes either, but there are some new API endpoints
in preparation for a future web management interface.
