package cbsd

import (
	"context"
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
			got, err := s.Command(context.Background(), tt.args.name, tt.args.arg...)
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
	Name              string `json:"name,omitempty"`
	StringValue       string `json:"string_value,omitempty"`
	EmptyValue        string `json:"empty_value,omitempty"`
	NoTagValue        string
	Bool              bool  `json:"bool,omitempty"`
	BoolFieldPtr      *bool `json:"bool_field,omitempty"`
	EmptyBoolFieldPtr *bool `json:"empty_bool_field,omitempty"`
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
			name: "Test Struct Bool Ptr",
			args: args{b: &TestStruct{
				Name:         "test",
				StringValue:  "test-value",
				NoTagValue:   "no-tag",
				BoolFieldPtr: Bool(true),
			}},
			want: []string{
				"name=test",
				"string_value=test-value",
				"bool=0",
				"bool_field=1",
			},
		},
		{
			name: "Test Struct Bool Ptr False",
			args: args{b: &TestStruct{
				Name:         "test",
				StringValue:  "test-value",
				NoTagValue:   "no-tag",
				BoolFieldPtr: Bool(false),
			}},
			want: []string{
				"name=test",
				"string_value=test-value",
				"bool=0",
				"bool_field=0",
			},
		},
		{
			name: "Test Struct Bool",
			args: args{b: &TestStruct{
				Name:        "test",
				StringValue: "test-value",
				NoTagValue:  "no-tag",
				Bool:        true,
			}},
			want: []string{
				"name=test",
				"string_value=test-value",
				"bool=1",
			},
		},
		{
			name: "Test Struct Bool False",
			args: args{b: &TestStruct{
				Name:        "test",
				StringValue: "test-value",
				NoTagValue:  "no-tag",
				Bool:        false,
			}},
			want: []string{
				"name=test",
				"string_value=test-value",
				"bool=0",
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

func Test_quote(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "Empty",
			s:    "",
			want: "''",
		},
		{
			name: "Normal",
			s:    "1",
			want: "1",
		},
		{
			name: "Escape",
			s:    "1 2",
			want: "'1 2'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := quote(tt.s); got != tt.want {
				t.Errorf("quote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShellExec_CommandWithInterface(t *testing.T) {
	type args struct {
		name string
		i    *TestStruct
		arg  []string
	}
	tests := []struct {
		name    string
		args    args
		wantStr string
		wantErr bool
	}{
		{
			name: "CommandWithInterface",
			args: args{
				name: "ls",
				i: &TestStruct{
					Name:        "test",
					StringValue: "1",
				},
				arg: nil,
			},
			wantStr: "ls name=test string_value=1 bool=0",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ShellExec{}
			_, err := s.CommandWithInterface(context.Background(), tt.args.name, tt.args.i, tt.args.arg...)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommandWithInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if s.String() != tt.wantStr {
				t.Errorf("CommandWithInterface() got = %v, want %v", s.String(), tt.wantStr)
			}
		})
	}
}
