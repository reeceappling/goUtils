//go:build !linux

package local

func IsLocal() bool {
	return true
}
