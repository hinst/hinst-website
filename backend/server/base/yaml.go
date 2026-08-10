package base

import (
	"github.com/goccy/go-yaml"
	"github.com/hinst/go-gophers"
)

var yamlEncodeOption = yaml.UseLiteralStyleIfMultiline(true)

func EncodeYaml(object any) []byte {
	return gophers.AssertResultError(yaml.MarshalWithOptions(object, yamlEncodeOption))
}
