<template>
  <div class="orders-container">
    <div class="container">
      <h1 class="page-title">Mes commandes</h1>

      <div v-if="loading" class="loading-indicator">
        Chargement de vos commandes...
      </div>

      <div v-else-if="error" class="error-message">
        {{ error }}
      </div>

      <div v-else-if="orders.length === 0" class="no-orders neo-border">
        <p>Vous n'avez pas encore passé de commande.</p>
        <router-link to="/" class="shop-button">
          Commencer mes achats
        </router-link>
      </div>

      <template v-else>
        <!-- Filtres (statut, date) -->
        <div class="filters">
          <div class="filter-group">
            <label for="status-filter">Filtrer par statut:</label>
            <select
              id="status-filter"
              v-model="statusFilter"
              class="filter-select"
            >
              <option value="">Tous les statuts</option>
              <option value="pending">En attente</option>
              <option value="confirmed">Confirmée</option>
              <option value="shipped">Expédiée</option>
              <option value="delivered">Livrée</option>
              <option value="cancelled">Annulée</option>
            </select>
          </div>

          <div class="filter-group">
            <label for="date-filter">Trier par:</label>
            <select id="date-filter" v-model="dateSort" class="filter-select">
              <option value="desc">Plus récentes d'abord</option>
              <option value="asc">Plus anciennes d'abord</option>
            </select>
          </div>
        </div>

        <!-- Liste des commandes -->
        <div class="orders-list">
          <div
            v-for="order in filteredOrders"
            :key="order.id"
            class="order-card neo-border"
            :class="{ 'order-cancelled': order.status === 'cancelled' }"
          >
            <div class="order-header">
              <div class="order-meta">
                <div class="order-number">Commande #{{ order.id }}</div>
                <div class="order-date">{{ formatDate(order.createdAt) }}</div>
              </div>
              <div class="order-status" :class="getStatusClass(order.status)">
                {{ formatStatus(order.status) }}
              </div>
            </div>

            <div class="order-details">
              <div class="order-summary">
                <div class="summary-item">
                  <span class="summary-label">Articles:</span>
                  <span class="summary-value">{{
                    order.items?.length || 0
                  }}</span>
                </div>
                <div class="summary-item">
                  <span class="summary-label">Total:</span>
                  <span class="summary-value total-price">{{
                    formatPrice(order.totalPrice)
                  }}</span>
                </div>
              </div>

              <!-- Liste des produits (affichage conditionnel) -->
              <div v-if="expandedOrder === order.id" class="order-items">
                <div class="items-header">
                  <span class="item-header-cell">Produit</span>
                  <span class="item-header-cell">Prix</span>
                  <span class="item-header-cell">Quantité</span>
                  <span class="item-header-cell">Sous-total</span>
                </div>

                <div
                  v-for="item in order.items"
                  :key="item.id"
                  class="order-item"
                >
                  <div class="item-cell item-product">
                    <div class="item-image" v-if="item.product?.image">
                      <img :src="item.product.image" :alt="item.product.name" />
                    </div>
                    <div class="item-name">
                      {{ item.product?.name || "Produit non disponible" }}
                    </div>
                  </div>
                  <div class="item-cell item-price">
                    {{ formatPrice(item.product?.price || 0) }}
                  </div>
                  <div class="item-cell item-quantity">
                    {{ item.quantity }}
                  </div>
                  <div class="item-cell item-subtotal">
                    {{
                      formatPrice((item.product?.price || 0) * item.quantity)
                    }}
                  </div>
                </div>

                <!-- Totaux de commande -->
                <div class="order-totals">
                  <div class="totals-row subtotal">
                    <span>Sous-total</span>
                    <span>{{ formatPrice(order.subtotalPrice || 0) }}</span>
                  </div>
                  <div class="totals-row shipping">
                    <span>Livraison</span>
                    <span>{{ formatPrice(order.shippingCost || 0) }}</span>
                  </div>
                  <div class="totals-row taxes">
                    <span>Taxes</span>
                    <span>{{ formatPrice(order.taxAmount || 0) }}</span>
                  </div>
                  <div class="totals-row total">
                    <span>Total</span>
                    <span>{{ formatPrice(order.totalPrice) }}</span>
                  </div>
                </div>
              </div>

              <div class="order-actions">
                <button
                  @click="toggleOrderDetails(order.id)"
                  class="action-button detail-button"
                >
                  {{
                    expandedOrder === order.id
                      ? "Masquer les détails"
                      : "Voir les détails"
                  }}
                </button>
                <button
                  v-if="order.status === 'pending'"
                  class="action-button cancel-button"
                  disabled
                >
                  Annuler la commande
                </button>
                <button class="action-button reorder-button" disabled>
                  Commander à nouveau
                </button>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useAuthStore } from "../../stores/authStore";

