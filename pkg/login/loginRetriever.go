package login

import (
	"fmt"
	"html/template"

	"github.com/kwhite17/Neighbors/pkg/retriever"
)

var loginTemplatePath = "login/login"

type LoginRetriever struct {
	retriever.TemplateRetriever
}

func (lr LoginRetriever) RetrieveCreateEntityTemplate() (*template.Template, error) {
	return nil, fmt.Errorf("UnsupportedOperation")
}

func (lr LoginRetriever) RetrieveSingleEntityTemplate() (*template.Template, error) {
	return retriever.RetrieveTemplate(loginTemplatePath)
}

func (lr LoginRetriever) RetrieveAllEntitiesTemplate() (*template.Template, error) {
	return nil, fmt.Errorf("UnsupportedOperation")
}

func (lr LoginRetriever) RetrieveEditEntityTemplate() (*template.Template, error) {
	return nil, fmt.Errorf("UnsupportedOperation")
}
