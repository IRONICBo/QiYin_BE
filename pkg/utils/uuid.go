package utils

import (
	"fmt"

	"github.com/gofrs/uuid"
)

// GenUUID generate uuid.
func GenUUID() string {
	return uuid.Must(uuid.NewV4()).String()
}

// GenUUIDWithoutHyphen generate uuid without hyphen.
func GenUUIDWithoutHyphen() string {
	return toString(uuid.Must(uuid.NewV4()))
}

// encodeCanonical encodes the canonical RFC-4122 form of UUID u into the
// first 36 bytes dst.
func encodeCanonical(dst []byte, u uuid.UUID) {
	const hextable = "0123456789abcdef"
	dst[8] = '-'
	dst[13] = '-'
	dst[18] = '-'
	dst[23] = '-'
	for i, x := range [16]byte{
		0, 2, 4, 6,
		9, 11,
		14, 16,
		19, 21,
		24, 26, 28, 30, 32, 34,
	} {
		c := u[i]
		dst[x] = hextable[c>>4]
		dst[x+1] = hextable[c&0x0f]
	}
}

// String returns a canonical RFC-4122 string representation of the UUID:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx, but delete separator -.
func toString(u uuid.UUID) string {
	var buf [36]byte
	encodeCanonical(buf[:], u)
	return fmt.Sprintf("%s%s%s%s%s", buf[0:7], buf[9:12], buf[14:17], buf[19:22], buf[24:])
}
