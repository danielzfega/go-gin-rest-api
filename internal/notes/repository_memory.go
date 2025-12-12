package notes

import (
    "errors"
    "sync"
    "time"
)

type MemoryRepository struct {
    mu     sync.RWMutex
    nextID int64
    items  map[int64]*Note
}

func NewMemoryRepository() *MemoryRepository {
    return &MemoryRepository{nextID: 1, items: make(map[int64]*Note)}
}

func (r *MemoryRepository) Create(n Note) (Note, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    id := r.nextID
    r.nextID++
    now := time.Now().UTC()
    n.ID = id
    n.CreatedAt = now
    n.UpdatedAt = now
    v := n
    r.items[id] = &v
    return v, nil
}

func (r *MemoryRepository) GetByID(id int64) (Note, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    v, ok := r.items[id]
    if !ok {
        return Note{}, errors.New("not_found")
    }
    return *v, nil
}

func (r *MemoryRepository) List() ([]Note, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    out := make([]Note, 0, len(r.items))
    for _, v := range r.items {
        out = append(out, *v)
    }
    return out, nil
}

func (r *MemoryRepository) Update(id int64, title *string, content *string) (Note, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    v, ok := r.items[id]
    if !ok {
        return Note{}, errors.New("not_found")
    }
    if title != nil {
        v.Title = *title
    }
    if content != nil {
        v.Content = *content
    }
    v.UpdatedAt = time.Now().UTC()
    return *v, nil
}

func (r *MemoryRepository) Delete(id int64) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.items[id]; !ok {
        return errors.New("not_found")
    }
    delete(r.items, id)
    return nil
}

