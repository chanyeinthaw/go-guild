# go-guild

Create bot-owned discord servers.

# Requirements

- discord developer account
- bot access token

# Building

Run `go mod tidy` to build `go-guild`.

# Usage

**Listing bot-owned guilds**

`go-guild --token=<bot access token> -op=ls`

Example output

```
Guilds: 2
Guild ID: 880324740949090315
Guild ID: 880324740949090316
```

**Creating a guild**

`go-guild --token=<bot access token> --op=cm --name=[<server name>]`

Example output

```
Guild   :  <server name>
Guild ID:  880324740949090315
Invite  :  Y8bPD3VWTp
OTP     :  B80704
Bot is now running.  Press CTRL-C to exit.
```

**Managing existing guild**

`go-guild --token=<bot access token> --op=cm --guild=<guild id>`

Example output

```
Guild   :  Server Name
Guild ID:  880324740949090315
Invite  :  Y8bPD3VWTp
OTP     :  B80704
Bot is now running.  Press CTRL-C to exit.
```

**Deleting a guild**

`go-guild --token=<bot access token> --op=del --guild=<guild id>`

Example output
```
Guilds: 2
Guild ID: 880324740949090315
Guild ID: 880324740949090316
```

# Bot Usage

`!help` - display help menu.

`!own <otp>` - get `@owner` role with **Admin** permission.

`!release <otp>` - release `@owner` role.

`!transfer <otp>` - transfer server ownership.
