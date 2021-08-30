package cmn

type LC struct {
	Lang string
}

func NewLC(lang string) *LC {
	return &LC{Lang: lang}
}

func (l *LC) T(text string) string {
	return text
}
