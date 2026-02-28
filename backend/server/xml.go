package server

import (
	"encoding/xml"

	"github.com/hinst/go-common"
)

func readXml(data []byte, v interface{}) {
	common.AssertError(xml.Unmarshal(data, v))
}
