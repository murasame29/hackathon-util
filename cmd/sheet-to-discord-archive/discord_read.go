package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/bwmarrin/discordgo"
)

type reader struct {
	dg        *discordgo.Session
	guildID   string
	nameCache map[string]string
}

// fetchMessages returns all messages of a channel in ascending order.
func (r *reader) fetchMessages(channelID string) ([]srcMessage, error) {
	var all []srcMessage
	before := ""
	for {
		u := discordgo.EndpointChannelMessages(channelID) + "?limit=100"
		if before != "" {
			u += "&before=" + before
		}
		body, err := r.dg.RequestWithBucketID("GET", u, nil, discordgo.EndpointChannelMessages(channelID))
		if err != nil {
			return nil, fmt.Errorf("fetch messages %s: %w", channelID, err)
		}
		var page []srcMessage
		if err := json.Unmarshal(body, &page); err != nil {
			return nil, fmt.Errorf("decode messages %s: %w", channelID, err)
		}
		all = append(all, page...)
		if len(page) < 100 {
			break
		}
		before = page[len(page)-1].ID
	}
	sort.Slice(all, func(i, j int) bool { return snowflakeLess(all[i].ID, all[j].ID) })
	return all, nil
}

func (r *reader) displayName(u srcUser) string {
	if n, ok := r.nameCache[u.ID]; ok {
		return n
	}
	nick := ""
	if m, err := r.dg.GuildMember(r.guildID, u.ID); err == nil && m != nil {
		nick = m.Nick
	}
	n := pickDisplayName(nick, u.GlobalName, u.Username)
	r.nameCache[u.ID] = n
	return n
}

// childThreads returns active + archived public + archived private threads of a channel.
func (r *reader) childThreads(channelID string, active []*discordgo.Channel) []*discordgo.Channel {
	seen := map[string]bool{}
	var out []*discordgo.Channel
	add := func(ts []*discordgo.Channel) {
		for _, t := range ts {
			if t.ParentID == channelID && !seen[t.ID] {
				seen[t.ID] = true
				out = append(out, t)
			}
		}
	}
	add(active)
	add(r.archivedThreads(channelID, false))
	add(r.archivedThreads(channelID, true))
	sort.Slice(out, func(i, j int) bool { return snowflakeLess(out[i].ID, out[j].ID) })
	return out
}

func (r *reader) archivedThreads(channelID string, private bool) []*discordgo.Channel {
	var out []*discordgo.Channel
	var before *time.Time
	for {
		var list *discordgo.ThreadsList
		var err error
		if private {
			list, err = r.dg.ThreadsPrivateArchived(channelID, before, 100)
		} else {
			list, err = r.dg.ThreadsArchived(channelID, before, 100)
		}
		if err != nil {
			kind := "public"
			if private {
				kind = "private"
			}
			log.Printf("[WARN] cannot list archived %s threads of %s: %v", kind, channelID, err)
			return out
		}
		out = append(out, list.Threads...)
		if !list.HasMore || len(list.Threads) == 0 {
			return out
		}
		last := list.Threads[len(list.Threads)-1]
		if last.ThreadMetadata == nil {
			return out
		}
		ts := last.ThreadMetadata.ArchiveTimestamp
		before = &ts
	}
}

// forumPostIndex maps post title -> thread ID for all active and archived public posts in the forum.
func (r *reader) forumPostIndex(forumID string, active []*discordgo.Channel) map[string]string {
	idx := map[string]string{}
	for _, t := range active {
		if t.ParentID == forumID {
			idx[t.Name] = t.ID
		}
	}
	for _, t := range r.archivedThreads(forumID, false) {
		idx[t.Name] = t.ID
	}
	return idx
}
