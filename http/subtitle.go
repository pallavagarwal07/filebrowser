package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	fb "github.com/pallavagarwal07/filebrowser"
)

func subtitlesHandler(c *fb.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	files, err := ReadDir(filepath.Dir(c.File.Path))
	if err != nil {
		return http.StatusInternalServerError, err
	}
	var subtitles = make([]map[string]string, 0)
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext == ".vtt" || ext == ".srt" {
			var sub map[string]string = make(map[string]string)
			sub["src"] = filepath.Dir(c.File.Path) + "/" + file.Name()
			sub["kind"] = "subtitles"
			sub["label"] = file.Name()
			subtitles = append(subtitles, sub)
		}
	}
	return renderJSON(w, subtitles)
}

func subtitleHandler(c *fb.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	str, err := CleanSubtitle(c.File.Path)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	file, err := os.Open(c.File.Path)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Header().Set("Content-Disposition", "inline")
	w.Header().Set("Content-Type", "text/vtt")
	http.ServeContent(w, r, stat.Name(), stat.ModTime(), bytes.NewReader([]byte(str)))

	return 0, nil

}

func CleanSubtitle(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	str := string(b) // convert content to a 'string'
	ext := filepath.Ext(filename)
	if ext == ".srt" {
		re := regexp.MustCompile("([0-9]{2}:[0-9]{2}:[0-9]{2}),([0-9]{3})")
		str = "WEBVTT\n\n" + re.ReplaceAllString(str, "$1.$2")
	}
	return str, err
}

func ReadDir(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	return list, nil
}
