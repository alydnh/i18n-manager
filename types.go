package i18n_manager

import (
	"gitlab.ptit365.com/utils"
	"i18n-manager/models"
	"strings"
)

type I18N struct {
	Items []*I18NItem `yaml:"i18n"`
}

func (i I18N) ToApiModel() []*models.LanguageAPIModel {
	return utils.Select(i.Items, func(item *I18NItem) *models.LanguageAPIModel { return item.ToApiModel() }).([]*models.LanguageAPIModel)
}

type I18NItem struct {
	Key      string            `yaml:"key"`
	Default  string            `yaml:"default"`
	Items    map[string]string `yaml:"items"`
	Deletion bool             `yaml:"deletion"`
}

func (i I18NItem) Match(language string, settled bool) (bool, string) {
	if strings.Compare(language, "all") == 0 {
		return true, ""
	}

	if text, exists := i.Items[language]; (!exists && settled) || (exists && !settled) {
		return false, text
	} else {
		return true, text
	}
}

func (i I18NItem) Clone() *I18NItem {
	item := &I18NItem{
		Key:     i.Key,
		Default: i.Default,
		Items:   make(map[string]string),
	}
	for key, value := range i.Items {
		item.Items[key] = value
	}

	return item
}

func (i I18NItem) ToApiModel() *models.LanguageAPIModel {
	model := &models.LanguageAPIModel{
		Default: i.Default,
		Items:   make([]*models.LanguageAPIModelItemsItems0, 0, len(i.Items)),
		Key:     i.Key,
	}

	for language, text := range i.Items {
		model.Items = append(model.Items, &models.LanguageAPIModelItemsItems0{
			Location: language,
			Text:     text,
		})
	}

	return model
}
