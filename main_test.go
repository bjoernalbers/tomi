package main

import "testing"

func TestTomedoDownloadURL(t *testing.T) {
	got := tomedoDownloadURL("tomedo.example.com")
	want := "http://tomedo.example.com:8080/tomedo_live/filebyname/serverinternal/tomedo.app.tar"
	if got != want {
		t.Errorf("tomedoDownloadURL()\ngot:\t%s\nwant:\t%s\n", got, want)
	}
}
