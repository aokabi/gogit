package pkg

import (
	"os"
	"path/filepath"
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
