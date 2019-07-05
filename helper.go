package gokendoparser

import (
	"regexp"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// BoolToString cast bool to string with extra format (usefull for inline func)
// BoolToString(true, "Yes", "No") will return "Yes"
func BoolToString(b bool, y string, n string) string {
	if y == "" {
		y = "Enable"
	}
	if n == "" {
		n = "Disable"
	}
	if b {
		return y
	}
	return n
}

// StringToBool cast string to bool with extra options default value. empty string will be default
// normally like js will use
// StringToBool("ya", false) will return true
// StringToBool("", false) will return false
// StringToBool("", true) will return true
func StringToBool(str string, def bool) bool {
	str = strings.ToLower(strings.TrimSpace(str))
	if !def {
		if str == "y" || str == "yes" || str == "true" || str == "1" || str == "ya" || str == "active" ||
			strings.Contains(str, "true") || strings.HasPrefix(str, "yes") {
			return true
		}
		return false
	}
	if str == "n" || str == "no" || str == "false" || str == "0" || str == "not" || str == "not active" || str == "inactive" ||
		strings.HasPrefix(str, "not") || strings.HasPrefix(str, "false") || strings.HasPrefix(str, "no") {
		return false
	}
	return true
}

// RegexCaseInsensitive Generate bson.RegEx for case insensitive
func RegexCaseInsensitive(value string) bson.RegEx {
	value = regexp.QuoteMeta(value)
	return bson.RegEx{Pattern: "^" + strings.ToLower(value) + "$", Options: "i"}
}

// RegexContains Generate bson.RegEx for contains
func RegexContains(value string, ignoreCase bool) bson.RegEx {
	value = regexp.QuoteMeta(value)
	if ignoreCase {
		return bson.RegEx{Pattern: "" + strings.ToLower(value) + "", Options: "i"}
	} else {
		return bson.RegEx{Pattern: "" + strings.ToLower(value) + "", Options: ""}
	}
}
