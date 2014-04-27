package vedis

/*
 #include <stdlib.h>

 #include "vedis.h"

 extern int
 kv_fetch_callback_func(void *data, unsigned int dataLen, void *userData);

 // you cannot export xConsumer directory because a function with a const pointer argument cannot be exported.
 typedef int (*xConsumer)(const void *pData, unsigned int iDataLen, void *pUserData);

 static inline int
 _vedis_kv_fetch_callback(vedis *pDb, const void *pKey,int nKeyLen, void *pUserData) {
   return vedis_kv_fetch_callback(pDb, pKey, nKeyLen, (xConsumer) kv_fetch_callback_func, pUserData);
 }
*/
import "C"

import (
	"unsafe"
)

type Vedis struct {
	db *C.vedis
}

// Open a database and return vedis object.
// If fileName is ":mem:", then a private, in-memory database is created for the connection.
// See: http://vedis.symisc.net/c_api/vedis_open.html
func OpenVedis(fileName string) (*Vedis, error) {
	cname := C.CString(fileName)
	defer C.free(unsafe.Pointer(cname))

	var db *C.vedis
	if rc := C.vedis_open(&db, cname); rc != C.VEDIS_OK {
		return nil, ErrCode(rc)
	}

	return &Vedis{db}, nil
}

// Close the database.
// See: http://vedis.symisc.net/c_api/vedis_close.html
func (u *Vedis) Close() error {
	if rc := C.vedis_close(u.db); rc != C.VEDIS_OK {
		return ErrCode(rc)
	}

	return nil
}

/*
 * Key-Value Store Interface
 */

// KvStore write a new record into the database. If the record does not exists, it is created. Otherwise, it is replaced.
// See: http://vedis.symisc.net/c_api/vedis_kv_store.html
func (u *Vedis) KvStore(key []byte, value []byte) error {
	if rc := C.vedis_kv_store(u.db, unsafe.Pointer(&key[0]), C.int(len(key)), unsafe.Pointer(&value[0]), C.vedis_int64(len(value))); rc != C.VEDIS_OK {
		return ErrCode(rc)
	}

	return nil
}

// KvStore write a new record into the database. If the record does not exists, it is created. Otherwise, the new data chunk is appended to the end of the old chunk.
// See: http://vedis.symisc.net/c_api/vedis_kv_append.html
func (u *Vedis) KvAppend(key []byte, value []byte) error {
	if rc := C.vedis_kv_append(u.db, unsafe.Pointer(&key[0]), C.int(len(key)), unsafe.Pointer(&value[0]), C.vedis_int64(len(value))); rc != C.VEDIS_OK {
		return ErrCode(rc)
	}

	return nil
}

// KvFetch a record from the database
// See: http://vedis.symisc.net/c_api/vedis_kv_fetch.html
func (u *Vedis) KvFetch(key []byte) ([]byte, error) {
	var n C.vedis_int64

	if rc := C.vedis_kv_fetch(u.db, unsafe.Pointer(&key[0]), C.int(len(key)), nil, &n); rc != C.VEDIS_OK {
		return nil, ErrCode(rc)
	}

	buf := make([]byte, int64(n))
	if rc := C.vedis_kv_fetch(u.db, unsafe.Pointer(&key[0]), C.int(len(key)), unsafe.Pointer(&buf[0]), &n); rc != C.VEDIS_OK {
		return nil, ErrCode(rc)
	}

	return buf, nil
}

type kvFetchCallbackFunc func([]byte) ErrCode

//export kv_fetch_callback_func
func kv_fetch_callback_func(data unsafe.Pointer, dataLen C.uint, userData unsafe.Pointer) C.int {
	f := *(*kvFetchCallbackFunc)(userData)
	d := C.GoBytes(data, C.int(int(uint(dataLen))))
	return C.int(f(d))
}

// TODO: temporary variables might be destroyed by GC?
func (u *Vedis) KvFetchCallback(key []byte, f kvFetchCallbackFunc) error {
	if rc := C._vedis_kv_fetch_callback(u.db, unsafe.Pointer(&key[0]), C.int(len(key)), unsafe.Pointer(&f)); rc != C.VEDIS_OK {
		return ErrCode(rc)
	}

	return nil
}

// KvDelete remove a particular record from the database, you can use this high-level thread-safe routine to perform the deletion.
// See: http://vedis.symisc.net/c_api/vedis_kv_delete.html
func (u *Vedis) KvDelete(key []byte) error {
	if rc := C.vedis_kv_delete(u.db, unsafe.Pointer(&key[0]), C.int(len(key))); rc != C.VEDIS_OK {
		return ErrCode(rc)
	}

	return nil
}
