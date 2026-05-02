package mealdomain

import "strings"

type GarnishOption struct {
	ID         string
	TemplateID string
	Name       string
}

func NewGarnishOption(templateID, name string) (*GarnishOption, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrEmptyGarnishName
	}
	return &GarnishOption{
		TemplateID: templateID,
		Name:       name,
	}, nil
}
