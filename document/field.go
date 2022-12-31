package document

import (
	"github.com/getumen/sakuin/fieldname"
	"github.com/getumen/sakuin/term"
)

type Field struct {
	name    fieldname.FieldName
	content []term.Term
}

func NewField(
	name fieldname.FieldName,
	content []term.Term,
) *Field {
	return &Field{
		name:    name,
		content: content,
	}
}

func (f Field) FieldName() fieldname.FieldName {
	return f.name
}

func (f Field) Content() []term.Term {
	return f.content
}

func (f Field) TermCount() int {
	return len(f.content)
}
