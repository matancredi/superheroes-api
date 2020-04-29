package model

import (
	"log"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/matancredi/superheroes-api/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestSaveSuper(t *testing.T) {

	err := refreshSuperTable()
	if err != nil {
		log.Fatal(err)
	}
	newSuper := models.Super{
		Name: "Super Chock",
	}
	savedSuper, err := newSuper.SaveSuper(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the supers: %v\n", err)
		return
	}
	assert.Equal(t, newSuper.Uuid, savedSuper.Uuid)
	assert.Equal(t, newSuper.Name, savedSuper.Name)
}

func TestFindAllSupers(t *testing.T) {

	err := refreshSuperTable()
	if err != nil {
		log.Fatal(err)
	}

	err = seedSupers()
	if err != nil {
		log.Fatal(err)
	}

	supers, err := superInstance.FindAllSupers(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the supers: %v\n", err)
		return
	}
	assert.Equal(t, len(*supers), 2)
}

func TestFindSuperByID(t *testing.T) {

	err := refreshSuperTable()
	if err != nil {
		log.Fatal(err)
	}

	super, err := seedOneSuper()
	if err != nil {
		log.Fatalf("cannot seed supers table: %v", err)
	}
	foundSuper, err := superInstance.FindSuperByID(server.DB, super.Uuid)
	if err != nil {
		t.Errorf("this is the error getting one super: %v\n", err)
		return
	}
	assert.Equal(t, foundSuper.Uuid, super.Uuid)
	assert.Equal(t, foundSuper.Name, super.Name)
}

func TestFindSuperByName(t *testing.T) {

	err := refreshSuperTable()
	if err != nil {
		log.Fatal(err)
	}

	super, err := seedOneSuper()
	if err != nil {
		log.Fatalf("cannot seed supers table: %v", err)
	}
	foundSuper, err := superInstance.FindSuperByName(server.DB, super.Name)
	if err != nil {
		t.Errorf("this is the error getting one super: %v\n", err)
		return
	}
	assert.Equal(t, foundSuper.Uuid, super.Uuid)
	assert.Equal(t, foundSuper.Name, super.Name)
}

func TestFindSuperByAlignment(t *testing.T) {

	err := refreshSuperTable()
	if err != nil {
		log.Fatal(err)
	}

	err = seedSupers()
	if err != nil {
		log.Fatal(err)
	}

	supers, err := superInstance.FindSuperByAlignment(server.DB, "good")
	if err != nil {
		t.Errorf("this is the error getting the supers: %v\n", err)
		return
	}
	assert.Equal(t, len(*supers), 1)
}

func TestDeleteASuper(t *testing.T) {

	err := refreshSuperTable()
	if err != nil {
		log.Fatal(err)
	}

	super, err := seedOneSuper()

	if err != nil {
		log.Fatalf("Cannot seed super: %v\n", err)
	}

	isDeleted, err := superInstance.DeleteASuper(server.DB, super.Uuid)
	if err != nil {
		t.Errorf("this is the error updating the super: %v\n", err)
		return
	}

	assert.Equal(t, isDeleted, int64(1))
}
