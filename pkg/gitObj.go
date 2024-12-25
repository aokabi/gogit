package pkg

import (
	"io"
	"strconv"
	"strings"
)

type header struct {
	objType string
	size    int
}

type gitObj struct {
	header
	content []byte
}

func (o *gitObj) GetContent() []byte {
	return o.content
}

func Parse(r io.Reader) (*gitObj, error) {
	all, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	objType := strings.Split(string(all), " ")[0]

	second := strings.Split(string(all), " ")[1]
	// split null byte
	size := strings.Split(second, "\x00")[0]
	content := strings.Split(second, "\x00")[1]
	sizeInt, err := strconv.Atoi(size)

	if err != nil {
		return nil, err
	}

	return &gitObj{
		header: header{
			objType: objType,
			size:    sizeInt,
		},
		content: []byte(content),
	}, nil
}
