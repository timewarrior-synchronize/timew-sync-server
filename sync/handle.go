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

package sync

import (
	"git.rwth-aachen.de/computer-aided-synthetic-biology/bachelorpraktika/2020-67-timewarrior-sync/timew-sync-server/data"
	"git.rwth-aachen.de/computer-aided-synthetic-biology/bachelorpraktika/2020-67-timewarrior-sync/timew-sync-server/storage"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// HandleSyncRequest receives sync requests and starts the sync
// process with the received data.
func HandleSyncRequest(w http.ResponseWriter, req *http.Request) {
	requestBody, reqError := ioutil.ReadAll(req.Body)
	if reqError != nil {
		log.Printf("Error reading sync request. Ignoring request.")
		return
	}
	requestData, parseError := data.ParseSyncRequest(string(requestBody))
	if parseError != nil {
		log.Printf("Error parsing sync request. Ignoring request.")
		return
	}
	syncData, _, err := Sync(requestData, storage.GlobalStorage)
	if err != nil {
		log.Printf("syncing failed")
	}
	responseBody, respError := data.ToJSON(syncData)
	if respError != nil {
		log.Printf("Error creating response JSON. Ignoring request.")
		return
	}
	sendResponse(w, responseBody)
}

// sendResponse writes data to response buffer
func sendResponse(w http.ResponseWriter, data string) {
	_, err := io.WriteString(w, data)
	if err != nil {
		log.Printf("sync/handle.go:sendResponse Error writing response to ResponseWriter")
	}
}
