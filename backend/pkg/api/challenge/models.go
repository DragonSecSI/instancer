package challenge

import (
	"github.com/DragonSecSI/instancer/backend/pkg/database/models"
)

type ChallengeListRequest struct {
	Page     int `in:"query=page,default=1"`
	Pagesize int `in:"query=pagesize,default=50"`
}

type ChallengeRequest struct {
	ID uint `in:"path=id"`
}

type Challenge struct {
	Name         string                   `json:"name"`
	Description  string                   `json:"description"`
	Category     string                   `json:"category"`
	Type         models.ChallengeType     `json:"type"`
	RemoteID     string                   `json:"remote_id"`
	Flag         string                   `json:"flag"`
	FlagType     models.ChallengeFlagType `json:"flag_type"`
	Duration     int                      `json:"duration"`
	Repository   string                   `json:"repository"`
	Chart        string                   `json:"chart"`
	ChartVersion string                   `json:"chart_version"`
	Values       string                   `json:"values"`
}

func (c Challenge) ToModel() models.Challenge {
	return models.Challenge{
		Name:         c.Name,
		Description:  c.Description,
		Category:     c.Category,
		Type:         c.Type,
		RemoteID:     c.RemoteID,
		Flag:         c.Flag,
		FlagType:     c.FlagType,
		Duration:     c.Duration,
		Repository:   c.Repository,
		Chart:        c.Chart,
		ChartVersion: c.ChartVersion,
		Values:       c.Values,
	}
}

type ChallengeNewRequest struct {
	Challenge *Challenge `in:"body;nonzero"`
}

type ChallengePutRequest struct {
	ID        uint       `in:"path=id"`
	Challenge *Challenge `in:"body;nonzero"`
}

type ChallengeResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Category    string               `json:"category"`
	Type        models.ChallengeType `json:"type"`
	Duration    int                  `json:"duration"`
}

func newChallengeResponseList(challs []models.Challenge) []ChallengeResponse {
	var result []ChallengeResponse = []ChallengeResponse{}
	for _, chall := range challs {
		result = append(result, ChallengeResponse{
			ID:          chall.ID,
			Name:        chall.Name,
			Description: chall.Description,
			Category:    chall.Category,
			Type:        chall.Type,
			Duration:    chall.Duration,
		})
	}
	return result
}

func newChallengeResponse(chall models.Challenge) ChallengeResponse {
	return ChallengeResponse{
		ID:          chall.ID,
		Name:        chall.Name,
		Description: chall.Description,
		Category:    chall.Category,
		Type:        chall.Type,
		Duration:    chall.Duration,
	}
}
