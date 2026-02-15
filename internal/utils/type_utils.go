package utils

import "strconv"

func ToUint(input string) uint {
	_int, err := strconv.Atoi(input)
	if err != nil {
		return 0
	}
	return uint(_int)
}
