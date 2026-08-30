# APIUtils

**APIUtils** is a collection of utility packages for enhancing and supporting usage of the **[Karotte128/KarotteAPI](https://github.com/karotte128/karotteapi)** Go API framework. It contains features that don’t belong in the core framework itself but are useful across API projects using KarotteAPI.

This repository does *not* contain a standalone server; it provides helpers for configuration, database integration, authentication providers, and common utilities for API applications.

---

## Features

Currently included (folder names):

- `config` — helpers to load configuration from file and expand environment variables.
- `database` — utilities for postgresql database integration.
- `dbauth` — database-driven authentication provider.
- `simpleauth` — basic authentication provider.
- Other utility packages as needed by API applications.

> See the code under each folder for exact function documentation and usage examples.

---

## Usage

### Config

`apiutils/config` contains a `.toml`-file config loader and environment variable expander.
It can be used to create the `config` for [Karotte128/KarotteAPI](https://github.com/karotte128/karotteapi).

Usage:
```go
err, rawConf := config.ReadConfigFromFile("config.toml") // Load the .toml file
	if err != nil {
		log.Fatal("failed loading config: " + err.Error())
	}

	conf := config.ExpandEnvConfig(rawConf) // Replace ${VAR} or ${VAR:-default} in the config with environment variables
```

---

### Database

#### `CreateConnection`

This can be used to create a new database connection from a connection string.

Usage:
```go
database.CreateConnection("postgres://example:example@localhost:5432/exampledb") // Alternatively, load the connection string from config.
```

#### `InsertStruct`

This can be used to insert a struct into a postgresql table.

#### `SelectStruct` and `SelectStructs`

These can be used to select structs from a postgresql table.

#### `UpdateStruct`

This can be used to update a struct in a postgresql table.

---

### dbauth

This is a simple authentication provider for `simpleauth` using a database.

`GetAuth`, `SetAuth` and `UpdateAuth` are used to interact with the authentication data in the database.

`GetAuthProvider` is the authentication provider to use with simpleauth.

Usage:
```go
simpleauth.Setup(dbauth.GetAuthProvider) // use dbauth as auth provider for simpleauth
```

---

### simpleauth

The `simpleauth` package provides a simple, API-Key based authentication system.

To use it, a `AuthProvider` is needed. This is a struct that contains two functions that must be implemented:
- `ReadAuthInfo  func(string) (AuthInfo, error)`
- `WriteAuthInfo func(AuthInfo) error`

APIUtils includes a default database auth provider, `dbauth`.

To use simpleauth, set it up before calling `api.InitAPI()`.

Example:
```go
database.CreateConnection(dbconn) // First, create the database connection.

simpleauth.Setup(dbauth.GetAuthProvider) // Set up simpleauth with the authentication provider. In this example the included dbauth provider is used.

api.InitAPI(config) // Start the KarotteAPI API server.
```

To use the authentication check inside of a module, include the `simpleauth.HasPermission(AuthInfo, string)` check in the module request handler.

```go
func Handler(w http.ResponseWriter, r *http.Request) {
	authInfo, ok := simpleauth.GetAuthInfo(r)
	if !ok {
		http.Error(w, "Unauthorized", 401)
		return
	}

	hasPermission := simpleauth.HasPermission(authInfo, "exampleperm") // Check for permission using the request context

    if !hasPermission {
        http.Error(w, "Unauthorized", 401) // Client does not have the permission
		return
    }

    // Client does have the permission, handle the request normally
}
```

---

## Example Setup

This is a fully functional example setup, using the config system, postgresql database as authentication provider and simpleauth for API-Key based authentication.

```go
//main.go
package main

import (
	"log"

	"github.com/karotte128/apiutils/config"
	"github.com/karotte128/apiutils/database"
	"github.com/karotte128/apiutils/dbauth"
	"github.com/karotte128/apiutils/simpleauth"
	"github.com/karotte128/karotteapi"
	"github.com/karotte128/karotteapi/api"
	"github.com/karotte128/karotteapi/core"
)

func main() {
	err, rawConf := config.ReadConfigFromFile("config.toml") // Load the .toml config file (apiutils/config)
	if err != nil {
		log.Fatal("failed loading config: " + err.Error())
	}

	conf := config.ExpandEnvConfig(rawConf) // Replace ENV variables in the config (apiutils/config)

	dbconn, ok := core.GetNestedValue[string](conf, "database", "connection") // Get the database connection string from the config (karotteapi/core)
	if !ok {
		log.Fatal("no database config!")
	}

	database.CreateConnection(dbconn) // Create the database connection (apiutils/database)

	simpleauth.Setup(dbauth.GetAuthProvider) // Set up simpleauth using the dbauth authenticaton provider (apiutils/simpleauth and apiutils/dbauth)

	api.InitAPI(conf) // Start the KarotteAPI server (karotteapi/api)
}
```