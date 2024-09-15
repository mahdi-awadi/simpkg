package format

import "testing"

func TestFormat(t *testing.T) {
	type args struct {
		format string
		i      []any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"TestFormat", args{"%v %v %v", []any{"test", "1", "2"}}, "test 1 2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Format(tt.args.format, tt.args.i...); got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	type args struct {
		i any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"TestString", args{[]any{"test", "1", "2"}}, "[test 1 2]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := String(tt.args.i); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
