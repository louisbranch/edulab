package html

import "testing"

func Test_marshal(t *testing.T) {
	type args struct {
		vals any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "empty", args: args{vals: []int{}}, want: "[]"},
		{name: "integers", args: args{vals: []int{1, 2, 3}}, want: "[1,2,3]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := marshal(tt.args.vals); string(got) != tt.want {
				t.Errorf("marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
