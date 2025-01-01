package pkg

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
)

type indexHeader struct {
	signature string
	version   int
	entryNum  int
}

type indexEntry struct {
	name string
	objectName string // hash
	permission string
	

}

// https://git-scm.com/docs/gitformat-index/2.40.0 を参考に実装
type index struct {
	indexHeader
	entries []indexEntry
}

func ReadIndexFile() *index {
	f, err := os.Open(".git/index")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	index := &index{}

	sig := make([]byte, 4)
	f.Read(sig)
	index.signature = string(sig)

	version := make([]byte, 4)
	f.Read(version)
	index.version = int(binary.BigEndian.Uint32(version))

	entryNum := make([]byte, 4)
	f.Read(entryNum)
	index.entryNum = int(binary.BigEndian.Uint32(entryNum))

	for range index.entryNum {
		entry := indexEntry{}

		// 32-bit ctime seconds, the last time a file's metadata changed
		f.Seek(4, 1)

		// 32-bit ctime nanosecond fractions
		f.Seek(4, 1)

		// 32-bit mtime seconds, the last time a file's data changed
		f.Seek(4, 1)

		// 32-bit mtime nanosecond fractions
		f.Seek(4, 1)

		// 32-bit dev
		f.Seek(4, 1)

		// 32-bit ino
		f.Seek(4, 1)

		// 32-bit mode, split into (high to low bits)
		// 16-bit unused, must be zero
		f.Seek(2, 1)

		// 4-bit object type
		// 3-bit unused
		// 9-bit unix permission. Only 0755 and 0644 are valid for regular files.
		f.Seek(2, 1)

		// 32-bit uid
		f.Seek(4, 1)

		// 32-bit gid
		f.Seek(4, 1)

		// 32-bit file size
		f.Seek(4, 1)

		// Object name for the represented object
		oName := make([]byte, 20)
		f.Read(oName)
		fmt.Println(hex.EncodeToString(oName))

		// A 16-bit 'flags' field split into (high to low bits)
		// 1-bit assume-valid flag
		// 1-bit extended flag (must be zero in version 2)
		// 2-bit stage (during merge)
		// 12-bit name length if the length is less than 0xFFF; otherwise 0xFFF is stored in this field.
		f.Seek(2, 1)

		name := make([]byte, 5)
		f.Read(name)
		entry.name = string(name)

		index.entries = append(index.entries, entry)
	}

	return index
}
