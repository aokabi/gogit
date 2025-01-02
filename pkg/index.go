package pkg

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"syscall"
	"time"
)

type indexHeader struct {
	signature string
	version   int
	entryNum  int
}

type indexEntry struct {
	ctimeSec   time.Time
	ctimeNsec  int64
	mtimeSec   time.Time
	mtimeNsec  int64
	dev        uint32
	inode      uint32
	mode       os.FileMode
	uid        uint32
	gid        uint32
	filesize   uint32
	objectName string // hash
	flags      flag   // 16 bit
	name       string
}

// 1-bit assume-valid flag
// 1-bit extended flag (must be zero in version 2)
// 2-bit stage (during merge)
// 12-bit name length if the length is less than 0xFFF; otherwise 0xFFF is stored in this field.
type flag struct {
	assumeValid bool
	extended    bool
	stage       int
	nameLength  int
}

// TREE
// NULL終端のパス文字列(ルートツリーだとnullだけ)
// この拡張でのエントリの数(ASCII数値)
// 0x20(空白)
// このツリーに含まれるサブツリーの数(ASCII数値)
// 0x0A(改行)
// 例：00 2d 31 20 30 0a
// 00 → パス文字列
// 2d 31 → エントリの数(-1)
// 20 → 空白
// 30 → サブツリーの数(0)
// 0a → 改行
// ASCII数値は数値を文字列としてバイト列に変換する
type extension struct {
	signature string
	size      uint32
	data      []byte // 32-bit
}

// https://git-scm.com/docs/gitformat-index/2.40.0 を参考に実装
type index struct {
	indexHeader
	entries   []indexEntry
	extension extension
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
		var byteSum int64 = 0

		// 32-bit ctime seconds, the last time a file's metadata changed
		ctimeSec := make([]byte, 4)
		n, _ := f.Read(ctimeSec)
		byteSum += int64(n)
		entry.ctimeSec = time.Unix(int64(binary.BigEndian.Uint32(ctimeSec)), 0)

		// 32-bit ctime nanosecond fractions
		ctimeNSec := make([]byte, 4)
		n, _ = f.Read(ctimeNSec)
		byteSum += int64(n)
		entry.ctimeNsec = int64(binary.BigEndian.Uint32(ctimeNSec))

		// 32-bit mtime seconds, the last time a file's data changed
		mtimeSec := make([]byte, 4)
		n, _ = f.Read(mtimeSec)
		byteSum += int64(n)
		entry.mtimeSec = time.Unix(int64(binary.BigEndian.Uint32(mtimeSec)), 0)

		// 32-bit mtime nanosecond fractions
		mtimeNSec := make([]byte, 4)
		n, _ = f.Read(mtimeNSec)
		byteSum += int64(n)
		entry.mtimeNsec = int64(binary.BigEndian.Uint32(mtimeNSec))

		// 32-bit dev
		dev := make([]byte, 4)
		n, _ = f.Read(dev)
		byteSum += int64(n)
		entry.dev = binary.BigEndian.Uint32(dev)

		// 32-bit ino
		inode := make([]byte, 4)
		n, _ = f.Read(inode)
		byteSum += int64(n)
		entry.inode = binary.BigEndian.Uint32(inode)

		// 32-bit mode, split into (high to low bits)
		// 16-bit unused, must be zero
		// 4-bit object type
		// 3-bit unused
		// 9-bit unix permission. Only 0755 and 0644 are valid for regular files.
		mode := make([]byte, 4)
		n, _ = f.Read(mode)
		byteSum += int64(n)
		modeInt := binary.BigEndian.Uint32(mode)
		entry.mode = os.FileMode(modeInt)

		// 32-bit uid
		uid := make([]byte, 4)
		n, _ = f.Read(uid)
		byteSum += int64(n)
		entry.uid = binary.BigEndian.Uint32(uid)

		// 32-bit gid
		gid := make([]byte, 4)
		n, _ = f.Read(gid)
		byteSum += int64(n)
		entry.gid = binary.BigEndian.Uint32(gid)

		// 32-bit file size
		filesize := make([]byte, 4)
		n, _ = f.Read(filesize)
		byteSum += int64(n)
		entry.filesize = binary.BigEndian.Uint32(filesize)

		// Object name for the represented object
		oName := make([]byte, 20)
		n, _ = f.Read(oName)
		byteSum += int64(n)
		entry.objectName = hex.EncodeToString(oName)

		// A 16-bit 'flags' field split into (high to low bits)
		// 1-bit assume-valid flag
		// 1-bit extended flag (must be zero in version 2)
		// 2-bit stage (during merge)
		// 12-bit name length if the length is less than 0xFFF; otherwise 0xFFF is stored in this field.
		flags := make([]byte, 2)
		n, _ = f.Read(flags)
		byteSum += int64(n)
		flagsInt := binary.BigEndian.Uint16(flags)
		length := int(flagsInt & 0xFFF)
		entry.flags = flag{assumeValid: (flagsInt>>15)&1 == 1, extended: (flagsInt>>14)&1 == 1, stage: int(flagsInt >> 12), nameLength: length}

		name := make([]byte, length)
		n, _ = f.Read(name)
		byteSum += int64(n)
		entry.name = string(name)

		// 1-8 nul bytes as necessary to pad the entry to a multiple of eight bytes while keeping the name NUL-terminated.
		f.Seek(8-byteSum%8, 1)

		index.entries = append(index.entries, entry)
	}

	// read extension
	signature := make([]byte, 4)
	_, err = f.Read(signature)
	if err != nil {
		panic(err)
	}

	size := make([]byte, 4)
	_, err = f.Read(size)
	if err != nil {
		panic(err)
	}
	sizeInt := binary.BigEndian.Uint32(size)

	data := make([]byte, sizeInt)
	_, err = f.Read(data)
	if err != nil {
		panic(err)
	}

	index.extension = extension{
		signature: string(signature),
		size:      sizeInt,
		data:      data,
	}

	return index
}

