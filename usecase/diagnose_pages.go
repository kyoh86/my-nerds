package usecase

import "github.com/kyoh86/my-nerds/model"

type PageDiagnostic struct {
	Index     int
	Landscape bool
}

func DiagnosePages(pages []model.Page) (diagnostics []PageDiagnostic, _ error) {
	for i, page := range pages {
		var d *PageDiagnostic
		prepare := func() {
			d = &PageDiagnostic{Index: i}
		}
		if page.Config.Width > page.Config.Height {
			prepare()
		}
		if d != nil {
			diagnostics = append(diagnostics, *d)
		}
	}
	return
}
