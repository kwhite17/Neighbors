package items

import (
	"html/template"

	"github.com/kwhite17/Neighbors/pkg/retriever"
)

var createItemTemplatePath = "items/new"
var getItemTemplatePath = "items/item"
var getItemsTemplatePath = "items/items"
var updateItemsTemplatePath = "items/edit"

type ItemRetriever struct {
	retriever.TemplateRetriever
}

func (ir ItemRetriever) RetrieveCreateEntityTemplate() (*template.Template, error) {
	return retriever.RetrieveTemplate(createItemTemplatePath)
}

func (ir ItemRetriever) RetrieveSingleEntityTemplate() (*template.Template, error) {
	return retriever.RetrieveTemplate(getItemTemplatePath)
}

func (ir ItemRetriever) RetrieveAllEntitiesTemplate() (*template.Template, error) {
	return retriever.RetrieveTemplate(getItemsTemplatePath)
}

func (ir ItemRetriever) RetrieveEditEntityTemplate() (*template.Template, error) {
	return retriever.RetrieveTemplate(updateItemsTemplatePath)
}
