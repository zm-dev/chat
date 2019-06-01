package model

type Page struct {
	Size    int          `json:"size"`
	Pages   int          `json:"pages"`
	Total   int          `json:"total"`
	Current int          `json:"current"`
	Records *interface{} `json:"records"`
}

func (p *Page) Offset() int {
	if p.Current > 0 {
		return (p.Current - 1) * p.Size
	}
	return 0
}

func (p *Page) SetPages() {
	if p.Size == 0 {
		p.Pages = 0
	}
	p.Pages = p.Total / p.Size
	if p.Total%p.Size != 0 {
		p.Pages = p.Pages + 1
	}
}

func (p *Page) setRecords(rs *interface{}) {
	p.Records = rs
}

func (p *Page) setTotal(t int) {
	p.Total = t
}
