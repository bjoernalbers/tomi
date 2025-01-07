package app

import "testing"

func TestTomedoName(t *testing.T) {
	tomedo := Tomedo{}
	got := tomedo.Name()
	want := "tomedo.app"
	if got != want {
		t.Fatalf("%#v.Name()\ngot:\t%q\nwant:\t%q", tomedo, got, want)
	}
}

func TestTomedoDownloadURL(t *testing.T) {
	tomedo := Tomedo{ServerURL: ServerURL()}
	want := "http://tomedo.example.com:8080/tomedo_live/filebyname/serverinternal/tomedo.app.tar"
	got, err := tomedo.DownloadURL()
	if err != nil {
		t.Fatalf("%#v.DownloadURL() error: %v", tomedo, err)
	}
	if got != want {
		t.Fatalf("%#v.DownloadURL():\ngot:\t%q\nwant:\t%q", tomedo, got, want)
	}
}