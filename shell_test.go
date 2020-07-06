package cbsd

import (
	"reflect"
	"strings"
	"testing"
)

func TestShellExec_Command(t *testing.T) {
	type fields struct {
		name  string
		value string
	}
	type args struct {
		name string
		arg  []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ls",
			fields: fields{
				name:  "NOCOLOR",
				value: "1",
			},
			args: args{
				name: "ls",
				arg:  []string{"-alh"},
			},
			want:    "cbsd.go",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ShellExec{}
			s.SetEnv(tt.fields.name, tt.fields.value)
			got, err := s.Command(tt.args.name, tt.args.arg...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Command() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(string(got), tt.want) {
				t.Errorf("Command() got = %v, want %v", string(got), tt.want)
			}
		})
	}
}

type TestStruct struct {
	Name        string `json:"name,omitempty"`
	StringValue string `json:"string_value,omitempty"`
	EmptyValue  string `json:"empty_value,omitempty"`
	NoTagValue  string
}

func Test_structToSlice(t *testing.T) {
	type args struct {
		b *TestStruct
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test Struct",
			args: args{b: &TestStruct{
				Name:        "test",
				StringValue: "test-value",
				NoTagValue:  "no-tag",
			}},
			want: []string{
				"name=test",
				"string_value=test-value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := structToSlice(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("structToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
