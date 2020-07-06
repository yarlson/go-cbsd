package cbsd

import (
	"fmt"
	"reflect"
	"testing"
)

type TestExec struct {
	env []string
}

func (s *TestExec) SetEnv(name, value string) {
	s.env = append(s.env, fmt.Sprintf("%s=%s", name, value))
}

func (s *TestExec) Command(name string, arg ...string) ([]byte, error) {
	output := ""
	if reflect.DeepEqual(arg, []string{"bls", "header=0", "display=jname,jid,vm_ram,vm_cpus,vm_os_type,status,vnc_port"}) {
		output = `build              45726  65536  12  linux    On   5910`
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

	if reflect.DeepEqual(arg, []string{"bstart", "inter=0", "jname=no-domain"}) {
		output = `No such domain: no-domain`
	}

	if reflect.DeepEqual(arg, []string{"bremove", "inter=0", "jname=no-domain"}) {
		output = `No such domain: no-domain
bremove done in 0 seconds`
	}

	if reflect.DeepEqual(arg, []string{"bremove", "inter=0", "jname=domain"}) {
		output = `Send SIGTERM to test. Soft timeout is 30 sec. 0 seconds left [..............................]
bstop done in 13 seconds
destroy parent zvol for test: crdata/test/dsk1.vhd
bremove done in 17 seconds`
	}

	return []byte(output), nil
}

func (s *TestExec) CommandWithInterface(name string, i interface{}, arg ...string) ([]byte, error) {
	return nil, nil
}

func TestBHyveService_List(t *testing.T) {
	type fields struct {
		exec Exec
	}
	tests := []struct {
		name    string
		fields  fields
		want    []BHyve
		wantErr bool
	}{
		{
			name: "BHyve List",
			fields: fields{
				exec: &TestExec{},
			},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BHyveService{
				exec: tt.fields.exec,
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

func TestBHyveService_Stop_No_Domain(t *testing.T) {
	type fields struct {
		exec Exec
	}
	type args struct {
		instanceId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Stop: No such domain",
			fields: fields{
				exec: &TestExec{},
			},
			args: args{instanceId: "no-domain"},
			want: "No such domain: no-domain",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BHyveService{
				exec: tt.fields.exec,
			}
			err := b.Stop(tt.args.instanceId)
			if err == nil {
				fmt.Println("Stop() error is empty!")
				return
			}
			if err.Error() != tt.want {
				t.Errorf("Stop() got = %v, want %v", err.Error(), tt.want)
			}
		})
	}
}

func TestBHyveService_Stop(t *testing.T) {
	type fields struct {
		exec Exec
	}
	type args struct {
		instanceId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Stop: Done",
			fields: fields{
				exec: &TestExec{},
			},
			args: args{instanceId: "domain"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BHyveService{
				exec: tt.fields.exec,
			}
			err := b.Stop(tt.args.instanceId)
			if err != nil {
				t.Errorf("Stop() got = %v, want nil", err.Error())
			}
		})
	}
}

func TestBHyveService_Start(t *testing.T) {
	type fields struct {
		exec Exec
	}
	type args struct {
		instanceId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Start: Done",
			fields: fields{
				exec: &TestExec{},
			},
			args: args{instanceId: "domain"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BHyveService{
				exec: tt.fields.exec,
			}
			err := b.Start(tt.args.instanceId)
			if err != nil {
				t.Errorf("Start() got = %v, want nil", err.Error())
			}
		})
	}
}

func TestBHyveService_Start_No_Domain(t *testing.T) {
	type fields struct {
		exec Exec
	}
	type args struct {
		instanceId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Start: No such domain",
			fields: fields{
				exec: &TestExec{},
			},
			args: args{instanceId: "no-domain"},
			want: "No such domain: no-domain",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BHyveService{
				exec: tt.fields.exec,
			}
			err := b.Start(tt.args.instanceId)
			if err == nil {
				fmt.Println("Start() error is empty!")
				return
			}
			if err.Error() != tt.want {
				t.Errorf("Start() got = %v, want %v", err.Error(), tt.want)
			}
		})
	}
}

func TestBHyveService_Remove_No_Domain(t *testing.T) {
	type fields struct {
		exec Exec
	}
	type args struct {
		instanceId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Remove: No such domain",
			fields: fields{
				exec: &TestExec{},
			},
			args: args{instanceId: "no-domain"},
			want: "No such domain: no-domain",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BHyveService{
				exec: tt.fields.exec,
			}
			err := b.Remove(tt.args.instanceId)
			if err == nil {
				fmt.Println("Remove() error is empty!")
				return
			}
			if err.Error() != tt.want {
				t.Errorf("Remove() got = %v, want %v", err.Error(), tt.want)
			}
		})
	}
}

func TestBHyveService_Remove(t *testing.T) {
	type fields struct {
		exec Exec
	}
	type args struct {
		instanceId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Remove: Done",
			fields: fields{
				exec: &TestExec{},
			},
			args: args{instanceId: "domain"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BHyveService{
				exec: tt.fields.exec,
			}
			err := b.Remove(tt.args.instanceId)
			if err != nil {
				t.Errorf("Remove() got = %v, want nil", err.Error())
			}
		})
	}
}
