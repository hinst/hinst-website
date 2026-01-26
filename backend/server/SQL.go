package server

import "strings"

type SQLite struct {
}

// Query: LIKE ? ESCAPE '\\'
// Parameter: "%"" + escapeLikeString(text) + "%"
func (SQLite) escapeLikeString(text string) string {
	text = strings.ReplaceAll(text, "\\", "\\\\")
	text = strings.ReplaceAll(text, "_", "\\_")
	text = strings.ReplaceAll(text, "%", "\\%")
	return text
}

// Find substring in field. Parameter matcherText will be escaped
func (SQLite) matchSubstring(fieldName string, matcherText string, parameters *[]any) string {
	*parameters = append(*parameters, "%"+SQLite{}.escapeLikeString(matcherText)+"%")
	return "UPPER(" + fieldName + ") LIKE ? ESCAPE '\\'"
}

// Find substring in multiple fields. Parameter matcherText will be escaped
func (SQLite) matchSubstringInFields(fieldNames []string, matcherText string, parameters *[]any) string {
	var fieldQueries []string
	for _, fieldName := range fieldNames {
		var query = SQLite{}.matchSubstring(fieldName, matcherText, parameters)
		fieldQueries = append(fieldQueries, "("+query+")")
	}
	return strings.Join(fieldQueries, " OR ")
}
