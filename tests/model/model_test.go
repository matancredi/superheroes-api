package model

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/matancredi/superheroes-api/api/controllers"
	"github.com/matancredi/superheroes-api/api/models"
)

var server = controllers.Server{}
var superInstance = models.Super{}
var groupInstance = models.Group{}

// Database and .env data
func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
}

// Resets database
func refreshTables() error {

	err := server.DB.DropTableIfExists(&models.Super{}, &models.Group{}, &models.SuperGroup{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Super{}, &models.Group{}, &models.SuperGroup{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

// Seeds one super
func seedOneSuper() (models.Super, error) {

	refreshTables()

	super := models.Super{
		Name:      "Homem Aranha",
		Biography: models.Biography{Alignment: "good"},
	}

	err := server.DB.Model(&models.Super{}).Create(&super).Error
	if err != nil {
		log.Fatalf("cannot seed supers table: %v", err)
	}
	return super, nil
}

// Seeds one group
func seedOneGroup() (models.Group, error) {

	refreshTables()

	group := models.Group{
		Name: "Avengers",
	}

	err := server.DB.Model(&models.Group{}).Create(&group).Error
	if err != nil {
		log.Fatalf("cannot seed group table: %v", err)
	}
	return group, nil
}

// Seeds more than one super
func seedSupers() error {

	supers := []models.Super{
		models.Super{
			Name:      "Pantera Negra",
			Biography: models.Biography{Alignment: "bad"},
		},
		models.Super{
			Name:      "Capitão América",
			Biography: models.Biography{Alignment: "good"},
		},
	}

	for i, _ := range supers {
		err := server.DB.Model(&models.Super{}).Create(&supers[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}
