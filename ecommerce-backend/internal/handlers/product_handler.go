package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thiomajid/ecommerce-backend/internal/database"
	"github.com/thiomajid/ecommerce-backend/internal/models"
)

// GetAllProducts récupère tous les produits
func GetAllProducts(c *gin.Context) {
	category := c.Query("category")
	query := c.Query("query")

	// Préparer la requête SQL de base
	sqlQuery := "SELECT id, name, price, description, category, image, stock FROM products"
	var args []interface{}
	var conditions []string

	// Ajouter des conditions à la requête si nécessaire
	if category != "" {
		conditions = append(conditions, "category = ?")
		args = append(args, category)
	}

	if query != "" {
		// Rechercher dans le nom et la description
		conditions = append(conditions, "(name LIKE ? OR description LIKE ?)")
		args = append(args, "%"+query+"%", "%"+query+"%")
	}

	// Compléter la requête avec les conditions WHERE
	if len(conditions) > 0 {
		sqlQuery += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			sqlQuery += " AND " + conditions[i]
		}
	}

	// Exécuter la requête
	rows, err := database.DB.Query(sqlQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des produits: " + err.Error()})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Category, &p.Image, &p.Stock); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du scan des produits: " + err.Error()})
			return
		}
		products = append(products, p)
	}

	c.JSON(http.StatusOK, products)
}

// GetProduct récupère un produit par son ID
func GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de produit invalide"})
		return
	}

	var p models.Product
	err = database.DB.QueryRow(
		"SELECT id, name, price, description, category, image, stock FROM products WHERE id = ?",
		id,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Category, &p.Image, &p.Stock)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produit non trouvé"})
		return
	}

	c.JSON(http.StatusOK, p)
}

// GetCategories récupère toutes les catégories de produits
func GetCategories(c *gin.Context) {
	rows, err := database.DB.Query("SELECT DISTINCT category FROM products")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des catégories"})
		return
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du scan des catégories"})
			return
		}
		categories = append(categories, category)
	}

	c.JSON(http.StatusOK, categories)
}

// CreateProduct crée un nouveau produit
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := database.DB.Exec(
		"INSERT INTO products (name, price, description, category, image, stock) VALUES (?, ?, ?, ?, ?, ?)",
		product.Name, product.Price, product.Description, product.Category, product.Image, product.Stock,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du produit: " + err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	product.ID = int(id)

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct met à jour un produit existant
func UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de produit invalide"})
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product.ID = id

	_, err = database.DB.Exec(
		"UPDATE products SET name = ?, price = ?, description = ?, category = ?, image = ?, stock = ? WHERE id = ?",
		product.Name, product.Price, product.Description, product.Category, product.Image, product.Stock, product.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du produit: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct supprime un produit
func DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de produit invalide"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression du produit: " + err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produit non trouvé"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produit supprimé avec succès"})
}
