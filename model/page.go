package model

type Page struct {
	Size    int32       `json:"size"`
	Pages   int32       `json:"pages"`
	Total   int32       `json:"total"`
	Current int32       `json:"current"`
	Records interface{} `json:"records"`
}

func (p *Page) Offset() int32 {
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
