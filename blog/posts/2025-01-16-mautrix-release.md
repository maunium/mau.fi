---
title: January 2025 releases // gomuks android
summary: A summary of the mautrix releases in January 2025
slug: 2025-01-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
- gomuks
---
Most bridges didn't have any significant changes over the holidays, so only
Signal and Twitter have new releases. I accidentally became an Android
developer and made a GeckoView wrapper for gomuks web.

| Bridge/library    | Version                                                          |
|-------------------|------------------------------------------------------------------|
| mautrix-twitter   | [v0.2.1](https://github.com/mautrix/twitter/releases/v0.2.1)     |
| mautrix-signal    | [v0.7.5](https://github.com/mautrix/signal/releases/v0.7.5)      |
| mautrix-go        | [v0.23.0](https://github.com/mautrix/go/releases/v0.23.0)        |
| go-util           | [v0.8.4](https://github.com/mautrix/go-util/releases/v0.8.4)     |

## Megabridge progress
Most bridges didn't have any notable changes over the holidays, so I decided
to skip their releases. However, Signal had some breaking changes which
necessitated a release, and the Twitter bridge had multiple bugfixes worth
releasing too.

Beeper is currently focusing on shipping the existing megabridges as local
bridges, so there probably won't be much progress on the remaining rewrites
until March or April.

## gomuks web
As with the last 2 months, gomuks web received a bunch of new features.
Contributors this time included [@nexy7574](https://github.com/nexy7574),
[@sumnerevans](https://github.com/sumnerevans), and
[@everypizza1](https://github.com/everypizza1)

### gomuks android
Last month gomuks web got some small changes for mobile optimization, which
ended up making it quite a good mobile client. The obvious next step was to add
push notifications and other integrations with native features.

I was initially considering using web push to keep it as a pure web app.
However, Firefox and iOS browsers don't support all the fancier options for
notifications, and even the fancy options in Chrome aren't as good as a native
app. Therefore, I decided to make a native wrapper app which uses GeckoView for
the web side, while handling push notifications from FCM natively in Kotlin.

The app can be found in the [gomuks/android](https://github.com/gomuks/android)
repo. It doesn't have any precompiled builds yet, but should work to some extent
if you build it yourself. It can already receive notifications as well as
auto-dismiss them when they're read from another device. The next features I'm
planning to add are showing media and avatars in notifications and then
providing a share target to send media using the OS share feature.

For now, gomuks android requires the backend to run elsewhere, because that's
how I use it. In the future, there may be an option to run a backend bundled
with the app. The backend running 24/7 allows notifications to be fully reliable
with minimal battery usage. All information is sent in an encrypted FCM payload,
which means the app doesn't need to fetch extra data to display notifications.

### User styles
Custom CSS is already supported, but this month gomuks web's custom CSS feature
gained support for `@import` statements. To make such statements easier to use,
I made a simple service where user styles can be shared: [css.gomuks.app](https://css.gomuks.app/).

### New features
A mostly complete list of new features in gomuks web:

* Spaces
* Rendering read receipts
* Joining rooms (via link and accepting invites) and leaving rooms
* Sticker picker to complete [MSC2545] support
  (custom emojis were already supported before)
* [MSC4095] URL preview rendering
* [MSC4144] per-message profile rendering
* Option for Discord-like compact replies
* Better message context menu for small screens (mobile)

[MSC2545]: https://github.com/matrix-org/matrix-spec-proposals/pull/2545
[MSC4095]: https://github.com/matrix-org/matrix-spec-proposals/pull/4095
[MSC4144]: https://github.com/matrix-org/matrix-spec-proposals/pull/4144
