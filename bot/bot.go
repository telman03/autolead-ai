package bot

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/telman03/autolead-ai/ai"
	"github.com/telman03/autolead-ai/matcher"
	"github.com/telman03/autolead-ai/parser"
	"github.com/telman03/autolead-ai/scraper"
)

func StartBot(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		userID := update.Message.From.ID

		if update.Message.IsCommand() {
			if update.Message.Command() == "start" {
				msg := tgbotapi.NewMessage(chatID, "👋 Welcome to AutoLead AI! Please send me your resume as a `.txt` or `.pdf` file.")
				bot.Send(msg)
			}
			continue
		}

		if update.Message.Document != nil {
			// ✅ Enforce Supabase-based usage limits
			allowed, _, err := CheckUserUsage(int64(userID))
			if err != nil {
				bot.Send(tgbotapi.NewMessage(chatID, "❌ Error checking usage limits. Try again later."))
				continue
			}
			if !allowed {
				msg := "⚠️ You've reached your daily limit of 3 cover letters. Upgrade to premium for unlimited access."
				bot.Send(tgbotapi.NewMessage(chatID, msg))
				continue
			}

			file := update.Message.Document
			fileURL, err := bot.GetFileDirectURL(file.FileID)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(chatID, "❌ Failed to get file URL."))
				continue
			}

			fileExt := strings.ToLower(filepath.Ext(file.FileName))
			localPath := "data/resume_from_user" + fileExt
			outputPath := "data/resume_from_user.txt"

			err = downloadFile(localPath, fileURL)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(chatID, "❌ Failed to download resume."))
				continue
			}

			if fileExt == ".pdf" {
				err = convertPDFToText(localPath, outputPath)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, "❌ Failed to extract text from PDF."))
					continue
				}
			} else if fileExt == ".txt" {
				outputPath = localPath
			} else {
				bot.Send(tgbotapi.NewMessage(chatID, "❌ Unsupported file type. Please send .txt or .pdf."))
				continue
			}

			bot.Send(tgbotapi.NewMessage(chatID, "📄 Resume received! Matching jobs and generating cover letters..."))

			resume, err := parser.ParseResume(outputPath)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(chatID, "❌ Couldn't parse the resume."))
				continue
			}

			jobs, err := scraper.ScrapeRemoteOK(10)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(chatID, "❌ Failed to fetch jobs."))
				continue
			}

			matches := matcher.MatchJobs(resume, jobs)
			topN := 3
			for i := 0; i < topN && i < len(matches); i++ {
				job := matches[i].Job
				letter, err := ai.GenerateCoverLetter(resume, job)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, "❌ Error generating cover letter."))
					continue
				}

				filename := fmt.Sprintf("output/%s_%s.txt", sanitize(job.Company), sanitize(job.Title))
				err = os.WriteFile(filename, []byte(letter), 0644)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, "❌ Failed to save cover letter."))
					continue
				}

				doc := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(filename))
				doc.Caption = fmt.Sprintf("📬 Cover Letter: %s at %s", job.Title, job.Company)
				bot.Send(doc)

				_ = IncrementUserUsage(int64(userID))
			}
		}
	}
}

func sanitize(name string) string {
	return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(name)), " ", "_")
}

func downloadFile(path string, url string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func convertPDFToText(pdfPath, txtPath string) error {
	cmd := exec.Command("pdftotext", pdfPath, txtPath)
	return cmd.Run()
}
