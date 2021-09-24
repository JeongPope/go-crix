package utils

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io/ioutil"
)

func DecompressFlate(data []byte) ([]byte, error) {
	return ioutil.ReadAll(flate.NewReader(bytes.NewReader(data)))
}

func DecompressGzip(data []byte) (resData []byte, err error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}
