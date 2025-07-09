package utils

import "strings"

func IsInListString(value string, list []string) bool {
	for _, item := range list {
		if strings.EqualFold(item, value) { 
			return true
		}
	}
	return false
}