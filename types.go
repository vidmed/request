package request

import (
	"bytes"
	"fmt"
)

// Request is a representation of request - sequence of two english letters
type Request [2]byte

// String returns the r converted to string
func (r Request) String() string {
	return string(r.Bytes())
}

// Bytes returns the r converted to []byte
func (r Request) Bytes() []byte {
	return r[:]
}

// Views is a representation of Request view history - this is a map Request to number of its views
type Views map[Request]uint

// String returns the v converted to string
func (v Views) String() string {
	return string(v.Bytes())
}

// Bytes returns the v converted to []byte
func (v Views) Bytes() []byte {
	keys := GetSortedKeys(v)

	b := new(bytes.Buffer)
	for _, k := range keys {
		fmt.Fprintf(b, "%s - %d\n", k.String(), v[k])
	}
	return b.Bytes()
}
