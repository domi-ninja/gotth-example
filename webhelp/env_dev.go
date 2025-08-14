//go:build !production

package webhelp

func DevMode() bool {
	return true
}
