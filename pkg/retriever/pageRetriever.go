package retriever

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
