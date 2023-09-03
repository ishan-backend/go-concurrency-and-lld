package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// singletonDB naming ensures, we are not using it directly, but using some kind of factory
type singletonDB struct {
	dBName string
	dbPass string
}

// Once ensures that any construct called within is called once and only once.
var once sync.Once

// sync.Once init() - thread safety

var instanceDB *singletonDB

// GetDBInstance is called for lazy initialization
func GetDBInstance() *singletonDB {
	once.Do(func() {
		readConfig, err := readData(".\\config.txt")
		if err != nil {
			panic(err)
		}
		db := singletonDB{dBName: readConfig["dbName"], dbPass: readConfig["dbPass"]}
		instanceDB = &db
	})

	return instanceDB
}

func main() {
	db := GetDBInstance()
	db.HealthCheck()
}

func readData(path string) (map[string]string, error) {
	exe, err := os.Executable() // returns path of current executable file
	if err != nil {
		panic(err)
	}

	exePath := filepath.Dir(exe) // beautifies path of exe, and returns its directory
	file, err := os.Open(exePath + path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file) // io.Reader expects a File to read
	scanner.Split(bufio.ScanLines)

	result := map[string]string{}
	for scanner.Scan() {
		k := scanner.Text()
		scanner.Scan()
		v := scanner.Bytes()
		result[k] = string(v)
	}

	return result, nil
}

func (i *singletonDB) HealthCheck() {
	fmt.Println("DB Instance is up")
}
