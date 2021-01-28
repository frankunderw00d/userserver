package user

type (
	LoginRequest struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}

	LoginResponse struct {
		Token    string    `json:"token"`
		Session  string    `json:"session"`
	}
)
