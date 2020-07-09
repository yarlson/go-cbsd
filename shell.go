package cbsd

import (
	"fmt"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
)

var pattern *regexp.Regexp

type Exec interface {
	SetEnv(name, value string)
	Command(name string, arg ...string) ([]byte, error)
	CommandWithInterface(name string, i interface{}, arg ...string) ([]byte, error)
}

type ShellExec struct {
	env         []string
	commandLine string
}

func (s *ShellExec) SetEnv(name, value string) {
	s.env = append(s.env, fmt.Sprintf("%s=%s", name, value))
}

func (s *ShellExec) Command(name string, arg ...string) ([]byte, error) {
	s.commandLine = name + " " + strings.Join(arg, " ")
	cmd := exec.Command(name, arg...)
	cmd.Env = append(cmd.Env, s.env...)

	return cmd.Output()
}

func (s *ShellExec) String() string {
	return s.commandLine
}

func (s *ShellExec) CommandWithInterface(name string, i interface{}, arg ...string) ([]byte, error) {
	arg = append(arg, structToSlice(i)...)

	return s.Command(name, arg...)
}

func init() {
	pattern = regexp.MustCompile(`[^\w@%+=:,./-]`)
}

func quote(s string) string {
	if len(s) == 0 {
		return "''"
	}
	if pattern.MatchString(s) {
		return "'" + strings.Replace(s, "'", "'\"'\"'", -1) + "'"
	}

	return s
}

func structToSlice(b interface{}) []string {
	iVal := reflect.ValueOf(b).Elem()
	typ := iVal.Type()

	var slice []string
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		f.Type()
		tag := typ.Field(i).Tag.Get("json")
		if tag == "" {
			continue
		}

		var value interface{}
		switch f.Kind() {
		case reflect.Ptr:
			{
				if f.IsNil() {
					continue
				}
				value = 0
				if f.Elem().Bool() {
					value = 1
				}
			}
		case reflect.String:
			{
				if f.String() == "" {
					continue
				}
				value = quote(f.String())
			}
		case reflect.Bool:
			{
				value = 0
				if f.Bool() {
					value = 1
				}
			}
		}

		fields := strings.Split(tag, ",")

		slice = append(slice, fmt.Sprintf("%s=%v", fields[0], value))
	}

	return slice
}
