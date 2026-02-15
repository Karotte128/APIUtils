# APIUtils

**APIUtils** is a collection of utility packages for enhancing and supporting usage of the **[Karotte128/KarotteAPI](https://github.com/karotte128/karotteapi)** Go API framework. It contains features that don’t belong in the core framework itself but are useful across API projects using KarotteAPI.

This repository does *not* contain a standalone server; it provides helpers for configuration, database integration, permission providers, and common utilities for API applications.

---

## Features

Currently included (folder names):

- `config` — helpers to load configuration from file and expand environment variables.
- `database` — utilities for postgresql database integration.
- `db_perm` — database-driven permission provider.
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

### db_perm

This is a simple permission provider for `simpleauth` using a database.

`GetPermission`, `SetPermission` and `UpdatePermission` are used to interact with the permission data in the database.

`GetPermissionWrapper` is the permission provider to use with simpleauth.

Usage:
```go
simpleauth.Setup(db_perm.GetPermissionWrapper) // use db_perm as permission  provider for simpleauth
```

---

### simpleauth

The `simpleauth` package provides a simple, API-Key based authentication system.

To use it, a `PermissionProvider` is needed. This is a `type PermissionProvider func(string) []string` that takes a `string` containing the API key and returns a `[]string` of permissions.

APIUtils includes a default database permission provider, `db_perm`.

To use simpleauth, set it up before calling `api.InitAPI()`.

Example:
```go
database.CreateConnection(dbconn) // First, create the database connection.

simpleauth.Setup(db_perm.GetPermissionWrapper) // Set up simpleauth with the permission provider. In this example the included db_perm provider is used.

api.InitAPI(details) // Start the KarotteAPI API server.
```

To use the authentication check inside of a module, include the `simpleauth.HasPermission(context.Context, string)` check in the module request handler.

```go
func Handler(w http.ResponseWriter, r *http.Request) {
	hasPermission := simpleauth.HasPermission(r.Context(), "exampleperm") // Check for permission using the request context

    if !hasPermission {
        http.Error(w, "Unauthorized", 401) // Client does not have the permission
		return
    }

    // Client does have the permission, handle the request normally
}
```

---

## Example Setup

This is a fully functional example setup, using the config system, postgresql database as permission provider and simpleauth for API-Key based authentication.

```go
//main.go
package main

import (
	"log"

	"github.com/karotte128/apiutils/config"
	"github.com/karotte128/apiutils/database"
	"github.com/karotte128/apiutils/db_perm"
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

	simpleauth.Setup(db_perm.GetPermissionWrapper) // Set up simpleauth using the db_perm permission provider (apiutils/simpleauth and apiutils/db_perm)

	var details karotteapi.ApiDetails // Create the API details (karotteapi)
	details.Config = conf // Set the config
	api.InitAPI(details) // Start the KarotteAPI server (karotteapi/api)
}
```