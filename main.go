package main

import (
	"flag"
	"fmt"
	"github.com/fruitgum/nexus-helm-push/archive"
	"github.com/fruitgum/nexus-helm-push/logging"
	"github.com/fruitgum/nexus-helm-push/server"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

func main() {

	username := flag.String("username", "", "Nexus account username. $NEXUS_HELM_USER")
	password := flag.String("password", "", "Nexus account password. $NEXUS_HELM_PASSWORD")
	repo := flag.String("repo", "", "Nexus helm repo name. $NEXUS_HELM_REPO")
	url := flag.String("url", "", "Nexus helm repo URL. $NEXUS_HELM_URL")
	debug := flag.Bool("debug", false, "Debug mode. $NEXUS_HELM_DEBUG")
	chart := flag.String("chart", "", "/path/to/chart/archive.tgz")
	maskedPassword := strings.Repeat("*", utf8.RuneCountInString(*password))
	flag.Parse()

	if *username == "" {
		if os.Getenv("NEXUS_HELM_USER") == "" {
			logging.Error("No Username provided")
			return
		} else {
			*username = os.Getenv("NEXUS_HELM_USER")
		}
	}

	if *password == "" {
		if os.Getenv("NEXUS_HELM_PASSWORD") == "" {
			logging.Error("No Password provided")
			return
		} else {
			*password = os.Getenv("NEXUS_HELM_PASSWORD")
		}
	}

	if *repo == "" {
		if os.Getenv("NEXUS_HELM_REPO") == "" {
			logging.Error("No repo name provided")
			return
		} else {
			*repo = os.Getenv("NEXUS_HELM_REPO")
		}
	}

	if *url == "" {
		if os.Getenv("NEXUS_HELM_URL") == "" {
			logging.Error("No url provided")
			return
		} else {
			*url = os.Getenv("NEXUS_HELM_URL")
		}
	}

	if *chart == "" {
		logging.Error("No file provided")
		return
	}

	absFilePath, _ := filepath.Abs(*chart)

	if *debug || os.Getenv("NEXUS_HELM_DEBUG") == "true" {
		logging.Debug("Username: %s", *username)
		logging.Debug("Password: %s", maskedPassword)
		logging.Debug("URL: %s", *url)
		logging.Debug("Repo: %s", *repo)
		logging.Debug("Chart: %s", *chart)
		logging.Debug("Chart absolute path: %s", absFilePath)
	}

	chartName := archive.BaseName(*chart)

	logging.Info("Checking connection")

	checkConnection, err := server.Ping(*url)
	if err != nil {
		logging.Error(checkConnection)
		return
	}

	scheme := server.CheckScheme(*url)

	if scheme != "https" {
		logging.Warn("Using insecure connection")
	}

	logging.Info("PING: OK")
	fileCheck, err := archive.CheckFileExists(*chart)

	if err != nil {
		logging.Error(fileCheck)
		return
	}

	logging.Info("Uploading " + *chartName + " to " + *repo + "repo")
	upload, success := server.Upload(*username, *password, *url, *repo, *chart)

	if success != true {
		logging.Error(upload)
		return
	}

	successMessage := *chartName + " uploaded to " + *url

	logging.Success(successMessage)

	asterisks := strings.Repeat("*", 30)

	links := "helm repo add nexus " + *url + "/repository/" + *repo
	links += "\n" + "helm install " + *chartName + " " + *repo + "/" + *chartName

	linkMessage := asterisks + "\n" + links + "\n" + asterisks

	fmt.Println(linkMessage)
}
