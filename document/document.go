package document

type Document struct {
	id     uint64
	fields []*Field
}

func NewDocument(
	id uint64,
	fields []*Field,
) *Document {
	return &Document{
		id:     id,
		fields: fields,
	}
}

func (d Document) ID() uint64 {
	return d.id
}

func (d Document) Fields() []*Field {
	return d.fields
}
