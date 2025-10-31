package utils

import (
    "crypto/sha512"
    "encoding/hex"
)

// SHA512Hex returns lowercase hex-encoded SHA512 of the input string
func SHA512Hex(input string) string {
    sum := sha512.Sum512([]byte(input))
    return hex.EncodeToString(sum[:])
}


