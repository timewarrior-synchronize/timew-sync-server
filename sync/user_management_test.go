/*
Copyright 2026 - timew-sync-server contributors

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
	"os"
	"path/filepath"
	"testing"
)

func TestGetUsedUserIDs_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	oldLocation := PublicKeyLocation
	PublicKeyLocation = tmpDir
	defer func() { PublicKeyLocation = oldLocation }()

	ids := GetUsedUserIDs()
	if len(ids) != 0 {
		t.Errorf("expected empty map, got %v", ids)
	}
}

func TestGetUsedUserIDs_WithFiles(t *testing.T) {
	tmpDir := t.TempDir()
	oldLocation := PublicKeyLocation
	PublicKeyLocation = tmpDir
	defer func() { PublicKeyLocation = oldLocation }()

	for _, name := range []string{"0_keys", "2_keys", "5_keys"} {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte("key"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "notakey"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	ids := GetUsedUserIDs()
	if len(ids) != 3 {
		t.Errorf("expected 3 ids, got %d", len(ids))
	}
	for _, id := range []int64{0, 2, 5} {
		if !ids[id] {
			t.Errorf("expected id %d to be present", id)
		}
	}
}

func TestGetFreeUserID_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	oldLocation := PublicKeyLocation
	PublicKeyLocation = tmpDir
	defer func() { PublicKeyLocation = oldLocation }()

	id := GetFreeUserID()
	if id != 0 {
		t.Errorf("expected 0, got %d", id)
	}
}

func TestGetFreeUserID_WithGap(t *testing.T) {
	tmpDir := t.TempDir()
	oldLocation := PublicKeyLocation
	PublicKeyLocation = tmpDir
	defer func() { PublicKeyLocation = oldLocation }()

	for _, name := range []string{"0_keys", "1_keys", "3_keys"} {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte("key"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	id := GetFreeUserID()
	if id != 2 {
		t.Errorf("expected 2, got %d", id)
	}
}

func TestGetFreeUserID_StartingFromZero(t *testing.T) {
	tmpDir := t.TempDir()
	oldLocation := PublicKeyLocation
	PublicKeyLocation = tmpDir
	defer func() { PublicKeyLocation = oldLocation }()

	for _, name := range []string{"1_keys", "2_keys"} {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte("key"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	id := GetFreeUserID()
	if id != 0 {
		t.Errorf("expected 0, got %d", id)
	}
}
