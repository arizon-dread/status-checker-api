package businesslayer

import "testing"

func Test_getContentType(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "json",
			args: args{"{\"json\": \"is true\"}"},
			want: "application/json"},
		{
			name: "xml",
			args: args{"<hello></hello>"},
			want: "text/xml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getContentType(tt.args.s); got != tt.want {
				t.Errorf("getContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}
