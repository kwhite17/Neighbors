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
var testStatus = "testStatus"

func TestRenderItemTemplate(t *testing.T) {
	testArray := make([]byte, 0)
	testBuffer := bytes.NewBuffer(testArray)
	testItem := generateItem()
	tmpl, err := itemRetriever.RetrieveSingleEntityTemplate()
	if err != nil {
		t.Fatal(err)
	}

	tmpl.Execute(testBuffer, testItem)
	htmlBytes, err := ioutil.ReadAll(testBuffer)
	if err != nil {
		t.Error(err)
	}

	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "strong") || !strings.Contains(htmlStr, testItem.Status) {
		t.Errorf("TestRenderItemTemplate Failure - Expected html to contain 'strong' or correct status, Actual: %s\n", testItem.Status)
	}
}

func TestRenderCreateItemTemplate(t *testing.T) {
	testArray := make([]byte, 0)
	testBuffer := bytes.NewBuffer(testArray)
	tmpl, err := itemRetriever.RetrieveCreateEntityTemplate()
	if err != nil {
		t.Fatal(err)
	}

	tmpl.Execute(testBuffer, nil)
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

	tmpl.Execute(testBuffer, []*managers.Item{testItem})
	htmlBytes, err := ioutil.ReadAll(testBuffer)
	if err != nil {
		t.Error(err)
	}

	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "table") || !strings.Contains(htmlStr, testItem.Status) {
		t.Errorf("TestRenderAllItemsTemplate Failure - Expected html to contain 'strong' or correct status: %s, Actual: %s\n", testItem.Status, htmlStr)
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
