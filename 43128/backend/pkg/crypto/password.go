package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(pwd string) (string, error) {
	h := sha256.Sum256([]byte(pwd + "event-platform-salt"))
	return hex.EncodeToString(h[:]), nil
}

func CheckPassword(hash, pwd string) bool {
	h := sha256.Sum256([]byte(pwd + "event-platform-salt"))
	return hex.EncodeToString(h[:]) == hash
}
