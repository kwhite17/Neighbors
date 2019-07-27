package retrievers

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/kwhite17/Neighbors/pkg/managers"
)

var testCity = "testCity"
var testCountry = "testCountry"
var testName = "testName"
var testPostalCode = "testPostalCode"
var testState = "testState"
var testStreet = "testStreet"

var shelterRetriever = &ShelterRetriever{}

func TestRenderSingleShelterTemplate(t *testing.T) {
	testArray := make([]byte, 0)
	testBuffer := bytes.NewBuffer(testArray)
	testShelter := generateShelter()
	tmpl, err := shelterRetriever.RetrieveSingleEntityTemplate()

	if err != nil {
		t.Fatal(err)
	}

	tmpl.Execute(testBuffer, map[string]interface{}{
		"Shelter": testShelter,
		"ShelterSession": nil,
	})
	htmlBytes, err := ioutil.ReadAll(testBuffer)

	if err != nil {
		t.Error(err)
	}

	htmlStr := string(htmlBytes)

	if !strings.Contains(htmlStr, "<h6 class=\"card-subtitle text-muted\">testStreet, testCity, testState,") || !strings.Contains(htmlStr, testShelter.Name) {
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

	tmpl.Execute(testBuffer, map[string]interface{}{
		"ShelterSession": nil,
	})
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

	tmpl.Execute(testBuffer, map[string]interface{}{
		"Shelters": []*managers.Shelter{testShelter},
		"ShelterSession": nil,
	})
	htmlBytes, err := ioutil.ReadAll(testBuffer)

	if err != nil {
		t.Error(err)
	}

	htmlStr := string(htmlBytes)

	if !strings.Contains(htmlStr, "table") || !strings.Contains(htmlStr, testShelter.City) {
		t.Errorf("GetAllShelters Failure - Expected html to contain 'strong' or correct city: %s, Actual: %s\n", testShelter.City, htmlStr)
	}
}

func generateShelter() *managers.Shelter {
	contactInfo := &managers.ContactInformation{
		City:       testCity,
		Country:    testCountry,
		Name:       testName,
		PostalCode: testPostalCode,
		State:      testState,
		Street:     testStreet,
	}

	return &managers.Shelter{ContactInformation: contactInfo}
}
