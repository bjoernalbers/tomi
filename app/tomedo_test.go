package app

import "testing"

func TestTomedoDownloadURL(t *testing.T) {
	tomedo := Tomedo{ServerURL: ServerURL()}
	want := "http://tomedo.example.com:8080/tomedo_live/filebyname/serverinternal/tomedo.app.tar"
	if got := tomedo.DownloadURL(); got != want {
		t.Fatalf("%#v.DownloadURL():\ngot:\t%q\nwant:\t%q", tomedo, got, want)
	}
}
