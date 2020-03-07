package retrievers

import (
	"fmt"
	"html/template"
)

var loginTemplatePath = "login/login"

type LoginRetriever struct {
	TemplateRetriever
}

func (lr LoginRetriever) RetrieveCreateEntityTemplate() (*template.Template, error) {
	return nil, fmt.Errorf("UnsupportedOperation")
}

func (lr LoginRetriever) RetrieveSingleEntityTemplate() (*template.Template, error) {
	return RetrieveMultiTemplate(layoutTemplatePath, loginTemplatePath)
}

func (lr LoginRetriever) RetrieveAllEntitiesTemplate() (*template.Template, error) {
	return nil, fmt.Errorf("UnsupportedOperation")
}

func (lr LoginRetriever) RetrieveEditEntityTemplate() (*template.Template, error) {
	return nil, fmt.Errorf("UnsupportedOperation")
}
