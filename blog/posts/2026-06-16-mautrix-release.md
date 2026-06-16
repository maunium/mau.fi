---
title: June 2026 releases // Full-text search in gomuks
summary: A summary of the mautrix releases in June 2026
slug: 2026-06-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
---
This month's development was focused on gomuks, with new features like full-text
search including encrypted rooms, a proper desktop wrapper and more.

| Software          | Version                                                           |
|-------------------|-------------------------------------------------------------------|
| mautrix-whatsapp  | [v26.06](https://github.com/mautrix/whatsapp/releases/v0.2606.0)  |
| mautrix-telegram  | [v26.06](https://github.com/mautrix/telegram/releases/v0.2606.0)  |
| mautrix-twitter   | [v26.06](https://github.com/mautrix/twitter/releases/v0.2606.0)    |
| mautrix-signal    | [v26.06](https://github.com/mautrix/signal/releases/v0.2606.0)    |
| mautrix-slack     | [v26.06](https://github.com/mautrix/slack/releases/v0.2606.0)     |
| mautrix-meta      | [v26.06](https://github.com/mautrix/meta/releases/v0.2606.0)      |
| meowlnir          | [v26.06](https://github.com/maunium/meowlnir/releases/v0.2606.0)  |
| gomuks            | [v26.06](https://github.com/gomuks/gomuks/releases/v0.2606.0)     |
| mautrix-go        | [v0.28.1](https://github.com/mautrix/go/releases/v0.28.1)         |
| go-util           | [v0.9.10](https://github.com/mautrix/go-util/releases/v0.9.10)    |

## gomuks

### Electron wrapper
Previously, gomuks had a desktop wrapper written with Wails. The primary issue
was that Wails uses WebKitGTK on Linux, which is too buggy to use in practice.
Since there's no way to use Firefox/Gecko in a standalone desktop app the way
GeckoView allows making Android apps, the only solution was to switch to
Electron.

In addition to having a browser engine that actually works, the new Electron
wrapper brings multiple improvements:

* macOS builds are properly signed to allow installing without workarounds.
  It also has a built-in auto-updater.
* The app can be closed to tray, which still keeps embedded backends running,
  including desktop notifications.
* To enable multiple accounts in one app, it can run multiple embedded backends
  and also connect to remote backends. However, there's currently no UI for
  managing those, so you have to edit the config file manually to add accounts.
* Screen sharing should work on Linux with Wayland and macOS. Windows support
  will be added later (Windows is behind and doesn't have an OS-level screen
  picker yet, only Linux and macOS have that).

### Search
The backend now has full-text search powered by SQLite's FTS5 extension. It
covers all events that the backend has encountered, so if you run the backend
24/7, it'll have all new messages. Fetching old messages is not yet automated,
but you can use the archivemuks script to do it in individual rooms.

The search UI in the web frontend supports both local backend search and the
standard search API that calls the homeserver. By default, it uses the
homeserver for searching in unencrypted rooms and backend in encrypted rooms,
but you can of course choose which to use.

### Other changes
The devtools in the web interface now include a fully featured push rule editor,
which allows you to inspect and modify push rules however you like.

For people who want to make their own gomuks frontends, there's now slightly
better documentation for the [RPC](https://spec.mau.fi/gomuks/rpc.html) and
[HTTP](https://spec.mau.fi/gomuks/http.html) APIs.

## Bridges
Bridges didn't have many changes this month, though the upcoming months might
bring bigger changes in Metaland. It appears that in addition to [abandoning
encrypted chats on Instagram](https://www.bbc.com/news/articles/clypzxl3lvqo),
they've ditched the entire shared API. The current version of Instagram's web
client no longer uses the same lightspeed API as Facebook Messenger, which
could mean that the bridges have to be split up into mautrix-facebook and
mautrix-instagram again.

Telegram added a new rich text message format, which the bridge now mostly
supports.
