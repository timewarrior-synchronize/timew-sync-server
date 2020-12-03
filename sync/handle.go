package sync

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// HandleSyncRequest receives sync requests and starts the sync
// process with the received data.
func HandleSyncRequest(w http.ResponseWriter, req *http.Request) {
	responseData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading sync request. Ignoring request.")
	}

	fmt.Println(string(responseData))
}
