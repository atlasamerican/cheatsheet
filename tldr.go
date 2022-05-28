package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	TldrRemoteUrl   = "https://github.com/tldr-pages/tldr-pages.github.io/"
	TldrRemotePath  = "raw/master/assets/tldr.zip"
	GithubStatusUrl = "https://www.githubstatus.com/api/v2/status.json"
)

var osMap = map[string]string{
	"linux":   "linux",
	"darwin":  "osx",
	"windows": "windows",
}

type TldrArchive struct {
	remoteUrl  string
	remotePath string
	statusUrl  string
	path       string
	zipPath    string
	revPath    string
	lang       string
	updating   chan bool
}

type TldrPage struct {
	name    string
	content string
}

func newTldrArchive(path string) *TldrArchive {
	a := &TldrArchive{
		remoteUrl:  TldrRemoteUrl,
		remotePath: TldrRemotePath,
		statusUrl:  GithubStatusUrl,
		path:       path,
		zipPath:    filepath.Join(path, "tldr.zip"),
		revPath:    filepath.Join(path, "rev"),
		lang:       "en",
		updating:   make(chan bool, 1),
	}

	a.updating <- true

	go func() {
		if a.checkUpdate() {
			if ok, err := a.update(); err != nil {
				if !ok {
					log.Fatal(err)
				}
				logger.Log("[error] %v", err)
			}
		}
		a.updating <- false
	}()

	return a
}

func (a *TldrArchive) waitForUpdate() {
	select {
	case state := <-a.updating:
		if state {
			<-a.updating
		}
	default:
		return
	}
}

func (a *TldrArchive) getRemoteRev() string {
	out, err := exec.Command("git", "ls-remote", a.remoteUrl, "HEAD").Output()
	if err != nil {
		log.Fatal(err)
	}
	var rev string
	_, err = fmt.Sscanf(string(out), "%s HEAD", &rev)
	if err != nil {
		log.Fatal(err)
	}
	return rev
}

func (a *TldrArchive) getRev() (string, error) {
	buf, err := ioutil.ReadFile(a.revPath)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (a *TldrArchive) checkStatus() bool {
	_, err := http.Get(a.statusUrl)
	return err == nil
}

func (a *TldrArchive) checkUpdate() bool {
	if !a.checkStatus() {
		logger.Log("[error] failed to get archive status; check your internet connection")
		return false
	}
	debugLogger.Log("[archive] checking for updates...")
	rev, err := a.getRev()
	if err != nil || rev != a.getRemoteRev() {
		return true
	}
	return false
}

func (a *TldrArchive) update() (bool, error) {
	debugLogger.Log("[archive] updating %s", a.zipPath)

	res, err := http.Get(a.remoteUrl + a.remotePath)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		msg := fmt.Sprintf("bad status code: %d", res.StatusCode)
		return false, errors.New(msg)
	}

	if err := os.MkdirAll(a.path, 0700); err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(a.zipPath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return false, err
	}

	err = ioutil.WriteFile(
		a.revPath,
		[]byte(a.getRemoteRev()),
		0600,
	)
	if err != nil {
		return true, err
	}

	return true, nil
}

func (a *TldrArchive) getPage(name string) (*TldrPage, error) {
	a.waitForUpdate()

	archive, err := zip.OpenReader(a.zipPath)
	if err != nil {
		log.Fatal(err)
	}
	defer archive.Close()

	var pages string
	if a.lang == "en" {
		pages = "pages"
	} else {
		pages = "pages." + a.lang
	}

	osName := osMap[runtime.GOOS]
	for i, dir := range []string{osName, "common"} {
		path := filepath.Join(pages, dir, name+".md")
		file, err := archive.Open(path)
		if err != nil {
			if i == 0 {
				continue
			}
			return nil, err
		}
		defer file.Close()
		buf, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
		return &TldrPage{name, string(buf)}, nil
	}
	return nil, nil
}
