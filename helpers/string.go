package helpers

import "strconv"

func StringToNumber(s string) int {
	atoi, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return atoi
}
