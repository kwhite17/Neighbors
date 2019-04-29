package retriever

import (
	"fmt"
	"html/template"
	"path/filepath"
)

type TemplateRetriever interface {
	RetrieveCreateEntityTemplate() (*template.Template, error)
	RetrieveSingleEntityTemplate() (*template.Template, error)
	RetrieveAllEntitiesTemplate() (*template.Template, error)
}

func RetrieveTemplate(templatePath string) (*template.Template, error) {
	fullPath := filepath.Join("templates", templatePath+".html")
	htmlBytes, err := Asset(fullPath)
	if err != nil {
		return nil, fmt.Errorf("ERROR - Couldn't Retrieve Asset From Path: %s\n", templatePath)
	}
	return template.New(templatePath).Parse(string(htmlBytes))
}
