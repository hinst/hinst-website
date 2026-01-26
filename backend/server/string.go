package server

import (
	"log"
	"os"
	"strconv"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/text/unicode/norm"
)

func getIntFromString(text string) int {
	return assertResultError(strconv.Atoi(text))
}

func getInt64FromString(text string) int64 {
	return assertResultError(strconv.ParseInt(text, 10, 64))
}

func getStringFromInt64(number int64) string {
	return strconv.FormatInt(number, 10)
}

func getStringFromInt(number int) string {
	return strconv.Itoa(number)
}

func getStringFromBool(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func getQuotedString(text string) string {
	return "\"" + text + "\""
}

func requireEnvVar(name string) string {
	var value = os.Getenv(name)
	if value == "" {
		log.Fatalln("Environment variable is required:", name)
	}
	return value
}

func normalizeString(text string) string {
	return norm.NFC.String(text)
}

func stripHtml(text string) string {
	return bluemonday.StrictPolicy().Sanitize(text)
}
