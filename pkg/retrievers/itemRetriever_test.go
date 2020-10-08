package retrievers

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"strings"
	"testing"

	"github.com/kwhite17/Neighbors/pkg/managers"
)

var itemRetriever = &ItemRetriever{}

var testCategory = "testCategory"
var testGender = "testGender"
var testQuantity = int8(rand.Int() % 127)
var testShelterID = rand.Int63()
var testSize = "testSize"
var testStatus = managers.CLAIMED

func TestRenderItemTemplate(t *testing.T) {
	testArray := make([]byte, 0)
	testBuffer := bytes.NewBuffer(testArray)
	testItem := generateItem()
	tmpl, err := itemRetriever.RetrieveSingleEntityTemplate()

	if err != nil {
		t.Fatal(err)
	}

	tmpl.Execute(testBuffer, map[string]interface{}{
		"Item":           testItem,
		"ShelterSession": nil,
	})
	htmlBytes, err := ioutil.ReadAll(testBuffer)

	if err != nil {
		t.Error(err)
	}

	htmlStr := string(htmlBytes)

	if !strings.Contains(htmlStr, "<p class=\"card-text\">Status: CLAIMED</p>") || !strings.Contains(htmlStr, StatusAsString(testItem.Status)) {
		t.Errorf("TestRenderItemTemplate Failure - Expected html to contain 'strong' or correct status, Actual: %s\n", htmlStr)
	}
}

func TestRenderCreateItemTemplate(t *testing.T) {
	testArray := make([]byte, 0)
	testBuffer := bytes.NewBuffer(testArray)
	tmpl, err := itemRetriever.RetrieveCreateEntityTemplate()

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

	if !strings.Contains(htmlStr, "form") || !strings.Contains(htmlStr, "Quantity") {
		t.Errorf("TestRenderCreateItemTemplate Failure - Expected html to contain 'form' or 'Quantity', Actual: %s\n", htmlStr)
	}
}

func TestRenderAllItemsTemplate(t *testing.T) {
	testArray := make([]byte, 0)
	testBuffer := bytes.NewBuffer(testArray)
	testItem := generateItem()
	tmpl, err := itemRetriever.RetrieveAllEntitiesTemplate()

	if err != nil {
		t.Fatal(err)
	}

	tmpl.Execute(testBuffer, map[string]interface{}{
		"Items":          []managers.Item{*testItem},
		"ShelterSession": nil,
	})
	htmlBytes, err := ioutil.ReadAll(testBuffer)

	if err != nil {
		t.Error(err)
	}

	htmlStr := string(htmlBytes)

	if !strings.Contains(htmlStr, "table") || !strings.Contains(htmlStr, StatusAsString(testItem.Status)) {
		t.Errorf("TestRenderAllItemsTemplate Failure - Expected html to contain 'strong' or correct status: %s, Actual: %s\n", StatusAsString(testItem.Status), htmlStr)
	}
}
func generateItem() *managers.Item {
	return &managers.Item{
		Category:  testCategory,
		Gender:    testGender,
		Quantity:  testQuantity,
		ShelterID: testShelterID,
		Size:      testSize,
		Status:    testStatus,
	}
}
