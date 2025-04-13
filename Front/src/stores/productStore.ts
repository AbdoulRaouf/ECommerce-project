import { defineStore } from "pinia";
import { computed, ref, watch } from "vue";
import { CartService, ProductService } from "../services/api";
import { useAuthStore } from "./authStore";

export interface Product {
  id: number;
  name: string;
  price: number;
  description: string;
  category: string;
  image: string;
  stock: number;
}

export interface CartItem {
  id?: number;
  userId?: string;
  productId: number;
  quantity: number;
  product?: Product;
}

export const useProductStore = defineStore("product", () => {
  // État des produits et du panier
  const products = ref<Product[]>([]);
  const cart = ref<CartItem[]>([]); // Initialisation avec un tableau vide au lieu de null
  const filterCategory = ref<string | null>(null);
  const searchQuery = ref("");
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  const authStore = useAuthStore();

  // Getters
  const filteredProducts = computed(() => {
    let result = products.value;

    // Filtrer par catégorie
    if (filterCategory.value) {
      result = result.filter(
        (product) => product.category === filterCategory.value
      );
    }

    // Filtrer par recherche
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      result = result.filter(
        (product) =>
          product.name.toLowerCase().includes(query) ||
          product.description.toLowerCase().includes(query)
      );
    }

    return result;
  });

  const categories = computed(() => {
    const uniqueCategories = new Set<string>();
    products.value.forEach((product) => uniqueCategories.add(product.category));
    return Array.from(uniqueCategories);
  });

  const cartTotal = computed(() => {
    // S'assurer que cart.value n'est jamais null
    if (!cart.value || cart.value.length === 0) return 0;

    return cart.value.reduce((total, item) => {
      const product = item.product;
      if (product) {
        return total + product.price * item.quantity;
      }
      return total;
    }, 0);
  });

  const cartItemsCount = computed(() => {
    // S'assurer que cart.value n'est jamais null
    if (!cart.value || cart.value.length === 0) return 0;

    return cart.value.reduce((count, item) => count + item.quantity, 0);
  });

  // Actions
  async function fetchProducts() {
    isLoading.value = true;
    error.value = null;
    try {
      const response = await ProductService.getAll();
      products.value = response.data;
    } catch (err) {
      console.error("Erreur lors de la récupération des produits:", err);
      error.value = "Impossible de charger les produits.";
    } finally {
      isLoading.value = false;
    }
  }

  async function fetchProduct(id: number) {
    isLoading.value = true;
    error.value = null;
    try {
      const response = await ProductService.getById(id);
      // Mettre à jour le produit dans la liste ou l'ajouter s'il n'existe pas
      const index = products.value.findIndex((p) => p.id === id);
      if (index !== -1) {
        products.value[index] = response.data;
      } else {
        products.value.push(response.data);
      }
      return response.data;
    } catch (err) {
      console.error(`Erreur lors de la récupération du produit ${id}:`, err);
      error.value = "Impossible de charger les détails du produit.";
      return null;
    } finally {
      isLoading.value = false;
    }
  }

  async function fetchCategories() {
    try {
      const response = await ProductService.getCategories();
      return response.data;
    } catch (err) {
      console.error("Erreur lors de la récupération des catégories:", err);
      return [];
    }
  }

  async function fetchCart() {
    // Ne récupérer le panier que si l'utilisateur est connecté
    if (!authStore.isLoggedIn) {
      cart.value = [];
      return;
    }

    isLoading.value = true;
    error.value = null;
    try {
      const response = await CartService.getCart();
      cart.value = response.data || []; // S'assurer que cart.value n'est jamais null
    } catch (err) {
      console.error("Erreur lors de la récupération du panier:", err);
      error.value = "Impossible de charger le panier.";
      cart.value = []; // Initialiser à un tableau vide en cas d'erreur
    } finally {
      isLoading.value = false;
    }
  }

  async function addToCart(product: Product, quantity = 1) {
    // Vérifier si l'utilisateur est connecté
    if (!authStore.isLoggedIn) {
      error.value =
        "Veuillez vous connecter pour ajouter des produits au panier.";
      return null;
    }

    isLoading.value = true;
    error.value = null;
    try {
      const response = await CartService.addToCart(product.id, quantity);
      await fetchCart(); // Recharger le panier pour avoir les données à jour
      return response.data;
    } catch (err) {
      console.error("Erreur lors de l'ajout au panier:", err);
      error.value = "Impossible d'ajouter le produit au panier.";
      return null;
    } finally {
      isLoading.value = false;
    }
  }

  async function removeFromCart(itemId: number) {
    isLoading.value = true;
    error.value = null;
    try {
      await CartService.removeFromCart(itemId);
      await fetchCart(); // Recharger le panier pour avoir les données à jour
    } catch (err) {
      console.error("Erreur lors de la suppression du panier:", err);
      error.value = "Impossible de supprimer le produit du panier.";
    } finally {
      isLoading.value = false;
    }
  }

  async function updateCartItemQuantity(itemId: number, quantity: number) {
    isLoading.value = true;
    error.value = null;
    try {
      if (quantity <= 0) {
        await removeFromCart(itemId);
      } else {
        await CartService.updateCartItem(itemId, quantity);
        await fetchCart(); // Recharger le panier pour avoir les données à jour
      }
    } catch (err) {
      console.error("Erreur lors de la mise à jour du panier:", err);
      error.value = "Impossible de mettre à jour la quantité.";
    } finally {
      isLoading.value = false;
    }
  }

  async function clearCart() {
    isLoading.value = true;
    error.value = null;
    try {
      await CartService.clearCart();
      cart.value = [];
    } catch (err) {
      console.error("Erreur lors du vidage du panier:", err);
      error.value = "Impossible de vider le panier.";
    } finally {
      isLoading.value = false;
    }
  }

  function setFilterCategory(category: string | null) {
    filterCategory.value = category;
  }

  function setSearchQuery(query: string) {
    searchQuery.value = query;
  }

  // Fonction pour surveiller les changements d'authentification
  function watchAuthState() {
    watch(
      () => authStore.isLoggedIn,
      (newIsLoggedIn) => {
        if (newIsLoggedIn) {
          // Si l'utilisateur vient de se connecter, charger son panier
          fetchCart();
        } else {
          // Si l'utilisateur s'est déconnecté, vider le panier local
          cart.value = [];
        }
      }
    );
  }

  // Initialiser le store en chargeant les produits
  fetchProducts();

  // Initialiser le panier si l'utilisateur est déjà connecté
  if (authStore.isLoggedIn) {
    fetchCart();
  }

  // Commencer à surveiller les changements d'état d'authentification
  watchAuthState();

  return {
    // État
    products,
    cart,
    isLoading,
    error,
    filterCategory,
    searchQuery,

    // Getters
    filteredProducts,
    categories,
    cartTotal,
    cartItemsCount,

    // Actions
    fetchProducts,
    fetchProduct,
    fetchCategories,
    fetchCart,
    addToCart,
    removeFromCart,
    updateCartItemQuantity,
    clearCart,
    setFilterCategory,
    setSearchQuery,
  };
});
