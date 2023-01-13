package document

import (
	"github.com/getumen/sakuin/analysis/token"
	"github.com/getumen/sakuin/fieldname"
)

type Field struct {
	name    fieldname.FieldName
	content token.TokenStream
}

func NewField(
	name fieldname.FieldName,
	content token.TokenStream,
) *Field {
	return &Field{
		name:    name,
		content: content,
	}
}

func (f Field) FieldName() fieldname.FieldName {
	return f.name
}

func (f Field) Content() token.TokenStream {
	return f.content
}

func (f Field) TermCount() int {
	return len(f.content)
}
