//go:build !production

package webapp

func DevMode() bool {
	return true
}
