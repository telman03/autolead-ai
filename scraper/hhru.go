package scraper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeHHRU(query string, max int) ([]Job, error) {
	var jobs []Job
	searchURL := fmt.Sprintf("https://hh.ru/search/vacancy?text=%s", strings.ReplaceAll(query, " ", "+"))

	// Create request with Russian user-agent
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/113")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch page: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	doc.Find("div.serp-item").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i >= max {
			return false
		}

		title := s.Find("a.bloko-link").Text()
		href, _ := s.Find("a.bloko-link").Attr("href")
		company := s.Find("div.vacancy-serp-item__meta-info-company").Text()
		description := s.Find("div.g-user-content").Text()

		jobs = append(jobs, Job{
			Title:       strings.TrimSpace(title),
			Link:        strings.TrimSpace(href),
			Company:     strings.TrimSpace(company),
			Description: strings.TrimSpace(description),
			Tags:        []string{}, // Initialize empty tags
		})

		return true
	})

	return jobs, nil
}
