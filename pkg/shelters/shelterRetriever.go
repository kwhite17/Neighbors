package shelters

import (
	"html/template"

	"github.com/kwhite17/Neighbors/pkg/retriever"
)

var createShelterTemplatePath = "shelters/new"
var getShelterTemplatePath = "shelters/shelter"
var getSheltersTemplatePath = "shelters/shelters"
var updateSheltersTemplatePath = "shelters/edit"

type ShelterRetriever struct {
	retriever.TemplateRetriever
}

func (sr ShelterRetriever) RetrieveCreateEntityTemplate() (*template.Template, error) {
	return retriever.RetrieveTemplate(createShelterTemplatePath)
}

func (sr ShelterRetriever) RetrieveSingleEntityTemplate() (*template.Template, error) {
	return retriever.RetrieveTemplate(getShelterTemplatePath)
}

func (sr ShelterRetriever) RetrieveAllEntitiesTemplate() (*template.Template, error) {
	return retriever.RetrieveTemplate(getSheltersTemplatePath)
}

func (sr ShelterRetriever) RetrieveEditEntityTemplate() (*template.Template, error) {
	return retriever.RetrieveTemplate(updateSheltersTemplatePath)
}
