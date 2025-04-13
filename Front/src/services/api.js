import axios from "axios";

const API_URL = "http://localhost:8080/api";

// Création d'une instance axios configurée
const apiClient = axios.create({
  baseURL: API_URL,
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 10000,
});

// Intercepteur pour ajouter l'ID utilisateur à toutes les requêtes
apiClient.interceptors.request.use((config) => {
  // Simuler un ID utilisateur (dans une vraie application, on récupérerait l'identifiant de l'utilisateur connecté)
  config.headers["User-ID"] = localStorage.getItem("userId") || "test-user";

  // Ajouter le token utilisateur si disponible
  const token = localStorage.getItem("token");
  if (token && !config.url.startsWith("/admin")) {
    config.headers["Authorization"] = `Bearer ${token}`;
  }

  // Ajouter le token admin si disponible
  const adminToken = localStorage.getItem("adminToken");
  if (adminToken && config.url.startsWith("/admin")) {
    config.headers["Authorization"] = `Bearer ${adminToken}`;
  }

  return config;
});

// Configuration de l'intercepteur pour les erreurs d'authentification
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      // Rediriger vers la page de connexion si token expiré ou invalide
      if (error.config.url.startsWith("/admin")) {
        localStorage.removeItem("adminToken");
        window.location.href = "/admin/login";
      } else {
        localStorage.removeItem("token");
        window.location.href = "/login";
      }
    }
    return Promise.reject(error);
  }
);

// Service pour les produits
export const ProductService = {
  // Récupérer tous les produits avec filtres optionnels
  getAll: (params = {}) => apiClient.get("/products", { params }),

  // Récupérer un produit spécifique par son ID
  getById: (id) => apiClient.get(`/products/${id}`),

  // Récupérer toutes les catégories disponibles
  getCategories: () => apiClient.get("/categories"),

  // Rechercher des produits
  searchProducts: (query) =>
    apiClient.get("/products/search", { params: { q: query } }),
};

// Service pour le panier d'achat
export const CartService = {
  // Récupérer le contenu du panier
  getCart: () => apiClient.get("/cart"),

  // Ajouter un produit au panier
  addToCart: (productId, quantity = 1) =>
    apiClient.post("/cart/items", { productId, quantity }),

  // Mettre à jour la quantité d'un produit dans le panier
  updateCartItem: (itemId, quantity) =>
    apiClient.put(`/cart/items/${itemId}`, { quantity }),

  // Supprimer un produit du panier
  removeFromCart: (itemId) => apiClient.delete(`/cart/items/${itemId}`),

  // Vider complètement le panier
  clearCart: () => apiClient.delete("/cart"),
};

// Service pour l'authentification
export const AuthService = {
  // Connexion utilisateur
  login: (email, password) =>
    apiClient.post("/auth/login", { email, password }),

  // Inscription utilisateur
  register: (userData) => apiClient.post("/auth/register", userData),

  // Vérifier si l'utilisateur est connecté
  isLoggedIn: () => !!localStorage.getItem("token"),

  // Déconnexion
  logout: () => {
    localStorage.removeItem("token");
  },

  // Obtenir le profil utilisateur
  getUserProfile: () => apiClient.get("/users/profile"),

  // Obtenir l'historique des commandes
  getUserOrders: () => apiClient.get("/users/orders"),
};

// Service pour les commandes
export const OrderService = {
  // Créer une nouvelle commande à partir du panier actuel
  createOrder: () => apiClient.post("/orders"),

  // Récupérer l'historique des commandes
  getOrders: () => apiClient.get("/orders"),

  // Récupérer les détails d'une commande spécifique
  getOrderDetails: (id) => apiClient.get(`/orders/${id}`),
};

// Service pour l'administration
export const AdminService = {
  // Authentification de l'administrateur
  login: (username, password) =>
    apiClient.post("/admin/login", { username, password }),

  // Vérifier si l'utilisateur est connecté en tant qu'admin
  isLoggedIn: () => !!localStorage.getItem("adminToken"),

  // Déconnexion de l'administrateur
  logout: () => {
    localStorage.removeItem("adminToken");
  },

  // Tableau de bord - statistiques
  getDashboardStats: () => apiClient.get("/admin/dashboard"),

  // Gestion des produits
  getProducts: () => apiClient.get("/admin/products"),

  createProduct: (product) =>
    apiClient.post("/admin/products", product, {
      headers: {
        "Content-Type": "application/json",
      },
    }),

  updateProduct: (id, product) =>
    apiClient.put(`/admin/products/${id}`, product, {
      headers: {
        "Content-Type": "application/json",
      },
    }),

  deleteProduct: (id) => apiClient.delete(`/admin/products/${id}`),

  // Gestion des commandes
  getOrders: (params = {}) => apiClient.get("/admin/orders", { params }),

  getOrderDetails: (id) => apiClient.get(`/admin/orders/${id}`),

  updateOrderStatus: (id, status) =>
    apiClient.patch(`/admin/orders/${id}/status`, { status }),

  deleteOrder: (id) => apiClient.delete(`/admin/orders/${id}`),
};

export default {
  ProductService,
  CartService,
  AuthService,
  OrderService,
  AdminService,
};
