import { defineStore } from "pinia";
import { computed, ref, watch } from "vue";
import { useAuthStore } from "./authStore";
import { Product } from "./productStore";

export interface WishlistItem {
  id?: number;
  userId?: string;
  productId: number;
  product?: Product;
  createdAt?: string;
}

export const useWishlistStore = defineStore("wishlist", () => {
  // État de la liste de souhaits
  const wishlist = ref<WishlistItem[]>([]);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  const authStore = useAuthStore();

  // Getters
  const wishlistCount = computed(() => wishlist.value.length);

  const isInWishlist = (productId: number) => {
    return wishlist.value.some((item) => item.productId === productId);
  };

  // Actions
  async function fetchWishlist() {
    // Ne récupérer la liste que si l'utilisateur est connecté
    if (!authStore.isLoggedIn) {
      wishlist.value = [];
      return;
    }

    isLoading.value = true;
    error.value = null;
    try {
      // Simulation API - à remplacer par un vrai appel API
      // const response = await WishlistService.getWishlist();
      // wishlist.value = response.data || [];

      // Pour le moment, utiliser localStorage comme stockage temporaire
      const storedWishlist = localStorage.getItem(
        `wishlist_${authStore.user?.id}`
      );
      wishlist.value = storedWishlist ? JSON.parse(storedWishlist) : [];
    } catch (err) {
      console.error(
        "Erreur lors de la récupération de la liste de souhaits:",
        err
      );
      error.value = "Impossible de charger la liste de souhaits.";
      wishlist.value = [];
    } finally {
      isLoading.value = false;
    }
  }

  async function addToWishlist(product: Product) {
    // Vérifier si l'utilisateur est connecté
    if (!authStore.isLoggedIn) {
      error.value =
        "Veuillez vous connecter pour ajouter des produits à votre liste.";
      return null;
    }

    // Vérifier si le produit est déjà dans la liste
    if (isInWishlist(product.id)) {
      return null;
    }

    isLoading.value = true;
    error.value = null;
    try {
      const newItem: WishlistItem = {
        id: Date.now(), // Générer un ID temporaire
        userId: authStore.user?.id,
        productId: product.id,
        product: product,
        createdAt: new Date().toISOString(),
      };

      // Simulation API - à remplacer par un vrai appel API
      // const response = await WishlistService.addToWishlist(product.id);
      // await fetchWishlist();

      // Pour le moment, utiliser localStorage comme stockage temporaire
      wishlist.value.push(newItem);
      saveWishlistToStorage();

      return newItem;
    } catch (err) {
      console.error("Erreur lors de l'ajout à la liste de souhaits:", err);
      error.value = "Impossible d'ajouter le produit à la liste.";
      return null;
    } finally {
      isLoading.value = false;
    }
  }

  async function removeFromWishlist(productId: number) {
    isLoading.value = true;
    error.value = null;
    try {
      // Simulation API - à remplacer par un vrai appel API
      // await WishlistService.removeFromWishlist(productId);

      // Pour le moment, utiliser localStorage comme stockage temporaire
      wishlist.value = wishlist.value.filter(
        (item) => item.productId !== productId
      );
      saveWishlistToStorage();
    } catch (err) {
      console.error(
        "Erreur lors de la suppression de la liste de souhaits:",
        err
      );
      error.value = "Impossible de supprimer le produit de la liste.";
    } finally {
      isLoading.value = false;
    }
  }

  async function clearWishlist() {
    isLoading.value = true;
    error.value = null;
    try {
      // Simulation API - à remplacer par un vrai appel API
      // await WishlistService.clearWishlist();

      // Pour le moment, utiliser localStorage comme stockage temporaire
      wishlist.value = [];
      saveWishlistToStorage();
    } catch (err) {
      console.error("Erreur lors du vidage de la liste de souhaits:", err);
      error.value = "Impossible de vider la liste de souhaits.";
    } finally {
      isLoading.value = false;
    }
  }

  // Helper pour sauvegarder dans localStorage
  function saveWishlistToStorage() {
    if (authStore.user?.id) {
      localStorage.setItem(
        `wishlist_${authStore.user.id}`,
        JSON.stringify(wishlist.value)
      );
    }
  }

  // Fonction pour surveiller les changements d'authentification
  function watchAuthState() {
    watch(
      () => authStore.isLoggedIn,
      (newIsLoggedIn) => {
        if (newIsLoggedIn) {
          // Si l'utilisateur vient de se connecter, charger sa liste
          fetchWishlist();
        } else {
          // Si l'utilisateur s'est déconnecté, vider la liste locale
          wishlist.value = [];
        }
      }
    );
  }

  // Initialiser la liste si l'utilisateur est déjà connecté
  if (authStore.isLoggedIn) {
    fetchWishlist();
  }

  // Commencer à surveiller les changements d'état d'authentification
  watchAuthState();

  return {
    // État
    wishlist,
    isLoading,
    error,

    // Getters
    wishlistCount,
    isInWishlist,

    // Actions
    fetchWishlist,
    addToWishlist,
    removeFromWishlist,
    clearWishlist,
  };
});
