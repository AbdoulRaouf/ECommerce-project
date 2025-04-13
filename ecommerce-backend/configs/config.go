package configs

import (
	"flag"
	"os"
	"path/filepath"
)

// Config représente la configuration de l'application
type Config struct {
	Port           string
	DatabasePath   string
	AllowedOrigins string
	ProductionMode bool
}

// LoadConfig charge la configuration depuis les arguments de ligne de commande ou les valeurs par défaut
func LoadConfig() Config {
	config := Config{}

	// Définir les drapeaux de ligne de commande
	flag.StringVar(&config.Port, "port", "8080", "Le port sur lequel le serveur écoute")
	flag.StringVar(&config.DatabasePath, "db", "./data/ecommerce.db", "Le chemin vers la base de données SQLite")
	flag.StringVar(&config.AllowedOrigins, "origins", "*", "Les origines autorisées pour CORS")
	flag.BoolVar(&config.ProductionMode, "prod", false, "Exécute le serveur en mode production")

	// Analyser les drapeaux
	flag.Parse()

	// Utiliser les variables d'environnement si elles existent
	if port := os.Getenv("PORT"); port != "" {
		config.Port = port
	}

	if dbPath := os.Getenv("DB_PATH"); dbPath != "" {
		config.DatabasePath = dbPath
	}

	if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
		config.AllowedOrigins = origins
	}

	if prod := os.Getenv("PRODUCTION"); prod == "true" {
		config.ProductionMode = true
	}

	// Convertir le chemin relatif en chemin absolu
	if !filepath.IsAbs(config.DatabasePath) {
		absPath, err := filepath.Abs(config.DatabasePath)
		if err == nil {
			config.DatabasePath = absPath
		}
	}

	return config
}
