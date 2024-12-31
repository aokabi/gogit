package pkg

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	NullByte = "\x00"
	HashSize = 20
)

type header struct {
	objType string
	size    int
}

type GitObj struct {
	header
	content []byte
}

func (o *GitObj) GetObjType() string {
	return o.objType
}

func Parse(r io.Reader) (*GitObj, error) {
	all, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	objType := strings.SplitN(string(all), " ", 2)[0]

	second := strings.SplitN(string(all), " ", 2)[1]
	// split null byte
	size := strings.SplitN(second, NullByte, 2)[0]
	content := strings.SplitN(second, NullByte, 2)[1]
	sizeInt, err := strconv.Atoi(size)

	if err != nil {
		return nil, err
	}

	return &GitObj{
		header: header{
			objType: objType,
			size:    sizeInt,
		},
		content: []byte(content),
	}, nil
}

// ハッシュはヘッダ込みで計算する
func (o *GitObj) Hash() string {
	store := fmt.Sprintf("%s %d\x00%s", o.objType, o.size, o.content)
	hash := sha1.New()
	hash.Write([]byte(store))
	return hex.EncodeToString(hash.Sum(nil))
}

func (o *GitObj) Store(w io.Writer) {
	store := fmt.Sprintf("%s %d\x00%s", o.objType, o.size, o.content)
	w.Write([]byte(store))
}
