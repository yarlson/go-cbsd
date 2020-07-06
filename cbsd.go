package cbsd

type CBSD struct {
	BHyve *BHyveService
	Jail  *JailService
	Xen   *XenService
}

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