func AddEntry(objects map[string]*GitObj) {
	index := ReadIndexFile()

	// truncate index file
	if err := os.Truncate(".git/index", 0); err != nil {
		panic(err)
	}

	f, err := os.OpenFile(".git/index", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for fileName, obj := range objects {
		addf, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}

		// create index entry
		fileinfo, err := addf.Stat()
		if err != nil {
			panic(err)
		}
		newEntry := indexEntry{
			ctimeSec:   fileinfo.ModTime(),
			mtimeSec:   fileinfo.ModTime(),
			dev:        uint32(fileinfo.Sys().(*syscall.Stat_t).Dev),
			inode:      uint32(fileinfo.Sys().(*syscall.Stat_t).Ino),
			mode:       fileinfo.Mode(),
			uid:        fileinfo.Sys().(*syscall.Stat_t).Uid,
			gid:        fileinfo.Sys().(*syscall.Stat_t).Gid,
			filesize:   uint32(fileinfo.Size()),
			ctimeNsec:  int64(fileinfo.ModTime().Nanosecond()),
			mtimeNsec:  int64(fileinfo.ModTime().Nanosecond()),
			objectName: obj.Hash(),
			flags:      flag{assumeValid: false, extended: false, stage: 0, nameLength: len(fileName)},
			name:       fileName,
		}

		index.entries = append(index.entries, newEntry)
		index.entryNum++
		addf.Close()
	}

	if _, err := f.Write(index.encodeBinary()); err != nil {
		panic(err)
	}
	// fmt.Println(hex.Dump(index.encodeBinary()))

}

