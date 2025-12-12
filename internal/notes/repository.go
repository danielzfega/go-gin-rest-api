package notes

type NoteRepository interface {
	Create(n Note) (Note, error)
	GetByID(id int64) (Note, error)
	List() ([]Note, error)
	Update(id int64, title *string, content *string) (Note, error)
	Delete(id int64) error
}
