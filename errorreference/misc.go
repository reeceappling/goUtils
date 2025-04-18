package errorreference

import "errors"

func FirstMatchingError(errToCheck error, errs ...error) error {
	for _, e := range errs {
		if errors.Is(errToCheck, e) {
			return e
		}
	}
	return nil
}

func ErrIsOneOf(errToCheck error, errs ...error) bool {
	return FirstMatchingError(errToCheck, errs...) != nil
}
