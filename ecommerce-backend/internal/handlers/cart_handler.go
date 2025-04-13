package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thiomajid/ecommerce-backend/internal/database"
	"github.com/thiomajid/ecommerce-backend/internal/models"
)

// GetCartItems récupère le contenu du panier d'un utilisateur
func GetCartItems(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		// Pour des fins de test, utiliser un ID par défaut si non fourni
		userID = "test-user"
	}

	// Récupérer les éléments du panier avec les informations des produits
	rows, err := database.DB.Query(`
		SELECT c.id, c.user_id, c.product_id, c.quantity, 
			p.id, p.name, p.price, p.description, p.category, p.image, p.stock
		FROM cart_items c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_id = ?
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération du panier: " + err.Error()})
		return
	}
	defer rows.Close()

	var cartItems []models.CartItem
	for rows.Next() {
		var ci models.CartItem
		var p models.Product
		if err := rows.Scan(
			&ci.ID, &ci.UserID, &ci.ProductID, &ci.Quantity,
			&p.ID, &p.Name, &p.Price, &p.Description, &p.Category, &p.Image, &p.Stock,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du scan des éléments du panier: " + err.Error()})
			return
		}
		ci.Product = p
		cartItems = append(cartItems, ci)
	}

	c.JSON(http.StatusOK, cartItems)
}

// AddToCart ajoute un produit au panier
func AddToCart(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		// Pour des fins de test, utiliser un ID par défaut si non fourni
		userID = "test-user"
	}

	var cartItem models.CartItem
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Vérifier si le produit existe et a suffisamment de stock
	var stock int
	err := database.DB.QueryRow("SELECT stock FROM products WHERE id = ?", cartItem.ProductID).Scan(&stock)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Produit non trouvé"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification du stock: " + err.Error()})
		}
		return
	}

	if stock < cartItem.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock insuffisant"})
		return
	}

	// Vérifier si le produit est déjà dans le panier
	var existingID int
	var existingQuantity int
	err = database.DB.QueryRow(
		"SELECT id, quantity FROM cart_items WHERE user_id = ? AND product_id = ?",
		userID, cartItem.ProductID,
	).Scan(&existingID, &existingQuantity)

	var result sql.Result
	if err == sql.ErrNoRows {
		// Le produit n'est pas dans le panier, l'ajouter
		result, err = database.DB.Exec(
			"INSERT INTO cart_items (user_id, product_id, quantity) VALUES (?, ?, ?)",
			userID, cartItem.ProductID, cartItem.Quantity,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'ajout au panier: " + err.Error()})
			return
		}

		id, _ := result.LastInsertId()
		cartItem.ID = int(id)
	} else if err != nil {
		// Une erreur s'est produite
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification du panier: " + err.Error()})
		return
	} else {
		// Le produit est déjà dans le panier, mettre à jour la quantité
		newQuantity := existingQuantity + cartItem.Quantity
		if newQuantity > stock {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stock insuffisant pour la quantité demandée"})
			return
		}

		_, err = database.DB.Exec(
			"UPDATE cart_items SET quantity = ? WHERE id = ?",
			newQuantity, existingID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du panier: " + err.Error()})
			return
		}

		cartItem.ID = existingID
		cartItem.Quantity = newQuantity
	}

	// Récupérer les informations complètes du produit
	var product models.Product
	err = database.DB.QueryRow(
		"SELECT id, name, price, description, category, image, stock FROM products WHERE id = ?",
		cartItem.ProductID,
	).Scan(&product.ID, &product.Name, &product.Price, &product.Description, &product.Category, &product.Image, &product.Stock)

	if err != nil {
		log.Printf("Erreur lors de la récupération des détails du produit: %v", err)
	} else {
		cartItem.Product = product
	}

	cartItem.UserID = userID
	c.JSON(http.StatusOK, cartItem)
}

// UpdateCartItem met à jour la quantité d'un produit dans le panier
func UpdateCartItem(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		userID = "test-user"
	}

	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID d'élément de panier invalide"})
		return
	}

	var cartItem models.CartItem
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Vérifier si l'élément du panier existe et appartient à l'utilisateur
	var productID int
	err = database.DB.QueryRow(
		"SELECT product_id FROM cart_items WHERE id = ? AND user_id = ?",
		itemID, userID,
	).Scan(&productID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Élément du panier non trouvé"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification de l'élément du panier: " + err.Error()})
		return
	}

	// Vérifier le stock disponible
	var stock int
	err = database.DB.QueryRow("SELECT stock FROM products WHERE id = ?", productID).Scan(&stock)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification du stock: " + err.Error()})
		return
	}

	if cartItem.Quantity > stock {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock insuffisant"})
		return
	}

	if cartItem.Quantity <= 0 {
		// Si la quantité est 0 ou négative, supprimer l'élément du panier
		_, err = database.DB.Exec("DELETE FROM cart_items WHERE id = ?", itemID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression de l'élément du panier: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Élément supprimé du panier"})
		return
	}

	// Mettre à jour la quantité
	_, err = database.DB.Exec(
		"UPDATE cart_items SET quantity = ? WHERE id = ?",
		cartItem.Quantity, itemID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour de l'élément du panier: " + err.Error()})
		return
	}

	// Récupérer l'élément mis à jour avec les détails du produit
	var updatedItem models.CartItem
	var product models.Product
	err = database.DB.QueryRow(`
		SELECT c.id, c.user_id, c.product_id, c.quantity, 
			p.id, p.name, p.price, p.description, p.category, p.image, p.stock
		FROM cart_items c
		JOIN products p ON c.product_id = p.id
		WHERE c.id = ?
	`, itemID).Scan(
		&updatedItem.ID, &updatedItem.UserID, &updatedItem.ProductID, &updatedItem.Quantity,
		&product.ID, &product.Name, &product.Price, &product.Description, &product.Category, &product.Image, &product.Stock,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération de l'élément mis à jour: " + err.Error()})
		return
	}

	updatedItem.Product = product
	c.JSON(http.StatusOK, updatedItem)
}

// RemoveFromCart supprime un élément du panier
func RemoveFromCart(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		userID = "test-user"
	}

	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID d'élément de panier invalide"})
		return
	}

	// Vérifier si l'élément existe et appartient à l'utilisateur
	var exists bool
	err = database.DB.QueryRow(
		"SELECT 1 FROM cart_items WHERE id = ? AND user_id = ?",
		itemID, userID,
	).Scan(&exists)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Élément du panier non trouvé"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification de l'élément du panier: " + err.Error()})
		return
	}

	// Supprimer l'élément du panier
	_, err = database.DB.Exec("DELETE FROM cart_items WHERE id = ?", itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression de l'élément du panier: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Élément supprimé du panier"})
}

// ClearCart vide le panier d'un utilisateur
func ClearCart(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		userID = "test-user"
	}

	_, err := database.DB.Exec("DELETE FROM cart_items WHERE user_id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression du panier: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Panier vidé avec succès"})
}
