import { AdminService } from "@/services/api";
import { defineStore } from "pinia";
import { computed, ref } from "vue";

export interface DashboardStats {
  totalOrders: number;
  totalProducts: number;
  totalRevenue: number;
  recentOrders: Order[];
  lowStockProducts: Product[];
}

export interface Order {
  id: number;
  userId: string;
  status: string;
  totalPrice: number;
  createdAt: string;
  items?: CartItem[];
}

export interface CartItem {
  id?: number;
  userId?: string;
  productId: number;
  quantity: number;
  product?: Product;
}

export interface Product {
  id?: number;
  name: string;
  price: number;
  description: string;
  category: string;
  image: string;
  stock: number;
}

export const useAdminStore = defineStore("admin", () => {
  // État
  const isLoading = ref(false);
  const error = ref<string | null>(null);
  const dashboardStats = ref<DashboardStats | null>(null);
  const products = ref<Product[]>([]);
  const orders = ref<Order[]>([]);
  const selectedOrder = ref<Order | null>(null);

  // Données calculées
  const isAuthenticated = computed(() => AdminService.isLoggedIn());

  // Actions - Authentification
  async function login(username: string, password: string): Promise<boolean> {
    isLoading.value = true;
    error.value = null;

    try {
      const response = await AdminService.login(username, password);
      const token = response.data.token;

      if (token) {
        localStorage.setItem("adminToken", token);
        return true;
      }

      error.value = "Token d'authentification non reçu";
      return false;
    } catch (err: any) {
      error.value = err.response?.data?.error || "Erreur lors de la connexion";
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  function logout(): void {
    AdminService.logout();
  }

  // Actions - Dashboard
  async function fetchDashboardStats(): Promise<void> {
    isLoading.value = true;
    error.value = null;

    try {
      const response = await AdminService.getDashboardStats();
      dashboardStats.value = response.data;
    } catch (err: any) {
      error.value =
        err.response?.data?.error ||
        "Erreur lors du chargement des statistiques";
    } finally {
      isLoading.value = false;
    }
  }

  // Actions - Produits
  async function fetchProducts(): Promise<void> {
    isLoading.value = true;
    error.value = null;

    try {
      const response = await AdminService.getProducts();
      products.value = response.data;
    } catch (err: any) {
      error.value =
        err.response?.data?.error || "Erreur lors du chargement des produits";
    } finally {
      isLoading.value = false;
    }
  }

  async function createProduct(product: Product): Promise<boolean> {
    isLoading.value = true;
    error.value = null;

    try {
      await AdminService.createProduct(product);
      await fetchProducts(); // Rafraîchir la liste des produits
      return true;
    } catch (err: any) {
      error.value =
        err.response?.data?.error || "Erreur lors de la création du produit";
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  async function updateProduct(
    id: number,
    productData: Product
  ): Promise<boolean> {
    isLoading.value = true;
    error.value = null;

    try {
      await AdminService.updateProduct(id, productData);
      await fetchProducts(); // Rafraîchir la liste des produits
      return true;
    } catch (err: any) {
      error.value =
        err.response?.data?.error || "Erreur lors de la mise à jour du produit";
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  async function deleteProduct(id: number): Promise<boolean> {
    isLoading.value = true;
    error.value = null;

    try {
      await AdminService.deleteProduct(id);
      await fetchProducts(); // Rafraîchir la liste des produits
      return true;
    } catch (err: any) {
      error.value =
        err.response?.data?.error || "Erreur lors de la suppression du produit";
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  // Actions - Commandes
  async function fetchOrders(status?: string): Promise<void> {
    isLoading.value = true;
    error.value = null;

    try {
      const params = status ? { status } : {};
      const response = await AdminService.getOrders(params);
      orders.value = response.data;
    } catch (err: any) {
      error.value =
        err.response?.data?.error || "Erreur lors du chargement des commandes";
    } finally {
      isLoading.value = false;
    }
  }

  async function fetchOrderDetails(id: number): Promise<void> {
    isLoading.value = true;
    error.value = null;

    try {
      const response = await AdminService.getOrderDetails(id);
      selectedOrder.value = response.data;
    } catch (err: any) {
      error.value =
        err.response?.data?.error ||
        "Erreur lors du chargement des détails de la commande";
    } finally {
      isLoading.value = false;
    }
  }

  async function updateOrderStatus(
    id: number,
    status: string
  ): Promise<boolean> {
    isLoading.value = true;
    error.value = null;

    try {
      await AdminService.updateOrderStatus(id, status);
      if (selectedOrder.value && selectedOrder.value.id === id) {
        selectedOrder.value.status = status;
      }
      await fetchOrders(); // Rafraîchir la liste des commandes
      return true;
    } catch (err: any) {
      error.value =
        err.response?.data?.error || "Erreur lors de la mise à jour du statut";
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  async function deleteOrder(id: number): Promise<boolean> {
    isLoading.value = true;
    error.value = null;

    try {
      await AdminService.deleteOrder(id);
      if (selectedOrder.value && selectedOrder.value.id === id) {
        selectedOrder.value = null;
      }
      await fetchOrders(); // Rafraîchir la liste des commandes
      return true;
    } catch (err: any) {
      error.value =
        err.response?.data?.error ||
        "Erreur lors de la suppression de la commande";
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  return {
    // État
    isLoading,
    error,
    dashboardStats,
    products,
    orders,
    selectedOrder,

    // Getters
    isAuthenticated,

    // Actions - Authentification
    login,
    logout,

    // Actions - Dashboard
    fetchDashboardStats,

    // Actions - Produits
    fetchProducts,
    createProduct,
    updateProduct,
    deleteProduct,

    // Actions - Commandes
    fetchOrders,
    fetchOrderDetails,
    updateOrderStatus,
    deleteOrder,
  };
});
