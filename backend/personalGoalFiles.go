package main

import (
	"regexp"
)

var goalIdStringMatcher = regexp.MustCompile(`^\d{1,10}$`)
var goalFileNameMatcher = regexp.MustCompile(`^\d\d\d\d-\d\d-\d\d`)
