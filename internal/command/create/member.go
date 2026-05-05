package create

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/sync/errgroup"
)

// lookupUserIDs resolves a list of usernames to Discord user IDs concurrently.
// Each username is searched via GuildMembersSearch and only exact (case-insensitive)
// matches are included in the result. Usernames that are not found are silently omitted.
func lookupUserIDs(dg *discordgo.Session, guildID string, usernames []string) (map[string]string, error) {
	var (
		mu     sync.Mutex
		result = make(map[string]string, len(usernames))
	)

	g := new(errgroup.Group)
	for _, raw := range usernames {
		username := strings.ToLower(strings.TrimSpace(raw))
		if username == "" {
			continue
		}

		g.Go(func() error {
			members, err := dg.GuildMembersSearch(guildID, username, 10)
			if err != nil {
				return fmt.Errorf("GuildMembersSearch %q: %w", username, err)
			}
			for _, m := range members {
				if m.User == nil {
					continue
				}
				if strings.ToLower(m.User.Username) == username {
					mu.Lock()
					result[username] = m.User.ID
					mu.Unlock()
					break
				}
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return result, nil
}
