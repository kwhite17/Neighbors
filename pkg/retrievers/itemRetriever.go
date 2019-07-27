package retrievers

import (
	"html/template"
)

var createItemTemplatePath = "items/new"
var getItemTemplatePath = "items/item"
var getItemsTemplatePath = "items/items"
var updateItemsTemplatePath = "items/edit"
var layoutTemplatePath = "home/layout"

type ItemRetriever struct {
	TemplateRetriever
}

func (ir ItemRetriever) RetrieveCreateEntityTemplate() (*template.Template, error) {
	return RetrieveMultiTemplate(layoutTemplatePath, createItemTemplatePath)
}

func (ir ItemRetriever) RetrieveSingleEntityTemplate() (*template.Template, error) {
	return RetrieveMultiTemplate(layoutTemplatePath, getItemTemplatePath)
}

func (ir ItemRetriever) RetrieveAllEntitiesTemplate() (*template.Template, error) {
	return RetrieveMultiTemplate(layoutTemplatePath, getItemsTemplatePath)
}

func (ir ItemRetriever) RetrieveEditEntityTemplate() (*template.Template, error) {
	return RetrieveMultiTemplate(layoutTemplatePath, updateItemsTemplatePath)
}
