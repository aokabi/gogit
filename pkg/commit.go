package pkg

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// commit objectのフォーマット
/*
tree {hash}
[parent {hash}]
author {name} <{email}> {timestamp} {timezone}
committer {name} <{email}> {timestamp} {timezone}

{commit message}
*/

type commit struct {
	treeHash           string
	parent             string
	author             string
	authorEmail        string
	authorTimestamp    time.Time
	committer          string
	committerEmail     string
	committerTimestamp time.Time
	message            string
}

func NewCommit(treeHash string, parent string, author string, authorEmail string, authorTimestamp time.Time, committer string, committerEmail string, committerTimestamp time.Time, message string) *commit {
	return &commit{
		treeHash:           treeHash,
		parent:             parent,
		author:             author,
		authorEmail:        authorEmail,
		authorTimestamp:    authorTimestamp,
		committer:          committer,
		committerEmail:     committerEmail,
		committerTimestamp: committerTimestamp,
		message:            message,
	}
}

func (c *commit) EncodeCommit() *GitObj {
	contents := make([]string, 0)

	contents = append(contents, fmt.Sprintf("tree %s\n", c.treeHash))

	if c.parent != "" {
		contents = append(contents, fmt.Sprintf("parent %s\n", c.parent))
	}

	contents = append(contents, fmt.Sprintf("author %s <%s> %d %s\n", c.author, c.authorEmail, c.authorTimestamp.Unix(), offsetStr(c.authorTimestamp)))

	contents = append(contents, fmt.Sprintf("committer %s <%s> %d %s\n\n", c.committer, c.committerEmail, c.committerTimestamp.Unix(), offsetStr(c.committerTimestamp)))

	contents = append(contents, c.message+"\n")

	return NewGitObj(COMMIT, []byte(strings.Join(contents, "")))

}

func DecodeCommit(o *GitObj) *commit {
	if o.objType != "commit" {
		panic(fmt.Sprintf("not commit: %s", o.objType))
	}

	tmp := strings.Split(string(o.content), "\n\n")
	headers := strings.Split(tmp[0], "\n")
	message := tmp[1]

	var treeHash string
	var parent string
	var authorLine []string
	var committerLine []string

	for _, h := range headers {
		prefix := strings.Split(h, " ")[0]
		switch prefix {
		case "tree":
			treeHash = strings.Split(h, " ")[1]	
		case "parent":
			parent = strings.Split(h, " ")[1]
		case "author":
			authorLine = strings.Split(h, " ")
		case "committer":
			committerLine = strings.Split(h, " ")
		
		}
	}

	return &commit{
		treeHash:           treeHash,
		parent:             parent,
		author:             authorLine[1],
		authorEmail:        authorLine[2],
		authorTimestamp:    parseTimestamp(authorLine[3], authorLine[4]),
		committer:          committerLine[1],
		committerEmail:     committerLine[2],
		committerTimestamp: parseTimestamp(committerLine[3], committerLine[4]),
		message:            message,
	}
}

func parseTimestamp(unixTimeStr string, offsetStr string) time.Time {
	unixtime, _ := strconv.Atoi(unixTimeStr)
	utcTime := time.Unix(int64(unixtime), 0).In(time.UTC)
	t, _ := time.Parse(time.RFC1123Z, fmt.Sprintf("%s %s", utcTime.Format(time.RFC1123), offsetStr))

	return t
}

func offsetStr(t time.Time) string {
	_, offset := t.Local().Zone()
	sign := "+"
	if offset < 0 {
		offset = -offset
		sign = "-"
	}
	hr := offset / 3600
	min := offset % 3600 / 60
	return fmt.Sprintf("%s%02d%02d", sign, hr, min)
}
