---
title: Writing a Twilio bridge
summary: Showcasing the new bridgev2 module in mautrix-go by writing a Twilio bridge
slug: megabridge-twilio
tags:
- mautrix
- Bridges
- Matrix
- Go
draft: true
---

## Getting started with a new Go project
The first step to make a new Go project is to create a new directory and run
`go mod init <module path>`, where `<module path>` is the import path of your
new module (for example, the GitHub repo). I have custom import paths, so even
though I use GitHub for hosting, I run `go mod init go.mau.fi/mautrix-twilio`.

When naming your bridge, please make up your own name and don't use `mautrix-*`.

The next step is creating the network connector itself. The connector is
effectively a library, so I put it in `pkg/connector/`.

## The connector itself
Create a file called `connector.go` and define a struct called `TwilioConnector`.

Then, I add a line like this:

```go
var _ bridgev2.NetworkConnector = (*TwilioConnector)(nil)
```

This is a conventional method of ensuring that a struct implements a specific
interface. The line is creating a value of type `*TwilioConnector`, then
assigning it to a variable of type `bridgev2.NetworkConnector`, which will only
compile if the value implements the interface. The variable name is the blank
identifier `_`, which means the value will be discarded after being evaluated.

If you're using a smart editor, it should complain that `*TwilioConnector` does
in fact not implement `bridgev2.NetworkConnector`, and possibly even offer you
a quick way to create stub methods to implement the interface.

### `Init`
The Init function is called when the bridge is initializing all types. It also
gives you access to the bridge struct, which you'll need to store for later.

This function should not do any kind of IO or other complicated operations, it
should just initialize the in-memory struct.

```go
type TwilioConnector struct {
	br *bridgev2.Bridge
}

func (tc *TwilioConnector) Init(bridge *bridgev2.Bridge) {
	tc.br = bridge
}
```

### `Start`
The Start function is called slightly later in the startup. This can be used
for bridge-wide IO operations, such as upgrading database schemas if the
connector needs its own database tables.

In the case of Twilio, there's no need for special database tables, but we do
need to register some routes, as receiving messages requires a webhook. Other
networks that receive events via websockets/polling/etc may not need to do
anything at all here.

```go
func (tc *TwilioConnector) Start(ctx context.Context) error {
	server, ok := tc.br.Matrix.(bridgev2.MatrixConnectorWithServer)
	if !ok {
		return fmt.Errorf("matrix connector does not implement MatrixConnectorWithServer")
	} else if server.GetPublicAddress() == "" {
		return fmt.Errorf("public address of bridge not configured")
	}
	r := server.GetRouter().PathPrefix("/_twilio").Subrouter()
	r.HandleFunc("/{loginID}/receive", tc.ReceiveMessage).Methods(http.MethodPost)
	return nil
}

func (tc *TwilioConnector) ReceiveMessage(w http.ResponseWriter, r *http.Request) {}
```

We'll come back to `ReceiveMessage` later

### `GetCapabilities`
The `GetCapabilities` function on the network connector is used to signal some
bridge-wide capabilities, like disappearing message support. Twilio doesn't
have any of the relevant features, so we'll just leave this empty:

```go
func (tc *TwilioConnector) GetCapabilities() *bridgev2.NetworkGeneralCapabilities {
	return &bridgev2.NetworkGeneralCapabilities{}
}
```

### `GetName`
The `GetName` function is used to customize the name of the bridge.

```go
func (tc *TwilioConnector) GetName() bridgev2.BridgeName {
	return bridgev2.BridgeName{
		DisplayName:      "Twilio",
		NetworkURL:       "https://twilio.com",
		NetworkIcon:      "mxc://maunium.net/FYuKJHaCrSeSpvBJfHwgYylP",
		NetworkID:        "twilio",
		BeeperBridgeType: "go.mau.fi/mautrix-twilio",
	}
}
```

### `GetConfig`
Network connectors can define their own config fields, which for normal bridges
using `mxmain` will be in the `network:` section of the config.

The `GetConfig` function returns all data that is needed to provide the config.

* `example` is the example config.
* `data` is a pointer to the object where the config should be decoded to.
* `upgrader` is a helper to perform config upgrades.

On startup, the bridge will read the user's config file as well as the example
config, then call the upgrader to copy values from the user's config into the
example, and finally overwrite the user's config with the example. Users can
disable the overwriting part if they don't like it, but the first two steps are
done in any case. There are two benefits to this system:

