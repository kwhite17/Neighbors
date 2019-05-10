package items

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

var itemRetriever = &ItemRetriever{}

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

	tmpl.Execute(testBuffer, []*Item{testItem})
	htmlBytes, err := ioutil.ReadAll(testBuffer)
	if err != nil {
		t.Error(err)
	}

	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "table") || !strings.Contains(htmlStr, testItem.Status) {
		t.Errorf("TestRenderAllItemsTemplate Failure - Expected html to contain 'strong' or correct status: %s, Actual: %s\n", testItem.Status, htmlStr)
	}
}
