package utils

import "strconv"

// IntToString convert int to string.
func IntToString(i interface{}) string {
	return strconv.FormatInt(int64(i.(int)), 10)
}

// StringToInt convert string to int.
func StringToInt(i string) int {
	j, _ := strconv.Atoi(i)
	return j
}
