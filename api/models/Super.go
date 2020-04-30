package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Biography struct {
	Fullname  string `gorm:"size:255;not null" json:"full-name"`
	Alignment string `json:"alignment"`
}

type Powerstats struct {
	Intelligence string `json:"intelligence"`
	Power        string `json:"power"`
}

type Work struct {
	Occupation string `json:"occupation"`
}

type Image struct {
	Url string `json:"url"`
}

type Connections struct {
	GroupAffiliation string `json:"group-affiliation"`
	Relatives        string `json:"relatives"`
}

type Super struct {
	Uuid            uint64      `gorm:"primary_key;auto_increment" json:"uuid"`
	Name            string      `gorm:"size:255;not null;unique" json:"name"`
	Biography       Biography   `gorm:"embedded" json:"biography"`
	Powerstats      Powerstats  `gorm:"embedded" json:"powerstats"`
	Work            Work        `gorm:"embedded" json:"work"`
	Image           Image       `gorm:"embedded" json:"image"`
	Connections     Connections `json:"connections"`
	RelativesNumber int         `json:"relatives_number"`
	CreatedAt       time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (s *Super) Prepare() {
	s.Uuid = 0
	s.Name = html.EscapeString(strings.TrimSpace(s.Name))
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
}

func (s *Super) Validate() error {

	if s.Name == "" {
		return errors.New("É necessário ter um nome")
	}
	return nil
}

// Saves a new super in database
func (s *Super) SaveSuper(db *gorm.DB) (*Super, error) {
	var err error
	err = db.Debug().Model(&Super{}).Create(&s).Error
	if err != nil {
		return &Super{}, err
	}
	return s, nil
}

// Searches for all the supers
func (s *Super) FindAllSupers(db *gorm.DB) (*[]Super, error) {
	var err error
	supers := []Super{}
	err = db.Debug().Model(&Super{}).Limit(100).Find(&supers).Error
	if err != nil {
		return &[]Super{}, err
	}
	return &supers, nil
}

// Searches for supers, filtering by id
func (s *Super) FindSuperByID(db *gorm.DB, uuid uint64) (*Super, error) {
	var err error
	err = db.Debug().Model(&Super{}).Where("uuid = ?", uuid).Take(&s).Error
	if err != nil {
		return &Super{}, err
	}

	//Searches for groups of super
	sg := SuperGroup{}
	groups, err := sg.FindSuperGroupBySuperId(db, uuid)

	var groupAffiliation string = ""

	for i, _ := range groups {
		groupAffiliation += groups[i].Name
	}

	s.Connections.GroupAffiliation = groupAffiliation

	return s, nil
}

// Searches for supers, filtering by name
func (s *Super) FindSuperByName(db *gorm.DB, name string) (*Super, error) {
	var err error
	err = db.Debug().Model(&Super{}).Where("name = ?", strings.Title(name)).Take(&s).Error
	if err != nil {
		return &Super{}, err
	}
	return s, nil
}

// Searches for supers, filtering by alignment: good or bad
func (s *Super) FindSuperByAlignment(db *gorm.DB, alignment string) (*[]Super, error) {
	var err error
	supers := []Super{}
	err = db.Debug().Model(&Super{}).Where("alignment = ?", alignment).Find(&supers).Error
	if err != nil {
		return &[]Super{}, err
	}
	return &supers, nil
}

// Removes a super from database
func (s *Super) DeleteASuper(db *gorm.DB, uuid uint64) (int64, error) {

	// Deletes from super group table
	sg := SuperGroup{}
	sg.DeleteSuperGroup(db, uuid)

	// Deletes from super table
	db = db.Debug().Model(&Super{}).Where("uuid = ?", uuid).Delete(&Super{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Super not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
