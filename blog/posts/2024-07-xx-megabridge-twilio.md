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
In this post, I'll go over all the steps necessary to build a Twilio bridge
using the new bridgev2 module in mautrix-go. The whole bridge can be found at
[github.com/mautrix/twilio](https://github.com/mautrix/twilio).

## Getting started with a new Go project
The first step to make a new Go project is to create a new directory and run
`go mod init <module path>`, where `<module path>` is the import path of your
new module (for example, the GitHub repo). In addition to that, we'll want to
add the mautrix-go and Twilio libraries as dependencies.

```shell
go mod init go.mau.fi/mautrix-twilio
go get maunium.net/go/mautrix
go get github.com/twilio/twilio-go
```

When naming your bridge, please make up your own name and don't use `mautrix-*`.

## The connector itself
The next step is creating the network connector itself. The connector is
effectively a library, so we'll put it in `pkg/connector/` and create a file
called `connector.go`. Because this is a minimal example, that file will also be
the only file in the connector package, but real connectors will probably want
to split up the parts that come later into different files.

Inside the file, let's start by defining a struct called `TwilioConnector`. This
struct is the main entrypoint to the connector. It is passed to the central
bridge module and is used to initialize other things.

```go
package connector

import (
    "maunium.net/go/mautrix/bridgev2"
)

type TwilioConnector struct {
	br *bridgev2.Bridge
}
```

Then, add a line like this:

```go
var _ bridgev2.NetworkConnector = (*TwilioConnector)(nil)
```

This is a conventional method of ensuring that a struct implements a specific
interface. The line is creating a value of type `*TwilioConnector`, then
assigning it to a variable of type `bridgev2.NetworkConnector`, which will only
compile if the value implements the interface. The variable name is the blank
identifier `_`, which means the value will be discarded after being evaluated.

If you're using a smart editor, it should complain that `*TwilioConnector` does
not in fact implement `bridgev2.NetworkConnector`, and possibly even offer you
a quick way to create stub methods to implement the interface.

### `Init`
The Init function is called when the bridge is initializing all types. It also
gives you access to the bridge struct, which will need to be stored for later.

This function should not do any kind of IO or other complicated operations, it
should just initialize the in-memory struct.

```go
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
networks that receive events via websockets/polling/etc. may not need to do
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

We'll come back to `ReceiveMessage` [later](#twilio--matrix).

### `GetCapabilities`
The `GetCapabilities` function on the network connector is used to signal some
bridge-wide capabilities, like disappearing message support. Twilio doesn't
have any of the relevant features, so we'll just leave this empty.

```go
func (tc *TwilioConnector) GetCapabilities() *bridgev2.NetworkGeneralCapabilities {
	return &bridgev2.NetworkGeneralCapabilities{}
}
```

### `GetName`
The `GetName` function is used to customize the name of the bridge.

* `DisplayName` is a simple human-readable name for the network. It doesn't have
  any particular rules. It usually starts with a capital letter. This is used in
  lots of places.
* `NetworkURL` is the website associated with the network.
  This is used in the `protocol` section of `m.bridge` events.
* `NetworkIcon` is a `mxc://` URI which contains the logo of the network.
  This is used in the `protocol` section of `m.bridge` events, as well as in the
  avatar of the bridge bot user.
* `NetworkID` is a string that uniquely identifies the network. If there are
  multiple bridge implementations for the same network, they should use the same
  ID. This is conventionally all lowercase.
* `BeeperBridgeType` identifies the specific bridge implementation. The Go
  module import path is a good option for this to ensure uniqueness, but bridges
  used by Beeper use shorter types (e.g. the Go rewrite of the Discord bridge
  used `discordgo`).
* `DefaultPort` can optionally be set to change the default port when generating
  the example config. It is not required and will default to `8008` when unset.
  All mautrix bridges use ports defined in [mau.fi/ports](https://mau.fi/ports).
* `DefaultCommandPrefix` can optionally be set to change the default command
  prefix when generating the example config. It is not required and will default
  to `!` followed by the `NetworkID`.

```go
func (tc *TwilioConnector) GetName() bridgev2.BridgeName {
	return bridgev2.BridgeName{
		DisplayName:      "Twilio",
		NetworkURL:       "https://twilio.com",
		NetworkIcon:      "mxc://maunium.net/FYuKJHaCrSeSpvBJfHwgYylP",
		NetworkID:        "twilio",
		BeeperBridgeType: "go.mau.fi/mautrix-twilio",
		DefaultPort:      29322,
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
values.

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

We'll initialize the go-twilio client here. We're also initializing a
`RequestValidator`. It's used for the webhooks, which we'll come back to later.

```go
func (tc *TwilioConnector) LoadUserLogin(ctx context.Context, login *bridgev2.UserLogin) error {
	accountSID := login.Metadata.Extra["account_sid"].(string)
	authToken := login.Metadata.Extra["auth_token"].(string)
	restClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username:   accountSID,
		Password:   authToken,
		AccountSid: accountSID,
	})
	validator := client.NewRequestValidator(authToken)
	login.Client = &TwilioClient{
		UserLogin:        login,
		Twilio:           restClient,
		RequestValidator: validator,
	}
	return nil
}
```

## The network API
Next, we'll need to actually define `TwilioClient` and implement the `NetworkAPI`.

```go
type TwilioClient struct {
	UserLogin        *bridgev2.UserLogin
	Twilio           *twilio.RestClient
	RequestValidator client.RequestValidator
}

