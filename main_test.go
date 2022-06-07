package main

import (
	"bytes"
	"reflect"
	"testing"
)

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

func Test_solLicenseInput_String(t *testing.T) {
	type fields struct {
		lib        string
		version    string
		license    string
		licenseUrl string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "Happy Case", fields: fields{lib: "lib", version: "v1.0.0", license: "license", licenseUrl: "myURL/LICENSE"}, want: ";lib;v1.0.0;myURL/LICENSE;license"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sli := &solLicenseInput{
				lib:        tt.fields.lib,
				version:    tt.fields.version,
				license:    tt.fields.license,
				licenseUrl: tt.fields.licenseUrl,
			}
			if got := sli.String(); got != tt.want {
				t.Errorf("solLicenseInput.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGoLicenseOutput(t *testing.T) {
	type args struct {
		raw string
	}
	tests := []struct {
		name string
		args args
		want goLicenseOutput
	}{

		{
			name: "Happy Case",
			args: args{raw: "github.com/gin-gonic/gin,https://github.com/gin-gonic/gin/blob/v1.7.4/LICENSE,MIT"},
			want: goLicenseOutput{lib: "github.com/gin-gonic/gin",
				licenseUrl: "https://github.com/gin-gonic/gin/blob/v1.7.4/LICENSE",
				license:    "MIT"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGoLicenseOutput(tt.args.raw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGoLicenseOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteFromSliSlice(t *testing.T) {
	type args struct {
		slis []solLicenseInput
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantOut string
	}{
		{
			name:    "HappyCase",
			args:    args{slis: []solLicenseInput{{lib: "lib", version: "v1.0.0", license: "MIT", licenseUrl: "URL/LICENSE"}, {lib: "lib", version: "v1.0.0", license: "MIT", licenseUrl: "URL/LICENSE"}}},
			want:    56,
			wantOut: ";lib;v1.0.0;URL/LICENSE;MIT\n;lib;v1.0.0;URL/LICENSE;MIT\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			if got := WriteFromSliSlice(tt.args.slis, out); got != tt.want {
				t.Errorf("WriteFromSliSlice() = %v, want %v", got, tt.want)
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("WriteFromSliSlice() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestNewSolLicenseInput(t *testing.T) {
	type args struct {
		glo goLicenseOutput
	}
	tests := []struct {
		name string
		args args
		want solLicenseInput
	}{
		{
			name: "Happy Case 1",
			args: args{
				glo: goLicenseOutput{lib: "lib", licenseUrl: "URL/v1.0.0/LICENSE", license: "MIT"},
			},
			want: solLicenseInput{lib: "lib", version: "v1.0.0", license: "MIT", licenseUrl: "URL/v1.0.0/LICENSE"},
		},
		{
			name: "Git commit sha version",
			args: args{
				glo: goLicenseOutput{lib: "lib", licenseUrl: "URL/665e8c73:LICENSE", license: "MIT"},
			},
			want: solLicenseInput{lib: "lib", version: "665e8c73", license: "MIT", licenseUrl: "URL/665e8c73:LICENSE"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSolLicenseInput(tt.args.glo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSolLicenseInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
