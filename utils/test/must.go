package test

func Must[T any](t T, err error) T {
	CheckErr(err)
	return t
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
