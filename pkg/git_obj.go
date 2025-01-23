package pkg

import (
	"compress/zlib"
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

type objectType string

const (
	BLOB   objectType = "blob"
	TREE   objectType = "tree"
	COMMIT objectType = "commit"
)

type header struct {
	objType objectType
	size    int
}

type GitObj struct {
	header
	content []byte
}

func NewGitObj(objType objectType, content []byte) *GitObj {
	return &GitObj{
		header: header{
			objType: objType,
			size:    len(content),
		},
		content: content,
	}
}

func (o *GitObj) GetObjType() objectType {
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
			objType: objectType(objType),
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

func (o *GitObj) Store() {
	hash := o.Hash()

	// save object
	dir := fmt.Sprintf("objects/%s", hash[:2])
	if IsNotExist(dir) {
		if err := CreateDir(dir); err != nil {
			panic(err)
		}
	}
	wf, err := CreateFile(fmt.Sprintf("objects/%s/%s", hash[:2], hash[2:]))
	if err != nil {
		panic(err)
	}
	defer wf.Close()

	zipWriter := zlib.NewWriter(wf)
	defer zipWriter.Close()

	store := fmt.Sprintf("%s %d\x00%s", o.objType, o.size, o.content)
	zipWriter.Write([]byte(store))

}
