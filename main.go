/*
Copyright 2020 - 2021, Jan Bormet, Anna-Felicitas Hausmann, Joachim Schmidt, Vincent Stollenwerk, Arne Turuc

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/timewarrior-synchronize/timew-sync-server/storage"
	"github.com/timewarrior-synchronize/timew-sync-server/sync"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var versionFlag bool
var configFilePath string
var portNumber int
var keyDirectoryPath string
var dbPath string
var noAuth bool
var sourcePath string
var userID int64

func main() {
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	addUserCmd := flag.NewFlagSet("add-user", flag.ExitOnError)
	addKeyCmd := flag.NewFlagSet("add-key", flag.ExitOnError)

	startCmd.StringVar(&configFilePath, "config-file", "", "[RESERVED, not used] Path to the configuration file")
	startCmd.IntVar(&portNumber, "port", 8080, "Port on which the server will listen for connections")
	startCmd.StringVar(&keyDirectoryPath, "keys-location", "authorized_keys", "Path to the users' public keys")
	startCmd.StringVar(&dbPath, "sqlite-db", "db.sqlite", "Path to the SQLite database")
	startCmd.BoolVar(&noAuth, "no-auth", false, "Run server without client authentication")

	addUserCmd.StringVar(&sourcePath, "path", "", "Supply the path to a PEM RSA key")
	addUserCmd.StringVar(&keyDirectoryPath, "keys-location", "authorized_keys", "Path to the users' public keys")

	addKeyCmd.StringVar(&sourcePath, "path", "", "Supply the path to a PEM RSA key")
	addKeyCmd.Int64Var(&userID, "id", -1, "Supply user id")
	addKeyCmd.StringVar(&keyDirectoryPath, "keys-location", "authorized_keys", "Path to the users' public keys")

	flag.BoolVar(&versionFlag, "version", false, "Print version information")

	if len(os.Args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Use commands start, add-user or add-key\n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "start":
		_ = startCmd.Parse(os.Args[2:])
		sync.PublicKeyLocation = keyDirectoryPath
	case "add-user":
		addUserCase(addUserCmd)
	case "add-key":
		addKeyCase(addKeyCmd)
	default:
		flag.Parse()
		if versionFlag {
			_, _ = fmt.Fprintf(os.Stderr, "timewarrior sync server version %v\n", "1.1.0")
			os.Exit(0)
		} else {
			log.Fatal("Use commands start, add-user or add-key")
		}
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error while opening SQLite database: %v", err)
	}
	defer db.Close()
	sqlStorage := &storage.Sql{DB: db}

	err = sqlStorage.Initialize()
	if err != nil {
		log.Fatalf("Error while initializing database: %v", err)
	}
	storage.GlobalStorage = sqlStorage

	handler := func(w http.ResponseWriter, req *http.Request) {
		sync.HandleSyncRequest(w, req, noAuth)
	}
	http.HandleFunc("/api/sync", handler)

	log.Printf("Listening on Port %v", portNumber)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", portNumber), nil))
}

// Subcommand for adding a new user
func addUserCase(addUserCmd *flag.FlagSet) {
	_ = addUserCmd.Parse(os.Args[2:])
	sync.PublicKeyLocation = keyDirectoryPath
	id := sync.GetFreeUserID()
	if sourcePath == "" {
		sync.AddKey(id, "")
	} else {
		key := sync.ReadKey(sourcePath)
		sync.AddKey(id, key)
	}
	_, _ = fmt.Fprintf(os.Stderr, "Successfully added new user %v", id)
	os.Exit(0)
}

// Subcommand for adding a new key
func addKeyCase(addKeyCmd *flag.FlagSet) {
	_ = addKeyCmd.Parse(os.Args[2:])
	sync.PublicKeyLocation = keyDirectoryPath
	if sourcePath == "" {
		log.Fatal("Provide a key file with --path [path-to-key-file]")
	}
	if userID < 0 {
		log.Fatal("Provide a non-negative user id with --id [user id]")
	}
	used := sync.GetUsedUserIDs()
	if !used[userID] {
		log.Fatalf("User %v does not exist", userID)
	}
	key := sync.ReadKey(sourcePath)
	sync.AddKey(userID, key)
	_, _ = fmt.Fprintf(os.Stderr, "Successfully added new key to user %v", userID)
	os.Exit(0)
}
