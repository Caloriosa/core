package tools

import (
	"encoding/hex"
	"github.com/golang/glog"
	"github.com/magical/argon2"
)

const encodeStd = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const PASSWORD_SALT_LENGTH = 32

func EncodeUInt64(n uint64) string {
	var (
		b   = make([]byte, 0)
		rem uint64
	)

	// Progressively divide by base, store remainder each time
	// Prepend as an additional character is the higher power
	for n > 0 {
		rem = n % 62
		n = n / 62
		b = append([]byte{encodeStd[rem]}, b...)
	}

	s := string(b)

	return s
}

func EncodeUserPassword(password, salt string) string {
	s, err := argon2.Key([]byte(password), []byte(salt), 5, 2, 16, 32)
	if err != nil {
		glog.Fatal("Error creating a hash: ", err)
	}

	return hex.EncodeToString(s)
}
