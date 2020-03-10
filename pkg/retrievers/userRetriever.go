package retrievers

import (
	"html/template"
)

var createShelterTemplatePath = "users/new"
var getShelterTemplatePath = "users/shelter"
var getSheltersTemplatePath = "users/shelters"
var getSamaritanSummaryTemplatePath = "users/samaritanSummary"
var getShelterSummaryTemplatePath = "users/shelterSummary"
var updateSheltersTemplatePath = "users/edit"

type ShelterRetriever struct {
	TemplateRetriever
}

func (sr ShelterRetriever) RetrieveCreateEntityTemplate() (*template.Template, error) {
	return RetrieveMultiTemplate(layoutTemplatePath, createShelterTemplatePath)
}

func (sr ShelterRetriever) RetrieveSingleEntityTemplate() (*template.Template, error) {
	return RetrieveMultiTemplate(layoutTemplatePath, getShelterTemplatePath, getSamaritanSummaryTemplatePath, getShelterSummaryTemplatePath)
}

func (sr ShelterRetriever) RetrieveAllEntitiesTemplate() (*template.Template, error) {
	return RetrieveMultiTemplate(layoutTemplatePath, getSheltersTemplatePath)
}

func (sr ShelterRetriever) RetrieveEditEntityTemplate() (*template.Template, error) {
	return RetrieveMultiTemplate(layoutTemplatePath, updateSheltersTemplatePath)
}
