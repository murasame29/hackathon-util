package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type channelState struct {
	PostThreadID        string `json:"postThreadID"`
	LastCopiedMessageID string `json:"lastCopiedMessageID"`
	PostedCount         int    `json:"postedCount"`
	Done                bool   `json:"done"`
}

type archiveState map[string]channelState

func loadState(path string) (archiveState, error) {
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return archiveState{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read state %s: %w", path, err)
	}
	st := archiveState{}
	if err := json.Unmarshal(b, &st); err != nil {
		return nil, fmt.Errorf("parse state %s: %w", path, err)
	}
	return st, nil
}

func saveState(path string, st archiveState) error {
	b, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return fmt.Errorf("write state: %w", err)
	}
	return os.Rename(tmp, path)
}

// with returns a copy of the state with key set to v.
func (s archiveState) with(key string, v channelState) archiveState {
	out := make(archiveState, len(s)+1)
	for k, x := range s {
		out[k] = x
	}
	out[key] = v
	return out
}

type channelAction int

const (
	actionCreate channelAction = iota
	actionResume
	actionSkipDone
	actionSkipExists
)

// decideAction determines what to do for a source channel given state and existing forum posts.
func decideAction(st archiveState, sourceID string, postExists, force bool) channelAction {
	if cs, ok := st[sourceID]; ok && cs.PostThreadID != "" {
		if cs.Done {
			return actionSkipDone
		}
		return actionResume
	}
	if postExists && !force {
		return actionSkipExists
	}
	return actionCreate
}
