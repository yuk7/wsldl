package download

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
	"github.com/yuk7/wsldl/lib/errutil"
)

var (
	downloadStatFile    = os.Stat
	downloadOpenForRead = func(path string) (io.ReadCloser, error) { return os.Open(path) }
	downloadCalcSHA256  = calculateSHA256
)

func DownloadFile(ctx context.Context, url, dest string, progressBarWidth int) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	resumeOffset := int64(0)
	info, err := downloadStatFile(dest)
	if err == nil {
		resumeOffset = info.Size()
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	if resumeOffset > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", resumeOffset))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusRequestedRangeNotSatisfiable && resumeOffset > 0 {
		return downloadCalcSHA256(dest)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	appendMode := resumeOffset > 0 && resp.StatusCode == http.StatusPartialContent
	if !appendMode {
		resumeOffset = 0
	}

	openFlags := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	if appendMode {
		openFlags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	}

	f, err := os.OpenFile(dest, openFlags, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()

	size := resp.ContentLength
	if appendMode && size > 0 {
		size += resumeOffset
	}

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
	if appendMode && resumeOffset > 0 && size > 0 {
		_ = bar.Add64(resumeOffset)
	}

	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
	if err != nil && err != io.EOF {
		return "", err
	}

	sha256String, err := downloadCalcSHA256(dest)
	if err != nil {
		return "", err
	}

	return sha256String, nil
}

func calculateSHA256(path string) (string, error) {
	f, err := downloadOpenForRead(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	sum := sha256.New()
	if _, err := io.Copy(sum, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(sum.Sum(nil)), nil
}
