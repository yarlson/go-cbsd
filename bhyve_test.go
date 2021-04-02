package cbsd

import (
	"context"
	"errors"
	"github.com/yarlson/go-cbsd/v2/mock"
	"reflect"
	"testing"
)

func TestBHyveService_List(t *testing.T) {
	mockedExec := &mock.ExecMock{
		CommandFunc: func(ctx context.Context, name string, arg ...string) ([]byte, error) {
			output := []byte("build               45726  65536  12  linux    On   5910")
			return output, nil
		},
		SetEnvFunc: func(name string, value string) {
			return
		},
	}
	mockedExecErr := &mock.ExecMock{
		CommandFunc: func(ctx context.Context, name string, arg ...string) ([]byte, error) {
			return nil, errors.New("error")
		},
		SetEnvFunc: func(name string, value string) {
			return
		},
	}
	tests := []struct {
		name    string
		exec    Exec
		want    []BHyve
		wantErr bool
	}{
		{
			name: "BHyve List",
			exec: mockedExec,
			want: []BHyve{
				{
					JName:    "build",
					JID:      45726,
					VmRam:    65536,
					VmCPUs:   12,
					VmOSType: "linux",
					Status:   "On",
					VNC:      "5910",
				},
			},
		},
		{
			name:    "BHyve List Error",
			exec:    mockedExecErr,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BHyveService{
				exec: tt.exec,
			}
			got, err := b.List(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBHyveService_Stop(t *testing.T) {
	mockedExec := &mock.ExecMock{
		CommandFunc: func(ctx context.Context, name string, arg ...string) ([]byte, error) {
			output := []byte(`Send SIGTERM to test Soft timeout is 30 sec. 0 seconds left [..............................]
bstop done in 11 seconds`)
			return output, nil
		},
		SetEnvFunc: func(name string, value string) {
			return
		},
	}
	mockedExecErr := &mock.ExecMock{
		CommandFunc: func(ctx context.Context, name string, arg ...string) ([]byte, error) {
			return nil, errors.New("error")
		},
		SetEnvFunc: func(name string, value string) {
			return
		},
	}
	mockedExeErrNoDomain := &mock.ExecMock{
		CommandFunc: func(ctx context.Context, name string, arg ...string) ([]byte, error) {
			output := []byte(`No such domain: no-domain`)
			return output, nil
		},
		SetEnvFunc: func(name string, value string) {
			return
		},
	}
	tests := []struct {
		name       string
		exec       Exec
		instanceId string
		wantErr    bool
	}{
		{
			name:       "Stop: Done",
			exec:       mockedExec,
			instanceId: "domain",
		},
		{
			name:       "Stop: Error",
			exec:       mockedExecErr,
			instanceId: "error",
			wantErr:    true,
		},
		{
			name:       "Stop: No such domain",
			exec:       mockedExeErrNoDomain,
			instanceId: "no-domain",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BHyveService{
				exec: tt.exec,
			}
			err := b.Stop(context.Background(), tt.instanceId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestBHyveService_Start(t *testing.T) {
	mockedExec := &mock.ExecMock{
		CommandFunc: func(ctx context.Context, name string, arg ...string) ([]byte, error) {
			output := []byte(`VRDP is enabled. VNC bind/port: 127.0.0.1:5912
For attach VM console, use: vncviewer 127.0.0.1:5912
Resolution: 1024x768.
bhyve renice: 1
Waiting for PID.
PID: 25681`)
			return output, nil
		},
		SetEnvFunc: func(name string, value string) {
			return
		},
	}
	mockedExecErr := &mock.ExecMock{
		CommandFunc: func(ctx context.Context, name string, arg ...string) ([]byte, error) {
			return nil, errors.New("error")
		},
		SetEnvFunc: func(name string, value string) {
			return
		},
	}
	mockedExeErrNoDomain := &mock.ExecMock{
		CommandFunc: func(ctx context.Context, name string, arg ...string) ([]byte, error) {
			output := []byte(`No such domain: no-domain`)
			return output, nil
		},
		SetEnvFunc: func(name string, value string) {
			return
		},
	}
	tests := []struct {
		name       string
		exec       Exec
		instanceId string
		wantErr    bool
	}{
		{
			name:       "Start: Done",
			exec:       mockedExec,
			instanceId: "domain",
		},
		{
			name:       "Start: Error",
			exec:       mockedExecErr,
			instanceId: "error",
			wantErr:    true,
		},
		{
			name:       "Start: No such domain",
			exec:       mockedExeErrNoDomain,
			instanceId: "no-domain",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BHyveService{
				exec: tt.exec,
			}
			err := b.Start(context.Background(), tt.instanceId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestBHyveService_Remove(t *testing.T) {
	mockedExec := &mock.ExecMock{
		CommandFunc: func(ctx context.Context, name string, arg ...string) ([]byte, error) {
			output := []byte(`Send SIGTERM to test. Soft timeout is 30 sec. 0 seconds left [..............................]
bstop done in 13 seconds
destroy parent zvol for test: crdata/test/dsk1.vhd
bremove done in 17 seconds`)
			return output, nil
		},
		SetEnvFunc: func(name string, value string) {
			return
		},
	}
	mockedExecErr := &mock.ExecMock{
		CommandFunc: func(ctx context.Context, name string, arg ...string) ([]byte, error) {
			return nil, errors.New("error")
		},
		SetEnvFunc: func(name string, value string) {
			return
		},
	}
	mockedExeErrNoDomain := &mock.ExecMock{
		CommandFunc: func(ctx context.Context, name string, arg ...string) ([]byte, error) {
			output := []byte(`No such domain: no-domain`)
			return output, nil
		},
		SetEnvFunc: func(name string, value string) {
			return
		},
	}
	tests := []struct {
		name       string
		exec       Exec
		instanceId string
		wantErr    bool
	}{
		{
			name:       "Remove: Done",
			exec:       mockedExec,
			instanceId: "domain",
		},
		{
			name:       "Remove: Error",
			exec:       mockedExecErr,
			instanceId: "error",
			wantErr:    true,
		},
		{
			name:       "Remove: No such domain",
			exec:       mockedExeErrNoDomain,
			instanceId: "no-domain",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BHyveService{
				exec: tt.exec,
			}
			err := b.Remove(context.Background(), tt.instanceId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
