package pkg

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type header struct {
	objType string
	size    int
}

type gitObj struct {
	header
	content []byte
}

func NewBlob(content []byte) *gitObj {
	return &gitObj{
		header: header{
			objType: "blob",
			size:    len(content),
		},
		content: content,
	}
}

func (o *gitObj) GetContent() []byte {
	return o.content
}

func Parse(r io.Reader) (*gitObj, error) {
	all, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	objType := strings.Split(string(all), " ")[0]

	second := strings.Split(string(all), " ")[1]
	// split null byte
	size := strings.Split(second, "\x00")[0]
	content := strings.Split(second, "\x00")[1]
	sizeInt, err := strconv.Atoi(size)

	if err != nil {
		return nil, err
	}

	return &gitObj{
		header: header{
			objType: objType,
			size:    sizeInt,
		},
		content: []byte(content),
	}, nil
}

// ハッシュはヘッダ込みで計算する
func (o *gitObj) Hash() string {
	store := fmt.Sprintf("%s %d\x00%s", o.objType, o.size, o.content)
	hash := sha1.New()
	hash.Write([]byte(store))
	return hex.EncodeToString(hash.Sum(nil))
}

func (o *gitObj) Store(w io.Writer) {
	store := fmt.Sprintf("%s %d\x00%s", o.objType, o.size, o.content)
	w.Write([]byte(store))
}