* the bridge doesn't need to have any backwards-compatibility for the config
  outside the upgrader function. The upgrader can simply copy fields from old
  locations into new ones.
* the user can easily get an upgraded config without having to manually figure
  out which fields have changed.

The Twilio connector doesn't need any special fields, so we can just return nil
values:

```go
func (tc *TwilioConnector) GetConfig() (example string, data any, upgrader configupgrade.Upgrader) {
	return "", nil, configupgrade.NoopUpgrader
}
```

If you did want config fields, the response would look something like this:

```go
type TwilioConnector struct {
	...
	Config Config
}

type Config struct {
	MyFancyField string `yaml:"my_fancy_field"`
}

//go:embed example-config.yaml
var ExampleConfig string

func upgradeConfig(helper configupgrade.Helper) {
	helper.Copy(configupgrade.Str, "my_fancy_field")
}

func (tc *TwilioConnector) GetConfig() (example string, data any, upgrader configupgrade.Upgrader) {
	return ExampleConfig, &tc.Config, configupgrade.SimpleUpgrader(upgradeConfig)
}
```

and you'd have `pkg/connector/example-config.yaml` with

```yaml
# A description of the field
my_fancy_field: this is the default value
```

### `LoadUserLogin`
`LoadUserLogin` is called when the bridge wants to prepare an existing login
for connection. This is where the `NetworkAPI` interface comes in: the primary
purpose is to fill the `Client` property of `UserLogin` with the network
client. This function should not do anything else, actually connecting to the
remote network (if applicable) happens later in `NetworkAPI.Connect`.

In the case of Twilio, we'll initialize the go-twilio client here:

```go
func (tc *TwilioConnector) LoadUserLogin(ctx context.Context, login *bridgev2.UserLogin) error {
	restClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username:   login.Metadata.Extra["api_key"].(string),
		Password:   login.Metadata.Extra["api_secret"].(string),
		AccountSid: login.Metadata.Extra["account_sid"].(string),
	})
	login.Client = &TwilioClient{
		UserLogin: login,
		Twilio:    restClient,
	}
	return nil
}
```

## The network API
Next, we'll need to actually define `TwilioClient` and implement the `NetworkAPI`.

```go
type TwilioClient struct {
	UserLogin *bridgev2.UserLogin
	Twilio    *twilio.RestClient
}

var _ bridgev2.NetworkAPI = (*TwilioClient)(nil)
```

Like with the network connector, we'll do the same interface implementation
assertion, In this case it's not technically necessary, as we're already
assigning a `&TwilioClient{}` to `UserLogin.Client`, which does the same
check. However, if you implement any of the optional extra interfaces, then the
explicit assertions become very useful, as there's nothing else ensuring the
correct functions are implemented.

### `Connect` and `Disconnect`
Since Twilio uses webhooks, these two methods don't need to do anything. For
networks with a persistent connection, this is where you'd set up the connection
and disconnect it.

```go
func (tc *TwilioClient) Connect(ctx context.Context) error { return nil }
func (tc *TwilioClient) Disconnect()                       {}
```

### `IsLoggedIn`
On some networks, logins can be invalidated, so this function is used to check
if the login is still valid. For Twilio, we'll just return true. Note that this
method is not meant to do any IO, it should just return cached values.

```go
func (tc *TwilioClient) IsLoggedIn() bool {
	return true
}
```

### `LogoutRemote`
This method is meant to invalidate remote network credentials and disconnect
from the network. Since Twilio doesn't have credentials that can be invalidated,
nor a persistent connection that can be disconnected, we don't need to do
anything here.

```go
func (tc *TwilioClient) LogoutRemote(ctx context.Context) {}
```

### `IsThisUser`
TODO

### `GetChatInfo`
TODO

### `GetUserInfo`
TODO

### `GetCapabilities`
TODO

## The login process
In the first section about the network connector, we skipped the `GetLoginFlows`
and `CreateLogin` functions, so let's get back to those.

`GetLoginFlows` returns the ways that can be used to log into the bridge. Just
an internal ID, a human-readable name and a brief description. The internal ID
of the flow the user picked is then passed to `CreateLogin`. The return type of
`CreateLogin` is the third and final primary interface, `LoginProcess`.

