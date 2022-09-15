package main

import (
	"archive/zip"
	"bufio"
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
	DataRemoteUrl   = "https://github.com/atlasamerican/cheatsheet/"
	DataRemotePath  = "raw/assets/data.zip"
	GithubStatusUrl = "https://www.githubstatus.com/api/v2/status.json"
)

var osMap = map[string]string{
	"linux":   "linux",
	"darwin":  "osx",
	"windows": "windows",
}

type TldrPage struct {
	name    string
	content string
}

type Archive[T TldrPage | Dataset] struct {
	remoteUrl  string
	remotePath string
	remoteRef  string
	statusUrl  string
	path       string
	zipPath    string
	revPath    string
	lang       string
	updating   chan bool
}

func newTldrArchive(path string) *Archive[TldrPage] {
	a := &Archive[TldrPage]{
		remoteUrl:  TldrRemoteUrl,
		remotePath: TldrRemotePath,
		remoteRef:  "HEAD",
		statusUrl:  GithubStatusUrl,
		path:       path,
		zipPath:    filepath.Join(path, "tldr.zip"),
		revPath:    filepath.Join(path, "tldr.zip.rev"),
		lang:       "en",
		updating:   make(chan bool, 1),
	}

	a.init()

	return a
}

func newDataArchive(path string) *Archive[Dataset] {
	a := &Archive[Dataset]{
		remoteUrl:  DataRemoteUrl,
		remotePath: DataRemotePath,
		remoteRef:  "refs/heads/assets",
		statusUrl:  GithubStatusUrl,
		path:       path,
		zipPath:    filepath.Join(path, "data.zip"),
		revPath:    filepath.Join(path, "data.zip.rev"),
		lang:       "en",
		updating:   make(chan bool, 1),
	}

	a.init()

	return a
}

func (a *Archive[T]) init() {
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
		close(a.updating)
	}()
}

func (a *Archive[T]) waitForUpdate() {
	_, ok := <-a.updating
	if !ok {
		return
	}
	<-a.updating
}

func (a *Archive[T]) getRemoteRev() string {
	cmd := exec.Command("git", "ls-remote", a.remoteUrl)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	var rev string
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if _, err := fmt.Sscanf(line, "%s "+a.remoteRef, &rev); err == nil {
			break
		}
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	if rev == "" {
		log.Fatalf("Failed to get revision for: %s", a.remoteUrl)
	}

	return rev
}

func (a *Archive[T]) getRev() (string, error) {
	buf, err := ioutil.ReadFile(a.revPath)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (a *Archive[T]) checkStatus() bool {
	_, err := http.Get(a.statusUrl)
	return err == nil
}

func (a *Archive[T]) checkUpdate() bool {
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

func (a *Archive[T]) update() (bool, error) {
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

func (a *Archive[T]) getPage(name string) (*TldrPage, error) {
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

func (a *Archive[T]) getCommands() []Command {
	a.waitForUpdate()

	cmds := make([]Command, 0)

	archive, err := zip.OpenReader(a.zipPath)
	if err != nil {
		log.Fatal(err)
	}
	defer archive.Close()

	for _, file := range archive.File {
		f, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		buf, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}

		readDataBuf(buf, &cmds)
	}

	return cmds
}
