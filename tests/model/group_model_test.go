package model

import (
	"log"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/go-playground/assert.v1"
)

func TestGroupByID(t *testing.T) {

	err := refreshTables()
	if err != nil {
		log.Fatal(err)
	}

	group, err := seedOneGroup()
	if err != nil {
		log.Fatalf("cannot seed group table: %v", err)
	}
	foundGroup, err := groupInstance.FindGroupById(server.DB, group.Uuid)
	if err != nil {
		t.Errorf("this is the error getting one group: %v\n", err)
		return
	}
	assert.Equal(t, foundGroup.Uuid, group.Uuid)
	assert.Equal(t, foundGroup.Name, group.Name)
}
