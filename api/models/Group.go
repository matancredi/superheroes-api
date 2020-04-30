package models

import (
	"errors"
	"html"
	"log"
	"strings"

	"github.com/jinzhu/gorm"
)

type Group struct {
	Uuid uint64 `gorm:"primary_key;auto_increment" json:"uuid"`
	Name string `gorm:"size:255;not null;unique" json:"name"`
}

func (g *Group) Prepare() {
	g.Uuid = 0
	g.Name = html.EscapeString(strings.TrimSpace(g.Name))
}

func (g *Group) Validate() error {

	if g.Name == "" {
		return errors.New("É necessário ter um nome")
	}
	return nil
}

func (g *Group) CreateGroup(db *gorm.DB, superUuid uint64, groups string) {
	groupsSeparated := strings.Split(groups, ",")

	for i, _ := range groupsSeparated {
		group := groupsSeparated[i]

		var err error
		gp := Group{}
		gp.Name = group

		//Checks if group is already registered
		err = db.Debug().Model(&Group{}).Where("name = ?", group).Take(&gp).Error

		if err != nil {
			// Group is not registered, so lets gonna register
			err = db.Debug().Model(&Group{}).Create(&gp).Error
			if err != nil {
				// Error when saving group
				log.Println(err)
			}
		}

		// Adds group id and super id into supergroup table
		sg := SuperGroup{}
		sg.CreateSuperGroup(db, superUuid, gp.Uuid)

	}
}

func (g *Group) FindGroupById(db *gorm.DB, GroupUuid uint64) (*Group, error) {

	var err error
	err = db.Debug().Model(&Group{}).Where("uuid = ?", GroupUuid).Take(&g).Error
	if err != nil {
		return &Group{}, err
	}

	return g, nil
}
