package pkg

import (
	"encoding/hex"
	"fmt"
	"iter"
	"strings"
)

type tree struct {
	entries []treeEntry
}

type treeEntry struct {
	perm     string
	objType  string
	hash     string
	filename string
}

func DecodeContent2Tree(o *GitObj) *tree {
	if o.objType != "tree" {
		panic(fmt.Sprintf("not tree: %s", o.objType))
	}

	// permission filename\x00hash というフォーマットで保存されている
	// 複数エントリがある場合はこれが区切り文字無しで連続している
	entries := make([]treeEntry, 0)
	left := string(o.content)
	for len(left) > 0 {
		strs := strings.SplitN(string(left), NullByte, 2)
		// sha1なので20byte
		hash := hex.EncodeToString([]byte(strs[1][:HashSize]))
		left = strs[1][HashSize:]

		// get entry type
		r, err := Decompress(ReadObjectFile(hash))
		if err != nil {
			panic(err)
		}
		defer r.Close()

		o, err := Parse(r)
		if err != nil {
			panic(err)
		}

		perm := strings.Split(strs[0], " ")[0]
		filename := strings.Split(strs[0], " ")[1]

		entries = append(entries, treeEntry{
			perm:     perm,
			objType:  o.objType,
			hash:     hash,
			filename: filename,
		})
	}

	return &tree{
		entries: entries,
	}
}

func (t *tree) Entries() iter.Seq[treeEntry] {
	return func(yield func(treeEntry) bool) {
		for _, e := range t.entries {
			if !yield(e) {
				return
			}
		}
	}
}

func (e treeEntry) GetPerm() string {
	return e.perm
}

func (e treeEntry) GetObjType() string {
	return e.objType
}

func (e treeEntry) GetHash() string {
	return e.hash
}

func (e treeEntry) GetFilename() string {
	return e.filename
}
