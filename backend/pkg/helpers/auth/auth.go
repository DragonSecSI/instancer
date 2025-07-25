package auth

type Auth struct {
	Admin AuthAdmin
	Token AuthToken
}

func NewAuthHelper() Auth {
	return Auth{
		Admin: AuthAdmin{
			IsAdmin: authAdminIsAdmin,
		},
		Token: AuthToken{
			GetTeam:       authTokenGetTeam,
			GenerateToken: authTokenGenerateToken,
		},
	}
}
