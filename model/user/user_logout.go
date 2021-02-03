package user

import "baseservice/middleware/authenticate"

type (
	LogoutRequest struct {
		authenticate.Request
	}

	LogoutResponse struct {
		authenticate.Response
	}
)

const ()

var ()
