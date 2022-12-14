package store

type Book struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Authors string `json:"authors"`
	Press   string `json:"press"`
}

type Store interface {
	Create(*Book) error
	Update(*Book) error
	Get(string) (Book, error)
	GetAll() ([]Book, error)
	Delete(string) error
}
