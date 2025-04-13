package models

// Product représente un produit dans notre boutique
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Image       string  `json:"image"`
	Stock       int     `json:"stock"`
}

// CartItem représente un élément du panier d'achat
type CartItem struct {
	ID        int     `json:"id"`
	UserID    string  `json:"userId"`
	ProductID int     `json:"productId"`
	Quantity  int     `json:"quantity"`
	Product   Product `json:"product"`
}

// Order représente une commande
type Order struct {
	ID         int        `json:"id"`
	UserID     string     `json:"userId"`
	Status     string     `json:"status"`
	TotalPrice float64    `json:"totalPrice"`
	CreatedAt  string     `json:"createdAt"`
	Items      []CartItem `json:"items"`
}
