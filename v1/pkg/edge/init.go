package edge

import (
	"github.com/CarsonSlovoka/go-pkg/v2/w32"
	"github.com/CarsonSlovoka/go-webview2/v1/dll"
	"log"
	"os"
)

func init() {
	r := dll.Ole.CoInitializeEx(0, w32.COINIT_APARTMENTTHREADED)
	if int(r) < 0 {
		log.Printf("Warning: CoInitializeEx call failed: E=%v", r)
		os.Exit(int(r))
	}
}
