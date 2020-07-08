package cbsd

import (
	"reflect"
	"testing"
)

func TestNewCBSD(t *testing.T) {
	shellExec := &ShellExec{}
	tests := []struct {
		name string
		want *CBSD
	}{
		{
			name: "CBSD",
			want: &CBSD{
				BHyve: &BHyveService{
					exec: shellExec,
				},
				Jail: &JailService{
					exec: shellExec,
				},
				Xen: &XenService{
					exec: shellExec,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCBSD(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCBSD() = %v, want %v", got, tt.want)
			}
		})
	}
}
