package problem

import (
	"encoding/xml"
	"errors"
	"net/http"
)

type RegisteredMembers struct {
	TypeURI       string `json:"type,omitempty"     xml:"type,omitempty"`
	TypeTitle     string `json:"title,omitempty"    xml:"title,omitempty"`
	StatusCode    int    `json:"status,omitempty"   xml:"status,omitempty"`
	DetailMessage string `json:"detail,omitempty"   xml:"detail,omitempty"`
	InstanceURI   string `json:"instance,omitempty" xml:"instance,omitempty"`

	XMLName xml.Name     `json:"-" xml:"problem"`
	XMLNS   XMLNamespace `json:"-" xml:"xmlns,attr"`
}

func (t RegisteredMembers) Type() string {
	return t.TypeURI
}

func (t RegisteredMembers) Status() int {
	return t.StatusCode
}

func (t RegisteredMembers) Title() string {
	return t.TypeTitle
}

func (t RegisteredMembers) Detail() string {
	return t.DetailMessage
}

func (t RegisteredMembers) Instance() string {
	return t.InstanceURI
}

type XMLNamespace struct{}

var (
	_ xml.MarshalerAttr   = XMLNamespace{}
	_ xml.UnmarshalerAttr = (*XMLNamespace)(nil)
)

func (_ XMLNamespace) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: xml.Name{Local: "xmlns"}, Value: problemNamespace}, nil
}

func (_ *XMLNamespace) UnmarshalXMLAttr(attr xml.Attr) error {
	if attr.Value != problemNamespace || attr.Name.Local != "xmlns" {
		return errors.New("problem-xml: attribute has not the correct namespace")
	}
	return nil
}

const problemNamespace = "urn:ietf:rfc:7807"

type Problem interface {
	Status() int
	http.Handler
}
