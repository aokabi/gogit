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
author {name} <{email}> {timestamp} {timezone}
committer {name} <{email}> {timestamp} {timezone}

{commit message}
*/

type commit struct {
	treeHash           string
	author             string
	authorEmail        string
	authorTimestamp    time.Time
	committer          string
	committerEmail     string
	committerTimestamp time.Time
	message            string
}

func NewCommit(treeHash string, author string, authorEmail string, authorTimestamp time.Time, committer string, committerEmail string, committerTimestamp time.Time, message string) *commit {
	return &commit{
		treeHash:           treeHash,
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

	contents = append(contents, fmt.Sprintf("author %s <%s> %d %s\n", c.author, c.authorEmail, c.authorTimestamp.Unix(), offsetStr(c.authorTimestamp)))

	contents = append(contents, fmt.Sprintf("committer %s <%s> %d %s\n\n", c.committer, c.committerEmail, c.committerTimestamp.Unix(), offsetStr(c.committerTimestamp)))

	contents = append(contents, c.message+"\n")

	return NewGitObj(COMMIT, []byte(strings.Join(contents, "")))

}

func DecodeCommit(o *GitObj) *commit {
	if o.objType != "commit" {
		panic(fmt.Sprintf("not commit: %s", o.objType))
	}

	lines := strings.Split(string(o.content), "\n")
	treeHash := strings.Split(lines[0], " ")[1]
	authorLine := strings.Split(lines[1], " ")
	committerLine := strings.Split(lines[2], " ")
	messageLine := strings.Join(lines[3:], "\n")

	return &commit{
		treeHash:           treeHash,
		author:             authorLine[1],
		authorEmail:        authorLine[2],
		authorTimestamp:    parseTimestamp(authorLine[3], authorLine[4]),
		committer:          committerLine[1],
		committerEmail:     committerLine[2],
		committerTimestamp: parseTimestamp(committerLine[3], committerLine[4]),
		message:            messageLine,
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
