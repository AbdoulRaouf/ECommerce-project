<template>
  <div class="wishlist-page container">
    <h1 class="page-title">Ma Liste de Souhaits</h1>

    <!-- État de chargement -->
    <div v-if="wishlistStore.isLoading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>Chargement de votre liste de souhaits...</p>
    </div>

    <!-- Gestion des erreurs -->
    <div v-else-if="wishlistStore.error" class="error-message">
      <p>{{ wishlistStore.error }}</p>
      <button @click="wishlistStore.fetchWishlist" class="retry-btn">
        Réessayer
      </button>
    </div>

    <div v-else-if="wishlistStore.wishlist.length === 0" class="empty-wishlist">
      <p>Votre liste de souhaits est vide</p>
      <router-link to="/" class="continue-shopping-btn"
        >DÉCOUVRIR DES PRODUITS</router-link
      >
    </div>

    <div v-else class="wishlist-contents">
      <div class="wishlist-items">
        <div
          v-for="item in wishlistStore.wishlist"
          :key="item.id"
          class="wishlist-item"
        >
          <div class="item-image">
            <img :src="item.product?.image" :alt="item.product?.name" />
          </div>

          <div class="item-details">
            <router-link :to="`/product/${item.productId}`" class="item-name">
              {{ item.product?.name }}
            </router-link>
            <p class="item-price">
              {{ formatPrice(item.product?.price || 0) }}
            </p>
            <p class="item-date">
              Ajouté le {{ formatDate(item.createdAt || "") }}
            </p>
          </div>

          <div class="item-actions">
            <button
              @click="addToCart(item.product!)"
              class="add-to-cart-btn"
              :disabled="productStore.isLoading"
            >
              Ajouter au panier
            </button>

            <button
              @click="removeFromWishlist(item.productId)"
              class="remove-btn"
              :disabled="wishlistStore.isLoading"
            >
              Supprimer
            </button>
          </div>
        </div>
      </div>

      <div class="wishlist-summary">
        <button
          @click="wishlistStore.clearWishlist()"
          class="clear-wishlist-btn"
          :disabled="wishlistStore.isLoading"
        >
          VIDER LA LISTE
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useNotificationStore } from "@/stores/notificationStore";
import { Product, useProductStore } from "@/stores/productStore";
import { useWishlistStore } from "@/stores/wishlistStore";
import { onMounted } from "vue";

const wishlistStore = useWishlistStore();
const productStore = useProductStore();
const notificationStore = useNotificationStore();

onMounted(() => {
  wishlistStore.fetchWishlist();
});

function formatPrice(price: number): string {
  return new Intl.NumberFormat("fr-FR", {
    style: "currency",
    currency: "EUR",
  }).format(price);
}

function formatDate(dateString: string): string {
  if (!dateString) return "";

  const date = new Date(dateString);
  return new Intl.DateTimeFormat("fr-FR", {
    day: "numeric",
    month: "long",
    year: "numeric",
  }).format(date);
}

async function addToCart(product: Product) {
  try {
    await productStore.addToCart(product, 1);
    notificationStore.show({
      message: `${product.name} a été ajouté au panier`,
      type: "success",
    });
  } catch (error) {
    console.error("Erreur lors de l'ajout au panier:", error);
    notificationStore.show({
      message: "Impossible d'ajouter le produit au panier",
      type: "error",
    });
  }
}

async function removeFromWishlist(productId: number) {
  try {
    await wishlistStore.removeFromWishlist(productId);
    notificationStore.show({
      message: "Produit supprimé de votre liste de souhaits",
      type: "success",
    });
  } catch (error) {
    console.error(
      "Erreur lors de la suppression de la liste de souhaits:",
      error
    );
    notificationStore.show({
      message: "Impossible de supprimer le produit de votre liste",
      type: "error",
    });
  }
}
</script>

<style scoped>
.wishlist-page {
  padding: 2rem 0;
}

.page-title {
  font-size: 2.5rem;
  margin-bottom: 2rem;
  position: relative;
  display: inline-block;
}

