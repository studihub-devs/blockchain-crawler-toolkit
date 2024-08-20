package utils

import "strings"

func ContainStrings(mainString string, substrings ...string) bool {
	var (
		matches         = 0
		isCompleteMatch = true
	)

	for _, substring := range substrings {
		if strings.Contains(mainString, substring) {
			matches += 1
		} else {
			isCompleteMatch = false
		}
	}

	return isCompleteMatch
}

func OneOfStrings(mainString string, otherString ...string) bool {
	for _, str := range otherString {
		if strings.EqualFold(mainString, str) {
			return true
		}
	}

	return false
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func RemoveDuplicates(strList []string) []string {
	list := []string{}
	for _, item := range strList {
		if !Contains(list, item) {
			list = append(list, item)
		}
	}
	return list
}