var _ bridgev2.NetworkAPI = (*TwilioClient)(nil)
```

Like with the network connector, we'll do the same interface implementation
assertion, In this case it's not technically necessary, as we're already
assigning a `&TwilioClient{}` to `UserLogin.Client`, which does the same
check. However, if you implement any of the optional extra interfaces, then the
explicit assertions become very useful, as there's nothing else ensuring the
correct functions are implemented.

### `Connect`
For most networks which use persistent connections, this is where you'd set up
the connection. Twilio doesn't use a persistent connection, so technically we
don't need to do anything here. However, we should still check access token
validity here.

```go
func (tc *TwilioClient) Connect(ctx context.Context) error {
	phoneNumbers, err := tc.Twilio.Api.ListIncomingPhoneNumber(nil)
	if err != nil {
		return fmt.Errorf("failed to list phone numbers: %w", err)
	}
	numberInUse := tc.UserLogin.Metadata.Extra["phone"].(string)
	var numberFound bool
	for _, number := range phoneNumbers {
		if number.PhoneNumber != nil && *number.PhoneNumber == numberInUse {
			numberFound = true
			break
		}
	}
	if !numberFound {
		return fmt.Errorf("phone number %s not found on account", numberInUse)
	}
	return nil
}
```

### `Disconnect`
For networks with persistent connections, Disconnect should tear down the
connection. Twilio doesn't have a persistent connection, so we don't need to do
anything here.

```go
func (tc *TwilioClient) Disconnect() {}
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

### `GetCapabilities`
This is similar to the network connector's GetCapabilities method, but is scoped
to a user login and a portal. Currently, these fields are only used to check
events before passing them to the network connector. Some of the fields are not
used at all yet. The plan is to also send these fields to the room as a state
event, so that clients could display limits directly to the user. The state
event will likely use [MSC4110] (or at least something similar).

[MSC4110]: https://github.com/matrix-org/matrix-spec-proposals/pull/4110

For now, we don't really need to define any fields, but let's include Twilio's
maximum message length.

```go
func (tc *TwilioClient) GetCapabilities(ctx context.Context, portal *bridgev2.Portal) *bridgev2.NetworkRoomCapabilities {
	return &bridgev2.NetworkRoomCapabilities{
		MaxTextLength: 1600,
	}
}
```

