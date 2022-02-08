package utils

import (
	"io"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

func DownloadFile(url, dest string, progressBarWidth int) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	size := resp.ContentLength

	if size == 0 {
		if progressBarWidth > 0 {
			ErrorRedPrintln("Failed to get total file size")
		}
		size = -1
	}

	bar := progressbar.NewOptions64(
		size,
		progressbar.OptionSetVisibility(false),
	)
	if progressBarWidth > 0 {
		bar = progressbar.NewOptions64(
			size,
			progressbar.OptionSetWidth(progressBarWidth),
			progressbar.OptionShowBytes(true),
			progressbar.OptionShowCount(),
		)
	} else if progressBarWidth < 0 {
		bar = progressbar.NewOptions64(
			-1,
			progressbar.OptionShowBytes(true),
			progressbar.OptionShowCount(),
		)
	}
	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}
