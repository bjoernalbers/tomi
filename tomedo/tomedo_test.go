package tomedo

import (
	"net/url"
	"testing"
)

const serverURL = "http://tomedo.example.com:8080/tomedo_live/"

func ServerURL() *url.URL {
	u, err := url.Parse(serverURL)
	if err != nil {
		panic(err)
	}
	return u
}

func TestTomedoName(t *testing.T) {
	app := Tomedo{}
	if got, want := app.Name(), "tomedo.app"; got != want {
		t.Fatalf("%#v.Name()\ngot:\t%q\nwant:\t%q", app, got, want)
	}
}

func TestTomedoPath(t *testing.T) {
	app := Tomedo{Dir: "/Users/gopher/Applications"}
	want := "/Users/gopher/Applications/tomedo.app"
	if got := app.Path(); got != want {
		t.Fatalf("%T.Path()\ngot:\t%q\nwant:\t%q", app, got, want)
	}
}

func TestTomedoURL(t *testing.T) {
	got := tomedoURL(DefaultServer().URL())
	want := "http://allgemeinmedizin.demo.tomedo.org:8080/tomedo_live/filebyname/serverinternal/tomedo.app.tar"
	if got != want {
		t.Fatalf("tomedoURL():\ngot:\t%q\nwant:\t%q", got, want)
	}
}
