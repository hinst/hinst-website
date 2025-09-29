package server

import "encoding/xml"

func readXml(data []byte, v interface{}) {
	assertError(xml.Unmarshal(data, v))
}
