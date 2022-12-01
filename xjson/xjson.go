// Package xjson implements JSON helpers.
package xjson

import (
	"bytes"
	"encoding/json"
)

// Marshal is like json.MarshalIndent but does not escape HTML.
// And does not return an error
func Marshal(v interface{}) []byte {
	var b bytes.Buffer
	e := json.NewEncoder(&b)
	// Allows < and > in JSON strings without escaping which we do with SrcArrow and
	// DstArrow. See https://stackoverflow.com/a/28596225
	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")
	err := e.Encode(v)
	if err != nil {
		buf, _ := json.Marshal(err.Error())
		return buf
	}
	buf := b.Bytes()
	// Remove trailing newline.
	return buf[:len(buf)-1]
}
