/*
Copyright 2021 - Jan Bormet, Anna-Felicitas Hausmann, Joachim Schmidt, Vincent Stollenwerk, Arne Turuc

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
	"git.rwth-aachen.de/computer-aided-synthetic-biology/bachelorpraktika/2020-67-timewarrior-sync/timew-sync-server/data"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

// Authenticate returns true iff the JWT specified in the http requests' Bearer token was signed by the correct user.
// If any step of the authentication process goes wrong or there is no matching public key Authenticate returns false
func Authenticate(r *http.Request, body data.SyncRequest) bool {
	keySet, err := GetKeySet(body.UserID)
	if err != nil {
		log.Printf("Error during Authentication. Unable to obtain keys for user %v", body.UserID)
		return false
	}

	for i := 0; i < keySet.Len(); i++ {
		key, ok := keySet.Get(i)
		if !ok {
			continue
		}

		token, err := jwt.ParseHeader(r.Header, "Authorization", jwt.WithValidate(true),
			jwt.WithVerify(jwa.RS256, key), jwt.WithAcceptableSkew(time.Duration(10e10)))
		if err != nil {
			continue
		}

		id, ok := token.Get("userID")
		if !ok {
			continue
		}

		presumedUserID, ok := id.(float64)
		if !ok || int(presumedUserID) != body.UserID {
			continue
		}
		return true
	}
	return false
}

// GetKeySet returns the key set of user with userId. Returns an error if the keys file of that user was not found
// or could not be parsed.
func GetKeySet(userId int) (jwk.Set, error) {
	filename := fmt.Sprintf("%d_keys", userId)
	path := filepath.Join(PublicKeyLocation, filename)

	keySet, err := jwk.ReadFile(path, jwk.WithPEM(true))
	if err != nil {
		log.Printf("Error parsing key set of user %d: %v", userId, err)
		return nil, err
	}

	return keySet, nil
}
