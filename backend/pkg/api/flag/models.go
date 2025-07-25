package flag

type FlagSubmit struct {
	Flag     string `json:"flag"`
	RemoteID string `json:"remote_id"`
}

type FlagSubmitRequest struct {
	Body *FlagSubmit `in:"body;nonzero"`
}

type FlagSubmitResponse struct {
	Correct        bool `json:"correct"`
	ActiveInstance bool `json:"active_instance"`
	WrongTeam      bool `json:"wrong_team"`
}
