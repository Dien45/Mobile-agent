package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/Dien45/Mobile-agent/internal/core"
)

type state struct {
	Sessions  map[string]core.Session         `json:"sessions"`
	Messages  map[string][]core.Message       `json:"messages"`
	Providers map[string]core.ProviderProfile `json:"providers"`
}

type Store struct {
	mu   sync.RWMutex
	path string
	data state
}

func Open(path string) (*Store, error) {
	s := &Store{path: path, data: state{
		Sessions: map[string]core.Session{}, Messages: map[string][]core.Message{}, Providers: map[string]core.ProviderProfile{},
	}}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil { return nil, err }
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) { return s, s.persistLocked() }
		return nil, err
	}
	if len(b) > 0 {
		if err := json.Unmarshal(b, &s.data); err != nil { return nil, fmt.Errorf("decode state: %w", err) }
	}
	if s.data.Sessions == nil { s.data.Sessions = map[string]core.Session{} }
	if s.data.Messages == nil { s.data.Messages = map[string][]core.Message{} }
	if s.data.Providers == nil { s.data.Providers = map[string]core.ProviderProfile{} }
	return s, nil
}

func (s *Store) persistLocked() error {
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil { return err }
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil { return err }
	backup := s.path + ".bak"
	_ = os.Remove(backup)
	if _, err := os.Stat(s.path); err == nil {
		if err := os.Rename(s.path, backup); err != nil { _ = os.Remove(tmp); return err }
	}
	if err := os.Rename(tmp, s.path); err != nil { _ = os.Rename(backup, s.path); return err }
	_ = os.Remove(backup)
	return nil
}

func (s *Store) ListSessions(includeDeleted bool) []core.Session {
	s.mu.RLock(); defer s.mu.RUnlock()
	out := make([]core.Session, 0, len(s.data.Sessions))
	for _, v := range s.data.Sessions {
		if !includeDeleted && v.DeletedAt != nil { continue }
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].UpdatedAt.After(out[j].UpdatedAt) })
	return out
}

func (s *Store) GetSession(id string) (core.Session, bool) {
	s.mu.RLock(); defer s.mu.RUnlock(); v, ok := s.data.Sessions[id]; return v, ok
}

func (s *Store) SaveSession(v core.Session) error {
	s.mu.Lock(); defer s.mu.Unlock(); s.data.Sessions[v.ID] = v; return s.persistLocked()
}

func (s *Store) RenameSession(id, title string) error {
	s.mu.Lock(); defer s.mu.Unlock()
	v, ok := s.data.Sessions[id]; if !ok { return os.ErrNotExist }
	v.Title, v.UpdatedAt = title, time.Now().UTC(); s.data.Sessions[id] = v; return s.persistLocked()
}

func (s *Store) DeleteSession(id string, purge bool) error {
	s.mu.Lock(); defer s.mu.Unlock()
	v, ok := s.data.Sessions[id]; if !ok { return os.ErrNotExist }
	if purge { delete(s.data.Sessions, id); delete(s.data.Messages, id) } else { now := time.Now().UTC(); v.DeletedAt = &now; v.UpdatedAt = now; s.data.Sessions[id] = v }
	return s.persistLocked()
}

func (s *Store) RestoreSession(id string) error {
	s.mu.Lock(); defer s.mu.Unlock(); v, ok := s.data.Sessions[id]; if !ok { return os.ErrNotExist }
	v.DeletedAt = nil; v.UpdatedAt = time.Now().UTC(); s.data.Sessions[id] = v; return s.persistLocked()
}

func (s *Store) Messages(id string) []core.Message {
	s.mu.RLock(); defer s.mu.RUnlock(); src := s.data.Messages[id]; out := make([]core.Message, len(src)); copy(out, src); return out
}

func (s *Store) AddMessage(v core.Message) error {
	s.mu.Lock(); defer s.mu.Unlock(); s.data.Messages[v.SessionID] = append(s.data.Messages[v.SessionID], v)
	if session, ok := s.data.Sessions[v.SessionID]; ok { session.UpdatedAt = time.Now().UTC(); s.data.Sessions[v.SessionID] = session }
	return s.persistLocked()
}

func (s *Store) ListProviders() []core.ProviderProfile {
	s.mu.RLock(); defer s.mu.RUnlock(); out := make([]core.ProviderProfile, 0, len(s.data.Providers))
	for _, p := range s.data.Providers { p.APIKey = ""; out = append(out, p) }
	sort.Slice(out, func(i,j int) bool { return out[i].Name < out[j].Name }); return out
}

func (s *Store) SaveProvider(v core.ProviderProfile) error {
	s.mu.Lock(); defer s.mu.Unlock()
	if old, ok := s.data.Providers[v.ID]; ok && v.APIKey == "" { v.APIKey = old.APIKey }
	s.data.Providers[v.ID] = v; return s.persistLocked()
}

func (s *Store) Provider(id string) (core.ProviderProfile, bool) { s.mu.RLock(); defer s.mu.RUnlock(); v, ok := s.data.Providers[id]; return v, ok }

func (s *Store) DeleteProvider(id string) error {
	s.mu.Lock(); defer s.mu.Unlock()
	if _, ok := s.data.Providers[id]; !ok { return os.ErrNotExist }
	delete(s.data.Providers, id)
	return s.persistLocked()
}
