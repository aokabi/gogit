package pkg


func (o *GitObj) DecodeContent2Blob() string {
	return string(o.content)
}

func NewBlob(content []byte) *GitObj {
	return &GitObj{
		header: header{
			objType: "blob",
			size:    len(content),
		},
		content: content,
	}
}