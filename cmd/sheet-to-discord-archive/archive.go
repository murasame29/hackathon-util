package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type target struct {
	SourceID       string
	Title          string
	ParentSourceID string // non-empty for thread posts
}

type config struct {
	DryRun    bool
	Force     bool
	MaxUpload int64
	StatePath string
}

type summary struct {
	Created, Resumed, SkippedExists, SkippedDone, Failed, Mismatch int
	MessagesPosted                                                 int
}

type archiver struct {
	cfg   config
	rd    *reader
	wh    *webhookClient
	httpc *http.Client
	state archiveState
	posts map[string]string // forum title -> thread id
	sum   summary
}

// postIDFor resolves the forum post thread ID for a source channel (state first, then forum index).
func (a *archiver) postIDFor(t target) string {
	if cs, ok := a.state[t.SourceID]; ok && cs.PostThreadID != "" {
		return cs.PostThreadID
	}
	return a.posts[t.Title]
}

func (a *archiver) run(t target) {
	msgs, err := a.rd.fetchMessages(t.SourceID)
	if err != nil {
		log.Printf("[ERROR] %s: %v", t.Title, err)
		a.sum.Failed++
		return
	}
	var postable []srcMessage
	var attCount int
	var attBytes int64
	for _, m := range msgs {
		if isPostable(m) {
			postable = append(postable, m)
			for _, at := range m.Attachments {
				attCount++
				attBytes += at.Size
			}
		}
	}
	_, exists := a.posts[t.Title]
	action := decideAction(a.state, t.SourceID, exists, a.cfg.Force)
	if a.cfg.DryRun {
		log.Printf("[DRY] %s (src %s): messages total %d / postable %d / skipped %d, attachments %d (%.1f KB), post exists=%v, action=%s",
			t.Title, t.SourceID, len(msgs), len(postable), len(msgs)-len(postable), attCount, float64(attBytes)/1024, exists, actionName(action))
		return
	}
	switch action {
	case actionSkipExists:
		log.Printf("[SKIP] already archived (post exists): %s", t.Title)
		a.sum.SkippedExists++
		return
	case actionSkipDone:
		log.Printf("[SKIP] already done (state): %s", t.Title)
		a.sum.SkippedDone++
		return
	}
	if err := a.copyChannel(t, action, postable); err != nil {
		log.Printf("[ERROR] %s: %v", t.Title, err)
		a.sum.Failed++
	}
}

func actionName(a channelAction) string {
	return [...]string{"create", "resume", "skip(done)", "skip(post exists)"}[a]
}

func (a *archiver) saveCS(id string, cs channelState) error {
	a.state = a.state.with(id, cs)
	return saveState(a.cfg.StatePath, a.state)
}

func (a *archiver) copyChannel(t target, action channelAction, postable []srcMessage) error {
	cs := a.state[t.SourceID]
	if action == actionCreate {
		header := fmt.Sprintf("from: <#%s>", t.SourceID)
		if t.ParentSourceID != "" {
			parentPost := a.postIDFor(target{SourceID: t.ParentSourceID, Title: parentTitle(t.Title)})
			if parentPost != "" {
				header += fmt.Sprintf("\nparent: <#%s>", parentPost)
			}
		}
		threadID, err := a.wh.execute(webhookPayload{Content: header, ThreadName: t.Title}, "", nil)
		if err != nil {
			return fmt.Errorf("create post: %w", err)
		}
		cs = channelState{PostThreadID: threadID}
		if err := a.saveCS(t.SourceID, cs); err != nil {
			return err
		}
		a.posts[t.Title] = threadID
		log.Printf("[OK] created post: %s (%s)", t.Title, threadID)
		a.sum.Created++
	} else {
		log.Printf("[RESUME] %s after %s", t.Title, cs.LastCopiedMessageID)
		a.sum.Resumed++
	}
	for _, m := range messagesAfter(postable, cs.LastCopiedMessageID) {
		if err := a.postMessage(cs.PostThreadID, m); err != nil {
			return fmt.Errorf("message %s: %w", m.ID, err)
		}
		cs.LastCopiedMessageID = m.ID
		cs.PostedCount++
		if err := a.saveCS(t.SourceID, cs); err != nil {
			return err
		}
		a.sum.MessagesPosted++
	}
	cs.Done = true
	if err := a.saveCS(t.SourceID, cs); err != nil {
		return err
	}
	if cs.PostedCount != len(postable) {
		log.Printf("[WARN] [VERIFY] %s: source %d / posted %d", t.Title, len(postable), cs.PostedCount)
		a.sum.Mismatch++
	} else {
		log.Printf("[VERIFY] %s: source %d / posted %d", t.Title, len(postable), cs.PostedCount)
	}
	return nil
}

func parentTitle(title string) string {
	i := strings.LastIndex(title, "/")
	if i < 0 {
		return title
	}
	return title[:i]
}

func (a *archiver) postMessage(threadID string, m srcMessage) error {
	upload, linked := splitUploads(m.Attachments, a.cfg.MaxUpload)
	var files []webhookFile
	for _, at := range upload {
		data, err := downloadFile(a.httpc, at.URL)
		if err != nil {
			log.Printf("[WARN] download %s failed (%v); linking instead", at.Filename, err)
			linked = append(linked, at)
			continue
		}
		files = append(files, webhookFile{Name: at.Filename, Data: data})
	}
	chunks := buildChunks(m, linked, a.rd.displayName)
	base := webhookPayload{
		Username:  sanitizeUsername(a.rd.displayName(m.Author)),
		AvatarURL: avatarURL(m.Author),
	}
	for i, c := range chunks {
		p := base
		p.Content = c
		var f []webhookFile
		if i == len(chunks)-1 {
			f = files
			p.Embeds = limitEmbeds(m.richEmbeds())
		}
		if _, err := a.wh.execute(p, threadID, f); err != nil {
			return err
		}
	}
	return nil
}

func limitEmbeds(e []json.RawMessage) []json.RawMessage {
	if len(e) > 10 {
		return e[:10]
	}
	return e
}
