package datamodels

type Page struct {
	CurrentPage int
	PageSize    int
	PageCount   int
}

func PageIndex(page *Page) int {
	if page.CurrentPage == 0 {
		page.CurrentPage = 1
	}
	return (page.CurrentPage - 1) * page.PageSize
}
