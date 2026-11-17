package main

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	maxTitleRunes    = 100
	maxContentRunes  = 2000
	maxUsernameRunes = 80
	replySnippetLen  = 40
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

// buildTitle joins path parts with "/" and truncates to Discord's 100-char thread name limit.
func buildTitle(parts ...string) string {
	return truncateRunes(strings.Join(parts, "/"), maxTitleRunes)
}

func truncateRunes(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	return string([]rune(s)[:n])
}

// pickDisplayName returns member nick > global name > username.
func pickDisplayName(nick, globalName, username string) string {
	for _, v := range []string{nick, globalName, username} {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return "unknown"
}

// sanitizeUsername makes a name acceptable as a webhook username
// (no "discord"/"clyde", 1-80 chars).
func sanitizeUsername(name string) string {
	out := replaceFold(name, "discord", "disc0rd")
	out = replaceFold(out, "clyde", "clyd3")
	out = strings.TrimSpace(out)
	if out == "" {
		return "unknown"
	}
	return truncateRunes(out, maxUsernameRunes)
}

func replaceFold(s, word, repl string) string {
	lower := strings.ToLower(s)
	var b strings.Builder
	i := 0
	for {
		j := strings.Index(lower[i:], word)
		if j < 0 {
			b.WriteString(s[i:])
			return b.String()
		}
		b.WriteString(s[i : i+j])
		b.WriteString(repl)
		i += j + len(word)
	}
}

func timestampSuffix(t time.Time) string {
	return "-# " + t.In(jst).Format("2006-01-02 15:04")
}

func replyPrefix(authorName, content string) string {
	snippet := strings.Join(strings.Fields(content), " ")
	snippet = truncateRunes(snippet, replySnippetLen)
	return fmt.Sprintf("-# ↪ %sへの返信：%s", authorName, snippet)
}

// splitContent splits s into chunks of at most max runes, preferring newline boundaries.
func splitContent(s string, max int) []string {
	if s == "" {
		return nil
	}
	var chunks []string
	rest := []rune(s)
	for len(rest) > max {
		cut := -1
		for i := max; i > 0; i-- {
			if rest[i-1] == '\n' {
				cut = i
				break
			}
		}
		if cut <= 0 {
			cut = max
		}
		chunks = append(chunks, strings.TrimRight(string(rest[:cut]), "\n"))
		rest = rest[cut:]
	}
	if len(rest) > 0 {
		chunks = append(chunks, string(rest))
	}
	return chunks
}
