package queue

type Page struct {
	index int
	data  []byte
}

func NewPage() *Page {
	return &Page{
		index: 0,
		data:  nil,
	}
}

func (p *Page) GetIndex() int {
	return p.index
}

func (p *Page) GetData() []byte {
	return p.data
}

func (p *Page) SetData(index int, data []byte) {
	p.index = index
	p.data = data
}
