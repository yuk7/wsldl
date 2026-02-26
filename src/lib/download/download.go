package download

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
	"github.com/yuk7/wsldl/lib/errutil"
)

func DownloadFile(url, dest string, progressBarWidth int) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()

	sum := sha256.New()

	size := resp.ContentLength

	if size == 0 {
		if progressBarWidth > 0 {
			errutil.ErrorRedPrintln("Failed to get total file size")
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
	_, err = io.Copy(io.MultiWriter(f, bar, sum), resp.Body)
	if err != nil && err != io.EOF {
		return "", err
	}

	sha256String := hex.EncodeToString(sum.Sum(nil))

	return sha256String, nil
}
