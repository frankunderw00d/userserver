package user

import "baseservice/middleware/authenticate"

type (
	AutoLoginRequest struct {
		authenticate.Request
	}

	AutoLoginResponse struct {
		authenticate.Response
	}
)
