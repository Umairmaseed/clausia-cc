package utils

import (
	"strings"

	"github.com/mozillazg/go-unidecode"
)

func ValidateAndCleanData(name string) string {
	// Remove spaces
	name = strings.Replace(name, " ", "", -1)

	// Uncapitalize
	name = strings.ToLower(name[:1]) + name[1:]

	// Remove accentuation
	name = unidecode.Unidecode(name)

	return name
}
