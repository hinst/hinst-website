package server

import "encoding/xml"

func readXml(data []byte, v interface{}) {
	AssertError(xml.Unmarshal(data, v))
}
