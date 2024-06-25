---
title: Project updates (and a blog)
summary: A summary of things that have happened in the mautrix projects
slug: 2024-h1-mautrix-updates
tags:
- mautrix
- hicli
- Bridges
- Stickerpicker
- Direct Media
- Matrix
- Go
---
I haven't been posting on <abbr title="This Week In Matrix">TWIM</abbr> for a
long time and too many things have piled up to post everything in TWIM. This
post recaps the potentially interesting things that I've worked on in the past
~half a year.

Also, since Automattic is basically a blogging company, I decided to try
making a blog. Can't promise there'll be any more posts though,
[last time](https://maunium.net/blog) I got to a grand total of three.

## Next generation of mautrix bridges
mautrix bridges have historically involved a lot of duplicate code across the
different bridges. Each bridge has had a very similar structure, but there were
lots of small network-specific differences, which made it difficult to share
code. This also made it difficult to build new bridges without having extensive
knowledge of how the existing bridges worked.

Recently, Beeper decided to move all bridges to run in clients instead of the
cloud. Client-side bridges have lots of similarities to cloud bridges, but they
also have some differences, quite similar to how different networks have slight
differences. Due to the differences, our initial proof-of-concept local signal
bridge was only able to reuse the message conversion code and none of the other
bridge code.

### Introducing Megabridge
To solve both problems at once, we started the megabridge project. The primary
output is the new [bridgev2 module] in mautrix-go, which contains all the
generic bridging code and defines interfaces for both remote network connectors
and Matrix connectors.

[bridgev2 module]: https://github.com/mautrix/go/tree/master/bridgev2

The name "megabridge" reflects the additional possibilities enabled by the
interface approach: one could build a binary that runs multiple bridges in one
process simply by creating multiple instances of the relevant interfaces. In the
old architecture, bridges were strictly singletons. Beeper also has a project
called "megahungry", which applies the same idea to hungryserv, our unfederated
Matrix homeserver. Beeper runs hungryserv and bridges in a single-tenant model,
but the new megahungry and megabridge projects will allow running multiple
single-tenant instances in one process.

### Matrix connectors
Two Matrix connectors have been written: one inside the Beeper app for local
bridges, and another in mautrix-go for connecting to a standard Matrix server as
an appservice. When using the appservice connector, the bridge will work exactly
like existing bridges do.

In the future, a third connector could use [MSC4144] instead of an appservice,
enabling self-hosting bridges against any homeserver. MSC4144 allows a bot to
specify the displayname and avatar separately for each message, which is the
main blocker for running bridges using a single bot account. There would still
be some limitations compared to appservices, e.g. you wouldn't be able to view
the remote member list or get read receipts, but it's better than nothing.

[MSC4144]: https://github.com/matrix-org/matrix-spec-proposals/pull/4144

### Network connectors
The current bridges will be replaced by implementations of the network connector
interface, which effectively means all existing bridges will be rewritten. The
Signal bridge already includes a v2 connector with all the basic functionality.

Rewrites of the Python bridges (Telegram, Twitter, Google Chat, LinkedIn) are
underway, and the other Go bridges will follow later. I'd expect somewhat
functional network connectors for all networks within a few months, although
achieving feature parity, especially for features that aren't used in Beeper,
will take longer.

### Writing new bridges
The interfaces are still evolving, but the current features are already
functional. If you want to build new bridges, it should already be much easier
to get started with the new architecture than the old one. See the
[unorganized-docs directory] in mautrix-go for more details and docs. The docs
will eventually be organized into docs.mau.fi and/or pkg.go.dev. You can also
join the [#go:maunium.net] room to ask questions.

At some point in the near-ish future I'll try writing a bridge to a simple
network (perhaps Twilio) as an example and make a blog post showcasing the
process.

[unorganized-docs directory]: https://github.com/mautrix/go/tree/master/bridgev2/unorganized-docs
[#go:maunium.net]: https://matrix.to/#/#go:maunium.net

## Old bridges

### mautrix-signal
The Go rewrite of mautrix-signal was finally released in February after being in
development for several months. The old bridge is fairly broken by now due to
Signal API changes (e.g. signald can't link new devices at all anymore), so I
assume most people have already updated.

### mautrix-meta
The legacy Facebook and Instagram bridges were replaced by a new Meta bridge
which uses the web app API. The API happens to be pretty much the same for both
FB and IG, which is why the bridges were merged. The old bridges may still work
to some extent, but they're not being maintained anymore, so they'll break more
and more as Meta stops supporting the old APIs.

### mautrix-slack
After Beeper joined Automattic, we had to start using Slack when talking to
Automattic people (we still use Matrix/Beeper for internal chats), which means
I had to actually care about the Slack bridge.

I was planning on [refactoring](https://github.com/mautrix/slack/commits/refactor)
the bridge to bring it up to the level of the other Go bridges, but
unfortunately didn't manage to finish the refactor before starting the
megabridge project. At this point, it is unlikely that the refactor would be
finished before the bridge is rewritten to use the new architecture.

mautrix-discord is similarly stuck on a slightly older mautrix-go version, so it
probably won't get any more releases until the megabridge rewrite.

### mautrix-python
With all bridges moving to Go, mautrix-python has seen less activity. It will
still be maintained, but the bridge module will likely be deprecated or even
removed once the Go rewrites are ready. maubot still uses mautrix-python and is
not going to be rewritten in Go. In fact, the [original version of maubot] was
written in Go, but it was rewritten in Python due to the lack of easily
reloadable plugins. As far as I know, the [plugin package of Go] has not
improved much since then.

[original version of maubot]: https://github.com/maubot/maubot/tree/f06c6dd7676771e740e665bb22ac5d42762fcc62
[plugin package of Go]: https://pkg.go.dev/plugin

## High-level client framework
I've recently started developing [hicli], an opinionated high-level module for
building clients. It manages room timeline storage, handles everything related
to end-to-end encryption, and more. The storage layer already works to some
extent, but more complex features like unread counts are not yet implemented.
Additionally, it hasn't been tested in an actual client yet. Since it's a side
project, it'll probably take a while to produce anything actually useful.

[hicli]: https://github.com/mautrix/go/tree/master/hicli

There are two reasons for this project. The first is that due to the atrocious
performance of Element Web, I've wanted to make my own web client for a while
now. Making a new Matrix SDK in JavaScript is a non-starter, so instead I looked
into ways to use Go. Unfortunately, running SQLite in WASM with OPFS doesn't
seem performant enough currently, so a pure web client is not feasible yet.
Instead, I plan to write a separate daemon that a web client can communicate
with over local HTTP, or perhaps through a webextension. The web client and
daemon could also be bundled into a desktop app using [Wails](https://wails.io/).

The other reason is that gomuks has been neglected for quite a while, primarily
because I don't actually use terminal clients, but also because the data layer
was grown organically and is rather bad. If I can make a good data layer for my
web client, I can hopefully also rewrite gomuks on top of that and maintain
both more actively.

Based on my experience of building an unfederated Matrix homeserver, SQLite is
very powerful and extremely fast if used correctly. If it can handle being the
database for a homeserver, it can certainly handle a client. To avoid some of
the mistakes in gomuks, hicli uses SQLite for storing everything and even uses
some of the more advanced SQL features like triggers.

## Gifs in maunium-stickerpicker
The sticker picker is more or less ready. It could still use a server component
for easy multi-user setup, but I'd hope that native stickers will happen within
a year or two, so spending a lot of extra effort on the sticker picker isn't
worth it. However, [@hedgenischay] recently contributed support for sending gifs
via Giphy, which means Element Web can now have a real gif picker.

[@hedgenischay]: https://github.com/hegdenischay

An interesting detail of the gif picker is that it doesn't reupload gifs.
Instead, it simply sends `mxc://giphy.mau.dev/*` URIs, which redirect to
`i.giphy.com`. More on that in the direct media section below.

The feature is enabled by default in the latest version and doesn't require any
extra configuration, although it is possible to override the gif proxy and
giphy API key if desired.

## Direct media access
mautrix-discord has had direct media access for a while now. The general idea
is that instead of reuploading media, the bridge will generate a `mxc://` URI
which actually redirects to the Discord CDN. This means less wasted disk space,
and if your homeserver caches it as remote media, it's easy to clean up.

The original implementation of direct media was simple enough that the redirects
could be implemented inside a reverse proxy. However, around February, Discord
started using expiring signed download links, which prevented a pure reverse
proxy implementation. A new version of direct media was built into the bridge
to support refreshing the download link if it expires. The new version also had
extra fun stuff, like implementing the `/_matrix/key/*` endpoints to pass the
federation tester.

More recently, [MSC3916] was accepted to add authentication to all media
download endpoints. The original implementation would've blocked direct media
redirects, but a late change shortly before accepting re-added proper redirect
support.

[MSC3916]: https://github.com/matrix-org/matrix-spec-proposals/pull/3916

Even though it allows redirects, MSC3916 still makes the federation download
endpoint more complicated, because the response must be a `multipart/mixed`.
This makes it harder to implement purely in a reverse proxy, so a small proxy
service is usually required. To make it easier to make such small services, I
moved the direct media code from mautrix-discord to mautrix-go:
<https://github.com/mautrix/go/blob/master/mediaproxy/mediaproxy.go>.

The `giphy.mau.dev` proxy mentioned in the sticker picker section uses the new
media proxy module in mautrix-go: <https://github.com/maunium/stickerpicker/blob/master/giphyproxy/main.go>.

In addition to mautrix-discord and the sticker/gif picker, the bridgev2 project
will include an interface to let network connectors implement direct media
access using the mediaproxy module.
