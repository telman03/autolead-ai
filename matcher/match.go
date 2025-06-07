package matcher

import (
	"strings"

	"github.com/telman03/autolead-ai/parser"
	"github.com/telman03/autolead-ai/scraper"
)

type ScoredJob struct {
	scraper.Job
	Score int
}

// clean & normalize text
func normalize(text string) string {
	return strings.ToLower(strings.TrimSpace(text))
}

// count how many resume skills are found in job tags/desc/title
func scoreJob(resume *parser.ResumeData, job scraper.Job) int {
	score := 0

	// Combine title + description + tags into one searchable text
	var tagsText string
	if job.Tags != nil {
		tagsText = strings.Join(job.Tags, " ")
	}
	fullText := normalize(job.Title + " " + job.Description + " " + tagsText)

	for _, skill := range resume.Skills {
		if skill == "" {
			continue
		}
		if strings.Contains(fullText, normalize(skill)) {
			score += 10 // skill match = +10 pts
		}
	}

	for _, exp := range resume.Experience {
		if exp == "" {
			continue
		}
		if strings.Contains(fullText, normalize(exp)) {
			score += 5 // experience line match = +5 pts
		}
	}

	if strings.Contains(normalize(job.Title), normalize(resume.Name)) {
		score -= 20 // avoid jobs accidentally matching name
	}

	if score > 100 {
		score = 100
	}

	return score
}

func MatchJobs(resume *parser.ResumeData, jobs []scraper.Job) []ScoredJob {
	var scored []ScoredJob
	for _, job := range jobs {
		score := scoreJob(resume, job)
		scored = append(scored, ScoredJob{
			Job:   job,
			Score: score,
		})
	}

	// Optional: sort by score descending
	for i := 0; i < len(scored); i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].Score > scored[i].Score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	return scored
}
