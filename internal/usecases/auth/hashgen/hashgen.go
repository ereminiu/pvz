package hashgen

import (
	"crypto/sha1"
	"fmt"
)

const (
	beforeSalt = "a[oCIRFkhnNtDMqKYWcQVlfnk"
	afterSalt  = "CPYGqBDVUaSPUAXoqLktpDVmO"
)

func GenerateHash(password string) string {
	hash := sha1.New()

	hash.Write([]byte(beforeSalt))
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(afterSalt)))
}
