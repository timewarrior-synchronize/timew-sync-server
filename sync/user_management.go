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
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// GetUsedUserIDs returns a map containing every user id with an existing file [user id]_keys
// in PublicKeyLocation directory
func GetUsedUserIDs() map[int64]bool {
	files, err := ioutil.ReadDir(PublicKeyLocation)
	if err != nil {
		log.Fatal("Error accessing keys-location directory")
	}
	used := make(map[int64]bool)

	for _, f := range files {
		s := strings.Split(f.Name(), "_")
		if len(s) != 2 {
			continue
		}
		i, err := strconv.ParseInt(s[0], 10, 64)
		if err != nil {
			continue
		}
		if s[1] == "keys" {
			used[i] = true
		}
	}
	return used
}

// GetFreeUserID returns the smallest valid unused user id
func GetFreeUserID() int64 {
	used := GetUsedUserIDs()
	for i := int64(0); i <= math.MaxInt64; i++ {
		if !used[i] {
			return i
		}
	}
	log.Fatal("Error obtaining free user id")
	return -1
}

// ReadKey reads the key from a file
func ReadKey(path string) string {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading key file at %v", path)
	}
	return string(key)
}

// AddKey adds the given key to the key file of the given user
func AddKey(userID int64, key string) {
	if userID < 0 {
		log.Fatal("Error adding key. Negative user id not allowed")
	}

	destFileName := fmt.Sprintf("%d_keys", userID)
	destFile, err := os.OpenFile(filepath.Join(PublicKeyLocation, destFileName), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error adding key. Unable to create new key file or write to existing key file with user id %v", userID)
	}
	defer destFile.Close()
	if key == "" {
		return
	}
	stat, err := destFile.Stat()
	if err != nil {
		log.Fatal("Unable to obtain kye file length")
	}
	if stat.Size() > 0 {
		key = "\n" + key
	}
	if _, err = destFile.WriteString(key); err != nil {
		destFile.Close()
		log.Fatalf("Error adding key. Unable to write to key file with user id %v", userID)
	}
}
