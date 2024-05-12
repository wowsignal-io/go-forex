package internal

import (
	"bytes"
	"encoding/csv"
	"io"
)

func SkipLineBytes(p []byte) ([]byte, error) {
	idx := bytes.IndexByte(p, '\n')
	if idx < 0 {
		return nil, io.EOF
	}
	return p[idx+1:], nil
}

func SkipLinesBytes(p []byte, n int) ([]byte, error) {
	var err error
	for i := 0; i < n; i++ {
		p, err = SkipLineBytes(p)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func SkipLinesCSV(cr *csv.Reader, n int) error {
	for i := 0; i < n; i++ {
		_, err := cr.Read()
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		if err != nil {
			return err
		}
	}
	return nil
}
