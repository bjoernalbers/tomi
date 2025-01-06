package app

import "testing"

func TestArzekoLatestVersionURL(t *testing.T) {
	arzeko := Arzeko{ServerURL: ServerURL()}
	want := "http://tomedo.example.com:8080/tomedo_live/arzeko/latestmac?version=0.0.0&aarch=intel"
	if got := arzeko.latestVersionURL(); got != want {
		t.Fatalf("%#v.latestVersionURL():\ngot:\t%q\nwant:\t%q", arzeko, got, want)
	}
}
