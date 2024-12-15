package simple_json_db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SimpleDB struct {
	data         map[string]interface{}
	filePath     string
	debugEnabled bool
	timeout      *time.Timer
}

func NewSimpleDB(options map[string]interface{}) *SimpleDB {
	db := &SimpleDB{
		data:         make(map[string]interface{}),
		filePath:     "./db.json",
		debugEnabled: false,
	}

	if file, ok := options["file"].(string); ok {
		db.filePath = file
	}
	if debug, ok := options["debug"].(bool); ok {
		db.debugEnabled = debug
	}

	db.init()
	return db
}

func (db *SimpleDB) init() {
	if db.debugEnabled {
		fmt.Println("[SimpleDB] starting...")
	}
	if !strings.HasSuffix(db.filePath, ".json") {
		db.filePath += ".json"
	}
	if !filepath.IsAbs(db.filePath) {
		wd, _ := os.Getwd()
		db.filePath = filepath.Join(wd, db.filePath)
	}
	if db.debugEnabled {
		fmt.Println("[SimpleDB] File path set")
	}
	fileContent, err := ioutil.ReadFile(db.filePath)
	if err == nil {
		json.Unmarshal(fileContent, &db.data)
	}
	if db.debugEnabled {
		fmt.Println("[SimpleDB] Database ready")
	}
}

func (db *SimpleDB) writeDb() {
	if db.debugEnabled {
		fmt.Println("[SimpleDB] Writing database to file")
	}
	fileContent, _ := json.Marshal(db.data)
	ioutil.WriteFile(db.filePath, fileContent, 0644)
	db.timeoutRemove()
}

func (db *SimpleDB) Set(key string, value interface{}) {
	if db.debugEnabled {
		fmt.Printf("[SimpleDB] Setting data for %s\n", key)
	}
	if db.timeout != nil {
		db.timeoutRemove()
	}
	db.timeoutSet()
	db.data[key] = value
}

func (db *SimpleDB) Get(key string) interface{} {
	if db.debugEnabled {
		fmt.Printf("[SimpleDB] Returning data for %s\n", key)
	}
	return db.data[key]
}

func (db *SimpleDB) Delete(key string) {
	if db.debugEnabled {
		fmt.Printf("[SimpleDB] Deleting data for %s\n", key)
	}
	if db.timeout != nil {
		db.timeoutRemove()
	}
	db.timeoutSet()
	delete(db.data, key)
}

func (db *SimpleDB) Keys() []string {
	if db.debugEnabled {
		fmt.Println("[SimpleDB] returning keys")
	}
	keys := make([]string, 0, len(db.data))
	for k := range db.data {
		keys = append(keys, k)
	}
	return keys
}

func (db *SimpleDB) timeoutSet() {
	if db.debugEnabled {
		fmt.Println("[SimpleDB] Setting timeout")
	}
	db.timeout = time.AfterFunc(500*time.Millisecond, db.writeDb)
}

func (db *SimpleDB) timeoutRemove() {
	if db.debugEnabled {
		fmt.Println("[SimpleDB] Removing timeout")
	}
	if db.timeout != nil {
		db.timeout.Stop()
		db.timeout = nil
	}
}
