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
	"crypto/rand"
	"crypto/rsa"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"net/http"
	"testing"
	"time"
)

func TestAuthenticateWithKeySet_positive(t *testing.T) {
	raw1, err1 := rsa.GenerateKey(rand.Reader, 1024)
	raw2, err2 := rsa.GenerateKey(rand.Reader, 1024)
	key1, err3 := jwk.New(raw1)
	key2, err4 := jwk.New(raw2)
	pub1, err5 := jwk.PublicKeyOf(key1)
	pub2, err6 := jwk.PublicKeyOf(key2)
	keySet := jwk.NewSet()
	keySet.Add(pub1)
	keySet.Add(pub2)
	token := jwt.New()
	token.Set("userID", 42)
	token.Set(jwt.ExpirationKey, time.Now().Add(time.Hour))
	payload, err7 := jwt.Sign(token, jwa.RS256, key2)
	bearer := "Bearer " + string(payload)
	req, err8 := http.NewRequest("POST", "", nil)
	req.Header.Add("Authorization", bearer)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil ||
		err8 != nil {
		t.Errorf("Failed to generate key set in preparation for testing")
	}
	b := AuthenticateWithKeySet(req, 42, keySet)
	if !b {
		t.Errorf("Failed to authenticate")
	}

}

func TestAuthenticateWithKeySet_negative(t *testing.T) {
	raw1, err1 := rsa.GenerateKey(rand.Reader, 1024)
	raw2, err2 := rsa.GenerateKey(rand.Reader, 1024)
	key1, err3 := jwk.New(raw1)
	key2, err4 := jwk.New(raw2)
	pub1, err5 := jwk.PublicKeyOf(key1)
	keySet := jwk.NewSet()
	keySet.Add(pub1)
	token := jwt.New()
	token.Set("userID", 42)
	token.Set(jwt.ExpirationKey, time.Now().Add(time.Hour))
	payload, err7 := jwt.Sign(token, jwa.RS256, key2)
	bearer := "Bearer " + string(payload)
	req, err8 := http.NewRequest("POST", "", nil)
	req.Header.Add("Authorization", bearer)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err7 != nil ||
		err8 != nil {
		t.Errorf("Failed to generate key set in preparation for testing")
	}
	b := AuthenticateWithKeySet(req, 42, keySet)
	if b {
		t.Errorf("Authenticated falsely")
	}
}

func TestAuthenticateWithKeySet_expired(t *testing.T) {
	raw1, err1 := rsa.GenerateKey(rand.Reader, 1024)
	raw2, err2 := rsa.GenerateKey(rand.Reader, 1024)
	key1, err3 := jwk.New(raw1)
	key2, err4 := jwk.New(raw2)
	pub1, err5 := jwk.PublicKeyOf(key1)
	pub2, err6 := jwk.PublicKeyOf(key2)
	keySet := jwk.NewSet()
	keySet.Add(pub1)
	keySet.Add(pub2)
	token := jwt.New()
	token.Set("userID", 42)
	token.Set(jwt.ExpirationKey, time.Now().Add(-time.Hour))
	payload, err7 := jwt.Sign(token, jwa.RS256, key2)
	bearer := "Bearer " + string(payload)
	req, err8 := http.NewRequest("POST", "", nil)
	req.Header.Add("Authorization", bearer)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil ||
		err8 != nil {
		t.Errorf("Failed to generate key set in preparation for testing")
	}
	b := AuthenticateWithKeySet(req, 42, keySet)
	if b {
		t.Errorf("Authenticated with expired jwt")
	}
}

func TestAuthenticateWithKeySet_IDMismatch(t *testing.T) {
	raw1, err1 := rsa.GenerateKey(rand.Reader, 1024)
	raw2, err2 := rsa.GenerateKey(rand.Reader, 1024)
	key1, err3 := jwk.New(raw1)
	key2, err4 := jwk.New(raw2)
	pub1, err5 := jwk.PublicKeyOf(key1)
	pub2, err6 := jwk.PublicKeyOf(key2)
	keySet := jwk.NewSet()
	keySet.Add(pub1)
	keySet.Add(pub2)
	token := jwt.New()
	token.Set("userID", 42)
	token.Set(jwt.ExpirationKey, time.Now().Add(time.Hour))
	payload, err7 := jwt.Sign(token, jwa.RS256, key2)
	bearer := "Bearer " + string(payload)
	req, err8 := http.NewRequest("POST", "", nil)
	req.Header.Add("Authorization", bearer)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil ||
		err8 != nil {
		t.Errorf("Failed to generate key set in preparation for testing")
	}
	b := AuthenticateWithKeySet(req, 0, keySet)
	if b {
		t.Errorf("Authenticated with mismatching userIDs")
	}

}
