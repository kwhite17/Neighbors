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
var testEmail = "test@test.com"

var shelterRetriever = &ShelterRetriever{}

func TestRenderSingleShelterTemplate(t *testing.T) {
	testArray := make([]byte, 0)
	testBuffer := bytes.NewBuffer(testArray)
	testShelter := generateShelter()
	testShelter.UserType = managers.SHELTER
	tmpl, err := shelterRetriever.RetrieveSingleEntityTemplate()

	if err != nil {
		t.Fatal(err)
	}

	tmpl.Execute(testBuffer, map[string]interface{}{
		"User":           testShelter,
		"ShelterSession": nil,
	})
	htmlBytes, err := ioutil.ReadAll(testBuffer)

	if err != nil {
		t.Error(err)
	}

	htmlStr := string(htmlBytes)

	if !strings.Contains(htmlStr, "testStreet, testCity, testState,") || !strings.Contains(htmlStr, testShelter.Name+" ("+testShelter.Email+")") {
		t.Errorf("GetSingleShelter Failure - Expected html to contain location info, shelter name, and shelter email. Actual: %s\n", htmlStr)
	}
}
func TestRenderSingleSamaritanTemplate(t *testing.T) {
	testArray := make([]byte, 0)
	testBuffer := bytes.NewBuffer(testArray)
	testShelter := generateShelter()
	testShelter.UserType = managers.SAMARITAN
	tmpl, err := shelterRetriever.RetrieveSingleEntityTemplate()

	if err != nil {
		t.Fatal(err)
	}

	tmpl.Execute(testBuffer, map[string]interface{}{
		"User":           testShelter,
		"ShelterSession": nil,
	})
	htmlBytes, err := ioutil.ReadAll(testBuffer)

	if err != nil {
		t.Error(err)
	}

	htmlStr := string(htmlBytes)

	if strings.Contains(htmlStr, "testStreet, testCity, testState,") || !strings.Contains(htmlStr, testShelter.Name+" ("+testShelter.Email+")") {
		t.Errorf("GetSingleSamaritan Failure - Expected html to contain samaritan name and email without location info, Actual: %s\n", htmlStr)
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
		"Users":          []*managers.User{testShelter},
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

func generateShelter() *managers.User {
	contactInfo := &managers.ContactInformation{
		City:       testCity,
		Name:       testName,
		PostalCode: testPostalCode,
		State:      testState,
		Street:     testStreet,
		Email:      testEmail,
	}

	return &managers.User{ContactInformation: contactInfo, UserType: managers.SHELTER}
}
