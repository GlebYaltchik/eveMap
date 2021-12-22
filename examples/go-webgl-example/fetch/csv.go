package fetch

import (
	"bytes"
	"compress/bzip2"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-humble/locstor"
	"github.com/gocarina/gocsv"
)

type Logger func(...interface{})

func CSV(url string, res interface{}, log Logger) error {
	r, err := fetchURL(url, log)
	if err != nil {
		return fmt.Errorf("csv: %w", err)
	}

	if err := gocsv.Unmarshal(r, res); err != nil {
		return fmt.Errorf("csv: unmarshal: %w", err)
	}

	return nil
}

func fetchURL(url string, log Logger) (io.Reader, error) {
	if v, _ := locstor.GetItem(url); len(v) != 0 {
		log("using cached value for", url)
		return strings.NewReader(v), nil
	}

	log("fetching:", url)

	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	log("request completed: ", resp.Status, "Content-Length:", resp.ContentLength)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch: unexpected status code: %d (%s)", resp.StatusCode, resp.Status)
	}

	var r io.Reader = resp.Body

	if resp.Header.Get("content-type") == "application/x-bzip2" {
		r = bzip2.NewReader(resp.Body)
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("fetch: can't read data: %w", err)
	}

	if err := locstor.SetItem(url, string(data)); err != nil {
		log("WARN:", err)
	}

	return bytes.NewReader(data), nil
}
