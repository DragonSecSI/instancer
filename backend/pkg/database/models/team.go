package models

import (
	goerrors "errors"
	"time"

	"gorm.io/gorm"

	"github.com/DragonSecSI/instancer/backend/pkg/database"
	"github.com/DragonSecSI/instancer/backend/pkg/errors"
)

type Team struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null;uniqueIndex"`
	RemoteID string `gorm:"not null;uniqueIndex"`
	Token    string `gorm:"not null;uniqueIndex"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func TeamGetList(db *gorm.DB, page int, pagesize int) ([]Team, error) {
	var teams []Team
	err := db.Scopes(database.Paginate(db, page, pagesize)).Find(&teams).Error
	if err != nil {
		return nil, &errors.DatabaseQueryError{
			Query: "TeamGetList",
			Err:   err,
		}
	}
	return teams, nil
}

func TeamGetByID(db *gorm.DB, id uint) (*Team, error) {
	var team Team
	err := db.First(&team, id).Error
	if err != nil {
		if goerrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, &errors.DatabaseQueryError{
			Query: "TeamGetByID",
			Err:   err,
		}
	}
	return &team, nil
}

func TeamGetByRemoteID(db *gorm.DB, remoteID string) (*Team, error) {
	var team Team
	err := db.Where("remote_id = ?", remoteID).First(&team).Error
	if err != nil {
		if goerrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, &errors.DatabaseQueryError{
			Query: "TeamGetByRemoteID",
			Err:   err,
		}
	}
	return &team, nil
}

func TeamGetByToken(db *gorm.DB, token string) (*Team, error) {
	var team Team
	err := db.Where("token = ?", token).First(&team).Error
	if err != nil {
		if goerrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, &errors.DatabaseQueryError{
			Query: "TeamGetByToken",
			Err:   err,
		}
	}
	return &team, nil
}

func TeamCreate(db *gorm.DB, team *Team) error {
	err := db.Create(team).Error
	if err != nil {
		return &errors.DatabaseQueryError{
			Query: "TeamCreate",
			Err:   err,
		}
	}
	return nil
}

func TeamUpdate(db *gorm.DB, team *Team) error {
	err := db.Save(team).Error
	if err != nil {
		return &errors.DatabaseQueryError{
			Query: "TeamUpdate",
			Err:   err,
		}
	}
	return nil
}

func TeamDelete(db *gorm.DB, id uint) error {
	err := db.Delete(&Team{}, id).Error
	if err != nil {
		return &errors.DatabaseQueryError{
			Query: "TeamDelete",
			Err:   err,
		}
	}
	return nil
}
