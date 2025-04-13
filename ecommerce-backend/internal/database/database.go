package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// DB est l'instance globale de la base de données
var DB *sql.DB

// InitDB initialise la connexion à la base de données et crée les tables si elles n'existent pas
func InitDB(dataSourceName string) error {
	var err error

	// S'assurer que le répertoire existe
	dir := filepath.Dir(dataSourceName)
	if dir != "." && dir != "/" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Ouvrir la connexion à la base de données
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}

	// Vérifier la connexion
	if err = DB.Ping(); err != nil {
		return err
	}

	log.Println("Connexion à la base de données établie")

	// Créer les tables si elles n'existent pas
	if err = createTables(DB); err != nil {
		return err
	}

	// Insérer des données de test si la base de données est vide
	if err = insertTestData(DB); err != nil {
		return err
	}

	return nil
}

// createTables crée les tables nécessaires dans la base de données
func createTables(db *sql.DB) error {
	// Table des produits
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		price REAL NOT NULL,
		description TEXT NOT NULL,
		category TEXT NOT NULL,
		image TEXT NOT NULL,
		stock INTEGER NOT NULL
	)`)
	if err != nil {
		return err
	}

	// Table du panier
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS cart_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT NOT NULL,
		product_id INTEGER NOT NULL,
		quantity INTEGER NOT NULL,
		FOREIGN KEY (product_id) REFERENCES products(id)
	)`)
	if err != nil {
		return err
	}

	// Table des commandes
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT NOT NULL,
		status TEXT NOT NULL,
		total_price REAL NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return err
	}

	// Table des produits dans les commandes
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS order_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		order_id INTEGER NOT NULL,
		product_id INTEGER NOT NULL,
		quantity INTEGER NOT NULL,
		price REAL NOT NULL,
		FOREIGN KEY (order_id) REFERENCES orders(id),
		FOREIGN KEY (product_id) REFERENCES products(id)
	)`)
	if err != nil {
		return err
	}

	log.Println("Tables créées avec succès")
	return nil
}

// insertTestData insère des données de test si la base est vide
func insertTestData(db *sql.DB) error {
	// Vérifier si des produits existent déjà
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		return err
	}

	// Si des produits existent déjà, ne rien faire
	if count > 0 {
		log.Println("Les données de test existent déjà")
		return nil
	}

	// Insérer des produits de test
	products := [][]interface{}{
		{"Chaise Brutaliste", 199.99, "Chaise moderne avec design néobrutaliste, robuste et élégante", "meubles", "https://images.unsplash.com/photo-1503602642458-232111445657", 10},
		{"Lampe Industrielle", 89.99, "Lampe de style industriel avec finition en métal brut", "luminaires", "https://images.unsplash.com/photo-1507473885765-e6ed057f782c", 15},
		{"Table Basse Béton", 299.99, "Table basse en béton brut avec piètement en acier", "meubles", "https://images.unsplash.com/photo-1533090368676-1fd25485db88", 5},
		{"Étagère Modulaire", 149.99, "Étagère modulaire en bois brut et métal", "meubles", "https://images.unsplash.com/photo-1594620302200-9a762244a156", 8},
		{"Horloge Minimaliste", 59.99, "Horloge murale avec design minimaliste et brut", "décoration", "https://images.unsplash.com/photo-1563861826100-9cb868fdbe1c", 20},
		{"Vase Géométrique", 39.99, "Vase avec formes géométriques et textures brutes", "décoration", "https://images.unsplash.com/photo-1612375867039-3e8162e75566", 12},
	}

	// Commencer une transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Préparer la requête d'insertion
	stmt, err := tx.Prepare(`
		INSERT INTO products (name, price, description, category, image, stock) 
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Insérer chaque produit
	for _, product := range products {
		_, err = stmt.Exec(product...)
		if err != nil {
			return err
		}
	}

	// Valider la transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	log.Println("Données de test insérées avec succès")
	return nil
}
