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

package sync

import (
	_ "github.com/lestrrat-go/jwx"
	"github.com/timewarrior-synchronize/timew-sync-server/data"
	"github.com/timewarrior-synchronize/timew-sync-server/storage"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var PublicKeyLocation string

// HandleSyncRequest receives sync requests and starts the sync
// process with the received data.
func HandleSyncRequest(w http.ResponseWriter, req *http.Request, noAuth bool) {
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading HTTP request, ignoring request: %v", err)
		return
	}

	requestData, err := data.ParseSyncRequest(string(requestBody))
	if err != nil {
		log.Printf("Error parsing sync request, ignoring request: %v", err)
		errorResponse := ErrorResponseBody{
			Message: "An error occurred while parsing the request",
			Details: err.Error(),
		}
		sendResponse(w, http.StatusBadRequest, errorResponse.ToString())
		return
	}

	// Authentication
	if !noAuth {
		authenticated := Authenticate(req, requestData)
		if !authenticated {
			errorResponse := ErrorResponseBody{
				Message: "An error occurred during authentication",
				Details: "",
			}
			sendResponse(w, http.StatusUnauthorized, errorResponse.ToString())
			return
		}
	}

	syncData, conflict, err := Sync(requestData, storage.GlobalStorage)
	if err != nil {
		log.Printf("Synchronization failed, ignoring request: %v", err)
		errorResponse := ErrorResponseBody{
			Message: "An error occurred while performing the synchronization",
			Details: err.Error(),
		}
		sendResponse(w, http.StatusInternalServerError, errorResponse.ToString())
		return
	}

	responseBody, err := data.ToJSON(syncData, conflict)
	if err != nil {
		log.Printf("Error creating response JSON, ignoring request: %v", err)
		errorResponse := ErrorResponseBody{
			Message: "An error occurred while creating the response",
			Details: err.Error(),
		}
		sendResponse(w, http.StatusInternalServerError, errorResponse.ToString())
		return
	}

	sendResponse(w, http.StatusOK, responseBody)
}

// sendResponse writes data to response buffer
func sendResponse(w http.ResponseWriter, statusCode int, data string) {
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	_, err := io.WriteString(w, data)
	if err != nil {
		log.Printf("Error writing response to ResponseWriter")
	}
}