## Identifiers
Before we get to the next functions in `NetworkAPI`, we'll need to cover network
identifiers. Network IDs are opaque identifiers for various things on the remote
network: logins, users, chats, messages, etc. Each identifier has its own type
in the networkid module, which ensures that you can't accidentally mix up types.
All the types are just strings behind the scenes.

All identifiers are generated by the network connector and won't be parsed by
any other component. Other components also will not make any assumptions about
different identifier types being similar, but the network connector itself is of
course allowed to define some types are equal. For example, most networks (but
not all) will define that `UserLoginID`s are the same as `UserID`s. However,
identifiers do have some uniqueness expectations that the network connector must
meet.

For Twilio, we'll define the identifiers as follows:

* `UserID`s are E.164 phone numbers without the leading `+`.
* `PortalID`s are equivalent to `UserID`s.
* `MessageID`s are Twilio message SIDs.
* `UserLoginID`s are account SID and phone SID joined with a `:`.

For convenience, we'll define some functions to cast strings into those types:

```go
func makeUserID(e164Phone string) networkid.UserID {
	return networkid.UserID(strings.TrimLeft(e164Phone, "+"))
}

func makePortalID(e164Phone string) networkid.PortalID {
	return networkid.PortalID(strings.TrimLeft(e164Phone, "+"))
}

func makeUserLoginID(accountSID, phoneSID string) networkid.UserLoginID {
	return networkid.UserLoginID(fmt.Sprintf("%s:%s", accountSID, phoneSID))
}
```

