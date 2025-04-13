package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thiomajid/ecommerce-backend/internal/database"
	"github.com/thiomajid/ecommerce-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

const AdminToken = "admin-secret-token" // En production, utilisez une méthode plus sécurisée

// AdminAuthMiddleware vérifie si la requête est authentifiée en tant qu'administrateur
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Vérifier d'abord l'en-tête Admin-Token (ancien format)
		token := c.GetHeader("Admin-Token")

		// Si l'Admin-Token n'est pas présent, vérifier l'en-tête Authorization (Bearer token)
		if token == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if token != AdminToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentification requise"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthMiddleware vérifie si la requête est authentifiée
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentification requise"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Vérifier si le token existe dans la base de données
		var userId string
		err := database.DB.QueryRow("SELECT user_id FROM auth_tokens WHERE token = ?", token).Scan(&userId)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
			c.Abort()
			return
		}

		// Ajouter l'ID utilisateur au contexte pour l'utiliser dans les handlers
		c.Set("userId", userId)

		c.Next()
	}
}

// LoginAdmin vérifie les identifiants administrateur et renvoie un token
func LoginAdmin(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données d'identification invalides"})
		return
	}

	// Vérification simple des identifiants (à remplacer par une authentification sécurisée en production)
	if credentials.Username == "admin" && credentials.Password == "admin123" {
		c.JSON(http.StatusOK, gin.H{
			"token":   AdminToken,
			"message": "Connexion réussie",
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Identifiants invalides"})
	}
}

// Fonction utilitaire pour générer un token aléatoire
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// RegisterUser crée un nouveau compte utilisateur
func RegisterUser(c *gin.Context) {
	var request models.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Vérifier si l'email existe déjà
	var exists bool
	err := database.DB.QueryRow("SELECT 1 FROM users WHERE email = ?", request.Email).Scan(&exists)
	if err != sql.ErrNoRows {
		c.JSON(http.StatusConflict, gin.H{"error": "Cette adresse email est déjà utilisée"})
		return
	}

	// Hasher le mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du traitement du mot de passe"})
		return
	}

	// Générer un ID utilisateur unique
	userId, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du compte"})
		return
	}

	// Préparer les données utilisateur
	now := time.Now().Format(time.RFC3339)

	// Insérer l'utilisateur dans la base de données
	_, err = database.DB.Exec(
		"INSERT INTO users (id, email, password, name, address, phone, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		userId, request.Email, string(hashedPassword), request.Name, request.Address, request.Phone, now,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du compte: " + err.Error()})
		return
	}

	// Générer un token d'authentification
	token, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la génération du token"})
		return
	}

	// Enregistrer le token dans la base de données
	_, err = database.DB.Exec(
		"INSERT INTO auth_tokens (token, user_id, created_at) VALUES (?, ?, ?)",
		token, userId, now,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'enregistrement du token"})
		return
	}

	// Préparer la réponse
	userResponse := models.UserResponse{
		ID:        userId,
		Email:     request.Email,
		Name:      request.Name,
		Address:   request.Address,
		Phone:     request.Phone,
		CreatedAt: now,
	}

	c.JSON(http.StatusCreated, models.LoginResponse{
		Token: token,
		User:  userResponse,
	})
}

// LoginUser authentifie un utilisateur existant
func LoginUser(c *gin.Context) {
	var request models.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Récupérer l'utilisateur depuis la base de données
	var user models.User
	err := database.DB.QueryRow(
		"SELECT id, email, password, name, address, phone, created_at FROM users WHERE email = ?",
		request.Email,
	).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, &user.Address, &user.Phone, &user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération de l'utilisateur"})
		return
	}

	// Vérifier le mot de passe
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	}

	// Générer un token d'authentification
	token, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la génération du token"})
		return
	}

	// Enregistrer le token dans la base de données
	now := time.Now().Format(time.RFC3339)
	_, err = database.DB.Exec(
		"INSERT INTO auth_tokens (token, user_id, created_at) VALUES (?, ?, ?)",
		token, user.ID, now,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'enregistrement du token"})
		return
	}

	// Préparer la réponse
	userResponse := models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Address:   user.Address,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User:  userResponse,
	})
}

// GetUserProfile récupère les informations du profil utilisateur
func GetUserProfile(c *gin.Context) {
	// Récupérer l'ID utilisateur du contexte
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	// Récupérer l'utilisateur depuis la base de données
	var user models.UserResponse
	err := database.DB.QueryRow(
		"SELECT id, email, name, address, phone, created_at FROM users WHERE id = ?",
		userId,
	).Scan(
		&user.ID, &user.Email, &user.Name, &user.Address, &user.Phone, &user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération du profil"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserOrders récupère l'historique des commandes de l'utilisateur
func GetUserOrders(c *gin.Context) {
	// Récupérer l'ID utilisateur du contexte
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	// Récupérer les commandes de l'utilisateur
	rows, err := database.DB.Query(
		"SELECT id, status, total_price, created_at FROM orders WHERE user_id = ? ORDER BY created_at DESC",
		userId,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des commandes"})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.Status, &order.TotalPrice, &order.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du scan des commandes"})
			return
		}
		order.UserID = userId.(string)
		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, orders)
}
