package internal

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Fetch returns the given resource data, handling URLs (including simple data
// URLs), as well as filesystem paths.
func Fetch(resource string) ([]byte, error) {
	u, err := url.Parse(resource)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "http", "https":
		return download(resource)
	case "data":
		return decode(u.Opaque)
	case "":
		return os.ReadFile(resource)
	default:
		return nil, fmt.Errorf("don't know how to handle URL scheme %s", u.Scheme)
	}
}

func download(uri string) ([]byte, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var b bytes.Buffer
	_, err = io.Copy(&b, resp.Body)
	return b.Bytes(), err
}

// decode handles common data URLs - it bails if it runs into anything at all
// complicated, but does support base64.
func decode(dataURL string) ([]byte, error) {
	idx := strings.IndexByte(dataURL, ',')
	if idx < 0 {
		return nil, fmt.Errorf("invalid data URL %q", dataURL)
	}

	spec := dataURL[:idx]
	payload, err := url.PathUnescape(dataURL[idx+1:])
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(spec, ";base64") {
		return base64.StdEncoding.DecodeString(payload)
	}

	return []byte(payload), nil
}
