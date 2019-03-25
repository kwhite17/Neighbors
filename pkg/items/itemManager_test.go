package items

import (
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestShouldWriteItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	log.Println(mock)
}
