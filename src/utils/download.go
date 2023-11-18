package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/schollz/progressbar/v3"
)

func DownloadFile(uri string, bar *progressbar.ProgressBar) (string, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return Empty, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Empty, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Empty, fmt.Errorf("failed to download file %s (%s)", uri, res.Status)
	}

	fileName := path.Base(res.Request.URL.Path)

	file, err := os.Create(fileName)
	if err != nil {
		return Empty, err
	}

	defer file.Close()

	if bar != nil {
		bar.ChangeMax(int(res.ContentLength))
		_, err = io.Copy(io.MultiWriter(file, bar), res.Body)
	} else {
		_, err = io.Copy(file, res.Body)
	}

	if err != nil {
		return Empty, err
	}

	return fileName, nil
}

func DownloadVersion(version string) (string, error) {
	bar := progressbar.NewOptions(
		114514,
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetDescription(" Downloading..."),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	uri := DownloadPrefix + version + DownloadSuffix
	
	return DownloadFile(uri, bar)
}
