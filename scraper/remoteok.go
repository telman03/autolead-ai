package scraper

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Job struct {
	Title       string   `json:"title"`
	Company     string   `json:"company"`
	Tags        []string `json:"tags"`
	Link        string   `json:"link"`
	Description string   `json:"description"`
}

func ScrapeRemoteOK(limit int) ([]Job, error) {
	var jobs []Job
	c := colly.NewCollector(
		colly.AllowedDomains("remoteok.com"),
	)

	c.OnHTML("tr.job", func(e *colly.HTMLElement) {
		if len(jobs) >= limit {
			return
		}
		title := e.ChildText("h2")
		company := e.ChildText(".companyLink > h3")
		link := "https://remoteok.com" + e.Attr("data-href")

		// Tags
		var tags []string
		e.ForEach(".tag", func(_ int, el *colly.HTMLElement) {
			tags = append(tags, strings.ToLower(el.Text))
		})

		// Job Description (can be extended by visiting link)
		description := e.ChildText(".description")

		job := Job{
			Title:       title,
			Company:     company,
			Tags:        tags,
			Link:        link,
			Description: description,
		}

		jobs = append(jobs, job)
	})

	err := c.Visit("https://remoteok.com/remote-dev-jobs")
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func SaveJobsToFile(jobs []Job, path string) error {
	data, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
