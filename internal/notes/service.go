package notes

import "errors"

type NoteService struct {
    repo NoteRepository
}

func NewNoteService(repo NoteRepository) *NoteService {
    return &NoteService{repo: repo}
}

func (s *NoteService) Create(title string, content string) (Note, error) {
    if len(title) == 0 {
        return Note{}, errors.New("validation_title_required")
    }
    n := Note{Title: title, Content: content}
    return s.repo.Create(n)
}

func (s *NoteService) Get(id int64) (Note, error) {
    return s.repo.GetByID(id)
}

func (s *NoteService) List() ([]Note, error) {
    return s.repo.List()
}

func (s *NoteService) Update(id int64, title *string, content *string) (Note, error) {
    if title != nil && len(*title) == 0 {
        return Note{}, errors.New("validation_title_required")
    }
    return s.repo.Update(id, title, content)
}

func (s *NoteService) Delete(id int64) error {
    return s.repo.Delete(id)
}

