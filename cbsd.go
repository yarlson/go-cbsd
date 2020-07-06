package cbsd

// CBSD is a service container
type CBSD struct {
	BHyve *BHyveService
	Jail  *JailService
	Xen   *XenService
}

// NewCBSD creates a new CBSD
func NewCBSD() *CBSD {
	shellExec := &ShellExec{}
	return &CBSD{
		BHyve: &BHyveService{
			exec: shellExec,
		},
		Jail: &JailService{
			exec: shellExec,
		},
		Xen: &XenService{
			exec: shellExec,
		},
	}
}
