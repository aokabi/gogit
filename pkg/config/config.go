package config

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aokabi/gogit/pkg"
)

type user struct {
	email string
	name  string
}

type remote struct {
	url string
}

type config struct {
	user
	remotes map[string]remote
}

func new() *config {
	return &config{
		remotes: map[string]remote{},
	}
}

func Read() *config {
	globalConf := readGlobal()
	localConf := readLocal()

	// 本当は適用の優先順とか考慮する
	return &config{
		user:    globalConf.user,
		remotes: localConf.remotes,
	}
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

func readLocal() *config {
	f, err := pkg.Open("config")
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

func (c *config) GetRemoteUrl(repo string) string {
	return c.remotes[repo].url
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

[remote "origin"]

	url = https://github.com/aokabi/gogit
*/
func decode(r io.Reader) *config {
	conf, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	// (?:...) 非キャプチャグループ
	sectionRegex := regexp.MustCompile(`^\[([^\s"\]]+)(?:\s+"([^"]+)")?\]$`)
	keyValueRegex := regexp.MustCompile(`^(\w+)\s*=\s*(.+)$`)

	lines := strings.Split(string(conf), "\n")
	c := new()
	var currentSection, currentSubSection string
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// start section
		if matches := sectionRegex.FindStringSubmatch(line); matches != nil {
			currentSection = matches[1]
			if len(matches) > 2 {
				currentSubSection = matches[2]
			} else {
				currentSubSection = ""
			}
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
		case "remote":
			if currentSubSection != "" {
				if matches := keyValueRegex.FindStringSubmatch(line); matches != nil {
					// 初期化
					if _, ok := c.remotes[currentSubSection]; !ok {
						c.remotes[currentSubSection] = remote{}
					}
					remote := c.remotes[currentSubSection]

					switch matches[1] {
					case "url":
						remote.url = matches[2]
					}

					c.remotes[currentSubSection] = remote
				}
			}
		}
	}

	return c
}
