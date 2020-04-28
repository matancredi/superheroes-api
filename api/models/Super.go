package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Super struct {
	Uuid      uint64    `gorm:"primary_key;auto_increment" json:"uuid"`
	Name      string    `gorm:"size:255;not null;unique" json:"name"`
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
		return errors.New("Required Name")
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
