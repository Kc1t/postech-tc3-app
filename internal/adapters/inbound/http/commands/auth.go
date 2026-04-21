package commands

type RegisterRequest struct {
	Name     string `json:"name"     binding:"required,max=100"`
	Email    string `json:"email"    binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,max=128"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required,max=255"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`  // segundos ate expirar
	TokenType    string `json:"token_type"`  // "Bearer"
}

type MessageResponse struct {
	Message string `json:"message"`
}
