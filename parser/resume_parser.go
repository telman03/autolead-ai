package parser

import (
	"os"
	"regexp"
	"strings"
)

type ResumeData struct {
	Name        string
	Email       string
	Phone       string
	Skills      []string
	Experience  []string
}

// Simple resume parser from .txt file
func ParseResume(filePath string) (*ResumeData, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	text := string(content)

	// Basic regex matchers
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	phoneRegex := regexp.MustCompile(`\+?\d[\d\s\-\(\)]{7,}\d`)

	email := emailRegex.FindString(text)
	phone := phoneRegex.FindString(text)

	// Simple field guesses
	lines := strings.Split(text, "\n")
	var name string
	if len(lines) > 0 {
		name = strings.TrimSpace(lines[0])
	}

	// Skills section
	var skills []string
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "skills") {
			skillLine := strings.ToLower(strings.ReplaceAll(line, "skills:", ""))
			skills = strings.Split(skillLine, ",")
			break
		}
	}

	// Experience (grab lines containing "experience" or job titles)
	var experience []string
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "experience") ||
			strings.Contains(strings.ToLower(line), "developer") ||
			strings.Contains(strings.ToLower(line), "engineer") {
			experience = append(experience, strings.TrimSpace(line))
		}
	}

	return &ResumeData{
		Name:       name,
		Email:      email,
		Phone:      phone,
		Skills:     skills,
		Experience: experience,
	}, nil
}
