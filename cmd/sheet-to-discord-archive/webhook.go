package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// parseWebhookURL extracts id and token from https://discord.com/api/webhooks/{id}/{token}.
func parseWebhookURL(raw string) (id, token string, err error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", "", fmt.Errorf("invalid webhook url")
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	for i := 0; i+2 < len(parts); i++ {
		if parts[i] == "webhooks" && parts[i+1] != "" && parts[i+2] != "" {
			return parts[i+1], parts[i+2], nil
		}
	}
	return "", "", fmt.Errorf("invalid webhook url: expected .../webhooks/{id}/{token}")
}

type webhookFile struct {
	Name string
	Data []byte
}

type webhookPayload struct {
	Content         string            `json:"content,omitempty"`
	Username        string            `json:"username,omitempty"`
	AvatarURL       string            `json:"avatar_url,omitempty"`
	ThreadName      string            `json:"thread_name,omitempty"`
	Embeds          []json.RawMessage `json:"embeds,omitempty"`
	AllowedMentions map[string]any    `json:"allowed_mentions"`
	Attachments     []map[string]any  `json:"attachments,omitempty"`
}

type webhookClient struct {
	id, token string
	http      *http.Client
	delay     time.Duration
}

// execute posts to the webhook and returns the created message's channel_id.
func (c *webhookClient) execute(p webhookPayload, threadID string, files []webhookFile) (string, error) {
	p.AllowedMentions = map[string]any{"parse": []string{}}
	q := url.Values{"wait": {"true"}}
	if threadID != "" {
		q.Set("thread_id", threadID)
	}
	endpoint := fmt.Sprintf("https://discord.com/api/v10/webhooks/%s/%s?%s", c.id, c.token, q.Encode())
	for i := range files {
		p.Attachments = append(p.Attachments, map[string]any{"id": i, "filename": files[i].Name})
	}
	for attempt := 0; attempt <= 5; attempt++ {
		body, ctype, err := encodeBody(p, files)
		if err != nil {
			return "", err
		}
		resp, err := c.http.Post(endpoint, ctype, body)
		if err != nil {
			// never surface the URL (contains token)
			return "", fmt.Errorf("webhook request failed")
		}
		data, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode == http.StatusTooManyRequests {
			var rl struct {
				RetryAfter float64 `json:"retry_after"`
			}
			_ = json.Unmarshal(data, &rl)
			wait := time.Duration(rl.RetryAfter*1000)*time.Millisecond + 100*time.Millisecond
			time.Sleep(wait)
			continue
		}
		time.Sleep(c.delay)
		if resp.StatusCode/100 != 2 {
			return "", fmt.Errorf("webhook status %d: %s", resp.StatusCode, truncateRunes(string(data), 300))
		}
		var msg struct {
			ChannelID string `json:"channel_id"`
		}
		if err := json.Unmarshal(data, &msg); err != nil {
			return "", fmt.Errorf("decode webhook response: %w", err)
		}
		return msg.ChannelID, nil
	}
	return "", fmt.Errorf("webhook rate limited: gave up after retries")
}

func encodeBody(p webhookPayload, files []webhookFile) (io.Reader, string, error) {
	js, err := json.Marshal(p)
	if err != nil {
		return nil, "", err
	}
	if len(files) == 0 {
		return bytes.NewReader(js), "application/json", nil
	}
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	if err := w.WriteField("payload_json", string(js)); err != nil {
		return nil, "", err
	}
	for i, f := range files {
		fw, err := w.CreateFormFile(fmt.Sprintf("files[%d]", i), f.Name)
		if err != nil {
			return nil, "", err
		}
		if _, err := fw.Write(f.Data); err != nil {
			return nil, "", err
		}
	}
	if err := w.Close(); err != nil {
		return nil, "", err
	}
	return &buf, w.FormDataContentType(), nil
}

func downloadFile(client *http.Client, u string) ([]byte, error) {
	resp, err := client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("download status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
