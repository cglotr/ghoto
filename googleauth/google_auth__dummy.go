package googleauth

import "net/http"

type Google_auth__dummy struct {
}

func Google_auth__dummy__new() *Google_auth__dummy {
	return &Google_auth__dummy{}
}

func (g *Google_auth__dummy) Get_client() (*http.Client, error) {
	return &http.Client{}, nil
}
