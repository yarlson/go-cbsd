package cbsd

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

var count int

type BHyveTestExec struct {
	env []string
}

func (s *BHyveTestExec) SetEnv(name, value string) {
	s.env = append(s.env, fmt.Sprintf("%s=%s", name, value))
}

func (s *BHyveTestExec) Command(name string, arg ...string) ([]byte, error) {
	output := ""
	if reflect.DeepEqual(arg, []string{"bls", "header=0", "display=jname,jid,vm_ram,vm_cpus,vm_os_type,status,vnc_port"}) {
		count += 1
		output = `a
build              45726  65536  12  linux    On   5910`
		if count == 2 {
			return nil, errors.New("error")
		}
	}

	if reflect.DeepEqual(arg, []string{"bstop", "inter=0", "jname=error"}) {
		return nil, errors.New("error")
	}

	if reflect.DeepEqual(arg, []string{"bstop", "inter=0", "jname=no-domain"}) {
		output = `No such domain: no-domain`
	}

	if reflect.DeepEqual(arg, []string{"bstop", "inter=0", "jname=domain"}) {
		output = `Send SIGTERM to test Soft timeout is 30 sec. 0 seconds left [..............................]
bstop done in 11 seconds`
	}

	if reflect.DeepEqual(arg, []string{"bstart", "inter=0", "jname=domain"}) {
		output = `VRDP is enabled. VNC bind/port: 127.0.0.1:5912
For attach VM console, use: vncviewer 127.0.0.1:5912
Resolution: 1024x768.
bhyve renice: 1
Waiting for PID.
PID: 25681`
	}

	if reflect.DeepEqual(arg, []string{"bstart", "inter=0", "jname=error"}) {
		return nil, errors.New("error")
	}

	if reflect.DeepEqual(arg, []string{"bstart", "inter=0", "jname=no-domain"}) {
		output = `No such domain: no-domain`
	}

	if reflect.DeepEqual(arg, []string{"bremove", "inter=0", "jname=no-domain"}) {
		output = `No such domain: no-domain
bremove done in 0 seconds`
	}

	if reflect.DeepEqual(arg, []string{"bremove", "inter=0", "jname=error"}) {
		return nil, errors.New("error")
	}

	if reflect.DeepEqual(arg, []string{"bremove", "inter=0", "jname=domain"}) {
		output = `Send SIGTERM to test. Soft timeout is 30 sec. 0 seconds left [..............................]
bstop done in 13 seconds
destroy parent zvol for test: crdata/test/dsk1.vhd
bremove done in 17 seconds`
	}

	return []byte(output), nil
}

func (s *BHyveTestExec) CommandWithInterface(name string, i interface{}, arg ...string) ([]byte, error) {
	return nil, nil
}

func TestBHyveService_List(t *testing.T) {
	tests := []struct {
		name    string
		exec    Exec
		want    []BHyve
		wantErr bool
	}{
		{
			name: "BHyve List",
			exec: &BHyveTestExec{},
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
			exec:    &BHyveTestExec{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BHyveService{
				exec: tt.exec,
			}
			got, err := b.List()
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
	tests := []struct {
		name       string
		exec       Exec
		instanceId string
		wantErr    bool
	}{
		{
			name:       "Stop: Done",
			exec:       &BHyveTestExec{},
			instanceId: "domain",
		},
		{
			name:       "Stop: Error",
			exec:       &BHyveTestExec{},
			instanceId: "error",
			wantErr:    true,
		},
		{
			name:       "Stop: No such domain",
			exec:       &BHyveTestExec{},
			instanceId: "no-domain",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BHyveService{
				exec: tt.exec,
			}
			err := b.Stop(tt.instanceId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestBHyveService_Start(t *testing.T) {
	tests := []struct {
		name       string
		exec       Exec
		instanceId string
		wantErr    bool
	}{
		{
			name:       "Start: Done",
			exec:       &BHyveTestExec{},
			instanceId: "domain",
		},
		{
			name:       "Start: Error",
			exec:       &BHyveTestExec{},
			instanceId: "error",
			wantErr:    true,
		},
		{
			name:       "Start: No such domain",
			exec:       &BHyveTestExec{},
			instanceId: "no-domain",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BHyveService{
				exec: tt.exec,
			}
			err := b.Start(tt.instanceId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestBHyveService_Remove(t *testing.T) {
	tests := []struct {
		name       string
		exec       Exec
		instanceId string
		wantErr    bool
	}{
		{
			name:       "Remove: Done",
			exec:       &BHyveTestExec{},
			instanceId: "domain",
		},
		{
			name:       "Remove: Error",
			exec:       &BHyveTestExec{},
			instanceId: "error",
			wantErr:    true,
		},
		{
			name:       "Remove: No such domain",
			exec:       &BHyveTestExec{},
			instanceId: "no-domain",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BHyveService{
				exec: tt.exec,
			}
			err := b.Remove(tt.instanceId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
