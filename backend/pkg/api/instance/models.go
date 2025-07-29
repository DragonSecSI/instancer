package instance

import "github.com/DragonSecSI/instancer/backend/pkg/database/models"

type InstanceListRequest struct {
	Page     int `in:"query=page;default=1"`
	Pagesize int `in:"query=pagesize;default=25"`
}

func (r *InstanceListRequest) Validate() bool {
	if r.Page < 1 {
		return false
	}
	if r.Pagesize < 1 || r.Pagesize > 100 {
		return false
	}
	return true
}

type InstanceRequest struct {
	ID uint `in:"path=id;required"`
}

type InstanceResponse struct {
	ID            int                  `json:"id"`
	Name          string               `json:"name"`
	ChallengeID   int                  `json:"challenge_id"`
	ChallengeType models.ChallengeType `json:"type"`
	CreatedAt     string               `json:"created_at"`
	Active        bool                 `json:"active"`
	Duration      int                  `json:"duration"`
}

func newInstanceResponseList(instances []models.Instance) []InstanceResponse {
	var responseList []InstanceResponse = []InstanceResponse{}
	for _, instance := range instances {
		responseList = append(responseList, InstanceResponse{
			ID:            int(instance.ID),
			Name:          instance.Name,
			ChallengeID:   int(instance.ChallengeID),
			ChallengeType: instance.ChallengeType,
			CreatedAt:     instance.CreatedAt.Format("2006-01-02 15:04:05"),
			Active:        instance.Active,
			Duration:      instance.Duration,
		})
	}
	return responseList
}

func newInstanceResponse(instance models.Instance) InstanceResponse {
	return InstanceResponse{
		ID:            int(instance.ID),
		Name:          instance.Name,
		ChallengeID:   int(instance.ChallengeID),
		ChallengeType: instance.ChallengeType,
		CreatedAt:     instance.CreatedAt.Format("2006-01-02 15:04:05"),
		Active:        instance.Active,
		Duration:      instance.Duration,
	}
}

type InstanceNewResponse struct {
	Name string `json:"name"`
}
