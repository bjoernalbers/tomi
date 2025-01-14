package tomedo

import (
	"net/url"
	"testing"
)

func TestDefaultServer(t *testing.T) {
	s := DefaultServer()
	if got, want := s.Scheme, "http"; got != want {
		t.Fatalf("DefaultServer().Scheme:\ngot:\t%q\nwant:\t%q", got, want)
	}
	if got, want := s.Addr, "allgemeinmedizin.demo.tomedo.org"; got != want {
		t.Fatalf("DefaultServer().Addr:\ngot:\t%q\nwant:\t%q", got, want)
	}
	if got, want := s.Port, "8080"; got != want {
		t.Fatalf("DefaultServer().Port:\ngot:\t%q\nwant:\t%q", got, want)
	}
	if got, want := s.Path, "/tomedo_live/"; got != want {
		t.Fatalf("DefaultServer().Path:\ngot:\t%q\nwant:\t%q", got, want)
	}
}

func TestServerURL(t *testing.T) {
	s := Server{Scheme: "https", Addr: "example.com", Port: "8443", Path: "/foo/"}
	want, _ := url.Parse("https://example.com:8443/foo/")
	if got := s.URL(); *got != *want {
		t.Fatalf("URL():\ngot:\t%q\nwant:\t%q", got, want)
	}
}

func TestServerString(t *testing.T) {
	s := Server{Scheme: "https", Addr: "example.com", Port: "8443", Path: "/foo/"}
	want := "https://example.com:8443/foo/"
	if got := s.String(); got != want {
		t.Fatalf("String():\ngot:\t%q\nwant:\t%q", got, want)
	}
}
