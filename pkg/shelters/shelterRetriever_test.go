package shelters

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

var shelterRetriever = &ShelterRetriever{}

func TestRenderSingleShelterTemplate(t *testing.T) {
	testArray := make([]byte, 0)
	testBuffer := bytes.NewBuffer(testArray)
	testShelter := generateShelter()
	tmpl, err := shelterRetriever.RetrieveSingleEntityTemplate()
	if err != nil {
		t.Fatal(err)
	}

	tmpl.Execute(testBuffer, testShelter)
	htmlBytes, err := ioutil.ReadAll(testBuffer)
	if err != nil {
		t.Error(err)
	}

	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "strong") || !strings.Contains(htmlStr, testShelter.Name) {
		t.Errorf("GetSingleShelter Failure - Expected html to contain 'strong' or correct ID, Actual: %s\n", htmlStr)
	}
}

func TestRenderCreateShelterTemplate(t *testing.T) {
	testArray := make([]byte, 0)
	testBuffer := bytes.NewBuffer(testArray)
	tmpl, err := shelterRetriever.RetrieveCreateEntityTemplate()
	if err != nil {
		t.Fatal(err)
	}

	tmpl.Execute(testBuffer, nil)
	htmlBytes, err := ioutil.ReadAll(testBuffer)
	if err != nil {
		t.Error(err)
	}

	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "form") || !strings.Contains(htmlStr, "Shelter Name") {
		t.Errorf("CreateSingleShelter Failure - Expected html to contain 'form' or 'Shelter Name', Actual: %s\n", htmlStr)
	}
}

func TestRenderAllSheltersTemplate(t *testing.T) {
	testArray := make([]byte, 0)
	testBuffer := bytes.NewBuffer(testArray)
	testShelter := generateShelter()
	tmpl, err := shelterRetriever.RetrieveAllEntitiesTemplate()
	if err != nil {
		t.Fatal(err)
	}

	tmpl.Execute(testBuffer, []*Shelter{testShelter})
	htmlBytes, err := ioutil.ReadAll(testBuffer)
	if err != nil {
		t.Error(err)
	}

	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "table") || !strings.Contains(htmlStr, testShelter.City) {
		t.Errorf("GetAllShelters Failure - Expected html to contain 'strong' or correct city: %s, Actual: %s\n", testShelter.City, htmlStr)
	}
}
