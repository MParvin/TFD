package tools

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

// StartBot initializes and starts the Telegram bot
func StartBot() error {
	// Read configuration
	token := viper.GetString("token")
	if token == "" {
		return fmt.Errorf("telegram bot token not configured")
	}

	adminChatID := viper.GetInt64("admin_chat_id")
	allowedUsers := viper.GetIntSlice("allowed_users")
	supportUser := viper.GetString("support_user")
	proxyURL := viper.GetString("proxy")

	// Create HTTP client with optional proxy
	client, err := GetHTTPClient(proxyURL)
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	// Create bot instance
	bot, err := tgbotapi.NewBotAPIWithClient(token, tgbotapi.APIEndpoint, client)
	if err != nil {
		return fmt.Errorf("failed to create bot: %w", err)
	}

	// Verify bot token
	me, err := bot.GetMe()
	if err != nil {
		return fmt.Errorf("failed to verify bot token: %w", err)
	}

	log.Printf("TFD Bot started successfully as @%s (ID: %d)", me.UserName, me.ID)

	// Create update config
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	// Start listening for updates
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		message := update.Message
		userID := int(message.From.ID)

		// Check if user is authorized
		if !isUserAllowed(userID, allowedUsers) {
			log.Printf("Unauthorized user %d (@%s) tried to use bot", userID, message.From.UserName)

			replyText := "❌ Sorry, you are not authorized to use this bot."
			if supportUser != "" {
				replyText += fmt.Sprintf(" Please contact @%s for access.", supportUser)
			}

			reply := tgbotapi.NewMessage(message.Chat.ID, replyText)
			if _, err := bot.Send(reply); err != nil {
				HandleError(fmt.Errorf("failed to send unauthorized message: %w", err))
			}
			continue
		}

		// Process the message
		if err := handleMessage(bot, message, client, adminChatID); err != nil {
			HandleError(err)
			SendErrorToAdmin(bot, adminChatID, err)

			// Send error reply to user
			reply := tgbotapi.NewMessage(message.Chat.ID, "❌ Sorry, an error occurred while processing your message.")
			if _, sendErr := bot.Send(reply); sendErr != nil {
				HandleError(fmt.Errorf("failed to send error reply: %w", sendErr))
			}
		}
	}

	return nil
}

// isUserAllowed checks if a user ID is in the allowed users list
func isUserAllowed(userID int, allowedUsers []int) bool {
	for _, allowedID := range allowedUsers {
		if userID == allowedID {
			return true
		}
	}
	return false
}

// handleMessage processes a single message based on its type
func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, client *http.Client, adminChatID int64) error {
	var err error
	var filename string
	var destDir string

	switch {
	case message.Photo != nil:
		filename, err = handlePhotoMessage(bot, message, client)
		destDir = viper.GetString("directories.photo")

	case message.Video != nil:
		filename, err = handleVideoMessage(bot, message, client)
		destDir = viper.GetString("directories.video")

	case message.Audio != nil:
		filename, err = handleAudioMessage(bot, message, client)
		destDir = viper.GetString("directories.music")

	case message.Voice != nil:
		filename, err = handleVoiceMessage(bot, message, client)
		destDir = viper.GetString("directories.voice")

	case message.Document != nil:
		filename, err = handleDocumentMessage(bot, message, client)
		destDir = viper.GetString("directories.document")

	case message.Text != "":
		filename, err = handleTextMessage(message.Text, client)
		destDir = viper.GetString("directories.other")

	default:
		log.Printf("Unsupported message type from user %d", message.From.ID)
		reply := tgbotapi.NewMessage(message.Chat.ID, "ℹ️ This message type is not supported. Please send photos, videos, documents, audio, or URLs.")
		_, sendErr := bot.Send(reply)
		return sendErr
	}

	if err != nil {
		return err
	}

	// Send confirmation to user
	if filename != "" {
		confirmText := fmt.Sprintf("✅ Downloaded successfully!\n📁 File: %s\n📂 Location: %s", filename, destDir)
		reply := tgbotapi.NewMessage(message.Chat.ID, confirmText)
		if _, err := bot.Send(reply); err != nil {
			return fmt.Errorf("failed to send confirmation message: %w", err)
		}
	}

	return nil
}

