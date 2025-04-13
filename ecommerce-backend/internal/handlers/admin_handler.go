package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thiomajid/ecommerce-backend/internal/database"
	"github.com/thiomajid/ecommerce-backend/internal/models"
)

// AdminGetAllOrders récupère toutes les commandes pour l'administrateur
func AdminGetAllOrders(c *gin.Context) {
	status := c.Query("status")

	var query string
	var args []interface{}

	query = `
		SELECT id, user_id, status, total_price, created_at
		FROM orders
	`

	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC"

	rows, err := database.DB.Query(query, args...)
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

// AdminGetOrderDetails récupère les détails d'une commande pour l'administrateur
func AdminGetOrderDetails(c *gin.Context) {
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
		WHERE id = ?
	`, orderID).Scan(&order.ID, &order.UserID, &order.Status, &order.TotalPrice, &order.CreatedAt)

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
		item.UserID = order.UserID

		items = append(items, item)
	}

	order.Items = items
	c.JSON(http.StatusOK, order)
}

// AdminUpdateOrderStatus met à jour le statut d'une commande
func AdminUpdateOrderStatus(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de commande invalide"})
		return
	}

	var updateData struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Vérifier que le statut est valide
	validStatuses := map[string]bool{
		"pending":   true,
		"confirmed": true,
		"shipped":   true,
		"delivered": true,
		"cancelled": true,
	}

	if !validStatuses[updateData.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Statut invalide"})
		return
	}

	// Vérifier que la commande existe
	var exists bool
	err = database.DB.QueryRow("SELECT 1 FROM orders WHERE id = ?", orderID).Scan(&exists)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Commande non trouvée"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la vérification de la commande: " + err.Error()})
		return
	}

	// Mise à jour du statut
	_, err = database.DB.Exec("UPDATE orders SET status = ? WHERE id = ?", updateData.Status, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du statut: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Statut mis à jour avec succès",
		"orderId": orderID,
		"status":  updateData.Status,
	})
}

// AdminDeleteOrder supprime une commande (à utiliser avec précaution)
func AdminDeleteOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de commande invalide"})
		return
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

	// Supprimer d'abord les éléments de la commande
	_, err = tx.Exec("DELETE FROM order_items WHERE order_id = ?", orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression des éléments de la commande: " + err.Error()})
		return
	}

	// Puis supprimer la commande
	result, err := tx.Exec("DELETE FROM orders WHERE id = ?", orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression de la commande: " + err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Commande non trouvée"})
		return
	}

	// Valider la transaction
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la finalisation de la suppression: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Commande supprimée avec succès"})
}

// AdminGetDashboardStats récupère des statistiques pour le tableau de bord admin
func AdminGetDashboardStats(c *gin.Context) {
	// Récupérer le nombre total de commandes
	var totalOrders int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM orders").Scan(&totalOrders)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des statistiques: " + err.Error()})
		return
	}

	// Récupérer le nombre total de produits
	var totalProducts int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&totalProducts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des statistiques: " + err.Error()})
		return
	}

	// Récupérer le chiffre d'affaires total
	var totalRevenue float64
	err = database.DB.QueryRow("SELECT SUM(total_price) FROM orders").Scan(&totalRevenue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des statistiques: " + err.Error()})
		return
	}

	// Récupérer les commandes récentes (5 dernières)
	rows, err := database.DB.Query(`
		SELECT id, user_id, status, total_price, created_at 
		FROM orders 
		ORDER BY created_at DESC LIMIT 5
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des commandes récentes: " + err.Error()})
		return
	}
	defer rows.Close()

	var recentOrders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.Status, &order.TotalPrice, &order.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du scan des commandes récentes: " + err.Error()})
			return
		}
		recentOrders = append(recentOrders, order)
	}

	// Récupérer les produits à faible stock (moins de 5 unités)
	rows, err = database.DB.Query(`
		SELECT id, name, price, description, category, image, stock 
		FROM products 
		WHERE stock < 5 
		ORDER BY stock ASC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des produits à faible stock: " + err.Error()})
		return
	}
	defer rows.Close()

	var lowStockProducts []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Description, &product.Category, &product.Image, &product.Stock); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du scan des produits à faible stock: " + err.Error()})
			return
		}
		lowStockProducts = append(lowStockProducts, product)
	}

	// Retourner les statistiques
	c.JSON(http.StatusOK, gin.H{
		"totalOrders":      totalOrders,
		"totalProducts":    totalProducts,
		"totalRevenue":     totalRevenue,
		"recentOrders":     recentOrders,
		"lowStockProducts": lowStockProducts,
	})
}
