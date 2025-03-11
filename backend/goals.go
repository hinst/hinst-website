package main

import "regexp"

var GoalFileNameMatcher = regexp.MustCompile(`^\d\d\d\d-\d\d-\d\d`)
