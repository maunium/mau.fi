---
title: December 2024 releases // Bluesky, gomuks desktop & mobile optimization
summary: A summary of the mautrix releases in December 2024
slug: 2024-12-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
---
There's a Wails wrapper for gomuks web now, plus mobile optimization and a lot
of features. A new Bluesky bridge was released, and the Twitter bridge was
upgraded to a megabridge.

| Bridge/library    | Version                                                          |
|-------------------|------------------------------------------------------------------|
| mautrix-gmessages | [v0.6.0](https://github.com/mautrix/gmessages/releases/v0.6.0)   |
| mautrix-whatsapp  | [v0.11.2](https://github.com/mautrix/whatsapp/releases/v0.11.2)  |
| mautrix-bluesky   | [v0.1.0](https://github.com/mautrix/bluesky/releases/v0.1.0)     |
| mautrix-twitter   | [v0.2.0](https://github.com/mautrix/twitter/releases/v0.2.0)     |
| mautrix-discord   | [v0.7.2](https://github.com/mautrix/discord/releases/v0.7.2)     |
| mautrix-signal    | [v0.7.4](https://github.com/mautrix/signal/releases/v0.7.4)      |
| mautrix-slack     | [v0.1.4](https://github.com/mautrix/slack/releases/v0.1.4)       |
| mautrix-meta      | [v0.4.3](https://github.com/mautrix/meta/releases/v0.4.3)        |
| mautrix-go        | [v0.22.1](https://github.com/mautrix/go/releases/v0.22.1)        |
| go-util           | [v0.8.3](https://github.com/mautrix/go-util/releases/v0.8.3)     |

## Megabridge progress
Bluesky is quickly gaining popularity and their current temporary chat protocol
is very simple, so I decided to go ahead and write a bridge: <https://bsky.app/profile/tulir.fi/post/3lbv3cd7q3223>.
It took a day to write, it's nearly fully featured and has around 1k lines of code.

Twitter testing was expected to be stalled like last month, but then Twitter
deleted some of the APIs the old bridge was using, so the new bridge has now
been tested and released.

Other than that, there haven't been any particularly exciting developments on
existing bridges. The Google Chat and LinkedIn rewrites are underway, but
progressing fairly slowly.

## gomuks web
gomuks web has received a bunch of new features again, including some
contributions from [@JadedBlueEyes](https://github.com/JadedBlueEyes)
and [@sumnerevans](https://github.com/sumnerevans).

### Windows builds with `zig cc`
The CI now builds binaries for Windows using `zig cc` for the C parts (libolm
and sqlite). I originally planned to use zig for all cgo cross-compiling, but
it turned out to be 2-3x as slow as my old musl.cc-based builder, so I reverted
to the old system for Linux and kept `zig cc` for Windows.

Overall, `zig cc` definitely makes cross-compiling cgo easier, but I prefer
speed and there was nothing wrong with the old builder for Linux.

### gomuks desktop & mobile
I decided to try wrapping the whole thing (backend and frontend) in a desktop
app using [Wails](https://wails.io/). It seems to work to some extent (at least
on Linux), so if you prefer desktop apps, give it a try. The CI builds desktop
binaries for all three OSes.

In addition to the desktop app, mobile support has been improved with features
like back button support, room switching animations and auto-reconnect. I'm
also considering adding web push support and/or making a full wrapper app for
the frontend.

### New features
A mostly complete list of new features:

* Custom CSS
* Preferences, like disabling read receipts/typing notifications, showing media
  previews and more
  * All preferences can be toggled per room and per device
* Auto-reconnect when losing connection to backend
* Auto-reload to update the frontend when the backend is updated
* Room topics in the room view header
* Viewing shared rooms (MSC2666) and device list in user info
* `matrix:` URI handling
* Rendering typing notifications (below a Discord-style floating composer)
* Support for Chrome and Safari (Firefox is still the primary target)
* Rendering blurhashes and media spoilers
* Sending and rendering location messages (with Leaflet+OSM or Google Maps)
* A gif picker (both Tenor and Giphy are supported)
* Navigation via browser back button
