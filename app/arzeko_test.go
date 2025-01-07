package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestArzekoName(t *testing.T) {
	arzeko := Arzeko{}
	got := arzeko.Name()
	want := "Arzeko.app"
	if got != want {
		t.Fatalf("%#v.Name()\ngot:\t%q\nwant:\t%q", arzeko, got, want)
	}
}

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

func TestArzekoDownloadURL(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		resp    string
		want    string
		wantErr bool
	}{
		{
			"not found",
			404,
			"",
			"",
			true,
		},
		{
			"empty response",
			200,
			"",
			"",
			true,
		},
		{
			"valid response",
			200,
			`{"url":"http://tomedo.example.com/Arzeko-1.2.3.zip","name":"1.2.3","notes":"Aarch: null","pub_date":null}`,
			"http://tomedo.example.com/Arzeko-1.2.3.zip",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				fmt.Fprintf(w, tt.resp)
			}))
			defer server.Close()
			u, err := url.Parse(server.URL)
			if err != nil {
				panic(err)
			}
			arzeko := Arzeko{ServerURL: u}
			got, err := arzeko.DownloadURL()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Arzeko.DownloadURL() error: %v, wantErr: %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("Arzeko.DownloadURL()\ngot:\t%q\nwant:\t%q", got, tt.want)
			}
		})
	}
}
