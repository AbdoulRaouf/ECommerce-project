<template>
  <div class="orders-container">
    <AdminNavbar />

    <div class="orders-content">
      <h1 class="page-title">Gestion des Commandes</h1>

      <!-- Order filters -->
      <div class="filters">
        <div class="filter-buttons">
          <button
            @click="loadOrders()"
            class="filter-btn"
            :class="{ active: !currentFilter }"
          >
            Toutes
          </button>
          <button
            @click="loadOrders('pending')"
            class="filter-btn"
            :class="{ active: currentFilter === 'pending' }"
          >
            En attente
          </button>
          <button
            @click="loadOrders('processing')"
            class="filter-btn"
            :class="{ active: currentFilter === 'processing' }"
          >
            En cours
          </button>
          <button
            @click="loadOrders('completed')"
            class="filter-btn"
            :class="{ active: currentFilter === 'completed' }"
          >
            Terminées
          </button>
          <button
            @click="loadOrders('cancelled')"
            class="filter-btn"
            :class="{ active: currentFilter === 'cancelled' }"
          >
            Annulées
          </button>
        </div>
      </div>

      <!-- Loading state -->
      <div v-if="adminStore.isLoading" class="loading">
        <div class="loading-spinner"></div>
        <p>Chargement des commandes...</p>
      </div>

      <!-- Error state -->
      <div v-else-if="adminStore.error" class="error-message">
        <h3>Erreur de chargement</h3>
        <p>{{ adminStore.error }}</p>
        <button @click="loadOrders(currentFilter)" class="retry-btn">
          Réessayer
        </button>
      </div>

      <!-- Orders table -->
      <div
        v-else-if="adminStore.orders.length"
        class="orders-table-container neo-border"
      >
        <table class="data-table orders-table">
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
              v-for="order in adminStore.orders"
              :key="order.id"
              @click="selectOrder(order.id)"
              :class="{ 'selected-row': selectedOrderId === order.id }"
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
                <div class="table-actions">
                  <button
                    @click.stop="selectOrder(order.id)"
                    class="action-btn view-btn"
                  >
                    Voir
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Empty state -->
      <div v-else class="empty-orders neo-border">
        <div class="empty-icon">🧾</div>
        <h2>Aucune commande trouvée</h2>
        <p>Il n'y a pas de commandes correspondant à vos critères</p>
      </div>
    </div>

    <!-- Order detail modal -->
    <div class="modal" :class="{ 'is-active': selectedOrderId !== null }">
      <div class="modal-background" @click="closeOrderDetails"></div>
      <div class="modal-content order-details-modal">
        <div v-if="adminStore.selectedOrder" class="order-details">
          <div class="modal-header">
            <h2 class="modal-title">
              Détails de la commande #{{ adminStore.selectedOrder.id }}
            </h2>
            <button @click="closeOrderDetails" class="close-btn">
              &times;
            </button>
          </div>

          <div class="order-info">
            <div class="info-row">
              <span class="label">Client:</span>
              <span class="value">{{ adminStore.selectedOrder.userId }}</span>
            </div>
            <div class="info-row">
              <span class="label">Date:</span>
              <span class="value">{{
                formatDate(adminStore.selectedOrder.createdAt)
              }}</span>
            </div>
            <div class="info-row">
              <span class="label">Statut:</span>
              <span
                class="value status-badge"
                :class="getStatusClass(adminStore.selectedOrder.status)"
              >
                {{ adminStore.selectedOrder.status }}
              </span>
            </div>
            <div class="info-row">
              <span class="label">Montant total:</span>
              <span class="value total-price">{{
                formatPrice(adminStore.selectedOrder.totalPrice)
              }}</span>
            </div>
          </div>

          <div class="status-update">
            <h3>Mettre à jour le statut</h3>
            <div class="status-buttons">
              <button
                v-for="status in [
                  'pending',
                  'processing',
                  'completed',
                  'cancelled',
                ]"
                :key="status"
                @click="updateOrderStatus(adminStore.selectedOrder.id, status)"
                class="status-update-btn"
                :class="[
                  `status-${status}`,
                  {
                    active:
                      adminStore.selectedOrder.status.toLowerCase() === status,
                  },
                ]"
                :disabled="
                  adminStore.isLoading ||
                  adminStore.selectedOrder.status.toLowerCase() === status
                "
              >
                {{ getStatusLabel(status) }}
              </button>
            </div>
          </div>

          <div class="order-items">
            <h3>Articles commandés</h3>
            <div
              v-if="
                adminStore.selectedOrder.items &&
                adminStore.selectedOrder.items.length
              "
              class="items-list"
            >
              <div
                v-for="item in adminStore.selectedOrder.items"
                :key="item.id"
                class="order-item neo-border"
              >
                <div class="item-image">
                  <img
                    v-if="item.product?.image"
                    :src="item.product.image"
                    :alt="item.product?.name"
                  />
                  <div v-else class="no-image">Pas d'image</div>
                </div>
                <div class="item-details">
                  <h4>{{ item.product?.name || "Produit inconnu" }}</h4>
                  <div class="item-meta">
                    <span class="quantity">{{ item.quantity }} × </span>
                    <span class="price">{{
                      formatPrice(item.product?.price || 0)
                    }}</span>
                  </div>
                </div>
                <div class="item-total">
                  {{ formatPrice((item.product?.price || 0) * item.quantity) }}
                </div>
              </div>
            </div>
            <div v-else class="empty-items">
              Aucun article dans cette commande
            </div>
          </div>

          <div class="modal-actions">
            <button @click="confirmDeleteOrder" class="delete-btn">
              Supprimer la commande
            </button>
            <button @click="closeOrderDetails" class="close-action-btn">
              Fermer
            </button>
          </div>
        </div>

        <div v-else class="order-loading">
          <div class="loading-spinner"></div>
          <p>Chargement des détails de la commande...</p>
        </div>
      </div>
    </div>

    <!-- Confirmation modal -->
    <div class="modal" :class="{ 'is-active': isConfirmDialogOpen }">
      <div class="modal-background" @click="isConfirmDialogOpen = false"></div>
      <div class="modal-content confirm-dialog">
        <h3>Confirmer la suppression</h3>
        <p>
          Êtes-vous sûr de vouloir supprimer cette commande ? Cette action est
          irréversible.
        </p>
        <div class="confirm-actions">
          <button @click="isConfirmDialogOpen = false" class="btn-cancel">
            Annuler
          </button>
          <button @click="deleteOrder" class="btn-confirm-delete">
            Supprimer
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AdminNavbar from "@/components/admin/AdminNavbar.vue";
import { useAdminStore } from "@/stores/adminStore";
import { onMounted, ref } from "vue";

