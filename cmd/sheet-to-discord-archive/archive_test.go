package main

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

func TestBuildTitle(t *testing.T) {
	long := strings.Repeat("あ", 120)
	tests := []struct {
		parts []string
		want  string
	}{
		{[]string{"タコスラブ", "やりとり"}, "タコスラブ/やりとり"},
		{[]string{"A", "b", "th"}, "A/b/th"},
		{[]string{long}, strings.Repeat("あ", 100)},
	}
	for _, tt := range tests {
		if got := buildTitle(tt.parts...); got != tt.want {
			t.Errorf("buildTitle(%v)=%q want %q", tt.parts, got, tt.want)
		}
	}
}

func TestPickDisplayName(t *testing.T) {
	tests := []struct{ nick, global, user, want string }{
		{"nick", "g", "u", "nick"},
		{"", "g", "u", "g"},
		{" ", "", "u", "u"},
		{"", "", "", "unknown"},
	}
	for _, tt := range tests {
		if got := pickDisplayName(tt.nick, tt.global, tt.user); got != tt.want {
			t.Errorf("got %q want %q", got, tt.want)
		}
	}
}

func TestSanitizeUsername(t *testing.T) {
	tests := []struct{ in, want string }{
		{"taro", "taro"},
		{"My Discord Bot", "My disc0rd Bot"},
		{"CLYDE fan", "clyd3 fan"},
		{"  ", "unknown"},
		{strings.Repeat("x", 90), strings.Repeat("x", 80)},
	}
	for _, tt := range tests {
		if got := sanitizeUsername(tt.in); got != tt.want {
			t.Errorf("sanitize(%q)=%q want %q", tt.in, got, tt.want)
		}
	}
}

func TestTimestampSuffix(t *testing.T) {
	ts := time.Date(2026, 9, 5, 9, 2, 30, 0, time.UTC)
	if got := timestampSuffix(ts); got != "-# 2026-09-05 18:02" {
		t.Errorf("got %q", got)
	}
}

func TestReplyPrefix(t *testing.T) {
	tests := []struct{ name, content, want string }{
		{"花子", "こんにちは\nよろしく", "-# ↪ 花子への返信：こんにちは よろしく"},
		{"a", strings.Repeat("x", 50), "-# ↪ aへの返信：" + strings.Repeat("x", 40)},
	}
	for _, tt := range tests {
		if got := replyPrefix(tt.name, tt.content); got != tt.want {
			t.Errorf("got %q want %q", got, tt.want)
		}
	}
}

func TestIsPostable(t *testing.T) {
	rich := json.RawMessage(`{"type":"rich","title":"x"}`)
	link := json.RawMessage(`{"type":"link","url":"x"}`)
	tests := []struct {
		name string
		m    srcMessage
		want bool
	}{
		{"default text", srcMessage{Type: 0, Content: "hi"}, true},
		{"reply", srcMessage{Type: 19, Content: "hi"}, true},
		{"slash", srcMessage{Type: 20, Content: "hi"}, true},
		{"context menu", srcMessage{Type: 23, Content: "hi"}, true},
		{"join", srcMessage{Type: 7, Content: "hi"}, false},
		{"thread created", srcMessage{Type: 18, Content: "t"}, false},
		{"empty", srcMessage{Type: 0}, false},
		{"attachment only", srcMessage{Type: 0, Attachments: []srcAttachment{{Filename: "a"}}}, true},
		{"rich embed only", srcMessage{Type: 0, Embeds: []json.RawMessage{rich}}, true},
		{"link embed only", srcMessage{Type: 0, Embeds: []json.RawMessage{link}}, false},
	}
	for _, tt := range tests {
		if got := isPostable(tt.m); got != tt.want {
			t.Errorf("%s: got %v", tt.name, got)
		}
	}
}

func TestSplitContent(t *testing.T) {
	line := strings.Repeat("a", 1500)
	tests := []struct {
		name  string
		in    string
		count int
	}{
		{"empty", "", 0},
		{"short", "abc", 1},
		{"exact", strings.Repeat("b", 2000), 1},
		{"newline split", line + "\n" + line, 2},
		{"hard split", strings.Repeat("c", 4500), 3},
	}
	for _, tt := range tests {
		got := splitContent(tt.in, 2000)
		if len(got) != tt.count {
			t.Errorf("%s: %d chunks want %d", tt.name, len(got), tt.count)
		}
		for _, c := range got {
			if utf8.RuneCountInString(c) > 2000 {
				t.Errorf("%s: chunk too long", tt.name)
			}
		}
	}
	if got := splitContent(line+"\n"+line, 2000); got[0] != line {
		t.Errorf("newline split should cut at newline")
	}
}