Login process is meant to be a simple abstraction over login flows to arbitrary
remote networks. It has three different step types that can hopefully be used to
build any login flow there is:

* User input: fairly self-explanatory, ask the user to give values for one or
  more fields.
* Cookies: display a webview to the user and extract cookies, localStorage or
  other things after completion. For non-graphical logins (like using the
  bridge bot), this will ask the user to manually go to the website and extract
  the relevant values. If only cookies are necessary, the extraction can be
  done by using the "Copy as cURL" feature in browser devtools and pasting the
  result to the bridge bot.
* Display and wait: display something to the user, then wait until the remote
  network returns a response. This is used for things like QR logins or other
  flows where the user has to do something on an existing login.

Every step also has an internal identifier (reverse java package naming style
is recommended), general instructions for the entire step, and type-specific
parameters.

In addition to the three real step types, there's a fourth special type
indicating the login was successful.

### `GetLoginFlows` & `CreateLogin`

In the case of Twilio, the user just needs to provide their API keys, so we'll
use the user input type. First, we'll implement the two functions in the network
connector:

```go
func (tc *TwilioConnector) GetLoginFlows() []bridgev2.LoginFlow {
	return []bridgev2.LoginFlow{{
		Name:        "API key & secret",
		Description: "Log in with your Twilio account SID, API key and API secret",
		ID:          "api-key-secret",
	}}
}

func (tc *TwilioConnector) CreateLogin(ctx context.Context, user *bridgev2.User, flowID string) (bridgev2.LoginProcess, error) {
	if flowID != "api-key-secret" {
		return nil, fmt.Errorf("unknown login flow ID")
	}
	return &TwilioLogin{User: user}, nil
}
```

### `TwilioLogin`

Then we need to define the actual `TwilioLogin` type that `CreateLogin` returns:

```go
type TwilioLogin struct {
	User *bridgev2.User
}
```

Here the interface implementation assertion is quite important. Returning
`TwilioLogin` from `CreateLogin` ensures that the interface implements
`bridgev2.LoginProcess`, but most login flows also need to implement one or more
of the step type specific interfaces. In this case, we're using the user input
type, so we want to make sure `bridgev2.LoginProcessUserInput` is implemented:

```go
var _ bridgev2.LoginProcessUserInput = (*TwilioLogin)(nil)
```

After that, we'll have three methods that need to be implemented: `Start`,
`SubmitUserInput` and `Cancel`.

### `Start`
Start returns the first step of the login process. For other networks that
require a connection, this is probably also where the connection would be
established. For Twilio, we don't have anything to connect to initially, we just
want the user to provide their API keys.

```go
func (tl *TwilioLogin) Start(ctx context.Context) (*bridgev2.LoginStep, error) {
	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeUserInput,
		StepID:       "fi.mau.twilio.enter_api_keys",
		Instructions: "",
		UserInputParams: &bridgev2.LoginUserInputParams{
			Fields: []bridgev2.LoginInputDataField{
				{
					Type:    bridgev2.LoginInputFieldTypeUsername,
					ID:      "account_sid",
					Name:    "Twilio account SID",
					Pattern: `^AC[0-9a-fA-F]{32}$`,
				},
				{
					Type: bridgev2.LoginInputFieldTypeUsername,
					ID:   "api_key",
					Name: "API key",
				},
				{
					Type: bridgev2.LoginInputFieldTypePassword,
					ID:   "api_secret",
					Name: "API secret",
				},
			},
		},
	}, nil
}
```

### `SubmitUserInput` and finishing the login
TODO this section

### `Cancel`
Cancel is called if the user cancels the login process. For networks that create
some sort of connection, you should tear it down here. Since Twilio doesn't have
any such connections, we don't need to do anything.

```go
func (tl *TwilioLogin) Cancel() {}
```

Note that this method is not called at the end of the login, nor if the login
process returns errors. In both of those cases, you need to disconnect yourself.
Errors returned by any step of the process are treated as fatal. If you want to
prompt the user to retry, you should return another login step with the
appropriate instructions. This is also how refreshing QR codes should be done.

## Bridging messages
With everything else out of the way, let's get to the main point: bridging
messages between Twilio and Matrix.

### Matrix → Twilio
To receive messages from Matrix, we need to implement the `HandleMatrixMessage`
function that we skipped over in the network API section.

