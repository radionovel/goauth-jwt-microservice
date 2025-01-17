package password

import (
	"crypto/md5"
	"encoding/hex"
)

func Generate(pass string, salt string) string {
	passHash := md5.Sum([]byte(pass))
	passHashWithSalt := hex.EncodeToString(passHash[:]) + salt
	result := md5.Sum([]byte(passHashWithSalt))

	return hex.EncodeToString(result[:])
}

func ComparePasswordHash(password string, hash string, salt string) bool {
	return Generate(password, salt) == hash
}
