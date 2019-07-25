package retrievers

import (
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/kwhite17/Neighbors/pkg/assets"
)

type TemplateRetriever interface {
	RetrieveCreateEntityTemplate() (*template.Template, error)
	RetrieveSingleEntityTemplate() (*template.Template, error)
	RetrieveAllEntitiesTemplate() (*template.Template, error)
}

func RetrieveTemplate(templatePath string) (*template.Template, error) {
	fullPath := filepath.Join("assets", "templates", templatePath+".html")
	htmlBytes, err := assets.Asset(fullPath)

	if err != nil {
		return nil, fmt.Errorf("ERROR - Couldn't Retrieve Asset From Path: %s\n", templatePath)
	}

	return template.New(templatePath).Parse(string(htmlBytes))
}

func RetrieveMultiTemplate(paths ...string) (*template.Template, error) {
	var tpl *template.Template

	for i, t := range paths {
		fp := filepath.Join("assets", "templates", t+".html")
		hb, err := assets.Asset(fp)

		if err != nil {
			return nil, fmt.Errorf("ERROR - Couldn't retrieve asset from path: %s\n", t)
		}

		if i == 0 {
			tpl, err = template.New(t).Parse(string(hb))
		} else {
			tpl, err = tpl.Parse(string(hb))
		}

		if err != nil {
			return nil, fmt.Errorf("ERROR - Failed to parse template: %s\n", t)
		}
	}

	return tpl, nil
}