TODO

### Twilio → Matrix
To receive messages from Twilio, we need to implement the `ReceiveMessages`
function that we created to handle HTTP webhooks.

TODO

## Main function
Now that we have a functional network connector, all that's left is to wrap it
up with a main function. The main function goes in `cmd/mautrix-twilio/main.go`,
because it's a command called mautrix-twilio rather than a part of the library.

The main file doesn't need to be particularly complicated. First, we define some
variables to store the version. These will be set at compile time using the `-X`
linker flag.

```go
var (
	Tag       = "unknown"
	Commit    = "unknown"
	BuildTime = "unknown"
)
```

Then, we make the actual main function, which just creates a `BridgeMain`, gives
it an instance of the connector, and runs the bridge.

```go
func main() {
	m := mxmain.BridgeMain{
		Name:        "mautrix-twilio",
		Description: "A Matrix-Twilio bridge",
		URL:         "https://github.com/mautrix/twilio",
		Version:     "0.1.0",
		Connector:   &connector.TwilioConnector{},
	}
	m.InitVersion(Tag, Commit, BuildTime)
	m.Run()
}
```

That's it. The `mxmain` module is designed to wrap all the parts together to
produce a traditional single-network bridge.

### Building the bridge
To build the bridge, you can simply use `go build ./cmd/mautrix-twilio` in the
repo root directory. However, to be slightly fancier, we also want to fill the
version info variables that we added to main. To do that, we'll make a script
called `build.sh` in the repo root.

```shell
#!/bin/sh
MAUTRIX_VERSION=$(cat go.mod | grep 'maunium.net/go/mautrix ' | awk '{ print $2 }' | head -n1)
GO_LDFLAGS="-s -w -X main.Tag=$(git describe --exact-match --tags 2>/dev/null) -X main.Commit=$(git rev-parse HEAD) -X 'main.BuildTime=`date -Iseconds`' -X 'maunium.net/go/mautrix.GoModVersion=$MAUTRIX_VERSION'"
go build -ldflags="$GO_LDFLAGS" ./cmd/mautrix-twilio "$@"
```

Let's break it down:

The first line gets the version of mautrix-go in use by somewhat crudely parsing
the `go.mod` file. Since there's a lot of code from mautrix-go being used, it's
useful to have the exact commit embedded rather than having to figure it out
based on the bridge version.

The second line defines all the linker flags. `-s` and `-w` are standard flags
to strip debug information and DWARF symbols, respectively. They make the binary
smaller, but also make it harder to debug using debuggers. Generally all you
need in production is stack traces, and fortunately those remain intact.

Each `-X` flag sets the value of a variable in the binary. We set four variables:

* If we're on a Git tag, we want to set `main.Tag` to the tag name. Otherwise,
  it's set to an empty string. To do this, we want both `--exact-match` (don't
  output anything unless we're on a tag) and `--tags` (consider all tags instead
  of only annotated ones). If you use annotated tags for releases, you may want
  to remove `--tags`. The `2>/dev/null` part is needed to suppress the error
  message when we're not on a tag.
* `Commit` is fairly straightforward, it's just the commit hash of the
  current commit (`HEAD`), which is easiest to find using `git rev-parse HEAD`.
* `BuildTime` is the current time in ISO 8601 format.
* Finally, we set `GoModVersion` inside mautrix-go to the version we extracted
  from `go.mod`.

With the linker flags defined, the last step is to actually call `go build`
with those flags and tell it to build our command. The `"$@"` at the end passes
any arguments given to the script to the `go build` command. For example, if you
wanted to output to a different path, you could use `./build.sh -o example.exe`.

### Running the bridge
Finally, we have a bridge, it's compiled, all that's left is to run and use it.
At this point, you can pretty much just follow [the docs](https://docs.mau.fi/bridges/go/setup.html)
starting from the "Configuring and running" part.

1. Generate the example config using `./mautrix-twilio -e`
   (it will be saved to `config.yaml`).
2. Edit the config like any other bridge.
3. Generate the appservice registration with `./mautrix-twilio -g`.
4. Pass the appservice registration to your homeserver.
5. Run the bridge with `./mautrix-twilio`.

After the bridge is running, start a chat with the bridge bot, and send `login`
to start the login process. Then send your API keys as instructed, and you're
good to go!
