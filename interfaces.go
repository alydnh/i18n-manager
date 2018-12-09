package i18n_manager

type Store interface {
	Save(items []*I18NItem) error
	LoadAll() ([]*I18NItem, error)
	Shutdown() error
}
