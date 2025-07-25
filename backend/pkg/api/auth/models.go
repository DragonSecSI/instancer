package auth

type Team struct {
	Name     string `json:"name"`
	RemoteID string `json:"remote_id"`
}

type TeamRegisterRequest struct {
	Team *Team `in:"body;nonzero"`
}

type TeamRegisterResponse struct {
	Token string `json:"token"`
}
