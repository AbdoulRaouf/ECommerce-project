package models

// User représente un utilisateur du site
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"-"` // Le tiret empêche ce champ d'être inclus dans les réponses JSON
	Name      string `json:"name"`
	Address   string `json:"address,omitempty"`
	Phone     string `json:"phone,omitempty"`
	CreatedAt string `json:"createdAt"`
}

// UserResponse est utilisé pour les réponses d'API sans données sensibles
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Address   string `json:"address,omitempty"`
	Phone     string `json:"phone,omitempty"`
	CreatedAt string `json:"createdAt"`
}

// LoginRequest représente une demande de connexion
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest représente une demande d'inscription
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
}

// LoginResponse est la réponse à une connexion réussie
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