func (index *index) encodeBinary() []byte {
	encoded := make([]byte, 0)
	// 文字列はbyte sliceにするだけ
	encoded = append(encoded, []byte(index.signature)...)
	// 数値はbinaryパッケージを使う
	encoded = binary.BigEndian.AppendUint32(encoded, uint32(index.version))
	encoded = binary.BigEndian.AppendUint32(encoded, uint32(index.entryNum))

	fmt.Println("entryNum:", len(index.entries))
	for _, entry := range index.entries {
		byteSum := 0
		// ctimeSec
		encoded = binary.BigEndian.AppendUint32(encoded, uint32(entry.ctimeSec.Unix()))
		byteSum += 4
		// ctimeNsec
		encoded = binary.BigEndian.AppendUint32(encoded, uint32(entry.ctimeNsec))
		byteSum += 4
		// mtimeSec
		encoded = binary.BigEndian.AppendUint32(encoded, uint32(entry.mtimeSec.Unix()))
		byteSum += 4
		// mtimeNsec
		encoded = binary.BigEndian.AppendUint32(encoded, uint32(entry.mtimeNsec))
		byteSum += 4
		// dev
		encoded = binary.BigEndian.AppendUint32(encoded, entry.dev)
		byteSum += 4
		// inode
		encoded = binary.BigEndian.AppendUint32(encoded, entry.inode)
		byteSum += 4
		// mode
		encoded = append(encoded, encodeBinary(entry.mode)...)
		byteSum += 4
		// uid
		encoded = binary.BigEndian.AppendUint32(encoded, entry.uid)
		byteSum += 4
		// gid
		encoded = binary.BigEndian.AppendUint32(encoded, entry.gid)
		byteSum += 4
		// filesize
		encoded = binary.BigEndian.AppendUint32(encoded, entry.filesize)
		byteSum += 4
		// objectName
		decodedObjectName, _ := hex.DecodeString(entry.objectName)
		encoded = append(encoded, decodedObjectName...)
		byteSum += 20
		// flags
		encoded = append(encoded, entry.flags.encodeBinary()...)
		byteSum += 2
		// name
		encoded = append(encoded, []byte(entry.name)...)
		byteSum += len(entry.name)

		// 1-8 nul bytes as necessary to pad the entry to a multiple of eight bytes while keeping the name NUL-terminated.
		encoded = append(encoded, make([]byte, 8-byteSum%8)...)
	}

	// extensions
	encoded = append(encoded, index.extension.encodeBinary()...)

	// hash checksum of the index file
	hash := sha1.New()
	hash.Write(encoded)
	encoded = append(encoded, hash.Sum(nil)...)

	return encoded
}

func (f *flag) encodeBinary() []byte {
	var tmp uint16
	if f.assumeValid {
		tmp = 1
	}

	tmp = tmp << 1
	if f.extended {
		tmp += 1
	}

	tmp = tmp << 2
	tmp += uint16(f.stage)

	tmp = tmp << 12
	tmp += uint16(f.nameLength)

	return binary.BigEndian.AppendUint16([]byte{}, tmp)
}

func encodeBinary(m os.FileMode) []byte {
	// 	16-bit unused, must be zero
	// 4-bit object type valid values in binary are 1000 (regular file), 1010 (symbolic link) and 1110 (gitlink)
	// 3-bit unused, must be zero
	// 9-bit unix permission. Only 0755 and 0644 are valid for regular files. Symbolic links and gitlinks have value 0 in this field.

	var tmp uint32
	if m.IsRegular() {
		tmp += 8
	} else {
		tmp += 10
	}
	// TODO gitlinksってなに？
	tmp = tmp << 3 // zero
	tmp = tmp << 9

	// regular fileの場合は0755か0644しかない
	// よくわかっていないのでとりあえず0644
	tmp += uint32(0644)
	ret := binary.BigEndian.AppendUint32([]byte{}, tmp)
	ret[0] = 0
	ret[1] = 0
	return ret
}

func (e *extension) encodeBinary() []byte {
	encoded := make([]byte, 0)

	encoded = append(encoded, []byte(e.signature)...)
	encoded = binary.BigEndian.AppendUint32(encoded, uint32(e.size))
	encoded = append(encoded, e.data...)

	return encoded
}
