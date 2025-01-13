package app

import "testing"

func TestTomedoName(t *testing.T) {
	tomedo := Tomedo{}
	tomedo.App.Path = "/Users/gopher/Applications/tomedo.app"
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

func TestTomedoPath(t *testing.T) {
	tomedo := Tomedo{}
	tomedo.App.Path = "/foo/bar/tomedo.app"
	want := "/foo/bar/tomedo.app"
	if got := tomedo.Path(); got != want {
		t.Fatalf("%T.Path()\ngot:\t%q\nwant:\t%q", tomedo, got, want)
	}
}
