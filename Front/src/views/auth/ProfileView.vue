<template>
  <div class="profile-container">
    <div class="container">
      <h1 class="page-title">Mon profil</h1>

      <div v-if="authStore.isLoading" class="loading-indicator">
        Chargement de vos informations...
      </div>

      <div v-else-if="authStore.error" class="error-message">
        {{ authStore.error }}
      </div>

      <template v-else-if="authStore.user">
        <div class="profile-section neo-border">
          <h2 class="section-title">Informations personnelles</h2>

          <div class="profile-info">
            <div class="info-row">
              <span class="info-label">Nom</span>
              <span class="info-value">{{ authStore.user.name }}</span>
            </div>

            <div class="info-row">
              <span class="info-label">Email</span>
              <span class="info-value">{{ authStore.user.email }}</span>
            </div>

            <div class="info-row">
              <span class="info-label">Adresse</span>
              <span class="info-value">{{
                authStore.user.address || "Non renseignée"
              }}</span>
            </div>

            <div class="info-row">
              <span class="info-label">Téléphone</span>
              <span class="info-value">{{
                authStore.user.phone || "Non renseigné"
              }}</span>
            </div>

            <div class="info-row">
              <span class="info-label">Membre depuis</span>
              <span class="info-value">{{
                formatDate(authStore.user.createdAt)
              }}</span>
            </div>
          </div>

          <!-- Bouton d'édition (future fonctionnalité) -->
          <div class="actions">
            <button class="neo-button" disabled>
              Modifier mes informations
            </button>
          </div>
        </div>

        <!-- Résumé des commandes récentes -->
        <div class="profile-section neo-border">
          <h2 class="section-title">Commandes récentes</h2>

          <div v-if="loading" class="loading-text">
            Chargement des commandes...
          </div>
          <div v-else-if="orders.length === 0" class="no-orders">
            Vous n'avez pas encore passé de commande.
          </div>
          <div v-else class="orders-summary">
            <!-- Liste des 3 dernières commandes -->
            <div
              v-for="order in orders.slice(0, 3)"
              :key="order.id"
              class="order-item"
            >
              <div class="order-header">
                <div>
                  <div class="order-number">Commande #{{ order.id }}</div>
                  <div class="order-date">
                    {{ formatDate(order.createdAt) }}
                  </div>
                </div>
                <div class="order-status" :class="getStatusClass(order.status)">
                  {{ formatStatus(order.status) }}
                </div>
              </div>
              <div class="order-total">
                Total: {{ formatPrice(order.totalPrice) }}
              </div>
            </div>

            <!-- Lien vers l'historique complet -->
            <div class="view-all">
              <router-link to="/orders" class="view-all-link">
                Voir toutes mes commandes
              </router-link>
            </div>
          </div>
        </div>

        <!-- Actions du compte -->
        <div class="profile-section neo-border">
          <h2 class="section-title">Actions du compte</h2>

          <div class="actions">
            <button @click="handleLogout" class="neo-button danger-button">
              Déconnexion
            </button>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "../../stores/authStore";
import { useProductStore } from "../../stores/productStore";

const authStore = useAuthStore();
const productStore = useProductStore();
const router = useRouter();

const orders = ref([]);
const loading = ref(false);

onMounted(async () => {
  // Récupérer le profil utilisateur s'il n'est pas déjà chargé
  if (!authStore.user) {
    await authStore.fetchUserProfile();
  }

  // Charger les commandes de l'utilisateur
  await fetchOrders();
});

async function fetchOrders() {
  loading.value = true;
  try {
    orders.value = await authStore.fetchUserOrders();
  } catch (error) {
    console.error("Erreur lors du chargement des commandes:", error);
  } finally {
    loading.value = false;
  }
}

function handleLogout() {
  authStore.logout();
  productStore.clearCart();
  router.push("/");
}

function formatDate(dateString) {
  if (!dateString) return "N/A";
  const date = new Date(dateString);
  return new Intl.DateTimeFormat("fr-FR", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
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
.profile-container {
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

.profile-section {
  background-color: white;
  border: 4px solid #000;
  box-shadow: 6px 6px 0 rgba(0, 0, 0, 0.8);
  padding: 2rem;
  margin-bottom: 2rem;
  transform: rotate(-0.5deg);
}

.section-title {
  font-family: var(--font-heading);
  font-size: 1.5rem;
  margin-bottom: 1.5rem;
  text-transform: uppercase;
}

.profile-info {
  transform: rotate(0.5deg);
}

.info-row {
  display: flex;
  padding: 0.75rem 0;
  border-bottom: 1px solid #eee;
}

.info-row:last-child {
  border-bottom: none;
}

.info-label {
  font-weight: 700;
  width: 180px;
  color: #555;
}

.info-value {
  flex-grow: 1;
}

.actions {
  margin-top: 1.5rem;
  display: flex;
  justify-content: flex-end;
}

.neo-button {
  background-color: var(--color-black);
  color: var(--color-white);
  border: none;
  padding: 0.75rem 1.5rem;
  font-family: var(--font-heading);
  font-size: 0.9rem;
  font-weight: 700;
  text-transform: uppercase;
  cursor: pointer;
  transform: skew(-5deg);
  transition: all 0.2s;
}

.neo-button:hover:not(:disabled) {
  transform: skew(-5deg) translateY(-3px);
  box-shadow: 4px 4px 0 rgba(0, 0, 0, 0.2);
}

.neo-button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.danger-button {
  background-color: #b91c1c;
}

.orders-summary {
  transform: rotate(0.5deg);
}

.no-orders {
  padding: 1.5rem 0;
  font-style: italic;
  color: #777;
}

.order-item {
  border: 2px solid #eee;
  padding: 1rem;
  margin-bottom: 1rem;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 0.5rem;
}

.order-number {
  font-weight: 700;
  font-size: 1.1rem;
}

.order-date {
  color: #777;
  font-size: 0.9rem;
}

.order-total {
  font-weight: 700;
  margin-top: 0.5rem;
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

.view-all {
  text-align: center;
  margin-top: 1.5rem;
}

.view-all-link {
  color: var(--color-black);
  font-weight: 700;
  text-decoration: underline;
  font-size: 0.9rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.loading-indicator,
.loading-text {
  text-align: center;
  padding: 2rem 0;
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
</style>
