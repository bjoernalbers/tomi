package server

import (
	"net/url"
	"testing"
)

func TestDefault(t *testing.T) {
	s := Default()
	if got, want := s.Scheme, "http"; got != want {
		t.Fatalf("Default().Scheme:\ngot:\t%q\nwant:\t%q", got, want)
	}
	if got, want := s.Addr, "allgemeinmedizin.demo.tomedo.org"; got != want {
		t.Fatalf("Default().Addr:\ngot:\t%q\nwant:\t%q", got, want)
	}
	if got, want := s.Port, "8080"; got != want {
		t.Fatalf("Default().Port:\ngot:\t%q\nwant:\t%q", got, want)
	}
	if got, want := s.Path, "/tomedo_live/"; got != want {
		t.Fatalf("Default().Path:\ngot:\t%q\nwant:\t%q", got, want)
	}
}

func TestURL(t *testing.T) {
	s := Server{Scheme: "https", Addr: "example.com", Port: "8443", Path: "/foo/"}
	want, _ := url.Parse("https://example.com:8443/foo/")
	if got := s.URL(); *got != *want {
		t.Fatalf("URL():\ngot:\t%q\nwant:\t%q", got, want)
	}
}

func TestString(t *testing.T) {
	s := Server{Scheme: "https", Addr: "example.com", Port: "8443", Path: "/foo/"}
	want := "https://example.com:8443/foo/"
	if got := s.String(); got != want {
		t.Fatalf("String():\ngot:\t%q\nwant:\t%q", got, want)
	}
}
