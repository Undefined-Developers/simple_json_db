package simple_json_db

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SimpleDB struct {
	data         map[string]string
	filePath     string
	debugEnabled bool
	timeout      *time.Timer
	delay        time.Duration
}

func NewSimpleDB(options map[string]interface{}) *SimpleDB {
	db := &SimpleDB{
		data:         make(map[string]string),
		filePath:     "./db.json",
		debugEnabled: false,
		delay:        5 * time.Second,
	}

	if file, ok := options["file"].(string); ok {
		db.filePath = file
	}
	if debug, ok := options["debug"].(bool); ok {
		db.debugEnabled = debug
	}
	if delay, ok := options["delay"].(int); ok {
		db.delay = time.Duration(delay) * time.Millisecond
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
	if err := os.MkdirAll(filepath.Dir(db.filePath), os.ModePerm); err != nil {
		if db.debugEnabled {
			fmt.Printf("[SimpleDB] Error creating directories: %v\n", err)
		}
	}
	if db.debugEnabled {
		fmt.Println("[SimpleDB] File path set")
		fmt.Printf("[SimpleDB] File path is: %s\n", db.filePath)
	}

	if _, err := os.Stat(db.filePath); os.IsNotExist(err) {
		if db.debugEnabled {
			fmt.Printf("[SimpleDB] File does not exist, creating new file: %s\n", db.filePath)
		}
		file, err := os.Create(db.filePath)
		if err != nil && db.debugEnabled {
			fmt.Printf("[SimpleDB] Error creating file: %v\n", err)
			return
		}
		err = file.Close()
		if err != nil {
			fmt.Printf("[SimpleDB] Error closing file after creation: %v\n", err)
			return
		}
	}

	fileContent, err := os.ReadFile(db.filePath)
	if err == nil {
		var tempData map[string]interface{}
		err := json.Unmarshal(fileContent, &tempData)
		if err != nil {
			if db.debugEnabled {
				fmt.Printf("[SimpleDB] Error unmarshalling file: %v\n", err)
			}
			return
		}

		for k, v := range tempData {
			db.data[k] = fmt.Sprintf("%v", v)
		}
		if db.debugEnabled {
			fmt.Println("[SimpleDB] Database ready")
		}
	} else {
		if db.debugEnabled {
			fmt.Printf("[SimpleDB] Error reading file: %v\n", err)
		}
	}
}

func (db *SimpleDB) writeDb() {
	if db.debugEnabled {
		fmt.Println("[SimpleDB] Writing database to file")
	}
	fileContent, _ := json.Marshal(db.data)
	err := os.WriteFile(db.filePath, fileContent, 0644)
	if err != nil {
		if db.debugEnabled {
			fmt.Printf("[SimpleDB] Error writing file: %v\n", err)
		}
		return
	}
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
	db.data[key] = fmt.Sprintf("%v", value)
}

func (db *SimpleDB) Get(key string) *string {
	if db.debugEnabled {
		fmt.Printf("[SimpleDB] Returning data for %s\n", key)
	}
	if value, exists := db.data[key]; exists {
		return &value
	}
	return nil
}

func (db *SimpleDB) Delete(key string) {
	if db.debugEnabled {
		fmt.Printf("[SimpleDB] Deleting data for %s\n", key)
	}
	if db.Has(key) {
		if db.timeout != nil {
			db.timeoutRemove()
		}
		db.timeoutSet()
		delete(db.data, key)
	}
}

func (db *SimpleDB) Has(key string) bool {
	if db.debugEnabled {
		fmt.Printf("[SimpleDB] Checking existence of %s\n", key)
	}
	_, exists := db.data[key]
	return exists
}

func (db *SimpleDB) Keys() []string {
	if db.debugEnabled {
		fmt.Println("[SimpleDB] Returning keys")
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
	db.timeout = time.AfterFunc(db.delay, db.writeDb)
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
