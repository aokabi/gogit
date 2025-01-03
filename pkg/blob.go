package pkg


func (o *GitObj) DecodeContent2Blob() string {
	return string(o.content)
}
