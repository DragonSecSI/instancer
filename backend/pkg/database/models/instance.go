package models

import (
	goerrors "errors"
	"time"

	"gorm.io/gorm"

	"github.com/DragonSecSI/instancer/backend/pkg/database"
	"github.com/DragonSecSI/instancer/backend/pkg/errors"
)

type Instance struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null;uniqueIndex"`
	Flag string `gorm:"not null;uniqueIndex"`

	TeamID        uint          `gorm:"not null;index"`
	Team          Team          `gorm:"foreignKey:TeamID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ChallengeID   uint          `gorm:"not null;index"`
	Challenge     Challenge     `gorm:"foreignKey:ChallengeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ChallengeType ChallengeType `gorm:"not null"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Active    bool      `gorm:"not null;index;default:true"`
	Duration  int       `gorm:"not null;default:1800"`
}

func InstanceGetList(db *gorm.DB, page int, pagesize int) ([]Instance, error) {
	var instances []Instance
	err := db.Scopes(database.Paginate(db, page, pagesize)).Order("id desc").Find(&instances).Error
	if err != nil {
		return nil, &errors.DatabaseQueryError{
			Query: "InstanceGetList",
			Err:   err,
		}
	}
	return instances, nil
}

func InstanceGetByID(db *gorm.DB, id uint) (*Instance, error) {
	var instance Instance
	err := db.First(&instance, id).Error
	if err != nil {
		if goerrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, &errors.DatabaseQueryError{
			Query: "InstanceGetByID",
			Err:   err,
		}
	}
	return &instance, nil
}

func InstanceGetByTeamID(db *gorm.DB, teamID uint, page int, pagesize int) ([]Instance, error) {
	var instances []Instance
	err := db.
		Scopes(database.Paginate(db, page, pagesize)).
		Where("team_id = ?", teamID).
		Order("id desc").
		Find(&instances).
		Error
	if err != nil {
		return nil, &errors.DatabaseQueryError{
			Query: "InstanceGetByTeamID",
			Err:   err,
		}
	}
	return instances, nil
}

func InstanceGetByName(db *gorm.DB, name string) (*Instance, error) {
	var instance Instance
	err := db.Where("name = ?", name).First(&instance).Error
	if err != nil {
		if goerrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, &errors.DatabaseQueryError{
			Query: "InstanceGetByName",
			Err:   err,
		}
	}
	return &instance, nil
}

func InstanceGetByFlag(db *gorm.DB, flag string) (*Instance, error) {
	var instance Instance
	err := db.Preload("Team").Preload("Challenge").Where("flag = ?", flag).First(&instance).Error
	if err != nil {
		if goerrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, &errors.DatabaseQueryError{
			Query: "InstanceGetByFlag",
			Err:   err,
		}
	}
	return &instance, nil
}

func InstanceGetActive(db *gorm.DB) ([]Instance, error) {
	var instances []Instance
	err := db.Where("active = ?", true).Find(&instances).Error
	if err != nil {
		return nil, &errors.DatabaseQueryError{
			Query: "InstanceGetActive",
			Err:   err,
		}
	}
	return instances, nil
}

func InstanceCreate(db *gorm.DB, instance *Instance) error {
	err := db.Create(instance).Error
	if err != nil {
		return &errors.DatabaseQueryError{
			Query: "InstanceCreate",
			Err:   err,
		}
	}
	return nil
}

func InstanceUpdate(db *gorm.DB, instance *Instance) error {
	err := db.Save(instance).Error
	if err != nil {
		return &errors.DatabaseQueryError{
			Query: "InstanceUpdate",
			Err:   err,
		}
	}
	return nil
}

func InstanceDelete(db *gorm.DB, id uint) error {
	err := db.Delete(&Instance{}, id).Error
	if err != nil {
		return &errors.DatabaseQueryError{
			Query: "InstanceDelete",
			Err:   err,
		}
	}
	return nil
}
