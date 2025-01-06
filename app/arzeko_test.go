package app

import "testing"

func TestArzekoLatestVersionURL(t *testing.T) {
	tests := []struct {
		arch string
		want string
	}{
		{
			"",
			"http://tomedo.example.com:8080/tomedo_live/arzeko/latestmac?version=0.0.0&aarch=intel",
		},
		{
			"intel",
			"http://tomedo.example.com:8080/tomedo_live/arzeko/latestmac?version=0.0.0&aarch=intel",
		},
		{
			"x86_64",
			"http://tomedo.example.com:8080/tomedo_live/arzeko/latestmac?version=0.0.0&aarch=intel",
		},
		{
			"amd64",
			"http://tomedo.example.com:8080/tomedo_live/arzeko/latestmac?version=0.0.0&aarch=intel",
		},
		{
			"arm",
			"http://tomedo.example.com:8080/tomedo_live/arzeko/latestmac?version=0.0.0&aarch=arm",
		},
		{
			"arm64",
			"http://tomedo.example.com:8080/tomedo_live/arzeko/latestmac?version=0.0.0&aarch=arm",
		},
	}
	for _, tt := range tests {
		arzeko := Arzeko{ServerURL: ServerURL(), Arch: tt.arch}
		if got := arzeko.latestVersionURL(); got != tt.want {
			t.Errorf("%#v.latestVersionURL():\ngot:\t%q\nwant:\t%q", arzeko, got, tt.want)
		}
	}
}
