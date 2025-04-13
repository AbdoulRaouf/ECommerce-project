<template>
  <div class="dashboard-container">
    <AdminNavbar />

    <div class="dashboard-content">
      <h1 class="page-title">Tableau de Bord</h1>

      <!-- Loading state -->
      <div v-if="adminStore.isLoading" class="loading">
        <div class="loading-spinner"></div>
        <p>Chargement des statistiques...</p>
      </div>

      <!-- Error state -->
      <div v-else-if="adminStore.error" class="error-message">
        <h3>Erreur de chargement</h3>
        <p>{{ adminStore.error }}</p>
        <button @click="loadDashboardData" class="retry-btn">Réessayer</button>
      </div>

      <!-- Dashboard content -->
      <div v-else-if="adminStore.dashboardStats" class="dashboard-layout">
        <!-- Key metrics -->
        <div class="stats-grid">
          <div class="stat-card neo-border neo-rotate-neg">
            <div class="stat-card-inner">
              <h3>Total des commandes</h3>
              <div class="stat-value">
                {{ adminStore.dashboardStats.totalOrders }}
              </div>
              <div class="stat-icon">📦</div>
            </div>
          </div>

          <div class="stat-card neo-border neo-rotate-pos">
            <div class="stat-card-inner">
              <h3>Total des produits</h3>
              <div class="stat-value">
                {{ adminStore.dashboardStats.totalProducts }}
              </div>
              <div class="stat-icon">🏷️</div>
            </div>
          </div>

          <div class="stat-card neo-border neo-rotate-neg">
            <div class="stat-card-inner">
              <h3>Chiffre d'affaires</h3>
              <div class="stat-value">
                {{ formatPrice(adminStore.dashboardStats.totalRevenue) }}
              </div>
              <div class="stat-icon">💰</div>
            </div>
          </div>
        </div>

        <!-- Recent orders table -->
        <div class="section-container neo-border">
          <h2 class="section-title">Commandes récentes</h2>

          <table
            class="data-table orders-table"
            v-if="
              adminStore.dashboardStats.recentOrders &&
              adminStore.dashboardStats.recentOrders.length
            "
          >
            <thead>
              <tr>
                <th>ID</th>
                <th>Client</th>
                <th>Date</th>
                <th>Montant</th>
                <th>Statut</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="order in adminStore.dashboardStats.recentOrders"
                :key="order.id"
              >
                <td># {{ order.id }}</td>
                <td>{{ order.userId }}</td>
                <td>{{ formatDate(order.createdAt) }}</td>
                <td>{{ formatPrice(order.totalPrice) }}</td>
                <td>
                  <span
                    class="status-badge"
                    :class="getStatusClass(order.status)"
                  >
                    {{ order.status }}
                  </span>
                </td>
                <td>
                  <router-link
                    :to="`/admin/orders?id=${order.id}`"
                    class="action-btn"
                  >
                    Détails
                  </router-link>
                </td>
              </tr>
            </tbody>
          </table>

          <div v-else class="empty-state">
            Aucune commande récente à afficher
          </div>

          <div class="section-footer">
            <router-link to="/admin/orders" class="view-all-btn">
              Voir toutes les commandes →
            </router-link>
          </div>
        </div>

        <!-- Low stock products table -->
        <div class="section-container neo-border">
          <h2 class="section-title">Produits à faible stock</h2>

          <table
            class="data-table products-table"
            v-if="
              adminStore.dashboardStats.lowStockProducts &&
              adminStore.dashboardStats.lowStockProducts.length
            "
          >
            <thead>
              <tr>
                <th>Image</th>
                <th>Produit</th>
                <th>Stock</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="product in adminStore.dashboardStats.lowStockProducts"
                :key="product.id"
              >
                <td class="product-image">
                  <img :src="product.image" :alt="product.name" />
                </td>
                <td>
                  <div class="product-name">{{ product.name }}</div>
                  <div class="product-category">{{ product.category }}</div>
                </td>
                <td>
                  <span
                    class="stock-count"
                    :class="{ critical: product.stock < 3 }"
                  >
                    {{ product.stock }}
                  </span>
                </td>
                <td>
                  <router-link
                    :to="`/admin/products?id=${product.id}`"
                    class="action-btn"
                  >
                    Modifier
                  </router-link>
                </td>
              </tr>
            </tbody>
          </table>

          <div v-else class="empty-state">
            Tous les produits ont un stock suffisant
          </div>

          <div class="section-footer">
            <router-link to="/admin/products" class="view-all-btn">
              Gérer les produits →
            </router-link>
          </div>
        </div>
      </div>

      <!-- Empty state -->
      <div v-else class="empty-dashboard">
        <div class="empty-icon">📊</div>
        <h2>Aucune donnée disponible</h2>
        <p>Impossible de charger les statistiques du tableau de bord</p>
        <button @click="loadDashboardData" class="retry-btn">Réessayer</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AdminNavbar from "@/components/admin/AdminNavbar.vue";
import { useAdminStore } from "@/stores/adminStore";
import { onMounted } from "vue";

const adminStore = useAdminStore();

// Fonction pour charger les données du dashboard
const loadDashboardData = async () => {
  await adminStore.fetchDashboardStats();
};

// Formatter les prix
const formatPrice = (price: number): string => {
  return new Intl.NumberFormat("fr-FR", {
    style: "currency",
    currency: "EUR",
  }).format(price);
};

// Formatter les dates
const formatDate = (dateString: string): string => {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat("fr-FR", {
    day: "numeric",
    month: "short",
    year: "numeric",
  }).format(date);
};

