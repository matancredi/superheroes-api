package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/matancredi/superheroes-api/api/models"
)

var supers = []models.Super{
	models.Super{
		Name: "Mulher Maravilha",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Super{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.Super{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i, _ := range supers {
		err = db.Debug().Model(&models.Super{}).Create(&supers[i]).Error
		if err != nil {
			log.Fatalf("cannot seed supers table: %v", err)
		}
	}
}
