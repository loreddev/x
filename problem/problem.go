package problem

import (
	"fmt"
)

type RegisteredMembers struct {
	Type     Type   `json:"type"               xml:"type"`
	Status   uint8  `json:"status"             xml:"status"`
	Title    string `json:"title"              xml:"title"`
	Detail   string `json:"detail,omitempty"   xml:"detail,omitempty"`
	Instance string `json:"instance,omitempty" xml:"instance,omitempty"`
}

type Members interface {
	Type() Type
}

func NewProblemType(URI, title string, status uint8, t Problem) Type {
	ty := Type{URI: URI, Title: title, Status: status}
	if _, ok := problemTypeList[URI]; ok {
		panic(fmt.Sprintf("problem: Problem type is already registered %q", URI))
	}
	problemTypeList[URI] = t
	return ty
}

type Type struct {
	URI    string
	Title  string
	Status uint8
}

type Problem interface {
	Type() Type
	SetType(Type)
}

var problemTypeList = map[string]Problem{}
