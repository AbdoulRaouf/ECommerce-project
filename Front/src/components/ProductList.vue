<template>
  <div class="product-list-container">
    <div class="filters">
      <div class="search-bar">
        <input
          type="text"
          v-model="searchQuery"
          @input="onSearchInput"
          placeholder="Rechercher des produits..."
          class="search-input"
        />
      </div>

      <div class="category-filters">
        <button
          class="category-btn"
          :class="{ active: selectedCategory === null }"
          @click="selectCategory(null)"
        >
          Tous
        </button>
        <button
          v-for="category in productStore.categories"
          :key="category"
          class="category-btn"
          :class="{ active: selectedCategory === category }"
          @click="selectCategory(category)"
        >
          {{ category }}
        </button>
      </div>
    </div>

    <!-- Affichage de l'état de chargement -->
    <div v-if="productStore.isLoading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>Chargement des produits...</p>
    </div>

    <!-- Affichage des erreurs -->
    <div v-else-if="productStore.error" class="error-message">
      <p>{{ productStore.error }}</p>
      <button @click="productStore.fetchProducts" class="retry-btn">
        Réessayer
      </button>
    </div>

    <div v-else class="product-list">
      <div
        v-if="productStore.filteredProducts.length === 0"
        class="no-products"
      >
        <h2>Aucun produit ne correspond à votre recherche</h2>
      </div>
      <div v-else class="products-grid">
        <ProductCard
          v-for="product in productStore.filteredProducts"
          :key="product.id"
          :product="product"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from "vue";
import { useProductStore } from "../stores/productStore";
import ProductCard from "./ProductCard.vue";

const productStore = useProductStore();
const selectedCategory = ref<string | null>(null);
const searchQuery = ref("");

function selectCategory(category: string | null) {
  selectedCategory.value = category;
  productStore.setFilterCategory(category);
}

function onSearchInput() {
  productStore.setSearchQuery(searchQuery.value);
}

// Rafraîchir les produits quand le composant est monté
onMounted(() => {
  productStore.fetchProducts();
});

// Observer les changements dans le store
watch(
  () => productStore.filterCategory,
  (newCategory) => {
    selectedCategory.value = newCategory;
  }
);

watch(
  () => productStore.searchQuery,
  (newQuery) => {
    searchQuery.value = newQuery;
  }
);
</script>

<style scoped>
.product-list-container {
  padding: 2rem 0;
}

.filters {
  margin-bottom: 2rem;
}

.search-bar {
  margin-bottom: 1rem;
}

.search-input {
  width: 100%;
  padding: 0.8rem;
  font-size: 1rem;
  border: 4px solid #000;
  background-color: #f1f1f1;
  transform: skew(-3deg);
  font-weight: 600;
}

.search-input:focus {
  outline: none;
  box-shadow: 5px 5px 0 rgba(0, 0, 0, 0.3);
}

.category-filters {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 1.5rem;
}

.category-btn {
  background-color: #e5e5e5;
  border: 3px solid #000;
  padding: 0.5rem 1rem;
  font-weight: 700;
  text-transform: uppercase;
  cursor: pointer;
  transform: skew(-5deg);
  transition: all 0.2s;
  font-size: 0.9rem;
}

.category-btn:hover {
  background-color: #ddd;
  transform: skew(-5deg) translateY(-2px);
  box-shadow: 3px 3px 0 #000;
}

.category-btn.active {
  background-color: #000;
  color: white;
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
}

.no-products {
  text-align: center;
  padding: 3rem;
  border: 4px dashed #000;
  font-weight: 800;
  text-transform: uppercase;
  transform: rotate(-1deg);
}

/* Styles pour le chargement */
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
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

/* Styles pour les messages d'erreur */
.error-message {
  text-align: center;
  padding: 2rem;
  border: 4px solid #ff3b30;
  background-color: #fff5f5;
  transform: rotate(-1deg);
  margin: 1rem 0;
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

@media (max-width: 768px) {
  .products-grid {
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 1.5rem;
  }
}
</style>