.page-title::after {
  content: "";
  position: absolute;
  bottom: -10px;
  left: 0;
  width: 100%;
  height: 5px;
  background-color: #000;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  border: 4px dashed #ccc;
  background-color: #f9f9f9;
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
  text-align: center;
  padding: 2rem;
  border: 4px solid #ff3b30;
  background-color: #fff5f5;
  transform: rotate(-1deg);
  margin-bottom: 1rem;
}

.retry-btn {
  background-color: #000;
  color: white;
  border: 0;
  padding: 0.8rem 1.5rem;
  font-weight: 700;
  text-transform: uppercase;
  cursor: pointer;
  margin-top: 1rem;
  transform: skew(-5deg);
}

.empty-wishlist {
  text-align: center;
  padding: 3rem;
  border: 4px dashed #000;
  font-size: 1.2rem;
}

.continue-shopping-btn {
  display: inline-block;
  margin-top: 1rem;
  background-color: #000;
  color: white;
  text-decoration: none;
  padding: 0.8rem 1.5rem;
  font-weight: 800;
  transform: skew(-5deg);
  transition: all 0.2s;
}

.continue-shopping-btn:hover {
  transform: skew(-5deg) translateY(-3px);
  box-shadow: 5px 5px 0 rgba(0, 0, 0, 0.1);
}

.wishlist-contents {
  border: 4px solid #000;
  background-color: #f5f5f5;
}

.wishlist-items {
  border-bottom: 4px solid #000;
}

.wishlist-item {
  display: grid;
  grid-template-columns: 120px 1fr auto;
  align-items: center;
  gap: 1.5rem;
  padding: 1.5rem;
  border-bottom: 2px solid #ddd;
}

.wishlist-item:last-child {
  border-bottom: none;
}

.item-image {
  height: 100px;
  width: 100px;
  overflow: hidden;
  border: 2px solid #000;
}

.item-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.item-name {
  font-weight: 800;
  margin: 0;
  font-size: 1.1rem;
  text-transform: uppercase;
  color: #000;
  text-decoration: none;
  display: inline-block;
}

.item-name:hover {
  text-decoration: underline;
}

.item-price {
  margin: 0.5rem 0;
  color: #333;
  font-weight: bold;
}

.item-date {
  font-size: 0.9rem;
  color: #666;
  margin: 0;
}

.item-actions {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.add-to-cart-btn,
.remove-btn,
.clear-wishlist-btn {
  font-weight: 700;
  text-transform: uppercase;
  cursor: pointer;
  border: 2px solid #000;
  transition: all 0.2s;
}

.add-to-cart-btn {
  background-color: #000;
  color: white;
  padding: 0.7rem 1rem;
}

.add-to-cart-btn:hover:not(:disabled) {
  background-color: #222;
  transform: translateY(-2px);
  box-shadow: 3px 3px 0 rgba(0, 0, 0, 0.2);
}

.remove-btn {
  background-color: #e5e5e5;
  color: #000;
  padding: 0.5rem;
}

.remove-btn:hover:not(:disabled) {
  background-color: #ddd;
}

.add-to-cart-btn:disabled,
.remove-btn:disabled,
.clear-wishlist-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.wishlist-summary {
  padding: 1.5rem;
  text-align: right;
}

.clear-wishlist-btn {
  background-color: #e5e5e5;
  color: #000;
  padding: 0.8rem 1.5rem;
  font-size: 0.9rem;
}

.clear-wishlist-btn:hover:not(:disabled) {
  background-color: #ddd;
}

@media (max-width: 768px) {
  .wishlist-item {
    grid-template-columns: 100px 1fr;
    grid-template-rows: auto auto;
  }

  .item-image {
    grid-row: 1 / span 2;
    width: 80px;
    height: 80px;
  }

  .item-actions {
    grid-column: 1 / span 2;
    flex-direction: row;
    margin-top: 1rem;
  }
}
</style>
