---
title: May 2025 releases // archivemuks and other gomuks web updates
summary: A summary of the mautrix releases in May 2025
slug: 2025-05-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
- Meowlnir
---
This month's releases mostly consist of some bridge bugfixes, a few new features
in Meowlnir and many features in gomuks web.

| Bridge/library    | Version                                                          |
|-------------------|------------------------------------------------------------------|
| mautrix-gmessages | [v0.6.2](https://github.com/mautrix/gmessages/releases/v0.6.2)   |
| mautrix-whatsapp  | [v0.12.1](https://github.com/mautrix/whatsapp/releases/v0.12.1)  |
| mautrix-twitter   | [v0.4.1](https://github.com/mautrix/twitter/releases/v0.4.1)     |
| mautrix-signal    | [v0.8.3](https://github.com/mautrix/signal/releases/v0.8.3)      |
| mautrix-go        | [v0.24.0](https://github.com/mautrix/go/releases/v0.24.0)        |
| meowlnir          | [v0.5.0](https://github.com/maunium/meowlnir/releases/v0.5.0)    |
| go-util           | [v0.8.7](https://github.com/mautrix/go-util/releases/v0.8.7)     |

## Bridges
Signal has basic support for direct media access now. It's a bit safer to use
than WhatsApp, but the links will only work for 45 days and will be permanently
broken after that, so for most use cases, traditional re-uploaded media is
likely still better.

The Twitter bridge got a few bug fixes as well as support for voice messages
in both directions.

The Signal and WhatsApp bridges turned out to have the exact same bug
implemented in different ways which broke backfilling exisitng portals after
re-logining. Backfilling old messages is not possible as there's no standard
mechanism in Matrix to do it, but backfilling missed messages when no new
messages have been sent should work again.

## Meowlnir
Meowlnir received some new tools for public servers, such as automatically
suspending local accounts when a ban policy is received. [@nexy7574] also made
a couple pull requests, like customizable auto-redact ban reasons and a fix for
the fallback redaction mechanism not deleting state events.

In addition to new features, the command handling system was redone on top of
mautrix-go's new generic command processing framework. There shouldn't be any
user-facing changes, but if you want to write your own bots in Go, the new
command framework will likely be useful.

## gomuks web
There have been a bunch of changes to gomuks web, like sending URL previews,
re-encoding uploads and mass-redacting messages. In addition to changes to web,
there's now a Go client package for using the gomuks websocket RPC API that the
web frontend uses.

Using the new RPC package, I made a simple CLI tool called [archivemuks] which
connects to the backend, asks it to fetch all messages in a single room, and
saves them to a JSON file. The long-term plan is to have background room history
crawling built into the backend to enable filling the local database for search
and export purposes, but the CLI tool works as a proof-of-concept. Because the
tool doesn't wait for decryption, you'll have to delete the export and re-run
the tool after the history has been downloaded once, otherwise your export will
have a bunch of undecrypted events.

[archivemuks]: https://github.com/tulir/gomuks/tree/main/cmd/archivemuks

Full changelog:

* ACL changes are now rendered in a `<details>` tag to avoid filling up the
  timeline for the first ACL event that adds lots of blocked servers
  (contributed by [@sumnerevans]).
* I also finally merged Sumner's PR for generating URL previews and bundling
  them in sent messages with [MSC4095]. Bundled previews can be seen by other
  gomuks users and will also be bridged by most mautrix bridges.
* [@nexy7574] contributed a button to remove all recent messages from a user.
  * The ban dialog now also offers removing messages as an option.
* A new media upload dialog was added, which allows previewing media and
  resizing/re-encoding it before uploading.
  * The dialog can be disabled in setting if you prefer the old immediate upload style.
  * Re-encoding videos requires ffmpeg to be installed and discoverable in `$PATH`.
  * Dumb formats like HEIC have re-encoding enabled by default to ensure you don't
    accidentally send images that nobody can view.
* Thumbnail generation for avatars and upload re-encoding will now apply JPEG
  orientation correctly (it won't preserve the raw metadata, but will rotate the
  pixels if necessary).
* The room view header buttons have been moved into a three dot menu on mobile,
  which means you can actually see the room name again.
* You can now click on reactions in the timeline to react with the same emoji.
  However, gomuks does not currently track whether you have reacted with a given
  emoji, so you can't un-react by clicking yet.
* There's a new button in settings for importing the entire server-side key
  backup. It may be useful if you want to export keys to file. By default,
  gomuks will only import keys as necessary rather than loading the entire
  backup.
* Sending mentions in captions has been fixed to actually include the mention
  in `m.mentions`.
* Apparently some weird people have dozens of emoji/sticker packs, so the
  emoji/sticker picker now has a scroll bar if there are more than 3 rows of
  packs.

[@sumnerevans]: https://github.com/sumnerevans
[@nexy7574]: https://github.com/nexy7574
[MSC4095]: https://github.com/matrix-org/matrix-spec-proposals/pull/4095
