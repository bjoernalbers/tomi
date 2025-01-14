package tomedo

import (
	"regexp"
	"testing"
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

func TestArzekoURLAutoUpdate(t *testing.T) {
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
		u := arzekoURL{ServerURL: DefaultServer().URL(), Arch: tt.arch}
		if got := u.autoUpdate(); got != tt.want {
			t.Errorf("%v.autoUpdate():\ngot:\t%q\nwant:\t%q", u, got, tt.want)
		}
	}
}

func TestArzekoURLString(t *testing.T) {
	u := arzekoURL{ServerURL: DefaultServer().URL()}
	want := regexp.MustCompile(`http://allgemeinmedizin.demo.tomedo.org:8080/tomedo_live/filebyname/serverinternalzip/arzeko/Arzeko-\d+.\d+.\d+-mac.zip`)
	got, err := u.String()
	if err != nil {
		t.Fatalf("%T.String(): error = %v", u, err)
	}
	if !want.MatchString(got) {
		t.Fatalf("%T.String():\ngot:\t%q\nwant:\t%q", u, got, want)
	}
}
