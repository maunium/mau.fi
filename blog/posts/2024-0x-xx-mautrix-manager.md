---
title: (Re-)Introducing mautrix-manager
summary: A desktop app to make logging into bridges easier
slug: introducing-mautrix-manager
tags:
- mautrix
- Bridges
draft: true
---
All mautrix bridges are primarily meant for puppeting, which means users need
to log into the bridges with their credentials. The difficulty of doing that
varies by network: some just require scanning a QR code, while others require
logging into the official website and grabbing cookies using browser devtools.

Long ago I made a web app called [mautrix-manager](https://github.com/tulir/mautrix-manager)
which could be used to log into many different bridges. It was also originally
embedded in the Beeper desktop app, but ended up being deprecated since Beeper
built bridge management directly into the desktop app and there wasn't much
demand for a management interface among self-hosting users at the time. It also
required effort to maintain, as each bridge had a slightly different API for
logging in. Only 3 of the 8 bridges it supported are still in use.

One of the many benefits of megabridge was that it added a new unified
abstraction for describing login flows. Instead of each bridge having its own
flavor of login commands and provisioning API, there's just one implementation
of both. The network connectors implement their login flows using three generic
step types: `user_input`, `cookies` and `display_and_wait`.

Due to the increase in the number of bridges requiring cookie extraction and
thanks to the new generic login APIs in megabridge, I decided to rewrite
mautrix-manager to support logging into any megabridge. Extracting cookies from
another website in a web app is obviously not possible, so the new manager is a
desktop app that uses Electron.

Electron specifically was chosen over more trendy alternatives to ensure that
the app can create login webviews in an high-quality browser engine. Tauri and
Wails use WebKitGTK or similar engines on Linux, which have lots of issues.
Of course a bundled Firefox webview would've been preferable, but unfortunately
no such thing exists.

You can find the new mautrix-manager at <https://github.com/mautrix/manager>.
