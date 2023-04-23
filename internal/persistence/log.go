package persistence

import "errors"

type LogEntry struct {
	Command []byte
	Term    int
}

type LogStore interface {
	Get(int) (*LogEntry, error)
	Save(int, *LogEntry) error
	DeleteFrom(int) error
}

type InMemoryLogStore struct {
	entries []*LogEntry
}

func (i *InMemoryLogStore) Get(idx int) (*LogEntry, error) {
	if idx >= len(i.entries) {
		return nil, errors.New("no such entry")
	}
	return i.entries[idx], nil
}

func (i *InMemoryLogStore) Append(entry *LogEntry) error {
	i.entries = append(i.entries, entry)
	return nil
}

func (i *InMemoryLogStore) DeleteFrom(idx int) error {
	i.entries = i.entries[:idx]
	return nil
}

func NewMemoryLogStore() *InMemoryLogStore {
	return &InMemoryLogStore{entries: []*LogEntry{}}
}
