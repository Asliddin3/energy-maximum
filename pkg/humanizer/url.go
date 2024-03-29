package humanizer

import "github.com/gosimple/slug"

type Url interface {
	Regenerate(text string) string
}

type ManagerHumanizer struct {
	Sub map[string]string
}

func NewHumanizer(sub map[string]string) *ManagerHumanizer {
	return &ManagerHumanizer{
		Sub: sub,
	}
}

func (h *ManagerHumanizer) Regenerate(text string) string {
	slug.CustomSub = h.Sub
	return slug.Make(text)
}
