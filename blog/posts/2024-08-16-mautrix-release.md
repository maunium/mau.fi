---
title: August 2024 releases & progress
summary: A summary of the mautrix releases in August 2024
slug: 2024-08-mautrix-release
tags:
- mautrix
- Bridges
- Matrix
---
This month's releases don't include as many bridges as last time, but the first
megabridge releases are here! Signal and Slack are now fully ported over to the
bridgev2 framework with the v0.7.0 and v0.1.0 releases respectively.

| Bridge/library | Version                                                          |
|----------------|------------------------------------------------------------------|
| mautrix-signal | [v0.7.0](https://github.com/mautrix/signal/releases/tag/v0.7.0)  |
| mautrix-slack  | [v0.1.0](https://github.com/mautrix/slack/releases/tag/v0.1.0)   |
| mautrix-go     | [v0.20.0](https://github.com/mautrix/go/releases/tag/v0.20.0)    |
| go-util        | [v0.7.0](https://github.com/mautrix/go-util/releases/tag/v0.7.0) |


Megabridge rewrites are primarily internal and for the most part bridges will
work exactly the same way as they did before. However, there are also a few new
features:

* Multi-account is supported on all megabridges. Users can simply log in
  multiple times and they'll get chats bridged from multiple accounts.
* Relay mode will disambiguate displaynames by default to prevent impersonation:
  if multiple users have displaynames that can be confused with each other,
  their user ID is appended to the name ("confusable" is defined by [UTS #39]).

[UTS #39]: https://www.unicode.org/reports/tr39/#Confusable_Detection

Additionally, the Slack bridge gained support for relay mode for the first time,
including support for logging in as a Slack app and using per-message names on
the Slack side.

[Go 1.23](https://go.dev/doc/go1.23) was also released a few days ago and
Go 1.21 is now EOL, which means that these latest releases require Go 1.22 or
higher to compile. As always, precompiled binaries don't have any dependencies.

## Megabridge progress
Telegram has unfortunately taken longer than expected and hasn't even reached
feature parity on Beeper, not to mention the features that Beeper doesn't use.

The Google Messages and Meta rewrites are currently being tested and will
likely be released next month.

WhatsApp, Twitter and Discord rewrites are also underway. This leaves Google
Chat, LinkedIn and iMessage. The rewrites for Google Chat and LinkedIn will be
started soon, but iMessage may get left behind for now as Beeper doesn't
support it anymore.

## Megabridge provisioning API
As [mentioned earlier](/blog/introducing-mautrix-manager/), I rewrote
mautrix-manager to use the new provisioning API in megabridge. The provisioning
API now also has a proper OpenAPI schema. The raw schema can be found in
[the repo](https://github.com/mautrix/go/blob/main/bridgev2/matrix/provisioning.yaml),
while a rendered version is available on spec.mau.fi: <https://spec.mau.fi/megabridge/>.

## Cryptography in mautrix-go
mautrix-go defaults to using libolm for implementing Matrix end-to-end
encryption. Since libolm is implemented in C, it's more difficult to compile
and generally not as nice as pure Go things, so the long-term goal has always
been to replace libolm with a pure Go implementation. That implementation
already exists in the form of goolm, but it's not yet ready for production,
as it still has issues where sessions get corrupted for unknown reasons.

In addition to the C-Go interface being suboptimal, libolm also has internal
issues. In particular, the implementations of cryptographic primitives (like
AES) used in libolm aren't cryptographically secure, meaning they're not
resistant against side-channel attacks.

The shortcomings in libolm were recently brought up again, including some
unwelcome FUD being spread prior to the actual disclosure. While the issues
aren't as critical as the fearmongering would suggest, they are real and libolm
has been formally deprecated now, so mautrix-go will also move away from it.
The plan is still to finish goolm and get it audited, but support for vodozemac
may be added in the nearer future as a short-term solution. You can follow
[mautrix/go#262](https://github.com/mautrix/go/issues/262) for updates.

mautrix-python will likely also switch to vodozemac at some point in the
future.