// handlePhotoMessage downloads the largest photo size
func handlePhotoMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, client *http.Client) (string, error) {
	// Get the largest photo size
	photos := message.Photo
	if len(photos) == 0 {
		return "", fmt.Errorf("no photo found in message")
	}

	largestPhoto := photos[len(photos)-1] // Last one is usually the largest
	destDir := viper.GetString("directories.photo")

	return downloadTelegramFile(bot, largestPhoto.FileID, destDir, "photo", client)
}

// handleVideoMessage downloads video files
func handleVideoMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, client *http.Client) (string, error) {
	destDir := viper.GetString("directories.video")
	return downloadTelegramFile(bot, message.Video.FileID, destDir, "video", client)
}

// handleAudioMessage downloads audio files
func handleAudioMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, client *http.Client) (string, error) {
	destDir := viper.GetString("directories.music")
	return downloadTelegramFile(bot, message.Audio.FileID, destDir, "audio", client)
}

// handleVoiceMessage downloads voice messages
func handleVoiceMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, client *http.Client) (string, error) {
	destDir := viper.GetString("directories.voice")
	return downloadTelegramFile(bot, message.Voice.FileID, destDir, "voice", client)
}

// handleDocumentMessage downloads document files
func handleDocumentMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, client *http.Client) (string, error) {
	destDir := viper.GetString("directories.document")
	return downloadTelegramFile(bot, message.Document.FileID, destDir, "document", client)
}

// handleTextMessage processes text messages and downloads URLs if found
func handleTextMessage(text string, client *http.Client) (string, error) {
	// Look for URLs in the text
	urlRegex := regexp.MustCompile(`https?://[^\s]+`)
	urls := urlRegex.FindAllString(text, -1)

	if len(urls) == 0 {
		return "", fmt.Errorf("no URLs found in text message")
	}

	destDir := viper.GetString("directories.other")

	// Download the first URL found
	url := urls[0]
	if err := DownloadFromURL(url, destDir, client); err != nil {
		return "", fmt.Errorf("failed to download from URL %s: %w", url, err)
	}

	// Extract filename from URL
	filename := filepath.Base(url)
	if filename == "" || filename == "." {
		filename = "downloaded_file"
	}

	return filename, nil
}

// downloadTelegramFile downloads a file from Telegram servers
func downloadTelegramFile(bot *tgbotapi.BotAPI, fileID, destDir, fileType string, client *http.Client) (string, error) {
	// Get file info from Telegram
	fileConfig := tgbotapi.FileConfig{FileID: fileID}
	file, err := bot.GetFile(fileConfig)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}

	// Get download URL
	fileURL := bot.GetFileDirectURL(fileID)

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s_%s", fileType, timestamp, filepath.Base(file.FilePath))
	if filepath.Base(file.FilePath) == "" {
		filename = fmt.Sprintf("%s_%s", fileType, timestamp)
	}

	// Sanitize filename
	filename = sanitizeFilename(filename)

	// Ensure destination directory exists
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	// Download the file
	fullPath := filepath.Join(destDir, filename)
	if err := downloadFileFromURL(fileURL, fullPath, client); err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}

	log.Printf("Downloaded %s: %s", fileType, filename)
	return filename, nil
}

// downloadFileFromURL downloads a file from a URL to a specific path
func downloadFileFromURL(url, destPath string, client *http.Client) error {
	if client == nil {
		client = http.DefaultClient
	}

	// Make HTTP request
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error %d when downloading file", resp.StatusCode)
	}

	// Create destination file
	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", destPath, err)
	}
	defer file.Close()

	// Copy file content
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file content: %w", err)
	}

	return nil
}

// sanitizeFilename removes dangerous characters from filenames (duplicate function for this package)
func sanitizeFilename(filename string) string {
	// Remove path separators and dangerous characters
	re := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)
	filename = re.ReplaceAllString(filename, "_")

	// Remove .. sequences
	filename = regexp.MustCompile(`\.\.+`).ReplaceAllString(filename, "_")

	// Limit length
	if len(filename) > 255 {
		filename = filename[:255]
	}

	// Ensure not empty
	if filename == "" {
		filename = "file"
	}

	return filename
}
