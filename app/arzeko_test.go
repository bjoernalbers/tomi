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

func TestArzekoUserSettings(t *testing.T) {
	arzeko := Arzeko{ServerURL: ServerURL()}
	want := `{
  "protocol": "http",
  "authPassword": "tomedo",
  "xApiKey": "tomedo",
  "ipAddress": "tomedo.example.com",
  "oneWaySSLPort": 9444,
  "authPortServerTools": 9442,
  "authPort": 9443,
  "httpsPort": 8443,
  "httpPort": 8080,
  "basePath": "tomedo_live",
  "autoUpdatePath": "arzeko/latestmac?version=",
  "calendarOptions": {
    "weekNumbers": true,
    "dayMaxEvents": 5,
    "minTime": "06:00:00",
    "maxTime": "23:00:00",
    "printSorting": "start,title",
    "snapDuration": "00:15:00",
    "slotDuration": "00:30:00",
    "allDaySlot": true,
    "showNamesMonth": false,
    "eventLimitWeek": 5,
    "showNamesInsteadOfKuerzel": true,
    "selectable": true,
    "restorePosition": false,
    "showAbwesenMonth": false,
    "hiddenDays": [],
    "showNutzerSchichtzeiten": true,
    "showPopovers": true,
    "pauseDuration": 0
  },
  "exportOptions": {
    "abwesenExport": false,
    "eigeneSchichtenOnly": false
  },
  "timezone": 1,
  "disableUpdates": false,
  "logLevel": "info"
}`
	if got := arzeko.userSettings(); got != want {
		t.Fatalf("%T.userSettings():\ngot:\t%s\nwant:\t%s", arzeko, got, want)
	}
}

func TestArzekoPath(t *testing.T) {
	arzeko := Arzeko{InstallDir: "/foo/bar"}
	want := "/foo/bar/Arzeko.app"
	if got := arzeko.Path(); got != want {
		t.Fatalf("%T.Path()\ngot:\t%q\nwant:\t%q", arzeko, got, want)
	}
}
