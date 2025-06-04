package bot

import (
	"fmt"
	"log"
	"time"
	"os"
	"github.com/joho/godotenv"
	supa "github.com/supabase-community/supabase-go"
)
var supabaseClient *supa.Client

func InitSupabase() error {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get environment variables
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	// Initialize client
	supabaseClient, err = supa.NewClient(supabaseURL, supabaseKey, nil)
	return err
}

type UserUsage struct {
	TelegramID int64     `json:"telegram_id"`
	DailyUses  int       `json:"daily_uses"`
	IsPremium  bool      `json:"is_premium"`
	LastReset  time.Time `json:"last_reset"`
}

const DailyLimit = 3

func CheckUserUsage(telegramID int64) (bool, bool, error) {
	var users []UserUsage

	_, err := supabaseClient.
		From("user_usage").
		Select("*", "", false).
		Eq("telegram_id", fmt.Sprintf("%d", telegramID)).
		ExecuteTo(&users)
	if err != nil {
		return false, false, err
	}

	if len(users) == 0 {
		newUser := UserUsage{
			TelegramID: telegramID,
			DailyUses:  0,
			IsPremium:  false,
			LastReset:  time.Now(),
		}
		_, _, err := supabaseClient.
			From("user_usage").
			Insert(newUser, false, "", "", "").
			Execute()
		if err != nil {
			return false, false, err
		}
		return true, false, nil
	}

	user := users[0]
	today := time.Now().Format("2006-01-02")
	last := user.LastReset.Format("2006-01-02")

	if today != last {
		_, _, err := supabaseClient.
			From("user_usage").
			Update(map[string]interface{}{
				"daily_uses": 0,
				"last_reset": time.Now(),
			}, "", "").
			Eq("telegram_id", fmt.Sprintf("%d", telegramID)).
			Execute()
		if err != nil {
			return false, user.IsPremium, err
		}
		user.DailyUses = 0
	}

	if user.IsPremium || user.DailyUses < DailyLimit {
		return true, user.IsPremium, nil
	}

	return false, user.IsPremium, nil
}

func IncrementUserUsage(telegramID int64) error {
	var users []UserUsage

	_, err := supabaseClient.
		From("user_usage").
		Select("*", "", false).
		Eq("telegram_id", fmt.Sprintf("%d", telegramID)).
		ExecuteTo(&users)
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return fmt.Errorf("user not found")
	}

	newCount := users[0].DailyUses + 1
	_, _, err = supabaseClient.
		From("user_usage").
		Update(map[string]interface{}{
			"daily_uses": newCount,
		}, "", "").
		Eq("telegram_id", fmt.Sprintf("%d", telegramID)).
		Execute()
	return err
}

func SetUserPremium(telegramID int64, premium bool) error {
	_, _, err := supabaseClient.
		From("user_usage").
		Update(map[string]interface{}{
			"is_premium": premium,
		}, "", "").
		Eq("telegram_id", fmt.Sprintf("%d", telegramID)).
		Execute()
	return err
}