package vedis

/*
 #include "vedis.h"
*/
import "C"

type ErrCode int

const (
	OK             ErrCode = C.VEDIS_OK             /* Successful result */
	NOMEM          ErrCode = C.VEDIS_NOMEM          /* Out of memory */
	ABORT          ErrCode = C.VEDIS_ABORT          /* Another thread have released this instance */
	IOERR          ErrCode = C.VEDIS_IOERR          /* IO error */
	CORRUPT        ErrCode = C.VEDIS_CORRUPT        /* Corrupt pointer */
	LOCKED         ErrCode = C.VEDIS_LOCKED         /* Forbidden Operation */
	BUSY           ErrCode = C.VEDIS_BUSY           /* The database file is locked */
	DONE           ErrCode = C.VEDIS_DONE           /* Operation done */
	PERM           ErrCode = C.VEDIS_PERM           /* Permission error */
	NOTIMPLEMENTED ErrCode = C.VEDIS_NOTIMPLEMENTED /* Method not implemented by the underlying Key/Value storage engine */
	NOTFOUND       ErrCode = C.VEDIS_NOTFOUND       /* No such record */
	NOOP           ErrCode = C.VEDIS_NOOP           /* No such method */
	INVALID        ErrCode = C.VEDIS_INVALID        /* Invalid parameter */
	EOF            ErrCode = C.VEDIS_EOF            /* End Of Input */
	UNKNOWN        ErrCode = C.VEDIS_UNKNOWN        /* Unknown configuration option */
	LIMIT          ErrCode = C.VEDIS_LIMIT          /* Database limit reached */
	EXISTS         ErrCode = C.VEDIS_EXISTS         /* Record exists */
	EMPTY          ErrCode = C.VEDIS_EMPTY          /* Empty record */
	FULL           ErrCode = C.VEDIS_FULL           /* Full database (unlikely) */
	CANTOPEN       ErrCode = C.VEDIS_CANTOPEN       /* Unable to open the database file */
	READ_ONLY      ErrCode = C.VEDIS_READ_ONLY      /* Read only Key/Value storage engine */
	LOCKERR        ErrCode = C.VEDIS_LOCKERR        /* Locking protocol error */
)

var errorString = map[ErrCode]string{
	OK:             "Successful result",
	NOMEM:          "Out of memory",
	ABORT:          "Another thread have released this instance",
	IOERR:          "IO error",
	CORRUPT:        "Corrupt pointer",
	LOCKED:         "Forbidden Operation",
	BUSY:           "The database file is locked",
	DONE:           "Operation done",
	PERM:           "Permission error",
	NOTIMPLEMENTED: "Method not implemented by the underlying Key/Value storage engine",
	NOTFOUND:       "No such record",
	NOOP:           "No such method",
	INVALID:        "Invalid parameter",
	EOF:            "End Of Input",
	UNKNOWN:        "Unknown configuration option",
	LIMIT:          "Database limit reached",
	EXISTS:         "Record exists",
	EMPTY:          "Empty record",
	FULL:           "Full database (unlikely)",
	CANTOPEN:       "Unable to open the database file",
	READ_ONLY:      "Read only Key/Value storage engine",
	LOCKERR:        "Locking protocol error",
}

func (err ErrCode) Error() string {
	return errorString[err]
}
