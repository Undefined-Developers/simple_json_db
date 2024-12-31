# SimpleDB

SimpleDB is a lightweight JSON-based database for Go. It allows you to store, retrieve, and manage data with ease.

## Installation

To install SimpleDB, use `go get`:

```sh
go get github.com/Undefined-Developers/simple_json_db
```

## Usage

Here is an example of how to use SimpleDB:

```go
package main

import (
	"fmt"
	"github.com/Undefined-Developers/simple_json_db"
)

func main() {
	// Create a new SimpleDB instance with options
	options := map[string]interface{}{
		"file":  "mydb.json",
		"debug": true,
		"delay": 2000, // Set delay to 2000 milliseconds (2 seconds)
	}
	db := simple_json_db.NewSimpleDB(options)

	// Set values
	db.Set("name", "John Doe")
	db.Set("age", 30)

	// Get values
	name := db.Get("name")
	age := db.Get("age")
	fmt.Printf("Name: %s, Age: %d\n", name, age)

	// Get all keys
	keys := db.Keys()
	fmt.Println("Keys:", keys)

	// Check if a key exists
	exists := db.Has("name")
	fmt.Printf("Does 'name' exist? %v\n", exists)

	// Delete a value
	db.Delete("age")

	// Check deletion
	age = db.Get("age")
	fmt.Printf("Age after deletion: %v\n", age)
}
```

## API

### `NewSimpleDB(options map[string]interface{}) *SimpleDB`

Creates a new SimpleDB instance. Options can include:
- `file`: The file path for the database (default: `./db.json`)
- `debug`: Enable debug mode (default: `false`)

### `Set(key string, value interface{})`

Sets a value in the database.

### `Get(key string) interface{}`

Gets a value from the database.

### `Delete(key string)`

Deletes a value from the database.

### `Keys() []string`

Returns all keys in the database.

### `Has(key string) bool`

Checks if a key exists in the database. Returns `true` if the key exists, otherwise `false`.

## License

This project is licensed under the GNU GENERAL PUBLIC LICENSE Version 3. See the `LICENSE` file for details.

