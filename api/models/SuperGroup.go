package models

import (
	"log"

	"github.com/jinzhu/gorm"
)

type SuperGroup struct {
	Uuid    uint64 `gorm:"primary_key;auto_increment" json:"uuid"`
	SuperId uint64 `gorm:"size:255;not null" json:"super_id"`
	GroupId uint64 `gorm:"size:255;not null" json:"group_id"`
}

func (sg *SuperGroup) Prepare() {
	sg.Uuid = 0
}

func (sg *SuperGroup) CreateSuperGroup(db *gorm.DB, superId uint64, groupId uint64) {
	sg.SuperId = superId
	sg.GroupId = groupId

	var err error

	err = db.Debug().Model(&SuperGroup{}).Create(&sg).Error
	if err != nil {
		// Error when saving supergroup
		log.Println(err)
	}
}

func (sg *SuperGroup) DeleteSuperGroup(db *gorm.DB, superId uint64) {
	db = db.Debug().Model(&SuperGroup{}).Where("super_id = ?", superId).Delete(&SuperGroup{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			log.Println("Super has no groups")
		}
	}
}

func (a *SuperGroup) FindSuperGroupBySuperId(db *gorm.DB, superUuid uint64) ([]Group, error) {

	var err error
	var groups []Group
	sg := []SuperGroup{}

	err = db.Debug().Model(&[]SuperGroup{}).Where("super_id = ?", superUuid).Find(&sg).Error
	if err != nil {
		return []Group{}, err
	}

	// For each id of groups, gets the group
	for i, _ := range sg {
		g := Group{}
		group, err := g.FindGroupById(db, sg[i].GroupId)
		if err != nil {
			return []Group{}, err
		}
		groups = append(groups, *group)
	}

	return groups, nil
}
