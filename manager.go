package i18n_manager

import (
	"gitlab.ptit365.com/utils"
	"strings"
	"sync"
)

const languagePlaceHolder = "<ENTER_TEXT_HERE>"

func CreateManager(store Store) (*Manager, error) {
	manager := &Manager{
		store: store,
		items: make(map[string]*I18NItem),
		lock:  &sync.RWMutex{},
	}

	if items, err := manager.store.LoadAll(); nil == err {
		for _, item := range items {
			manager.items[item.Key] = item
		}

		return manager, nil
	} else {
		return nil, err
	}
}

type Manager struct {
	store Store
	items map[string]*I18NItem
	lock  *sync.RWMutex
}

func (m *Manager) SaveOrUpdate(keys []string) (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	items := make([]*I18NItem, 0, len(keys))
	for _, key := range keys {
		if _, ok := m.items[key]; !ok {
			item := &I18NItem{
				Key:     key,
				Default: "",
				Items:   make(map[string]string),
			}
			items = append(items, item)
		}
	}

	if utils.EmptyArray(items) {
		return nil
	}

	if err = m.store.Save(items); nil == err {
		for _, item := range items {
			m.items[item.Key] = item
		}
	}

	return
}

func (m *Manager) Query(language, status string, languageTemplates []string) *I18N {
	m.lock.RLock()
	defer m.lock.RUnlock()

	language = strings.ToLower(language)
	isAllLanguage := strings.Compare(language, "all") == 0

	status = strings.ToLower(status)
	mustUnSet := strings.Compare(status, "unset") == 0
	mustSettled := !mustUnSet && strings.Compare(status, "settled") == 0
	templateMap := make(map[string]bool)
	hasTemplates := !utils.EmptyArray(languageTemplates)
	if hasTemplates {
		for _, language := range languageTemplates {
			templateMap[language] = true
		}
	}

	items := make([]*I18NItem, 0, len(m.items))
	for _, item := range m.items {
		match, text := item.Match(language, mustSettled)
		if !isAllLanguage && !match {
			continue
		}

		textIsEmptyOrWhiteSpace := utils.EmptyOrWhiteSpace(text)
		var newItem *I18NItem = nil
		if isAllLanguage {
			newItem = item.Clone()
		} else if mustUnSet && textIsEmptyOrWhiteSpace {
			newItem = &I18NItem{
				Key:     item.Key,
				Default: item.Default,
				Items:   map[string]string{language: languagePlaceHolder},
			}
		} else if mustSettled && !textIsEmptyOrWhiteSpace {
			newItem = &I18NItem{
				Key:     item.Key,
				Default: item.Default,
				Items:   map[string]string{language: text},
			}
		} else {
			newItem = &I18NItem{
				Key:     item.Key,
				Default: item.Default,
				Items:   map[string]string{language: text},
			}

			if utils.EmptyOrWhiteSpace(newItem.Items[language]) {
				newItem.Items[language] = languagePlaceHolder
			}
		}

		if nil != newItem {
			if hasTemplates {
				for language := range templateMap {
					if text, exists := newItem.Items[language]; !exists || utils.EmptyOrWhiteSpace(text) {
						newItem.Items[language] = languagePlaceHolder
					}
				}
				for language := range newItem.Items {
					if _, exists := templateMap[language]; !exists {
						delete(newItem.Items, language)
					}
				}
			}

			if utils.EmptyOrWhiteSpace(newItem.Default) {
				newItem.Default = languagePlaceHolder
			}

			items = append(items, newItem)
		}
	}

	return &I18N{
		Items: items,
	}
}
