package repository

import "regexp"

func IsDuplicateKeyError(err error) (bool, error) {
	return regexp.MatchString("Error 1062", err.Error())
}