// Obtenir la classe CSS pour le statut de commande
const getStatusClass = (status: string): string => {
  switch (status.toLowerCase()) {
    case "completed":
    case "terminée":
      return "status-completed";
    case "processing":
    case "en cours":
      return "status-processing";
    case "pending":
    case "en attente":
      return "status-pending";
    case "cancelled":
    case "annulée":
      return "status-cancelled";
    default:
      return "";
  }
};

// Charger les données au chargement de la page
onMounted(loadDashboardData);
</script>

<style scoped>
.dashboard-container {
  min-height: 100vh;
  background-color: #f9f9f9;
}

.dashboard-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.page-title {
  font-family: var(--font-heading);
  font-size: 3rem;
  font-weight: 900;
  margin-bottom: 2rem;
  text-transform: uppercase;
  letter-spacing: -1px;
  position: relative;
  display: inline-block;
}

.page-title::after {
  content: "";
  position: absolute;
  bottom: -10px;
  left: 0;
  width: 100px;
  height: 6px;
  background-color: #000;
}

.dashboard-layout {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-bottom: 1rem;
}

.stat-card {
  background-color: white;
  border: 4px solid #000;
  box-shadow: 8px 8px 0 rgba(0, 0, 0, 0.8);
  padding: 1.5rem;
  position: relative;
  overflow: hidden;
}

.neo-rotate-neg {
  transform: rotate(-1deg);
}

.neo-rotate-pos {
  transform: rotate(1deg);
}

.stat-card-inner {
  transform: rotate(0deg);
}

.stat-card h3 {
  font-family: var(--font-heading);
  font-size: 1.2rem;
  text-transform: uppercase;
  margin-bottom: 1rem;
  font-weight: 700;
}

.stat-value {
  font-family: var(--font-heading);
  font-size: 2.5rem;
  font-weight: 900;
  line-height: 1;
}

.stat-icon {
  position: absolute;
  top: 1rem;
  right: 1rem;
  font-size: 2rem;
  opacity: 0.2;
}

.section-container {
  background-color: white;
  border: 4px solid #000;
  box-shadow: 10px 10px 0 rgba(0, 0, 0, 0.8);
  margin-bottom: 2rem;
  transform: rotate(0.5deg);
  overflow: hidden;
}

.section-title {
  background-color: #000;
  color: white;
  padding: 1rem 1.5rem;
  font-family: var(--font-heading);
  font-size: 1.5rem;
  text-transform: uppercase;
  margin: 0;
  font-weight: 700;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th {
  border-bottom: 2px solid #000;
  padding: 1rem;
  text-align: left;
  font-family: var(--font-heading);
  text-transform: uppercase;
  font-weight: 700;
  font-size: 0.9rem;
}

.data-table td {
  padding: 1rem;
  border-bottom: 1px solid #eee;
}

.data-table tr:last-child td {
  border-bottom: none;
}

.status-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 3px;
  font-weight: 700;
  font-size: 0.8rem;
  text-transform: uppercase;
}

.status-completed {
  background-color: #d1fae5;
  color: #065f46;
}

.status-processing {
  background-color: #e0f2fe;
  color: #0369a1;
}

.status-pending {
  background-color: #fef3c7;
  color: #92400e;
}

.status-cancelled {
  background-color: #fee2e2;
  color: #b91c1c;
}

.action-btn {
  display: inline-block;
  padding: 0.5rem 1rem;
  background-color: #000;
  color: white;
  text-decoration: none;
  font-weight: 700;
  font-size: 0.9rem;
  transform: skew(-3deg);
  transition: all 0.2s;
}

.action-btn:hover {
  transform: skew(-3deg) translateY(-2px);
  box-shadow: 3px 3px 0 rgba(0, 0, 0, 0.3);
}

.product-image img {
  width: 50px;
  height: 50px;
  object-fit: cover;
  border: 2px solid #000;
}

.product-name {
  font-weight: 700;
  margin-bottom: 0.25rem;
}

.product-category {
  font-size: 0.85rem;
  color: #666;
}

.stock-count {
  font-weight: 700;
}

.critical {
  color: #dc2626;
}

.section-footer {
  padding: 1rem;
  text-align: right;
  border-top: 2px solid #eee;
}

.view-all-btn {
  display: inline-block;
  font-family: var(--font-heading);
  font-weight: 700;
  color: #000;
  text-decoration: none;
}

.view-all-btn:hover {
  text-decoration: underline;
}

.empty-state {
  padding: 2rem;
  text-align: center;
  color: #666;
  font-style: italic;
}

.empty-dashboard {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 50vh;
  text-align: center;
}

.empty-icon {
  font-size: 5rem;
  margin-bottom: 1rem;
  opacity: 0.3;
}

.retry-btn {
  display: inline-block;
  margin-top: 1rem;
  background-color: #000;
  color: white;
  border: 3px solid #000;
  padding: 0.75rem 1.5rem;
  font-family: var(--font-heading);
  font-weight: 700;
  text-transform: uppercase;
  cursor: pointer;
  transform: skew(-3deg);
  transition: all 0.2s;
}

.retry-btn:hover {
  transform: skew(-3deg) translateY(-3px);
  box-shadow: 5px 5px 0 rgba(0, 0, 0, 0.3);
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 50vh;
}

.loading-spinner {
  width: 60px;
  height: 60px;
  border: 6px solid #e5e5e5;
  border-top: 6px solid #000;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.error-message {
  background-color: #fee2e2;
  border-left: 6px solid #ef4444;
  padding: 1.5rem;
  margin: 2rem 0;
}

.error-message h3 {
  margin-top: 0;
  font-weight: 700;
  color: #b91c1c;
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .page-title {
    font-size: 2rem;
  }

  .data-table {
    display: block;
    overflow-x: auto;
  }
}
</style>
