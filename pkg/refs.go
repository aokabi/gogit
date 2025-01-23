package pkg

import (
	"io"
	"strings"
)

func UpdateRefs(ref string, newValue string) {
	f, err := CreateFile(ref)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.WriteString(newValue + "\n"); err != nil {
		panic(err)
	}
}

func ReadHEAD() string {
	f, err := Open("HEAD")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	ref := strings.Split(string(content), " ")
	return strings.Trim(ref[1], "\n")
}

func ReadRef(ref string) string {
	f, err := Open(ref)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	return string(content)
}
