package main

import "testing"

func Test_extractVersionFromLicenseUrl(t *testing.T) {
	type args struct {
		lurl string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "happy case 1", args: args{lurl: "https://github.com/tinylib/msgp/blob/v1.1.6/LICENSE"}, want: "v1.1.6"},

		{name: "happy case 1", args: args{lurl: "https://cs.opensource.google/go/x/sys/+/665e8c73:LICENSE"}, want: "665e8c73"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractVersionFromLicenseUrl(tt.args.lurl); got != tt.want {
				t.Errorf("extractVersionFromLicenseUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
