package edge

import "github.com/CarsonSlovoka/go-pkg/v2/w32"

type iUnknownImpl interface {
	// QueryInterface https://learn.microsoft.com/en-us/windows/win32/api/unknwn/nf-unknwn-iunknown-queryinterface(refiid_void)
	QueryInterface(rIID *w32.GUID, object uintptr) uintptr

	// AddRef https://learn.microsoft.com/en-us/windows/win32/api/unknwn/nf-unknwn-iunknown-addref
	// AddRef() int32 會遇到 expected function with one uintptr-sized result
	AddRef() uintptr

	// Release https://learn.microsoft.com/en-us/windows/win32/api/unknwn/nf-unknwn-iunknown-release
	Release() uintptr
}

type iUnknownVTbl struct {
	queryInterface uintptr
	addRef         uintptr
	release        uintptr
}
