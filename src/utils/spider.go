package utils

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	DownloadPage = "https://www.minecraft.net/en-us/download/server/bedrock"
	DownloadPrefix = "https://minecraft.azureedge.net/bin-linux/bedrock-server-"
	DownloadSuffix = ".zip"
)

const (
	Empty = ""
)

func fetchDownloadPage() (string, error) {
	header := make(http.Header)
	header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	header.Set("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	header.Set("Cache-Control", "max-age=0")
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	header.Set("Sec-Fetch-Dest", "document")
	header.Set("Sec-Fetch-Mode", "navigate")
	header.Set("Sec-Fetch-Site", "none")
	header.Set("Sec-Fetch-User", "?1")
	header.Set("Upgrade-Insecure-Requests", "1")

	req, err := http.NewRequest("GET", DownloadPage, nil)
	if err != nil {
		return Empty, err
	}

	req.Header = header
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Empty, err
	}

	if res.StatusCode != http.StatusOK {
		return Empty, fmt.Errorf("failed to fetch download page: %s", res.Status)
	}

	defer res.Body.Close()

	buf := make([]byte, 1024)
	var sb strings.Builder
	for {
		n, err := res.Body.Read(buf)
		sb.Write(buf[:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			return Empty, err
		}
	}

	return sb.String(), nil
}

func GetLatestVersion() (string, error) {
	body, err := fetchDownloadPage()
	if err != nil {
		return Empty, err
	}

	prefixIndex := strings.Index(body, DownloadPrefix)
	if prefixIndex == -1 {
		return Empty, fmt.Errorf("failed to find release download prefix")
	}

	body = body[prefixIndex+len(DownloadPrefix):]
	suffixIndex := strings.Index(body, DownloadSuffix)
	if suffixIndex == -1 {
		return Empty, fmt.Errorf("failed to find release download suffix")
	}

	version := body[:suffixIndex]

	return version, nil
}