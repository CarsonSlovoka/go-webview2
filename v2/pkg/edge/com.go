package edge

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
)

type iUnknownImpl interface {
	// QueryInterface https://learn.microsoft.com/en-us/windows/win32/api/unknwn/nf-unknwn-iunknown-queryinterface(refiid_void)
	QueryInterface(rIID *w32.GUID, object uintptr) w32.HRESULT

	// AddRef https://learn.microsoft.com/en-us/windows/win32/api/unknwn/nf-unknwn-iunknown-addref
	AddRef() int32

	// Release https://learn.microsoft.com/en-us/windows/win32/api/unknwn/nf-unknwn-iunknown-release
	Release() uint32
}

type iUnknownVTbl struct {
	queryInterface uintptr
	addRef         uintptr
	release        uintptr
}
