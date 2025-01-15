package pkg

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func UpdateRefs(ref string, newValue string) {
	filename := filepath.Join(".git/", ref)
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.WriteString(newValue + "\n"); err != nil {
		panic(err)
	}
}

func ReadHEAD() string {
	f, err := os.Open(".git/HEAD")
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
	f, err := os.Open(filepath.Join(".git/", ref))
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