func TestBuildChunks(t *testing.T) {
	ts := time.Date(2026, 9, 5, 9, 2, 0, 0, time.UTC)
	name := func(u srcUser) string { return u.Username }
	m := srcMessage{Type: 19, Content: "了解", Timestamp: ts,
		Reference: &struct {
			Type      int    `json:"type"`
			MessageID string `json:"message_id"`
		}{0, "1"},
		Referenced: &srcMessage{Content: "元", Author: srcUser{Username: "bob"}},
		Stickers: []struct {
			Name string `json:"name"`
		}{{"wave"}},
	}
	got := buildChunks(m, []srcAttachment{{URL: "https://x/big.mp4"}}, name)
	want := "-# ↪ bobへの返信：元\n了解\n[sticker: wave]\nhttps://x/big.mp4\n-# 2026-09-05 18:02"
	if len(got) != 1 || got[0] != want {
		t.Errorf("got %q", got)
	}
	long := buildChunks(srcMessage{Type: 0, Content: strings.Repeat("z", 1995), Timestamp: ts}, nil, name)
	if len(long) != 2 || !strings.HasSuffix(long[1], "18:02") {
		t.Errorf("long split wrong: %d", len(long))
	}
}

func TestSplitUploads(t *testing.T) {
	atts := []srcAttachment{{Size: 6}, {Size: 5}}
	if up, ln := splitUploads(atts, 20); len(up) != 2 || len(ln) != 0 {
		t.Errorf("under cap")
	}
	if up, ln := splitUploads(atts, 10); len(up) != 0 || len(ln) != 2 {
		t.Errorf("over cap")
	}
}

func TestParseWebhookURL(t *testing.T) {
	tests := []struct {
		in, id, token string
		ok            bool
	}{
		{"https://discord.com/api/webhooks/123/abc-DEF", "123", "abc-DEF", true},
		{"https://discord.com/api/v10/webhooks/123/tok?x=1", "123", "tok", true},
		{"https://discord.com/api/webhooks/123", "", "", false},
		{"nonsense", "", "", false},
	}
	for _, tt := range tests {
		id, tok, err := parseWebhookURL(tt.in)
		if (err == nil) != tt.ok || id != tt.id || tok != tt.token {
			t.Errorf("%q: %q %q %v", tt.in, id, tok, err)
		}
		if err != nil && strings.Contains(err.Error(), "123") {
			t.Errorf("error leaks url")
		}
	}
}

func TestStateLoadSave(t *testing.T) {
	path := filepath.Join(t.TempDir(), "s.json")
	st, err := loadState(path)
	if err != nil || len(st) != 0 {
		t.Fatalf("missing file should give empty state: %v", err)
	}
	st2 := st.with("c1", channelState{PostThreadID: "p", LastCopiedMessageID: "5", PostedCount: 2})
	if len(st) != 0 {
		t.Errorf("with must not mutate")
	}
	if err := saveState(path, st2); err != nil {
		t.Fatal(err)
	}
	got, err := loadState(path)
	if err != nil || got["c1"].LastCopiedMessageID != "5" || got["c1"].PostedCount != 2 {
		t.Errorf("roundtrip: %+v %v", got, err)
	}
}

func TestDecideAction(t *testing.T) {
	st := archiveState{
		"inprog": {PostThreadID: "p1", LastCopiedMessageID: "9"},
		"done":   {PostThreadID: "p2", Done: true},
	}
	tests := []struct {
		id            string
		exists, force bool
		want          channelAction
	}{
		{"inprog", true, false, actionResume},
		{"done", true, false, actionSkipDone},
		{"new", false, false, actionCreate},
		{"new", true, false, actionSkipExists},
		{"new", true, true, actionCreate},
	}
	for _, tt := range tests {
		if got := decideAction(st, tt.id, tt.exists, tt.force); got != tt.want {
			t.Errorf("%s exists=%v force=%v: got %v want %v", tt.id, tt.exists, tt.force, got, tt.want)
		}
	}
}

func TestMessagesAfter(t *testing.T) {
	msgs := []srcMessage{{ID: "9"}, {ID: "10"}, {ID: "100"}}
	tests := []struct {
		last string
		want int
	}{{"", 3}, {"9", 2}, {"10", 1}, {"100", 0}}
	for _, tt := range tests {
		if got := messagesAfter(msgs, tt.last); len(got) != tt.want {
			t.Errorf("after %q: %d want %d", tt.last, len(got), tt.want)
		}
	}
}
