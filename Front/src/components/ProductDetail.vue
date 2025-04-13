<template>
  <div class="product-detail" v-if="product">
    <div class="product-detail-container">
      <div class="product-image">
        <div class="image-wrapper">
          <img :src="product.image" :alt="product.name" />
        </div>
      </div>

      <div class="product-info">
        <h1 class="product-name">{{ product.name }}</h1>
        <p class="product-price">{{ formatPrice(product.price) }}</p>
        <div class="product-category">{{ product.category }}</div>

        <div class="product-description">
          <h2>Description</h2>
          <p>{{ product.description }}</p>
        </div>

        <div class="product-stock">
          <p v-if="product.stock > 0" class="in-stock">
            En stock ({{ product.stock }})
          </p>
          <p v-else class="out-of-stock">Rupture de stock</p>
        </div>

        <div class="product-actions">
          <div class="quantity-selector">
            <button
              @click="decrementQuantity"
              class="quantity-btn"
              :disabled="isLoading"
            >
              -
            </button>
            <span class="quantity">{{ quantity }}</span>
            <button
              @click="incrementQuantity"
              class="quantity-btn"
              :disabled="isLoading || quantity >= product.stock"
            >
              +
            </button>
          </div>

          <button
            @click="addToCart"
            class="add-to-cart-btn"
            :disabled="product.stock === 0 || isLoading"
          >
            {{ isLoading ? "AJOUT EN COURS..." : "AJOUTER AU PANIER" }}
          </button>
        </div>

        <div v-if="error" class="error-message">
          {{ error }}
        </div>

        <div class="product-details">
          <h2>Caractéristiques</h2>
          <ul>
            <li>Design néobrutaliste</li>
            <li>Matériaux de haute qualité</li>
            <li>Fabrication artisanale</li>
            <li>Garantie 2 ans</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
  <div v-else-if="isLoading" class="loading-container">
    <div class="loading-spinner"></div>
    <p>Chargement du produit...</p>
  </div>
  <div v-else-if="error" class="error-container">
    <h2>Erreur</h2>
    <p>{{ error }}</p>
    <button @click="loadProduct" class="retry-btn">Réessayer</button>
    <router-link to="/" class="back-link">Retourner à l'accueil</router-link>
  </div>
  <div v-else class="product-not-found">
    <h1>Produit non trouvé</h1>
    <router-link to="/" class="back-link">Retourner à l'accueil</router-link>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import { useProductStore, type Product } from "@/stores/productStore";

const props = defineProps<{
  productId?: number;
}>();

const productStore = useProductStore();
const quantity = ref(1);
const isLoading = ref(false);
const error = ref<string | null>(null);
const product = ref<Product | null>(null);

// Fonction pour charger les détails du produit
async function loadProduct() {
  if (!props.productId) return;

  isLoading.value = true;
  error.value = null;

  try {
    const fetchedProduct = await productStore.fetchProduct(props.productId);
    if (fetchedProduct) {
      product.value = fetchedProduct;
      // Réinitialiser la quantité quand on change de produit
      quantity.value = 1;
    } else {
      error.value = "Produit non trouvé";
    }
  } catch (err) {
    console.error("Erreur lors du chargement du produit:", err);
    error.value =
      "Impossible de charger les détails du produit. Veuillez réessayer plus tard.";
  } finally {
    isLoading.value = false;
  }
}

// Charger le produit au montage du composant
onMounted(() => {
  loadProduct();
});

// Observer les changements de l'ID du produit pour recharger les données
watch(
  () => props.productId,
  (newId) => {
    if (newId) {
      loadProduct();
    }
  }
);

function incrementQuantity() {
  if (product.value && quantity.value < product.value.stock) {
    quantity.value++;
  }
}

function decrementQuantity() {
  if (quantity.value > 1) {
    quantity.value--;
  }
}

async function addToCart() {
  if (!product.value || isLoading.value) return;

  isLoading.value = true;
  error.value = null;

  try {
    await productStore.addToCart(product.value, quantity.value);
    // Afficher un message de succès ici si nécessaire
  } catch (err) {
    console.error("Erreur lors de l'ajout au panier:", err);
    error.value =
      "Impossible d'ajouter le produit au panier. Veuillez réessayer.";
  } finally {
    isLoading.value = false;
  }
}

function formatPrice(price: number): string {
  return new Intl.NumberFormat("fr-FR", {
    style: "currency",
    currency: "EUR",
  }).format(price);
}
</script>

