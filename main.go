/*
Copyright 2020 - Jan Bormet, Anna-Felicitas Hausmann, Joachim Schmidt, Vincent Stollenwerk, Arne Turuc

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
	"git.rwth-aachen.de/computer-aided-synthetic-biology/bachelorpraktika/2020-67-timewarrior-sync/timew-sync-server/storage"
	"git.rwth-aachen.de/computer-aided-synthetic-biology/bachelorpraktika/2020-67-timewarrior-sync/timew-sync-server/sync"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var versionFlag bool
var configFilePath string
var portNumber int

func main() {
	flag.BoolVar(&versionFlag, "version", false, "Print version information")
	flag.StringVar(&configFilePath, "config-file", "", "Path to the configuration file")
	flag.IntVar(&portNumber, "port", 8080, "Port on which the server will listen for connections")
	flag.Parse()

	if versionFlag {
		_, _ = fmt.Fprintf(os.Stderr, "timewarrior sync server version %v\n", "unreleased")
		os.Exit(0)
	}

	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlStorage := &storage.Sql{DB: db}
	sqlStorage.Initialize()
	storage.GlobalStorage = sqlStorage

	http.HandleFunc("/api/sync", sync.HandleSyncRequest)

	log.Printf("Listening on Port %v", portNumber)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", portNumber), nil))
}
