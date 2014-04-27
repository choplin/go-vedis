package vedis

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TempFilename() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), "vedis"+hex.EncodeToString(randBytes)+".db")
}

func TestOpenVedis(t *testing.T) {
	tempFilename := TempFilename()
	db, err := OpenVedis(tempFilename)
	if err != nil {
		t.Fatal("Failed to open database: ", err)
	}
	defer os.Remove(tempFilename)
	defer db.Close()

	db.KvStore([]byte("k"), []byte("v"))

	if stat, err := os.Stat(tempFilename); err != nil || stat.IsDir() {
		t.Error("Failed to create a database file", err)
	}
}

func TestKvStore(t *testing.T) {
	tempFilename := TempFilename()
	db, err := OpenVedis(tempFilename)
	if err != nil {
		t.Fatal("Failed to open database: ", err)
	}
	defer os.Remove(tempFilename)
	defer db.Close()

	k := []byte("k")
	v1 := []byte("v1")
	v2 := []byte("v2")

	if err := db.KvStore(k, v1); err != nil {
		t.Fatal("Failed to store: ", err)
	}

	r1, err := db.KvFetch(k)
	if err != nil {
		t.Fatal("Failed to fetch: ", err)
	}

	if !bytes.Equal(v1, r1) {
		t.Errorf("Fetched result is different from stored value. stored:%s fetched:%s", v1, r1)
	}

	if err := db.KvStore(k, v2); err != nil {
		t.Fatal("Failed to store: ", err)
	}

	r2, err := db.KvFetch(k)
	if err != nil {
		t.Fatal("Failed to fetch: ", err)
	}
	if !bytes.Equal(v2, r2) {
		t.Errorf("Fetched result is not overwritten with new value. stored:%s fetched:%s", v2, r2)
	}
}

func TestKvAppend(t *testing.T) {
	tempFilename := TempFilename()
	db, err := OpenVedis(tempFilename)
	if err != nil {
		t.Fatal("Failed to open database: ", err)
	}
	defer os.Remove(tempFilename)
	defer db.Close()

	k := []byte("k")
	v1 := []byte("v1")
	v2 := []byte("v2")
	appended := append(v1, v2...)

	if err := db.KvStore(k, v1); err != nil {
		t.Fatal("Failed to store: ", err)
	}

	if err := db.KvAppend(k, v2); err != nil {
		t.Fatal("Failed to append: ", err)
	}

	r, err := db.KvFetch(k)
	if err != nil {
		t.Fatal("Failed to fetch: ", err)
	}

	if !bytes.Equal(appended, r) {
		t.Errorf("Fetched result is different from appended value. stored1:%s stored2:%s fetched:%s", v1, v2, r)
	}
}

func TestKvDelete(t *testing.T) {
	tempFilename := TempFilename()
	db, err := OpenVedis(tempFilename)
	if err != nil {
		t.Fatal("Failed to open database: ", err)
	}
	defer os.Remove(tempFilename)
	defer db.Close()

	k := []byte("k")
	v := []byte("v")

	if err := db.KvStore(k, v); err != nil {
		t.Fatal("Failed to store: ", err)
	}

	if err := db.KvDelete(k); err != nil {
		t.Fatal("Failed to delete: ", err)
	}

	r, err := db.KvFetch(k)
	if r != nil {
		t.Error("The record has not been deleted")
	}

	if err.(ErrCode) != NOTFOUND {
		t.Error("Failed to delete a record")
	}
}

func TestKvStoreCallback(t *testing.T) {
	tempFilename := TempFilename()
	db, err := OpenVedis(tempFilename)
	if err != nil {
		t.Fatal("Failed to open database: ", err)
	}
	defer os.Remove(tempFilename)
	defer db.Close()

	k := []byte("k")
	v := []byte("value")

	if err := db.KvStore(k, v); err != nil {
		t.Fatal("Failed to store: ", err)
	}

	f := func(data []byte) ErrCode {
		if !bytes.Equal(v, data) {
			t.Error("the Stored value and the value passed for callback function is different")
		}
		return OK
	}

	if err = db.KvFetchCallback(k, f); err != nil {
		t.Error("An error occured when KvFetchCallback is called: ", err)
	}
}
