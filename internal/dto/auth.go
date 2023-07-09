package dto

type LoginRequest struct {
	Username string `json:"username" from:"username" valid:"required"`
	Password string `json:"password" from:"password" valid:"required"`
}

type LoginResponse struct {
	StatusCode uint32 `json:"statusCode" from:"statusCode"`
	Token      string `json:"token" from:"token"`
}

type CreateStaffRequest struct {
	Name     string `json:"name" from:"name" valid:"type(string)"`
	Email    string `json:"email" from:"email" valid:"type(string),email"`
	Username string `json:"username" from:"username" valid:"type(string)"`
	Password string `json:"password" from:"password" valid:"type(string)"`
	ImageUrl string `json:"imageUrl" from:"imageUrl" valid:"type(string),url"`
	Role     string `json:"role" from:"role" valid:"type(string),in(STAFF|ADMIN)"`
}
