package core

import (
	"encoding/xml"
	"io"
)

type Xml struct {
	c *Core
}

func (c *Core) Xml() Xml {
	return Xml{c}
}

// Shortcut to encoding/xml.Marshal
func (x Xml) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

// Shortcut to encoding/xml.MarshalIndent
func (x Xml) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return xml.MarshalIndent(v, prefix, indent)
}

// Shortcut to encoding/xml.Unmarshal
func (x Xml) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

// Shortcut to encoding/xml.NewDecoder
func (x Xml) NewDecoder(r io.Reader) *xml.Decoder {
	return xml.NewDecoder(r)
}

// Shortcut to encoding/xml.NewEncoder
func (x Xml) NewEncoder(w io.Writer) *xml.Encoder {
	return xml.NewEncoder(w)
}

// Output in XML
func (x Xml) Send(v interface{}) {
	w := x.c.Pub.Writers["gzip"]
	if w == nil {
		w = x.c
	}
	xml.NewEncoder(w).Encode(v)
}

// Decode Request Body
func (x Xml) DecodeReqBody(v interface{}) error {
	return x.NewDecoder(x.c.Req.Body).Decode(v)
}
