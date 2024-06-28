package server

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
)

type Writer struct {
	*multipart.Writer
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func (w *Writer) CreateFormFile(fieldname, filename, contentType string) (io.Writer, error) {

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", contentType)
	return w.CreatePart(h)
}

func Upload(username string, password string, url string, repo string, chart string) (string, bool) {

	uploadLink := url + "/service/rest/v1/components?repository=" + repo

	archive, _ := os.Open(chart)

	var requestBody bytes.Buffer

	writer := &Writer{multipart.NewWriter(&requestBody)}

	formFile, err := writer.CreateFormFile("r.asset", archive.Name(), "application/x-compressed")

	if err != nil {
		return "Create form file: " + err.Error(), false
	}

	_, err = io.Copy(formFile, archive)
	if err != nil {
		return "Copy file: " + err.Error(), false
	}

	wrc := writer.Close()
	if wrc != nil {
		return wrc.Error(), false
	}
	req, err := http.NewRequest("POST", uploadLink, &requestBody)
	if err != nil {
		return req.Response.Status, false
	}
	req.SetBasicAuth(username, password)

	req.Header.Set("Content-Type", "multipart/form-data; boundary="+writer.Boundary())
	req.Header.Set("Accept", "application/json")

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
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return resp.Status, false
		}
		return string(respBody), false
	}

	return "", true

}