<style scoped>
.product-detail {
  padding: 2rem 0;
  max-width: 1200px;
  margin: 0 auto;
}

.product-detail-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 3rem;
  padding: 0 2rem;
}

.product-image {
  position: relative;
}

.image-wrapper {
  border: 4px solid #000;
  overflow: hidden;
  box-shadow: 10px 10px 0 #000;
  transform: rotate(-1deg);
  height: 500px;
}

.image-wrapper img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.product-info {
  padding: 2rem;
  background-color: #f5f5f5;
  border: 4px solid #000;
  box-shadow: -10px 10px 0 #000;
  transform: rotate(1deg);
}

.product-name {
  font-size: 2.5rem;
  font-weight: 900;
  margin: 0 0 1rem;
  line-height: 1.2;
  text-transform: uppercase;
  letter-spacing: -1px;
}

.product-price {
  font-size: 2rem;
  font-weight: 800;
  margin: 1rem 0;
}

.product-category {
  display: inline-block;
  background-color: #000;
  color: white;
  padding: 0.5rem 1rem;
  text-transform: uppercase;
  font-weight: 700;
  letter-spacing: 1px;
  margin-bottom: 1.5rem;
}

.product-description {
  margin: 1.5rem 0;
}

.product-description h2 {
  font-size: 1.5rem;
  font-weight: 800;
  margin-bottom: 0.5rem;
  text-transform: uppercase;
}

.product-description p {
  font-size: 1.1rem;
  line-height: 1.6;
}

.product-stock {
  margin: 1.5rem 0;
  font-weight: 700;
}

.in-stock {
  color: #2ecc71;
}

.out-of-stock {
  color: #e74c3c;
}

.product-actions {
  margin: 2rem 0;
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.quantity-selector {
  display: flex;
  align-items: center;
  border: 3px solid #000;
  overflow: hidden;
}

.quantity-btn {
  background-color: #e5e5e5;
  border: none;
  width: 40px;
  height: 40px;
  font-size: 1.5rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.quantity-btn:hover:not(:disabled) {
  background-color: #ddd;
}

.quantity-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.quantity {
  width: 50px;
  text-align: center;
  font-weight: bold;
  font-size: 1.2rem;
}

.add-to-cart-btn {
  flex-grow: 1;
  background-color: #000;
  color: white;
  border: 0;
  padding: 0 1.5rem;
  font-weight: 800;
  text-transform: uppercase;
  cursor: pointer;
  letter-spacing: 1px;
  transition: all 0.2s;
  transform: skew(-5deg);
  font-size: 1rem;
  min-height: 46px;
}

.add-to-cart-btn:hover:not(:disabled) {
  background-color: #333;
  transform: skew(-5deg) translateY(-3px);
  box-shadow: 3px 3px 0 rgba(0, 0, 0, 0.2);
}

.add-to-cart-btn:disabled {
  background-color: #ccc;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.product-details {
  margin-top: 2rem;
}

.product-details h2 {
  font-size: 1.5rem;
  font-weight: 800;
  margin-bottom: 1rem;
  text-transform: uppercase;
}

.product-details ul {
  list-style-type: none;
  padding: 0;
}

.product-details li {
  position: relative;
  padding-left: 1.5rem;
  margin-bottom: 0.5rem;
  font-size: 1.1rem;
}

.product-details li::before {
  content: "—";
  position: absolute;
  left: 0;
}

.product-not-found,
.loading-container,
.error-container {
  text-align: center;
  padding: 5rem 2rem;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
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

.product-not-found h1,
.error-container h2 {
  font-size: 2.5rem;
  margin-bottom: 2rem;
}

.back-link,
.retry-btn {
  display: inline-block;
  background-color: #000;
  color: white;
  text-decoration: none;
  padding: 1rem 2rem;
  font-weight: 700;
  text-transform: uppercase;
  transform: skew(-5deg);
  margin: 0 0.5rem;
}

.error-message {
  background-color: #fff5f5;
  border-left: 4px solid #ff3b30;
  padding: 1rem;
  margin: 1rem 0;
  color: #e74c3c;
}

@media (max-width: 900px) {
  .product-detail-container {
    grid-template-columns: 1fr;
  }

  .image-wrapper {
    height: 400px;
  }

  .product-actions {
    flex-direction: column;
  }

  .quantity-selector {
    width: 100%;
    justify-content: center;
  }
}
</style>
