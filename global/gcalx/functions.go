package gcalx

import (
	"crypto/md5"
	"encoding/hex"
	"unsafe"
)

// md5
func Md5(s string) string {
	m := md5.Sum([]byte(s))
	return hex.EncodeToString(m[:])
}

func IsNil(obj interface{}) bool {
	type eFace struct {
		data unsafe.Pointer
	}
	if obj == nil {
		return true
	}
	return (*eFace)(unsafe.Pointer(&obj)).data == nil
}
