package config

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type user struct {
	email string
	name string
}

type config struct {
	user
}

func Read() *config {
	return readGlobal()
}

func readGlobal() *config {
	homedir := os.Getenv("HOME")
	f, err := os.Open(filepath.Join(homedir, ".gitconfig"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	return decode(f)
}

func (c *config) GetEmail() string {
	return c.user.email
}

func (c *config) GetName() string {
	return c.user.name
}

/*
https://git-scm.com/docs/git-config#_configuration_file

例
[fetch]

	prune = true

[rebase]

	autosquash = true

[include]

	path = ~/.gitconfig.local

[user]

	email = aokabit@gmail.com
	name = aokabi
*/
func decode(r io.Reader) *config {
	conf, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	sectionRegex := regexp.MustCompile(`^\[(.+)]$`)
	keyValueRegex := regexp.MustCompile(`^(\w+)\s*=\s*(.+)$`)

	lines := strings.Split(string(conf), "\n")
	c := &config{}
	var currentSection string
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// start section
		if matches := sectionRegex.FindStringSubmatch(line); matches != nil {
			currentSection = matches[1]
			continue
		}

		switch currentSection {
		case "user":
			if matches := keyValueRegex.FindStringSubmatch(line); matches != nil {
				switch matches[1] {
				case "email":
					c.user.email = matches[2]
				case "name":
					c.user.name = matches[2]
				}
			}
		}
	}

	return c
}
