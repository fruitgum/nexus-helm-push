package server

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	//"path/filepath"
)

func Upload(username string, password string, url string, repo string, absFilePath string) (string, bool) {

	uploadLink := url + "/service/rest/v1/components?repository=" + repo

	file, _ := os.Open(absFilePath)

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	formFile, err := writer.CreateFormFile("r.asset", absFilePath)

	if err != nil {
		return "Create form file: " + err.Error(), false
	}

	writer.WriteField("type", "application/x-compressed")

	_, err = io.Copy(formFile, file)
	if err != nil {
		return "Copy file: " + err.Error(), false
	}

	writer.Close()

	req, err := http.NewRequest("POST", uploadLink, &requestBody)
	if err != nil {
		return req.Response.Status, false
	}
	req.SetBasicAuth(username, password)

	req.Header = map[string][]string{
		"Content-Type": {"multipart/form-data"},
		"Accept":       {"application/json"},
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err.Error(), false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return "Insufficient permissions to upload a component", false
	}

	if resp.StatusCode != 204 {
		return resp.Status, false
	}

	return "", true

}