const authStore = useAuthStore();

const orders = ref([]);
const loading = ref(true);
const error = ref(null);
const expandedOrder = ref(null);
const statusFilter = ref("");
const dateSort = ref("desc");

// Récupération des commandes au chargement du composant
onMounted(async () => {
  await fetchOrders();
});

async function fetchOrders() {
  loading.value = true;
  error.value = null;

  try {
    orders.value = await authStore.fetchUserOrders();
  } catch (err) {
    console.error("Erreur lors du chargement des commandes:", err);
    error.value =
      "Impossible de charger vos commandes. Veuillez réessayer ultérieurement.";
  } finally {
    loading.value = false;
  }
}

// Filtrage des commandes selon les critères sélectionnés
const filteredOrders = computed(() => {
  let result = [...orders.value];

  // Filtrer par statut si un statut est sélectionné
  if (statusFilter.value) {
    result = result.filter((order) => order.status === statusFilter.value);
  }

  // Trier par date
  result.sort((a, b) => {
    const dateA = new Date(a.createdAt).getTime();
    const dateB = new Date(b.createdAt).getTime();
    return dateSort.value === "desc" ? dateB - dateA : dateA - dateB;
  });

  return result;
});

function toggleOrderDetails(orderId) {
  if (expandedOrder.value === orderId) {
    expandedOrder.value = null;
  } else {
    expandedOrder.value = orderId;
  }
}

function formatDate(dateString) {
  if (!dateString) return "N/A";
  const date = new Date(dateString);
  return new Intl.DateTimeFormat("fr-FR", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
}

function formatPrice(price) {
  return new Intl.NumberFormat("fr-FR", {
    style: "currency",
    currency: "EUR",
  }).format(price);
}

function formatStatus(status) {
  const statuses = {
    pending: "En attente",
    confirmed: "Confirmée",
    shipped: "Expédiée",
    delivered: "Livrée",
    cancelled: "Annulée",
  };
  return statuses[status] || status;
}

function getStatusClass(status) {
  return `status-${status}`;
}
</script>

<style scoped>
.orders-container {
  padding: 3rem 1rem;
}

.page-title {
  font-family: var(--font-heading);
  font-size: 2.5rem;
  font-weight: 900;
  margin-bottom: 2rem;
  text-transform: uppercase;
  position: relative;
  display: inline-block;
}

.page-title::after {
  content: "";
  position: absolute;
  bottom: -10px;
  left: 0;
  width: 80px;
  height: 4px;
  background-color: #000;
}

.neo-border {
  border: 4px solid #000;
  box-shadow: 6px 6px 0 rgba(0, 0, 0, 0.8);
  background-color: white;
}

.filters {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  margin-bottom: 2rem;
}

.filter-group {
  flex-grow: 1;
  max-width: 300px;
}

.filter-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 700;
  font-size: 0.9rem;
}

.filter-select {
  width: 100%;
  padding: 0.5rem;
  border: 2px solid #000;
  background-color: #f9f9f9;
  font-family: var(--font-body);
  font-size: 0.9rem;
}

.no-orders {
  padding: 3rem;
  text-align: center;
  margin-bottom: 2rem;
}

.no-orders p {
  margin-bottom: 1.5rem;
  font-size: 1.2rem;
}

.shop-button {
  display: inline-block;
  background-color: var(--color-black);
  color: var(--color-white);
  padding: 0.75rem 1.5rem;
  text-decoration: none;
  font-weight: 700;
  text-transform: uppercase;
  transform: skew(-5deg);
  transition: transform 0.2s;
}

.shop-button:hover {
  transform: skew(-5deg) translateY(-3px);
}

