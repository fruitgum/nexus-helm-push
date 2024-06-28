package server

import (
	"net/http"
	"net/url"
	"strconv"
)

func CheckScheme(nexusURL string) string {
	u, _ := url.Parse(nexusURL)

	scheme := u.Scheme
	return scheme
}

func Ping(nexusURL string) (string, error) {

	apiEndpoint := nexusURL + "/v1/status"

	resp, err := http.Post(apiEndpoint, "application/json", nil)
	if err != nil {
		return resp.Status, err
	}

	if resp.StatusCode != 200 {
		return "Could not connect to server, HTTP " + strconv.Itoa(resp.StatusCode), err
	}

	return "", nil

}
