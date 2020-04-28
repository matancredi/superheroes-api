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

type Super struct {
	Uuid       uint64     `gorm:"primary_key;auto_increment" json:"uuid"`
	Name       string     `gorm:"size:255;not null" json:"name"`
	Biography  Biography  `gorm:"embedded" json:"biography"`
	Powerstats Powerstats `gorm:"embedded" json:"powerstats"`
	Work       Work       `gorm:"embedded" json:"work"`
	Image      Image      `gorm:"embedded" json:"image"`
	//grupos
	//numero de parentes
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
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

func (s *Super) SaveSuper(db *gorm.DB) (*Super, error) {
	var err error
	err = db.Debug().Model(&Super{}).Create(&s).Error
	if err != nil {
		return &Super{}, err
	}
	return s, nil
}

func (s *Super) FindAllSupers(db *gorm.DB) (*[]Super, error) {
	var err error
	supers := []Super{}
	err = db.Debug().Model(&Super{}).Limit(100).Find(&supers).Error
	if err != nil {
		return &[]Super{}, err
	}
	return &supers, nil
}

func (s *Super) DeleteASuper(db *gorm.DB, pid uint64) (int64, error) {

	db = db.Debug().Model(&Super{}).Where("uuid = ?", pid).Delete(&Super{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Super not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
