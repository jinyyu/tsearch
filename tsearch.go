package tsearch

type TextSearch struct {
	separator Separator
	storage   Storage
}

func NewTextSearch(separator Separator, storage Storage) *TextSearch {
	return &TextSearch{
		separator: separator,
		storage:   storage,
	}
}

func (t *TextSearch) UpdateText(id uint32, text string) (err error) {
	return err
}
