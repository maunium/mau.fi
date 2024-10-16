---
title: October 2024 releases & progress
summary: A summary of the mautrix releases in October 2024
slug: 2024-10-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
- Meowlnir
---
As expected, WhatsApp has been upgraded to a megabridge. Other bridges have
received minor bugfixes, but nothing too exciting. In non-bridge news, hicli
has finally progressed again.

| Bridge/library    | Version                                                          |
|-------------------|------------------------------------------------------------------|
| mautrix-gmessages | [v0.5.1](https://github.com/mautrix/gmessages/releases/v0.5.1)   |
| mautrix-whatsapp  | [v0.11.0](https://github.com/mautrix/whatsapp/releases/v0.11.0)  |
| mautrix-signal    | [v0.7.2](https://github.com/mautrix/signal/releases/v0.7.2)      |
| mautrix-slack     | [v0.1.2](https://github.com/mautrix/slack/releases/v0.1.2)       |
| mautrix-meta      | [v0.4.1](https://github.com/mautrix/meta/releases/v0.4.1)        |
| mautrix-go        | [v0.21.1](https://github.com/mautrix/go/releases/v0.21.1)        |
| go-util           | [v0.8.1](https://github.com/mautrix/go-util/releases/v0.8.1)     |
| meowlnir          | [v0.2.0](https://github.com/maunium/meowlnir/releases/v0.2.0)    |

## Megabridge progress
The Telegram bridge rewrite has mostly reached feature parity with the features
that Beeper uses, but there hasn't been any progress on non-Beeper features,
such as relaybots, so there's no ETA for the full release yet.

The Twitter rewrite has been resumed and may be ready by next month.
Unfortunately Discord hasn't progressed at all, so it'll likely take until
December at least, possibly even next year.

## gomuks / hicli
As mentioned in my [first post], I'm working on a high-level client framework
to build a new web client and to replace the internals of gomuks. Back when I
wrote that, Element Web was using 2-3 GB of ram. Now, it uses 3-6 GB and
freezes completely several times a day, so making a replacement has become more
urgent.

[first post]: https://mau.fi/blog/2024-h1-mautrix-updates/#high-level-client-framework

The new web client now exists. Inventing a new name was too much effort, so I
decided to call it gomuks web. You can find it in the [webmuks branch] (CI
binaries are available on mau.dev as always).

[webmuks branch]: https://github.com/tulir/gomuks/tree/webmuks

Like I said in the first post, the new client isn't exactly a true web client,
as it communicates with a Go backend that's running locally (although
technically nobody forces it to be local: you could perfectly well run the
backend on a raspberry pi or something).

The frontend is stateless: it receives the room list from the backend on
connect and then receives new events over websocket as the backend receives
them from sync. Since the backend knows the exact order of rooms and their
preview messages, it only takes a fraction of a second to send the first 100
rooms to the frontend, which means there's no waiting time when reopening the
frotend.

It's significantly faster than Element Web and it needs less than 250 MB of ram
for both the backend and the browser tab combined to handle my main account. In
terms of features, I still need to implement unread counts, reactions and maybe
spaces, but otherwise it's already usable as a basic client.

The old version of gomuks is now called gomuks legacy. At some point in the
future, I'll copy the old terminal UI over to the hicli-based version and call
it gomuks terminal. I might also eventually wrap the web client in a desktop
and/or mobile app. The web UI is not mobile optimized, but performance-wise it
seems decent even on mobile browsers.

## Meowlnir
Meowlnir now supports receiving reports via the standard Matrix [/report endpoınt].
In addition to just forwarding reports to the management room, bot admins can
also send commands to the bot in the report message, which allows quickly
banning spammers without having to switch to the management room.

Other than that, there were a couple of bug fixes around evaluating policies,
as well as a feature to notify the management room if a user pings the bot in
a protected room.

[/report endpoint]: https://spec.matrix.org/v1.12/client-server-api/#post_matrixclientv3roomsroomidreporteventid
