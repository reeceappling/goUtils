package this

import (
	"os"
	"path"
	"runtime"
)

// ThisDir returns a path with a trailing slash
func Dir() string {
	_, f, _, _ := runtime.Caller(1)
	return path.Dir(f) + string(os.PathSeparator)
}

func File() string {
	_, f, _, _ := runtime.Caller(1)
	return f
}
