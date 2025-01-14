package app

import (
	"net/http"
	"testing"

	"github.com/bjoernalbers/tomi/server"
)

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

func TestTomedoDownloaderURL(t *testing.T) {
	d := tomedoDownloader{server.Default().URL()}
	want := "http://allgemeinmedizin.demo.tomedo.org:8080/tomedo_live/filebyname/serverinternal/tomedo.app.tar"
	if got := d.URL(); got != want {
		t.Fatalf("%T.URL():\ngot:\t%q\nwant:\t%q", d, got, want)
	}
}

func TestTomedoDownloaderGet(t *testing.T) {
	d := tomedoDownloader{server.Default().URL()}
	resp, err := d.Get()
	if err != nil {
		t.Fatalf("%T.Get(): error = %v", d, err)
	}
	defer resp.Body.Close()
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("%T.Get(): HTTP Status = %d, want: %d", d, got, want)
	}
}
