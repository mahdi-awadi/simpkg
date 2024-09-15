package image

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// encode is our main function for
// base64 encoding a passed []byte
func encode(bin []byte) []byte {
	e64 := base64.StdEncoding

	maxEncLen := e64.EncodedLen(len(bin))
	encBuf := make([]byte, maxEncLen)

	e64.Encode(encBuf, bin)
	return encBuf
}

// FromBuffer accepts a buffer and returns a
// base64 encoded string.
func FromBuffer(buf bytes.Buffer) string {
	enc := encode(buf.Bytes())
	mime := http.DetectContentType(buf.Bytes())

	return format(enc, mime)
}

// FromLocal reads a local file and returns
// the base64 encoded version.
func FromLocal(fname string) (string, error) {
	var b bytes.Buffer

	fileExists, _ := exists(fname)
	if !fileExists {
		return "", fmt.Errorf("File does not exist\n")
	}

	file, err := os.Open(fname)
	if err != nil {
		return "", fmt.Errorf("err opening file\n")
	}

	_, err = b.ReadFrom(file)
	if err != nil {
		return "", fmt.Errorf("err reading file to buffer\n")
	}

	return FromBuffer(b), nil
}

// format is an abstraction of the mime switch to create the
// acceptable base64 string needed for browsers.
func format(enc []byte, mime string) string {
	switch mime {
	case "image/gif", "image/jpeg", "image/pjpeg", "image/png", "image/tiff":
		return fmt.Sprintf("data:%s;base64,%s", mime, enc)
	default:
	}

	return fmt.Sprintf("data:image/png;base64,%s", enc)
}

// cleanUrl converts whitespace in urls to %20
func cleanUrl(s string) string {
	return strings.Replace(s, " ", "%20", -1)
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
