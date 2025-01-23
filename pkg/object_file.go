package pkg

import (
	"compress/zlib"
	"io"
	"os"
	"path/filepath"
)

func ReadObjectFile(hash string) *os.File {
	f, err := Open(filepath.Join("objects", hash[:2], hash[2:]))
	if err != nil {
		panic(err)
	}
	return f
}

func Decompress(f *os.File) (io.ReadCloser, error) {
	r, err := zlib.NewReader(f)
	if err != nil {
		panic(err)
	}

	return r, nil
}
