package app

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/bjoernalbers/tomi/server"
)

func TestArzekoName(t *testing.T) {
	arzeko := Arzeko{}
	arzeko.App.Path = "/Users/gopher/Applications/Arzeko.app"
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
	arzeko := Arzeko{}
	arzeko.App.Path = "/foo/bar/Arzeko.app"
	want := "/foo/bar/Arzeko.app"
	if got := arzeko.Path(); got != want {
		t.Fatalf("%T.Path()\ngot:\t%q\nwant:\t%q", arzeko, got, want)
	}
}

func TestArzekoDownloaderAutoUpdateURL(t *testing.T) {
	tests := []struct {
		arch string
		want string
	}{
		{
			"",
			"http://allgemeinmedizin.demo.tomedo.org:8080/tomedo_live/arzeko/latestmac?version=0.0.0&aarch=intel",
		},
		{
			"intel",
			"http://allgemeinmedizin.demo.tomedo.org:8080/tomedo_live/arzeko/latestmac?version=0.0.0&aarch=intel",
		},
		{
			"x86_64",
			"http://allgemeinmedizin.demo.tomedo.org:8080/tomedo_live/arzeko/latestmac?version=0.0.0&aarch=intel",
		},
		{
			"amd64",
			"http://allgemeinmedizin.demo.tomedo.org:8080/tomedo_live/arzeko/latestmac?version=0.0.0&aarch=intel",
		},
		{
			"arm",
			"http://allgemeinmedizin.demo.tomedo.org:8080/tomedo_live/arzeko/latestmac?version=0.0.0&aarch=arm",
		},
		{
			"arm64",
			"http://allgemeinmedizin.demo.tomedo.org:8080/tomedo_live/arzeko/latestmac?version=0.0.0&aarch=arm",
		},
	}
	for _, tt := range tests {
		d := arzekoDownloader{ServerURL: server.Default().URL(), Arch: tt.arch}
		if got := d.AutoUpdateURL(); got != tt.want {
			t.Errorf("%v.AutoUpdateURL():\ngot:\t%q\nwant:\t%q", d, got, tt.want)
		}
	}
}

func TestArzekoDownloaderURL(t *testing.T) {
	d := arzekoDownloader{ServerURL: server.Default().URL()}
	want := regexp.MustCompile(`http://allgemeinmedizin.demo.tomedo.org:8080/tomedo_live/filebyname/serverinternalzip/arzeko/Arzeko-\d+.\d+.\d+-mac.zip`)
	got, err := d.URL()
	if err != nil {
		t.Fatalf("%T.URL(): error = %v", d, err)
	}
	if !want.MatchString(got) {
		t.Fatalf("%T.URL():\ngot:\t%q\nwant:\t%q", d, got, want)
	}
}

func TestArzekoDownloaderGet(t *testing.T) {
	d := arzekoDownloader{ServerURL: server.Default().URL()}
	resp, err := d.Get()
	if err != nil {
		t.Fatalf("%T.Get(): error = %v", d, err)
	}
	defer resp.Body.Close()
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("%T.Get(): HTTP Status = %d, want: %d", d, got, want)
	}
}