const adminStore = useAdminStore();
const currentFilter = ref<string | null>(null);
const selectedOrderId = ref<number | null>(null);
const isConfirmDialogOpen = ref(false);

// Charger les commandes
const loadOrders = async (status?: string) => {
  currentFilter.value = status || null;
  await adminStore.fetchOrders(status);
};

// Sélectionner une commande pour voir les détails
const selectOrder = async (id: number) => {
  selectedOrderId.value = id;
  await adminStore.fetchOrderDetails(id);
};

// Fermer les détails de la commande
const closeOrderDetails = () => {
  selectedOrderId.value = null;
};

// Mettre à jour le statut d'une commande
const updateOrderStatus = async (id: number, status: string) => {
  const success = await adminStore.updateOrderStatus(id, status);
  if (success) {
    // Recharger les détails de la commande
    await adminStore.fetchOrderDetails(id);
    // Rafraîchir la liste des commandes si on filtre par statut
    if (currentFilter.value) {
      await loadOrders(currentFilter.value);
    }
  }
};

// Confirmer la suppression
const confirmDeleteOrder = () => {
  isConfirmDialogOpen.value = true;
};

// Supprimer une commande
const deleteOrder = async () => {
  if (selectedOrderId.value) {
    const success = await adminStore.deleteOrder(selectedOrderId.value);
    if (success) {
      closeOrderDetails();
      await loadOrders(currentFilter.value);
      isConfirmDialogOpen.value = false;
    }
  }
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
    hour: "2-digit",
    minute: "2-digit",
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

// Obtenir le label du statut
const getStatusLabel = (status: string): string => {
  switch (status) {
    case "completed":
      return "Terminée";
    case "processing":
      return "En cours";
    case "pending":
      return "En attente";
    case "cancelled":
      return "Annulée";
    default:
      return status;
  }
};

// Charger les commandes au chargement de la page
onMounted(() => {
  loadOrders();
});
</script>

<style scoped>
.orders-container {
  min-height: 100vh;
  background-color: #f9f9f9;
}

.orders-content {
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

.filters {
  margin-bottom: 2rem;
}

.filter-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.filter-btn {
  padding: 0.75rem 1.25rem;
  background-color: #e5e5e5;
  border: 3px solid #000;
  font-family: var(--font-heading);
  font-weight: 700;
  text-transform: uppercase;
  cursor: pointer;
  transform: skew(-3deg);
  transition: all 0.2s;
}

.filter-btn:hover {
  background-color: #d5d5d5;
  transform: skew(-3deg) translateY(-2px);
  box-shadow: 3px 3px 0 rgba(0, 0, 0, 0.2);
}

.filter-btn.active {
  background-color: #000;
  color: white;
}

.orders-table-container {
  background-color: white;
  border: 4px solid #000;
  box-shadow: 10px 10px 0 rgba(0, 0, 0, 0.8);
  margin-bottom: 2rem;
  overflow: hidden;
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

.data-table tr {
  cursor: pointer;
  transition: background-color 0.2s;
}

.data-table tr:hover {
  background-color: #f9f9f9;
}

.selected-row {
  background-color: #f0f0f0;
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

.table-actions {
  display: flex;
  gap: 0.5rem;
}

.action-btn {
  padding: 0.4rem 0.8rem;
  background-color: #000;
  color: white;
  border: none;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  transform: skew(-3deg);
  font-size: 0.9rem;
}

.view-btn:hover {
  transform: skew(-3deg) translateY(-2px);
  box-shadow: 2px 2px 0 rgba(0, 0, 0, 0.3);
}

.empty-orders {
  background-color: white;
  border: 4px solid #000;
  box-shadow: 10px 10px 0 rgba(0, 0, 0, 0.8);
  padding: 3rem;
  text-align: center;
  margin-bottom: 2rem;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
  opacity: 0.3;
}

.empty-orders h2 {
  font-family: var(--font-heading);
  font-size: 1.8rem;
  margin-bottom: 1rem;
}

.empty-orders p {
  color: #666;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
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

/* Modal styles */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 100;
  visibility: hidden;
  opacity: 0;
  transition: all 0.3s ease;
}

.modal.is-active {
  visibility: visible;
  opacity: 1;
}

.modal-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.7);
  z-index: 101;
}

.modal-content {
  position: relative;
  z-index: 102;
  width: 90%;
  max-width: 800px;
  max-height: 90vh;
  overflow-y: auto;
  background-color: #fff;
  border: 4px solid #000;
  box-shadow: 15px 15px 0 rgba(0, 0, 0, 0.8);
}

.order-details-modal {
  transform: rotate(-0.5deg);
}

.order-details {
  transform: rotate(0.5deg);
  padding: 0;
}

.modal-header {
  background-color: #000;
  color: white;
  padding: 1rem 1.5rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-title {
  font-family: var(--font-heading);
  font-size: 1.5rem;
  text-transform: uppercase;
  margin: 0;
  font-weight: 700;
}

.close-btn {
  background: none;
  border: none;
  color: white;
  font-size: 2rem;
  cursor: pointer;
  line-height: 1;
  padding: 0;
}

.order-info {
  padding: 1.5rem;
  border-bottom: 2px solid #eee;
}

.info-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.75rem;
}

.info-row:last-child {
  margin-bottom: 0;
}

.label {
  font-weight: 700;
}

.total-price {
  font-size: 1.2rem;
  font-weight: 900;
}

.status-update {
  padding: 1.5rem;
  border-bottom: 2px solid #eee;
}

.status-update h3 {
  font-family: var(--font-heading);
  font-weight: 700;
  margin-top: 0;
  margin-bottom: 1rem;
  font-size: 1.2rem;
  text-transform: uppercase;
}

.status-buttons {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 0.75rem;
}

.status-update-btn {
  padding: 0.5rem;
  border: 2px solid;
  background-color: transparent;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  text-transform: uppercase;
  font-size: 0.8rem;
}

.status-update-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 3px 3px 0 rgba(0, 0, 0, 0.1);
}

.status-update-btn.active {
  border-width: 2px;
}

.status-update-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.status-pending {
  color: #92400e;
  border-color: #fef3c7;
  background-color: #fffbeb;
}

.status-pending.active {
  background-color: #fef3c7;
  border-color: #92400e;
}

.status-processing {
  color: #0369a1;
  border-color: #e0f2fe;
  background-color: #f0f9ff;
}

.status-processing.active {
  background-color: #e0f2fe;
  border-color: #0369a1;
}

.status-completed {
  color: #065f46;
  border-color: #d1fae5;
  background-color: #ecfdf5;
}

.status-completed.active {
  background-color: #d1fae5;
  border-color: #065f46;
}

.status-cancelled {
  color: #b91c1c;
  border-color: #fee2e2;
  background-color: #fef2f2;
}

.status-cancelled.active {
  background-color: #fee2e2;
  border-color: #b91c1c;
}

.order-items {
  padding: 1.5rem;
}

.order-items h3 {
  font-family: var(--font-heading);
  font-weight: 700;
  margin-top: 0;
  margin-bottom: 1.5rem;
  font-size: 1.2rem;
  text-transform: uppercase;
}

.items-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.order-item {
  display: grid;
  grid-template-columns: 60px 1fr auto;
  gap: 1rem;
  padding: 1rem;
  align-items: center;
  border: 2px solid #eee;
}

.item-image {
  width: 60px;
  height: 60px;
  overflow: hidden;
  border: 2px solid #000;
}

.item-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.no-image {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f0f0f0;
  font-size: 0.7rem;
  text-align: center;
  padding: 0.25rem;
}

.item-details h4 {
  margin: 0 0 0.5rem 0;
  font-weight: 700;
  font-size: 1rem;
}

.item-meta {
  font-size: 0.9rem;
}

.quantity {
  color: #666;
}

.item-total {
  font-weight: 700;
}

.empty-items {
  text-align: center;
  padding: 2rem 0;
  color: #666;
  font-style: italic;
}

.modal-actions {
  padding: 1.5rem;
  display: flex;
  justify-content: space-between;
  border-top: 2px solid #eee;
}

.delete-btn {
  padding: 0.75rem 1.25rem;
  background-color: #fee2e2;
  border: 2px solid #b91c1c;
  color: #b91c1c;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  text-transform: uppercase;
  transform: skew(-3deg);
}

.delete-btn:hover {
  background-color: #fecaca;
  transform: skew(-3deg) translateY(-2px);
  box-shadow: 3px 3px 0 rgba(185, 28, 28, 0.2);
}

.close-action-btn {
  padding: 0.75rem 1.25rem;
  background-color: #000;
  border: 2px solid #000;
  color: white;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  text-transform: uppercase;
  transform: skew(-3deg);
}

.close-action-btn:hover {
  background-color: #222;
  transform: skew(-3deg) translateY(-2px);
  box-shadow: 3px 3px 0 rgba(0, 0, 0, 0.3);
}

.order-loading {
  padding: 3rem;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.confirm-dialog {
  background-color: white;
  padding: 1.5rem;
  max-width: 400px;
  transform: rotate(-1deg);
}

.confirm-dialog h3 {
  font-family: var(--font-heading);
  font-weight: 700;
  margin-top: 0;
  font-size: 1.5rem;
  margin-bottom: 1rem;
}

.confirm-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 1.5rem;
}

.btn-cancel {
  padding: 0.75rem 1.25rem;
  background-color: #e5e5e5;
  border: 3px solid #000;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  transform: skew(-3deg);
}

.btn-cancel:hover {
  background-color: #d5d5d5;
  transform: skew(-3deg) translateY(-2px);
  box-shadow: 3px 3px 0 rgba(0, 0, 0, 0.2);
}

.btn-confirm-delete {
  padding: 0.75rem 1.25rem;
  background-color: #b91c1c;
  border: 3px solid #b91c1c;
  color: white;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  transform: skew(-3deg);
}

.btn-confirm-delete:hover {
  background-color: #991b1b;
  transform: skew(-3deg) translateY(-2px);
  box-shadow: 3px 3px 0 rgba(185, 28, 28, 0.4);
}

@media (max-width: 768px) {
  .page-title {
    font-size: 2rem;
  }

  .data-table {
    display: block;
    overflow-x: auto;
  }

  .filter-buttons {
    overflow-x: auto;
    padding-bottom: 0.5rem;
    justify-content: flex-start;
    flex-wrap: nowrap;
  }
}
</style>