.orders-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.order-card {
  padding: 1.5rem;
  transform: rotate(-0.5deg);
  transition: transform 0.2s;
}

.order-card:hover {
  transform: rotate(-0.5deg) translateY(-3px);
}

.order-cancelled {
  border-color: #b91c1c;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  border-bottom: 1px solid #eee;
  padding-bottom: 1rem;
}

.order-number {
  font-weight: 900;
  font-family: var(--font-heading);
  font-size: 1.2rem;
}

.order-date {
  font-size: 0.9rem;
  color: #555;
  margin-top: 0.25rem;
}

.order-status {
  font-size: 0.9rem;
  padding: 0.25rem 0.75rem;
  border-radius: 4px;
  font-weight: 700;
  text-transform: uppercase;
}

.status-pending {
  background-color: #fef3c7;
  color: #92400e;
}

.status-confirmed {
  background-color: #e0f2fe;
  color: #075985;
}

.status-shipped {
  background-color: #d1fae5;
  color: #065f46;
}

.status-delivered {
  background-color: #dcfce7;
  color: #15803d;
}

.status-cancelled {
  background-color: #fee2e2;
  color: #b91c1c;
}

.order-details {
  transform: rotate(0.5deg);
}

.order-summary {
  display: flex;
  flex-wrap: wrap;
  gap: 2rem;
  margin-bottom: 1.5rem;
}

.summary-item {
  display: flex;
  flex-direction: column;
}

.summary-label {
  font-size: 0.9rem;
  color: #555;
}

.summary-value {
  font-weight: 700;
  font-size: 1.1rem;
}

.total-price {
  font-size: 1.2rem;
  color: #000;
}

.order-items {
  border: 2px solid #eee;
  margin: 1rem 0;
  padding: 1rem;
}

.items-header {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr;
  gap: 1rem;
  padding-bottom: 0.5rem;
  margin-bottom: 0.5rem;
  border-bottom: 1px solid #eee;
  font-weight: 700;
  font-size: 0.9rem;
  text-transform: uppercase;
}

.order-item {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr;
  gap: 1rem;
  padding: 1rem 0;
  border-bottom: 1px solid #f5f5f5;
}

.item-product {
  display: flex;
  align-items: center;
}

.item-image {
  width: 50px;
  height: 50px;
  margin-right: 1rem;
  background-color: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #eee;
}

.item-image img {
  max-width: 100%;
  max-height: 100%;
  object-fit: cover;
}

.item-name {
  font-weight: 500;
}

.item-cell {
  display: flex;
  align-items: center;
}

.order-totals {
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid #eee;
}

.totals-row {
  display: flex;
  justify-content: space-between;
  padding: 0.5rem 0;
}

.total {
  font-weight: 700;
  font-size: 1.1rem;
  border-top: 1px solid #ddd;
  padding-top: 0.5rem;
  margin-top: 0.5rem;
}

.order-actions {
  margin-top: 1.5rem;
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.action-button {
  padding: 0.5rem 1rem;
  background: none;
  border: 2px solid #000;
  font-family: var(--font-body);
  font-weight: 700;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.2s;
}

.detail-button {
  background-color: #f9f9f9;
}

.detail-button:hover {
  background-color: #eee;
}

.cancel-button {
  color: #b91c1c;
  border-color: #b91c1c;
}

.cancel-button:hover:not(:disabled) {
  background-color: #fee2e2;
}

.reorder-button {
  background-color: #000;
  color: #fff;
}

.reorder-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 3px 3px 0 rgba(0, 0, 0, 0.2);
}

.action-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  border-color: #ccc;
}

.loading-indicator {
  text-align: center;
  padding: 3rem 0;
  font-style: italic;
  color: #555;
}

.error-message {
  background-color: #fee2e2;
  border-left: 4px solid #b91c1c;
  color: #b91c1c;
  padding: 1rem;
  margin-bottom: 1.5rem;
  font-weight: 500;
}

@media (max-width: 768px) {
  .order-item {
    grid-template-columns: 1fr;
    gap: 0.5rem;
  }

  .items-header {
    display: none;
  }

  .item-cell {
    padding: 0.25rem 0;
  }

  .item-cell::before {
    content: attr(data-label);
    font-weight: 700;
    margin-right: 0.5rem;
  }
}
</style>
