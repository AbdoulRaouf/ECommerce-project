import { defineStore } from "pinia";
import { computed, ref } from "vue";

export const useAuthStore = defineStore("auth", () => {
  // État
  const user = ref(null);
  const token = ref(localStorage.getItem("auth_token") || null);
  const isLoading = ref(false);
  const error = ref(null);

  // Getters
  const isLoggedIn = computed(() => !!token.value);
  const username = computed(
    () => user.value?.name?.split(" ")[0] || "Utilisateur"
  );

  // Actions
  async function login(email, password) {
    isLoading.value = true;
    error.value = null;

    try {
      // Simulation d'une API de connexion (à remplacer par de vraies requêtes API)
      await new Promise((resolve) => setTimeout(resolve, 800));

      // Vérifier les identifiants (exemple simplifié)
      if (email === "user@example.com" && password === "password") {
        const userData = {
          id: 1,
          name: "John Doe",
          email: email,
          address: "123 Rue de Paris, 75000 Paris",
          phone: "+33123456789",
          createdAt: "2023-09-15T14:30:00Z",
        };

        // Sauvegarder les données utilisateur et le token
        user.value = userData;
        const authToken =
          "simulated_jwt_token_" + Math.random().toString(36).substr(2);
        token.value = authToken;
        localStorage.setItem("auth_token", authToken);

        return true;
      } else {
        throw new Error("Email ou mot de passe incorrect");
      }
    } catch (err) {
      console.error("Erreur de connexion:", err);
      error.value =
        err.message || "Erreur lors de la connexion. Veuillez réessayer.";
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  async function register(userData) {
    isLoading.value = true;
    error.value = null;

    try {
      // Simulation d'une API d'inscription (à remplacer par de vraies requêtes API)
      await new Promise((resolve) => setTimeout(resolve, 1000));

      // Simulation d'une création de compte réussie
      const newUser = {
        id: 2,
        name: userData.name,
        email: userData.email,
        address: userData.address || null,
        phone: userData.phone || null,
        createdAt: new Date().toISOString(),
      };

      // Sauvegarder les données utilisateur et le token
      user.value = newUser;
      const authToken =
        "simulated_jwt_token_" + Math.random().toString(36).substr(2);
      token.value = authToken;
      localStorage.setItem("auth_token", authToken);

      return true;
    } catch (err) {
      console.error("Erreur d'inscription:", err);
      error.value =
        err.message || "Erreur lors de l'inscription. Veuillez réessayer.";
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  function logout() {
    user.value = null;
    token.value = null;
    localStorage.removeItem("auth_token");
  }

  async function fetchUserProfile() {
    if (!token.value) return null;

    isLoading.value = true;
    try {
      // Simulation d'une API pour récupérer le profil utilisateur
      await new Promise((resolve) => setTimeout(resolve, 500));

      // Exemple de profil utilisateur (à remplacer par de vraies données API)
      user.value = {
        id: 1,
        name: "John Doe",
        email: "user@example.com",
        address: "123 Rue de Paris, 75000 Paris",
        phone: "+33123456789",
        createdAt: "2023-09-15T14:30:00Z",
      };

      return user.value;
    } catch (err) {
      console.error("Erreur de récupération du profil:", err);
      error.value = "Impossible de charger votre profil. Veuillez réessayer.";
      return null;
    } finally {
      isLoading.value = false;
    }
  }

  async function fetchUserOrders() {
    if (!token.value) return [];

    // Simulation d'une API pour récupérer les commandes de l'utilisateur
    await new Promise((resolve) => setTimeout(resolve, 700));

    // Exemple de commandes (à remplacer par de vraies données API)
    return [
      {
        id: 1001,
        status: "delivered",
        createdAt: "2025-03-15T10:30:00Z",
        totalPrice: 129.99,
        subtotalPrice: 119.99,
        shippingCost: 5.0,
        taxAmount: 5.0,
        items: [
          {
            id: 1,
            product: {
              id: 101,
              name: "Chaise design",
              price: 59.99,
              image: "https://via.placeholder.com/150",
            },
            quantity: 2,
          },
        ],
      },
      {
        id: 1002,
        status: "shipped",
        createdAt: "2025-04-01T14:20:00Z",
        totalPrice: 249.5,
        subtotalPrice: 239.5,
        shippingCost: 0.0,
        taxAmount: 10.0,
        items: [
          {
            id: 2,
            product: {
              id: 102,
              name: "Table basse",
              price: 149.5,
              image: "https://via.placeholder.com/150",
            },
            quantity: 1,
          },
          {
            id: 3,
            product: {
              id: 103,
              name: "Lampe de bureau",
              price: 45.0,
              image: "https://via.placeholder.com/150",
            },
            quantity: 2,
          },
        ],
      },
      {
        id: 1003,
        status: "pending",
        createdAt: "2025-04-10T09:15:00Z",
        totalPrice: 75.9,
        subtotalPrice: 69.9,
        shippingCost: 3.0,
        taxAmount: 3.0,
        items: [
          {
            id: 4,
            product: {
              id: 104,
              name: "Coussin décoratif",
              price: 23.3,
              image: "https://via.placeholder.com/150",
            },
            quantity: 3,
          },
        ],
      },
    ];
  }

  // Initialiser le profil utilisateur si déjà connecté
  function init() {
    if (token.value) {
      fetchUserProfile();
    }
  }

  // Appeler init au démarrage
  init();

  return {
    user,
    token,
    isLoading,
    error,
    isLoggedIn,
    username,
    login,
    register,
    logout,
    fetchUserProfile,
    fetchUserOrders,
  };
});
