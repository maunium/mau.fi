---
title: (Re-)Introducing mautrix-manager
summary: A desktop app to make logging into bridges easier
slug: introducing-mautrix-manager
tags:
- mautrix
- Bridges
---
All mautrix bridges are primarily meant for puppeting, which means users need
to log into the bridges with their credentials. The difficulty of doing that
varies by network: some just require scanning a QR code, while others require
logging into the official website and grabbing cookies using browser devtools.

Long ago I made a web app called [mautrix-manager](https://github.com/tulir/mautrix-manager)
which could be used to log into many different bridges. It used to be embedded
in the Beeper desktop app, but later Beeper built bridge management directly
into the desktop app. There also wasn't much demand for a management interface
among self-hosting users at the time, and it required effort to maintain, as
each bridge had a slightly different API for logging in. Only 3 of the 8
bridges it supported are still in use.

## Enter megabridge
One of the many benefits of megabridge was that it added a new unified
abstraction for describing login flows. Instead of each bridge having its own
flavor of login commands and provisioning API, there's just one implementation
of both. The network connectors implement their login flows using three generic
step types: `user_input`, `cookies` and `display_and_wait`.

By combining those steps, even complicated login flows can be achieved.
For example, the Google account pairing in Google Messages uses `cookies`
(Google account login) followed by `display_and_wait` (tap emoji on phone).
Bridges can also offer multiple login flows, Google Messages offers QR and
Google account as separate flows.

## The new manager
Due to the increase in the number of bridges requiring cookie extraction and
thanks to the new generic login APIs in megabridge, I decided to rewrite
mautrix-manager to support logging into any megabridge.

I chose Electron as the framework for the app for two reasons:

1. A web app can't extract cookies. A desktop app or a browser extension is
   required, and a desktop app is easier.
2. Electron alternatives like Tauri and Wails use the system webview, which is
   often worse than Chromium. Some login pages get picky about the browser used
   to log in, so a high-quality webview is important.
   * (if there was an Electron equivalent built on Firefox, I would've used it)

Note that the usefulness of the new manager is still somewhat limited, as it
only works with megabridges. Once all the bridge rewrites are complete, it will
work with any mautrix bridge. For now, the Signal and Slack megabridges are
ready, Google Messages mostly works and Meta is coming up next.

You can find the new mautrix-manager at <https://github.com/mautrix/manager>.

![preview](/blog/res/mautrix-manager.png)
