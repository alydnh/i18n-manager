package utils

import (
	"strings"
)

const EmptyString = ""

func EmptyOrWhiteSpace(text interface{}) bool {
	switch text.(type) {
	case string:
		return strings.TrimSpace(text.(string)) == EmptyString
	case *string:
		value := text.(*string)
		return nil == value || EmptyOrWhiteSpace(*value)
	}

	return nil == text
}

func Compare(left, right *string, ignoreCase bool) int {
	if left == right {
		return 0
	}

	if nil == left && nil != right {
		return -1
	}

	if nil != left && nil == right {
		return 1
	}

	if ignoreCase {
		return strings.Compare(strings.ToUpper(*left), strings.ToUpper(*right))
	}

	return strings.Compare(*left, *right)
}

func Ptr2Str(text *string) (string, bool) {
	if nil == text {
		return EmptyString, false
	}

	return *text, true
}

func StrPtr2Ptr(getter func() (string, bool)) *string {
	if value, ok := getter(); ok {
		return &value
	}

	return nil

}