See the [networkid module godocs](https://pkg.go.dev/maunium.net/go/mautrix/bridgev2/networkid)
for docs on all the different types of identifiers.

### `IsThisUser`
Since `UserID`s and `UserLoginID`s are not interchangeable, we need to provide
some way for the bridge to determine if a given user ID belongs to a user login.
For most networks where user and login IDs are the same, you can just check for
equality.

For this bridge, we're segregating different logins to have their own portals,
which means this function is not actually necessary, and we could just hardcode
it to return `false`. It's not hard to implement though, so let's do it anyway.

```go
func (tc *TwilioClient) IsThisUser(ctx context.Context, userID networkid.UserID) bool {
	phoneNum := tc.UserLogin.Metadata.Extra["phone"].(string)
	return makeUserID(phoneNum) == userID
}
```

If you were to define UserLoginID and UserID the same way, you could have an
even simpler check:

```go
func (tc *NotTwilioClient) IsThisUser(ctx context.Context, userID networkid.UserID) bool {
	return networkid.UserID(tc.UserLogin.ID) == userID
}
```

### `GetChatInfo`
`GetChatInfo` returns the info for a given chat. All the values in the response
struct are pointers, which means they can be omitted to tell the bridge that the
corresponding room state event shouldn't be modified. For example, DMs generally
don't have names, topics or avatars. However, even DMs do have members. The
member list should always include all participants, so both the Matrix user and
the remote user in DMs.

We're only handling DMs in this bridge, so we don't need to return anything
other than the member list.

```go
func (tc *TwilioClient) GetChatInfo(ctx context.Context, portal *bridgev2.Portal) (*bridgev2.ChatInfo, error) {
	return &bridgev2.ChatInfo{
		Members: &bridgev2.ChatMemberList{
			IsFull: true,
			Members: []bridgev2.ChatMember{
				{
					EventSender: bridgev2.EventSender{
						IsFromMe: true,
						Sender:   makeUserID(tc.UserLogin.Metadata.Extra["phone"].(string)),
					},
					// This could be omitted, but leave it in to be explicit.
					Membership: event.MembershipJoin,
					// Make the user moderator, so they can adjust the room metadata if they want to.
					PowerLevel: 50,
				},
				{
					EventSender: bridgev2.EventSender{
						Sender: networkid.UserID(portal.ID),
					},
					Membership: event.MembershipJoin,
					PowerLevel: 50,
				},
			},
		},
	}, nil
}
```

### `GetUserInfo`
`GetUserInfo` is basically the same as `GetChatInfo`, but it returns info for a
given user instead of a chat. The returned info will be applied as the ghost
user's profile on Matrix.

Because we're bridging SMS and don't have a contact lits, we don't really have
any other info than the phone number itself. If we wanted to be fancy, we could
format the phone number nicely for `Name`, but I couldn't find any convenient
libraries similar to phonenumbers for Python.

In addition to `Name`, we set `Identifiers` which is a list of URIs that
represent the user. In this case, we're using the `tel:` scheme, but you could
also include network-specific @usernames with a custom scheme here.

```go
func (tc *TwilioClient) GetUserInfo(ctx context.Context, ghost *bridgev2.Ghost) (*bridgev2.UserInfo, error) {
	return &bridgev2.UserInfo{
		Identifiers: []string{fmt.Sprintf("tel:+%s", ghost.ID)},
		Name:        ptr.Ptr(fmt.Sprintf("+%s", ghost.ID)),
	}, nil
}
```

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
connector.

```go
func (tc *TwilioConnector) GetLoginFlows() []bridgev2.LoginFlow {
	return []bridgev2.LoginFlow{{
		Name:        "Auth token",
		Description: "Log in with your Twilio account SID and auth token",
		ID:          "auth-token",
	}}
}

func (tc *TwilioConnector) CreateLogin(ctx context.Context, user *bridgev2.User, flowID string) (bridgev2.LoginProcess, error) {
	if flowID != "auth-token" {
		return nil, fmt.Errorf("unknown login flow ID")
	}
	return &TwilioLogin{User: user}, nil
}
```

### `TwilioLogin`
Then we need to define the actual `TwilioLogin` type that `CreateLogin` returns:

```go
type TwilioLogin struct {
	User         *bridgev2.User
	Client       *twilio.RestClient
	PhoneNumbers []twilioPhoneNumber
	AccountSID   string
	AuthToken    string
}
```

We have a bunch of extra fields in addition to the `User`. They are used to
store data when there are multiple login steps. Specifically, if the Twilio
account has more than one phone number, we'll return a second step asking which
one to use.

Here the interface implementation assertion is quite important. Returning
`TwilioLogin` from `CreateLogin` ensures that the interface implements
`bridgev2.LoginProcess`, but most login flows also need to implement one or more
of the step type specific interfaces. In this case, we're using the user input
type, so we want to make sure `bridgev2.LoginProcessUserInput` is implemented.

```go
var _ bridgev2.LoginProcessUserInput = (*TwilioLogin)(nil)
```

After that, we'll have three methods that need to be implemented: `Start`,
`SubmitUserInput` and `Cancel`.

### `Start`
Start returns the first step of the login process. For other networks that
require a connection, this is probably also where the connection would be
established. For Twilio, we don't have anything to connect to initially, we just
want the user to provide their account SID and auth token.

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
					Type:    bridgev2.LoginInputFieldTypePassword,
					ID:      "auth_token",
					Name:    "Twilio auth token",
					Pattern: "^[0-9a-f]{32}$",
				},
			},
		},
	}, nil
}
```

### `SubmitUserInput` and finishing the login
After the user provides the values, we'll get a call to `SubmitUserInput`.

This will be a more complicated function. First, we need to validate the
credentials and get a list of phone numbers available on the Twilio account.
After that, we either finish the login if there's only one number, or ask the
user which one to use if there are multiple. If we ask the user, then we'll get
another call to `SubmitUserInput`, which means we need to remember the data from
the first call. After a successful login, we prepare the `UserLogin` instance.

Let's split up the function to keep it more readable. First, `SubmitUserInput`
itself. We have two paths, so we'll just split it into two calls. If `Client` is
not set in `TwilioLogin`, we're in the first step where we want API keys. If it
is set, we want to choose a phone number.

```go
func (tl *TwilioLogin) SubmitUserInput(ctx context.Context, input map[string]string) (*bridgev2.LoginStep, error) {
	if tl.Client == nil {
		return tl.submitAPIKeys(ctx, input)
	} else {
		return tl.submitChosenPhoneNumber(ctx, input)
	}
}
```

Then the API key submit function.

```go
type twilioPhoneNumber struct {
	SID          string
	Number       string
	PrettyNumber string
}

