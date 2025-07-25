package api

type Api struct {
	Pagination ApiPagination
	Response   ApiResponse
}

func NewApiHelper() Api {
	return Api{
		Pagination: ApiPagination{
			GetPagination: getPagination,
		},
		Response: ApiResponse{
			Json:      responseJson,
			JsonError: responseJsonError,
		},
	}
}
