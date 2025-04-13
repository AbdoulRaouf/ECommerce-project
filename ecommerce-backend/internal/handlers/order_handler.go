package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thiomajid/ecommerce-backend/internal/database"
	"github.com/thiomajid/ecommerce-backend/internal/models"
)

// CreateOrder crée une nouvelle commande à partir du panier de l'utilisateur
func CreateOrder(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		userID = "test-user"
	}

	// Vérifier si le panier contient des éléments
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM cart_items WHERE user_id = ?", userID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification du panier: " + err.Error()})
		return
	}

	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le panier est vide"})
		return
	}

	// Récupérer les éléments du panier
	rows, err := database.DB.Query(`
		SELECT c.product_id, c.quantity, p.price, p.stock
		FROM cart_items c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_id = ?
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération du panier: " + err.Error()})
		return
	}
	defer rows.Close()

	type CartItemWithPrice struct {
		ProductID int
		Quantity  int
		Price     float64
		Stock     int
	}

	var items []CartItemWithPrice
	var totalPrice float64 = 0

	// Vérifier le stock pour chaque produit
	for rows.Next() {
		var item CartItemWithPrice
		if err := rows.Scan(&item.ProductID, &item.Quantity, &item.Price, &item.Stock); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du scan des éléments du panier: " + err.Error()})
			return
		}

		// Vérifier si le stock est suffisant
		if item.Quantity > item.Stock {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Stock insuffisant",
				"productID": item.ProductID,
				"requested": item.Quantity,
				"available": item.Stock,
			})
			return
		}

		items = append(items, item)
		totalPrice += item.Price * float64(item.Quantity)
	}

	// Commencer une transaction
	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'initialisation de la transaction: " + err.Error()})
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Créer la commande
	result, err := tx.Exec(
		"INSERT INTO orders (user_id, status, total_price, created_at) VALUES (?, ?, ?, ?)",
		userID, "pending", totalPrice, time.Now().Format(time.RFC3339),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la commande: " + err.Error()})
		return
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération de l'ID de commande: " + err.Error()})
		return
	}

	// Ajouter les éléments à la commande et mettre à jour le stock
	for _, item := range items {
		_, err = tx.Exec(
			"INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)",
			orderID, item.ProductID, item.Quantity, item.Price,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'ajout des éléments à la commande: " + err.Error()})
			return
		}

		// Mettre à jour le stock
		_, err = tx.Exec(
			"UPDATE products SET stock = stock - ? WHERE id = ?",
			item.Quantity, item.ProductID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du stock: " + err.Error()})
			return
		}
	}

	// Vider le panier
	_, err = tx.Exec("DELETE FROM cart_items WHERE user_id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression du panier: " + err.Error()})
		return
	}

	// Valider la transaction
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la finalisation de la commande: " + err.Error()})
		return
	}

	// Récupérer la commande complète
	order := models.Order{
		ID:         int(orderID),
		UserID:     userID,
		Status:     "pending",
		TotalPrice: totalPrice,
		CreatedAt:  time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, order)
}

// GetOrders récupère toutes les commandes d'un utilisateur
func GetOrders(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		userID = "test-user"
	}

	rows, err := database.DB.Query(`
		SELECT id, user_id, status, total_price, created_at
		FROM orders
		WHERE user_id = ?
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des commandes: " + err.Error()})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.Status, &order.TotalPrice, &order.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du scan des commandes: " + err.Error()})
			return
		}
		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrderDetails récupère les détails d'une commande
func GetOrderDetails(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		userID = "test-user"
	}

	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de commande invalide"})
		return
	}

	// Récupérer les informations de la commande
	var order models.Order
	err = database.DB.QueryRow(`
		SELECT id, user_id, status, total_price, created_at
		FROM orders
		WHERE id = ? AND user_id = ?
	`, orderID, userID).Scan(&order.ID, &order.UserID, &order.Status, &order.TotalPrice, &order.CreatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Commande non trouvée"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération de la commande: " + err.Error()})
		return
	}

	// Récupérer les éléments de la commande
	rows, err := database.DB.Query(`
		SELECT oi.id, oi.product_id, oi.quantity, oi.price,
			p.name, p.description, p.category, p.image, p.stock
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		WHERE oi.order_id = ?
	`, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des éléments de la commande: " + err.Error()})
		return
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		var product models.Product

		if err := rows.Scan(
			&item.ID, &item.ProductID, &item.Quantity, &product.Price,
			&product.Name, &product.Description, &product.Category, &product.Image, &product.Stock,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du scan des éléments de la commande: " + err.Error()})
			return
		}

		product.ID = item.ProductID
		item.Product = product
		item.UserID = userID

		items = append(items, item)
	}

	order.Items = items
	c.JSON(http.StatusOK, order)
}
