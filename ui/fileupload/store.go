package fileupload

import "sync"

// FileMeta holds metadata about an uploaded file.
type FileMeta struct {
	ID       string
	Name     string
	Size     int64
	MimeType string
}

// Store is the interface for uploaded file metadata storage.
type Store interface {
	Add(key string, meta FileMeta)
	Remove(key, fileID string)
	List(key string) []FileMeta
	Clear(key string)
}

// MemoryStore is a thread-safe in-memory implementation of Store,
// keyed by "sessionID:componentID".
type MemoryStore struct {
	mu    sync.Mutex
	files map[string][]FileMeta
}

// compile-time check
var _ Store = (*MemoryStore)(nil)

// NewStore creates a new empty in-memory file metadata store.
func NewStore() *MemoryStore {
	return &MemoryStore{files: make(map[string][]FileMeta)}
}

// Add appends a file to the store under the given key.
func (s *MemoryStore) Add(key string, meta FileMeta) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.files[key] = append(s.files[key], meta)
}

// Remove deletes a file by ID from the store.
func (s *MemoryStore) Remove(key, fileID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	files := s.files[key]
	for i, f := range files {
		if f.ID == fileID {
			s.files[key] = append(files[:i], files[i+1:]...)
			return
		}
	}
}

// List returns all files stored under the given key.
func (s *MemoryStore) List(key string) []FileMeta {
	s.mu.Lock()
	defer s.mu.Unlock()
	dst := make([]FileMeta, len(s.files[key]))
	copy(dst, s.files[key])
	return dst
}

// Clear removes all files under the given key.
func (s *MemoryStore) Clear(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.files, key)
}
