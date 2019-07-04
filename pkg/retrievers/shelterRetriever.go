package retrievers

import (
	"html/template"
)

var createShelterTemplatePath = "shelters/new"
var getShelterTemplatePath = "shelters/shelter"
var getSheltersTemplatePath = "shelters/shelters"
var updateSheltersTemplatePath = "shelters/edit"

type ShelterRetriever struct {
	TemplateRetriever
}

func (sr ShelterRetriever) RetrieveCreateEntityTemplate() (*template.Template, error) {
	return RetrieveTemplate(createShelterTemplatePath)
}

func (sr ShelterRetriever) RetrieveSingleEntityTemplate() (*template.Template, error) {
	return RetrieveTemplate(getShelterTemplatePath)
}

func (sr ShelterRetriever) RetrieveAllEntitiesTemplate() (*template.Template, error) {
	return RetrieveTemplate(getSheltersTemplatePath)
}

func (sr ShelterRetriever) RetrieveEditEntityTemplate() (*template.Template, error) {
	return RetrieveTemplate(updateSheltersTemplatePath)
}