func (tl *TwilioLogin) submitAPIKeys(ctx context.Context, input map[string]string) (*bridgev2.LoginStep, error) {
	tl.AccountSID = input["account_sid"]
	tl.AuthToken = input["auth_token"]
	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username:   tl.AccountSID,
		Password:   tl.AuthToken,
		AccountSid: tl.AccountSID,
	})
	// Get the list of phone numbers. This doubles as a way to verify the credentials are valid.
	phoneNumbers, err := twilioClient.Api.ListIncomingPhoneNumber(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list phone numbers: %w", err)
	}
	var numbers []twilioPhoneNumber
	for _, number := range phoneNumbers {
		if number.Status == nil || number.PhoneNumber == nil || *number.Status != "in-use" {
			continue
		}
		numbers = append(numbers, twilioPhoneNumber{
			SID:          *number.Sid,
			Number:       *number.PhoneNumber,
			PrettyNumber: *number.FriendlyName,
		})
	}
	tl.Client = twilioClient
	tl.PhoneNumbers = numbers
	if len(numbers) == 0 {
		return nil, fmt.Errorf("no active phone numbers found")
	} else if len(numbers) == 1 {
		return tl.finishLogin(ctx, numbers[0])
	} else {
		phoneNumberList := make([]string, len(numbers))
		for i, number := range numbers {
			phoneNumberList[i] = fmt.Sprintf("* %s", number.Number)
		}
		return &bridgev2.LoginStep{
			Type:         bridgev2.LoginStepTypeUserInput,
			StepID:       "fi.mau.twilio.choose_number",
			Instructions: "Your Twilio account has multiple phone numbers. Please choose one:\n\n" + strings.Join(phoneNumberList, "\n"),
			UserInputParams: &bridgev2.LoginUserInputParams{
				Fields: []bridgev2.LoginInputDataField{{
					Type: bridgev2.LoginInputFieldTypePhoneNumber,
					ID:   "chosen_number",
					Name: "Phone number",
				}},
			},
		}, nil
	}
}
```

Choosing the phone number is fairly simple, as we already have a valid token and
have fetched the list of phone numbers. We just need to find the phone number
the user chose.

```go
func (tl *TwilioLogin) submitChosenPhoneNumber(ctx context.Context, input map[string]string) (*bridgev2.LoginStep, error) {
	numberIdx := slices.IndexFunc(tl.PhoneNumbers, func(e twilioPhoneNumber) bool {
		return e.Number == input["chosen_number"]
	})
	if numberIdx == -1 {
		// We could also return a new LoginStep here if we wanted to allow the user to retry.
		// Errors are always fatal, so returning an error here will cancel the login process.
		return nil, fmt.Errorf("invalid phone number")
	}
	return tl.finishLogin(ctx, tl.PhoneNumbers[numberIdx])
}
```

Finally, the finish function, which can be called from either path and creates
the `UserLogin` object. In addition to creating the object, we also send our
webhook URL to Twilio. We'll define the `GetWebhookURL` function later when
implementing `ReceiveMessages`.

```go
func (tl *TwilioLogin) finishLogin(ctx context.Context, phoneNumber twilioPhoneNumber) (*bridgev2.LoginStep, error) {
	ul, err := tl.User.NewLogin(ctx, &database.UserLogin{
		ID: makeUserLoginID(tl.AccountSID, phoneNumber.SID),
		Metadata: database.UserLoginMetadata{
			StandardUserLoginMetadata: database.StandardUserLoginMetadata{
				RemoteName: phoneNumber.PrettyNumber,
			},
			Extra: map[string]any{
				"phone":       phoneNumber.Number,
				"phone_sid":   phoneNumber.SID,
				"auth_token":  tl.AuthToken,
				"account_sid": tl.AccountSID,
			},
		},
	}, &bridgev2.NewLoginParams{
		LoadUserLogin: func(ctx context.Context, login *bridgev2.UserLogin) error {
			login.Client = &TwilioClient{
				UserLogin:        login,
				Twilio:           tl.Client,
				RequestValidator: tclient.NewRequestValidator(tl.AuthToken),
			}
			return nil
		},
	})
	if err != nil {
		return nil, err
	}
	tc := ul.Client.(*TwilioClient)
	// In addition to creating the UserLogin, we'll also want to set the webhook URL for the phone number.
	_, err = tc.Twilio.Api.UpdateIncomingPhoneNumber(phoneNumber.SID, &openapi.UpdateIncomingPhoneNumberParams{
		SmsMethod: ptr.Ptr(http.MethodPost),
		SmsUrl:    ptr.Ptr(tc.GetWebhookURL()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to set webhook URL for phone number: %w", err)
	}
	// Finally, return the special complete step indicating the login was successful.
	// It doesn't have any params other than the UserLogin we just created.
	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeComplete,
		StepID:       "fi.mau.twilio.complete",
		Instructions: "Successfully logged in",
		CompleteParams: &bridgev2.LoginCompleteParams{
			UserLoginID: ul.ID,
			UserLogin:   ul,
		},
	}, nil
}
```

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

### Twilio → Matrix
To receive messages from Twilio, we need to implement the `ReceiveMessages`
function that we created to handle HTTP webhooks. Before that, let's define the
webhook URLs, which was used in the login section.

We don't need to check if `Matrix` implements `MatrixConnectorWithServer`,
because we already validated that in `Start`. We can just cast it to access
`GetPublicAddress` and then append our path. We include the user login ID in the
in order to correctly route incoming webhooks.

```go
func (tc *TwilioClient) GetWebhookURL() string {
	server := tc.UserLogin.Bridge.Matrix.(bridgev2.MatrixConnectorWithServer)
	return fmt.Sprintf("%s/_twilio/%s/receive", server.GetPublicAddress(), tc.UserLogin.ID)
}
```

The `ReceiveMessages` function contains a lot of boilerplate code that the
Twilio library could handle, but doesn't. The main bridge-specific code is
finding the user login based on the path parameter.

```go
func (tc *TwilioConnector) ReceiveMessage(w http.ResponseWriter, r *http.Request) {
	// First make sure the signature header is present and that the request body is valid form data.
	sig := r.Header.Get("X-Twilio-Signature")
	if sig == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Missing signature header\n"))
		return
	}

	params := make(map[string]string)
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Failed to parse form data\n"))
		return
	}
	for key, value := range r.PostForm {
		params[key] = value[0]
	}

	// Get the user login based on the path. We need it to find the right token
	// to use for validating the request signature.
	loginID := mux.Vars(r)["loginID"]
	login := tc.br.GetCachedUserLoginByID(networkid.UserLoginID(loginID))
	if login == nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Unrecognized login ID in request path\n"))
		return
	}
	client := login.Client.(*TwilioClient)

	// Now that we have the client, validate the request.
	if !client.RequestValidator.Validate(client.GetWebhookURL(), params, sig) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("Invalid signature\n"))
		return
	}

	// Pass the request to the client for handling. This is where everything actually happens.
	client.HandleWebhook(r.Context(), params)

	// We don't want to respond immediately, so just send a blank TwiML response.
	twimlResult, err := twiml.Messages(nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "text/xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(twimlResult))
	}
}
```

Finally, we need the actual handling function. All we really need to do is pass
the event to the central bridge. To do that, we need to extract metadata like
the portal and message IDs, and provide a converter function to actually convert
the message into a Matrix event.

Simple connectors can use the `bridgev2.SimpleRemoteEvent` struct as the remote
event, but for more complicated connectors, it often makes sense to create an
interface which implements `bridgev2.RemoteEvent`. That way, the interface
methods can figure out the appropriate data to return, instead of having to fill
a struct for every different event type.

```go
func (tc *TwilioClient) HandleWebhook(ctx context.Context, params map[string]string) {
	tc.UserLogin.Bridge.QueueRemoteEvent(tc.UserLogin, &bridgev2.SimpleRemoteEvent[map[string]string]{
		Type: bridgev2.RemoteEventMessage,
		LogContext: func(c zerolog.Context) zerolog.Context {
			return c.
				Str("from", params["From"]).
				Str("message_id", params["MessageSid"])
		},
		PortalKey: networkid.PortalKey{
			ID:       makePortalID(params["From"]),
			Receiver: tc.UserLogin.ID,
		},
		Data:         params,
		CreatePortal: true,
		ID:           networkid.MessageID(params["MessageSid"]),
		Sender: bridgev2.EventSender{
			Sender: makeUserID(params["From"]),
		},
		Timestamp:          time.Now(),
		ConvertMessageFunc: tc.convertMessage,
	})
}
```

Let's go over each of the fields we're filling:

* `Type` is the event type. It's a normal message.
* `LogContext` is a function that adds structured fields to the event handler's
  zerolog logger. By default, the logger only has the portal key and user login
  ID, so other things should be added here.
* `PortalKey` is the ID of the chat. This is a combination of a portal ID and an
  optional "receiver". Receivers can be used to segregate portals, so that if
  multiple logged-in users have the same chat, they'll still get separate portal
  rooms. Most networks should use receivers for DMs, but it is also possible to
  use them for all rooms if you don't want any portals to be shared. If there's
  no receiver, then users will be added to the same Matrix room.
* `Data` is the event data itself. This is only here so that it can be passed to
  the message convert function.
* `CreatePortal` tells the central bridge module that we want it to create a
  portal room if one doesn't already exist for the given portal key. The bridge
  will then call `GetChatInfo` to get the info of the chat to create.
* `ID` is the message ID.
* `Sender` is the sender of the message. For networks where the user can send
  messages from other clients, you should also fill `IsFromMe` and/or
  `SenderLogin` appropriately. For Twilio, we'll just assume you can't send
  messages from other clients (we don't support receiving those anyway), so we
  don't need to fill anything else than `Sender`.
* `Timestamp` is the message timestamp. Twilio doesn't seem to provide
  timestamps, so we just declare that the message was sent now.
* `ConvertMessageFunc` is a function that gets `Data`, the `Portal` object as
  well as a `MatrixAPI` and returns the Matrix events that should be sent.

The convert message function is very simple, as we only support plain text
messages for now. If you wanted to bridge media, you'd download it from the
remote network and reupload it to Matrix using `intent.UploadMedia`.

```go
func (tc *TwilioClient) convertMessage(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, data map[string]string) (*bridgev2.ConvertedMessage, error) {
	return &bridgev2.ConvertedMessage{
		Parts: []*bridgev2.ConvertedMessagePart{{
			Type: event.EventMessage,
			Content: &event.MessageEventContent{
				MsgType: event.MsgText,
				Body:    data["Body"],
			},
		}},
	}, nil
}
```

### Matrix → Twilio
To receive messages from Matrix, we need to implement the `HandleMatrixMessage`
function that we skipped over in the network API section. Responding is very
simple, we just call the Twilio API and return the message ID.

```go
func (tc *TwilioClient) HandleMatrixMessage(ctx context.Context, msg *bridgev2.MatrixMessage) (message *bridgev2.MatrixMessageResponse, err error) {
	resp, err := tc.Twilio.Api.CreateMessage(&openapi.CreateMessageParams{
		To:   ptr.Ptr(fmt.Sprintf("+%s", msg.Portal.ID)),
		From: ptr.Ptr(tc.UserLogin.Metadata.Extra["phone"].(string)),
		Body: ptr.Ptr(msg.Content.Body),
	})
	if err != nil {
		return nil, err
	}
	return &bridgev2.MatrixMessageResponse{
		DB: &database.Message{
			ID:       networkid.MessageID(*resp.Sid),
			SenderID: makeUserID(*resp.From),
		},
	}, nil
}
```

## Bonus feature: starting chats
We've implemented everything that's strictly necessary for a bridge to work, but
let's add one optional feature on top: creating new portal rooms. To do this,
we'll add another interface assertion for `TwilioClient`:

```go
var _ bridgev2.IdentifierResolvingNetworkAPI = (*TwilioClient)(nil)
```

The interface requires us to implement the `ResolveIdentifier` method, which is
used for both checking if an identifier is reachable and actually starting a
direct chat. There are further optional interfaces for creating group chats with
resolved identifiers, but we don't support group chats at all here, so let's
stick to DMs.

The function just gets a raw string which is provided by the user. If you wanted
to make a fancy Twilio bridge, you'd probably use the lookup API to get more
info about the phone number, but we'll just do basic validation to make sure the
input is a number.

After validating the number, we'll get the ghost and portal objects as well as
their info. The info for both is the same shape as the `GetUserInfo` and
`GetChatInfo` methods, so we'll just call them instead of duplicating the same
behavior.

We don't actually care about the `createChat` parameter here, because Twilio
doesn't require creating chats explicitly. For networks which do require
creating chats, you'd need to use the bool to decide whether it should be
created or not. The `Chat` field in the response is mandatory when `createChat`
is true, but can be omitted when it's false.

We also don't create the portal room here: the central bridge module takes care
of that using the info we return.

```go
func (tc *TwilioClient) ResolveIdentifier(ctx context.Context, identifier string, createChat bool) (*bridgev2.ResolveIdentifierResponse, error) {
	e164Number, err := bridgev2.CleanPhoneNumber(identifier)
	if err != nil {
		return nil, err
	}
	userID := makeUserID(e164Number)
	portalID := networkid.PortalKey{
		ID:       makePortalID(e164Number),
		Receiver: tc.UserLogin.ID,
	}
	ghost, err := tc.UserLogin.Bridge.GetGhostByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ghost: %w", err)
	}
	portal, err := tc.UserLogin.Bridge.GetPortalByID(ctx, portalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get portal: %w", err)
	}
	ghostInfo, _ := tc.GetUserInfo(ctx, ghost)
	portalInfo, _ := tc.GetChatInfo(ctx, portal)
	return &bridgev2.ResolveIdentifierResponse{
		Ghost:    ghost,
		UserID:   userID,
		UserInfo: ghostInfo,
		Chat: &bridgev2.CreateChatResponse{
			Portal:     portal,
			PortalID:   portalID,
			PortalInfo: portalInfo,
		},
	}, nil
}
```

That's everything needed from the connector to enable starting chats. With that
function implemented, the `resolve-identifier` and `start-chat` bot commands
as well as the corresponding provisioning APIs will work.

## Main function
Now that we have a functional network connector, all that's left is to wrap it
up with a main function. The main function goes in `cmd/mautrix-twilio/main.go`,
because it's a command called mautrix-twilio rather than a part of the library.

The main file doesn't need to be particularly complicated. First, we define some
variables to store the version. These will be set at compile time using the `-X`
linker flag. We'll go over the exact flags in the next section.

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
* `Commit` is fairly straightforward, it's just the commit hash of the current
  commit (`HEAD`), which is easiest to find using `git rev-parse HEAD`.
* `BuildTime` is the current time in ISO 8601/RFC3339 format.
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

## Conclusion
If you want to ask anything related to mautrix-go or this post, feel free to
join [#go:maunium.net]. [This post] also accepts pull requests.

[#go:maunium.net]: https://matrix.to/#/#go:maunium.net
[this post]: https://github.com/maunium/mau.fi/blob/main/blog/posts/2024-07-xx-megabridge-twilio.md
