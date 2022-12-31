package document

type Document struct {
	id     int64
	fields []*Field
}

func NewDocument(
	id int64,
	fields []*Field,
) *Document {
	return &Document{
		id:     id,
		fields: fields,
	}
}

func (d Document) ID() int64 {
	return d.id
}

func (d Document) Fields() []*Field {
	return d.fields
}
