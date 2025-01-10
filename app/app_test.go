package app

import (
	"net/url"
)

const serverURL = "http://tomedo.example.com:8080/tomedo_live/"

func ServerURL() *url.URL {
	u, err := url.Parse(serverURL)
	if err != nil {
		panic(err)
	}
	return u
}
