package scanner

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/imroc/req/v3"
	sitemap "github.com/oxffaa/gopher-parse-sitemap"
	"golang.org/x/sync/errgroup"
)

const (
	MAX_SITEMAPS_INDEX = 10
)

func (s *ScanClient) isSitemapIndex(body io.Reader) bool {
	bodyStart := make([]byte, 256)
	_, err := body.Read(bodyStart)
	if err != nil {
		return false
	}

	bodyStartString := string(bodyStart)
	return strings.Contains(bodyStartString, "<sitemapindex")

}

func (s *ScanClient) getSitemapUrlsByUrl(url string) []string {
	var urls []string
	err := sitemap.ParseFromSite(url, func(e sitemap.Entry) error {
		urls = append(urls, e.GetLocation())
		return nil
	})
	if err != nil {
		return urls
	}

	return urls
}

func (s *ScanClient) getSitemap(bodyBuffer io.Reader, originUrl string) SitemapInfo {
	var currentSitemap SitemapInfo
	currentSitemap.LocationUrl = originUrl

	err := sitemap.Parse(bodyBuffer, func(e sitemap.Entry) error {
		currentSitemap.Urls = append(currentSitemap.Urls, e.GetLocation())
		return nil
	})
	if err != nil {
		return currentSitemap
	}

	return currentSitemap
}

func (s *ScanClient) getSitemapIndex(bodyBuffer io.Reader) SitemapIndex {
	var currentIndex SitemapIndex
	var sitemapsUrls []string

	err := sitemap.ParseIndex(bodyBuffer, func(e sitemap.IndexEntry) error {
		sitemapsUrls = append(sitemapsUrls, e.GetLocation())
		return nil
	})
	if err != nil {
		return currentIndex
	}

	g := new(errgroup.Group)
	for _, url := range sitemapsUrls[:MAX_SITEMAPS_INDEX] {

		url := url
		g.Go(func() error {
			currentIndex.Sitemaps = append(currentIndex.Sitemaps, SitemapInfo{
				LocationUrl: url,
				Urls: s.getSitemapUrlsByUrl(url),
			})
			return nil
		})

	}

	if err := g.Wait(); err == nil {
		fmt.Println("An error occurred while fetching sitemaps")
	}
	return currentIndex
}

func (s *ScanClient) sitemapExists(sitemapUrl string) bool {
	resp, err := req.C().R().Get(sitemapUrl)
	if err != nil || resp.IsErrorState() {
		fmt.Println("Sitemap Exists", err, resp.IsErrorState())
		fmt.Println(sitemapUrl, resp.StatusCode)
		return false
	}

	return true
}
 
func (s *ScanClient) GetSiteMaps() ([]SitemapIndex, error) {
	g := new(errgroup.Group)

	// Get the file
	if cap(s.sitemapUrls) == 0 {
		sitemapUrlString := fmt.Sprintf("%s/%s", s.baseUrl, "sitemap.xml")
		if !s.sitemapExists(sitemapUrlString) {
			return make([]SitemapIndex, 0), nil
		}

		// else continue on
		s.sitemapUrls = append(s.sitemapUrls, sitemapUrlString)
	}

	var sitemapIndexes []SitemapIndex

	for _, sitemapUrl := range s.sitemapUrls {

		sitemapUrl := sitemapUrl
		g.Go(func() error {
			resp, err := req.C().R().Get(sitemapUrl)
			if err != nil || resp.IsErrorState() {
				return err
			}
				
			// Read Body into Memory
			// This might be dangerous
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			
			bodyBuffer := bytes.NewBuffer(body)
			bodyBuffer2 := bytes.NewBuffer(body)
		
			if s.isSitemapIndex(bodyBuffer2) {
				sitemapIndexes = append(sitemapIndexes, s.getSitemapIndex(bodyBuffer))
			} else {
				// todo: maybe eventually check if is acuallty a sitemap
				sitemaps := make([]SitemapInfo, 0)
				sitemaps = append(sitemaps, s.getSitemap(bodyBuffer, sitemapUrl), )
				sitemapIndexes = append(sitemapIndexes, SitemapIndex{
					Sitemaps: sitemaps,
				})
			}

			return nil
		})

	}


	if err := g.Wait(); err == nil {
		fmt.Println("Error")
	}

	return sitemapIndexes, nil
}