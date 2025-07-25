package models

import (
	goerrors "errors"

	"gorm.io/gorm"

	"github.com/DragonSecSI/instancer/backend/pkg/database"
	"github.com/DragonSecSI/instancer/backend/pkg/errors"
)

type ChallengeFlagType int

const (
	ChallengeFlagTypeStatic ChallengeFlagType = iota
	ChallengeFlagTypeSuffix
	ChallengeFlagTypeLeetify
	ChallengeFlagTypeCapitalize
	ChallengeFlagTypeCombined
)

type Challenge struct {
	ID          uint              `gorm:"primaryKey"`
	Name        string            `gorm:"not null;uniqueIndex"`
	Description string            `gorm:"not null"`
	Flag        string            `gorm:"not null"`
	FlagType    ChallengeFlagType `gorm:"not null"`
	Duration    int               `gorm:"not null"`

	Repository   string `gorm:"not null"`
	Chart        string `gorm:"not null"`
	ChartVersion string `gorm:"not null"`
	Values       string `gorm:"not null"`
}

func ChallengeGetList(db *gorm.DB, page int, pagesize int) ([]Challenge, error) {
	var challenges []Challenge
	err := db.Scopes(database.Paginate(db, page, pagesize)).Find(&challenges).Error
	if err != nil {
		return nil, &errors.DatabaseQueryError{
			Query: "ChallengeGetList",
			Err:   err,
		}
	}
	return challenges, nil
}

func ChallengeGetByID(db *gorm.DB, id uint) (*Challenge, error) {
	var challenge Challenge
	err := db.First(&challenge, id).Error
	if err != nil {
		if goerrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, &errors.DatabaseQueryError{
			Query: "ChallengeGetByID",
			Err:   err,
		}
	}
	return &challenge, nil
}

func ChallengeCreate(db *gorm.DB, challenge *Challenge) error {
	err := db.Create(challenge).Error
	if err != nil {
		return &errors.DatabaseQueryError{
			Query: "ChallengeCreate",
			Err:   err,
		}
	}
	return nil
}

func ChallengeUpdate(db *gorm.DB, challenge *Challenge) error {
	err := db.Save(challenge).Error
	if err != nil {
		return &errors.DatabaseQueryError{
			Query: "ChallengeUpdate",
			Err:   err,
		}
	}
	return nil
}

func ChallengeDelete(db *gorm.DB, id uint) error {
	err := db.Delete(&Challenge{}, id).Error
	if err != nil {
		return &errors.DatabaseQueryError{
			Query: "ChallengeDelete",
			Err:   err,
		}
	}
	return nil
}
