package instance

import "github.com/DragonSecSI/instancer/backend/pkg/database/models"

type InstanceRequest struct {
	ID uint `in:"path=id;required"`
}

type InstanceResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ChallengeID int    `json:"challenge_id"`
	CreatedAt   string `json:"created_at"`
	Active      bool   `json:"active"`
}

func newInstanceResponseList(instances []models.Instance) []InstanceResponse {
	var responseList []InstanceResponse = []InstanceResponse{}
	for _, instance := range instances {
		responseList = append(responseList, InstanceResponse{
			ID:          int(instance.ID),
			Name:        instance.Name,
			ChallengeID: int(instance.ChallengeID),
			CreatedAt:   instance.CreatedAt.Format("2006-01-02 15:04:05"),
			Active:      instance.Active,
		})
	}
	return responseList
}

func newInstanceResponse(instance models.Instance) InstanceResponse {
	return InstanceResponse{
		ID:          int(instance.ID),
		Name:        instance.Name,
		ChallengeID: int(instance.ChallengeID),
		CreatedAt:   instance.CreatedAt.Format("2006-01-02 15:04:05"),
		Active:      instance.Active,
	}
}

type InstanceNewResponse struct {
	Name string `json:"name"`
}
