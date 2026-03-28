package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Bot struct {
	token      string
	chatID     string
	client     *http.Client
	OnApproval func(blockID, decision string)
}

func New(token, chatID string) *Bot {
	return &Bot{
		token:  token,
		chatID: chatID,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (b *Bot) SendWorkBlockApproval(blockID, title, goal, transition string) error {
	label := "Activate"
	if transition == "ready_to_ship" {
		label = "Ship"
	}

	text := fmt.Sprintf("*Work Block: %s*\nGoal: %s\nTransition: %s", escapeMarkdown(title), escapeMarkdown(goal), transition)

	payload := map[string]any{
		"chat_id":    b.chatID,
		"text":       text,
		"parse_mode": "Markdown",
		"reply_markup": map[string]any{
			"inline_keyboard": [][]map[string]string{
				{
					{"text": label, "callback_data": "approve:" + blockID},
					{"text": "Reject", "callback_data": "reject:" + blockID},
				},
			},
		},
	}

	body, _ := json.Marshal(payload)
	resp, err := b.client.Post(b.apiURL("sendMessage"), "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram sendMessage: status %d", resp.StatusCode)
	}
	return nil
}

func (b *Bot) SendMessage(text string) error {
	payload := map[string]any{
		"chat_id":    b.chatID,
		"text":       text,
		"parse_mode": "Markdown",
	}
	body, _ := json.Marshal(payload)
	resp, err := b.client.Post(b.apiURL("sendMessage"), "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (b *Bot) StartPolling(ctx context.Context) {
	offset := 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		updates, newOffset, err := b.getUpdates(offset)
		if err != nil {
			log.WithError(err).Warn("telegram: poll error")
			time.Sleep(5 * time.Second)
			continue
		}
		offset = newOffset

		for _, u := range updates {
			if u.CallbackQuery == nil {
				continue
			}
			data := u.CallbackQuery.Data
			parts := strings.SplitN(data, ":", 2)
			if len(parts) != 2 {
				continue
			}
			action, blockID := parts[0], parts[1]
			if action != "approve" && action != "reject" {
				continue
			}

			// Answer the callback to dismiss the spinner
			b.answerCallback(u.CallbackQuery.ID, action)

			if b.OnApproval != nil {
				b.OnApproval(blockID, action)
			}
		}
	}
}

type update struct {
	UpdateID      int            `json:"update_id"`
	CallbackQuery *callbackQuery `json:"callback_query"`
}

type callbackQuery struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

func (b *Bot) getUpdates(offset int) ([]update, int, error) {
	url := fmt.Sprintf("%s?offset=%d&timeout=30", b.apiURL("getUpdates"), offset)
	resp, err := b.client.Get(url)
	if err != nil {
		return nil, offset, err
	}
	defer resp.Body.Close()

	var result struct {
		OK     bool     `json:"ok"`
		Result []update `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, offset, err
	}

	newOffset := offset
	for _, u := range result.Result {
		if u.UpdateID >= newOffset {
			newOffset = u.UpdateID + 1
		}
	}
	return result.Result, newOffset, nil
}

func (b *Bot) answerCallback(callbackID, action string) {
	text := "Approved"
	if action == "reject" {
		text = "Rejected"
	}
	payload := map[string]any{
		"callback_query_id": callbackID,
		"text":              text,
	}
	body, _ := json.Marshal(payload)
	b.client.Post(b.apiURL("answerCallbackQuery"), "application/json", bytes.NewReader(body))
}

func (b *Bot) apiURL(method string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", b.token, method)
}

func escapeMarkdown(s string) string {
	r := strings.NewReplacer("*", "\\*", "_", "\\_", "`", "\\`", "[", "\\[")
	return r.Replace(s)
}
