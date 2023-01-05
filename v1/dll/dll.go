//go:build windows

package dll

import "github.com/CarsonSlovoka/go-pkg/v2/w32"

var (
	User   *w32.User32DLL
	Kernel *w32.Kernel32DLL
	Ole    *w32.Ole32DLL
)

func init() {
	User = w32.NewUser32DLL()
	Kernel = w32.NewKernel32DLL()
	Ole = w32.NewOle32DLL()
}
