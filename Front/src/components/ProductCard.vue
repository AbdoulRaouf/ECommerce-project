<template>
  <div class="product-card">
    <router-link :to="`/product/${product.id}`" class="product-link">
      <div class="product-image-container">
        <img :src="product.image" :alt="product.name" class="product-image" />
      </div>

      <div class="product-info">
        <h3 class="product-name">{{ product.name }}</h3>
        <p class="product-price">{{ formatPrice(product.price) }}</p>
        <p class="product-category">{{ product.category }}</p>
      </div>
    </router-link>

    <div class="product-actions">
      <button @click.prevent="addToCart" class="add-to-cart-btn">
        AJOUTER AU PANIER
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useProductStore, type Product } from "@/stores/productStore";

const props = defineProps<{
  product: Product;
}>();

const productStore = useProductStore();

function addToCart() {
  productStore.addToCart(props.product);
}

function formatPrice(price: number): string {
  return new Intl.NumberFormat("fr-FR", {
    style: "currency",
    currency: "EUR",
  }).format(price);
}
</script>

<style scoped>
.product-card {
  background-color: #f1f1f1;
  border: 4px solid #000;
  overflow: hidden;
  transition: transform 0.3s, box-shadow 0.3s;
  position: relative;
  box-shadow: 6px 6px 0 #000;
  transform: rotate(-1deg);
  height: 100%;
  display: flex;
  flex-direction: column;
}

.product-card:hover {
  transform: rotate(0deg) translateY(-5px);
  box-shadow: 10px 10px 0 #000;
}

.product-link {
  text-decoration: none;
  color: inherit;
  display: flex;
  flex-direction: column;
  flex-grow: 1;
}

.product-image-container {
  overflow: hidden;
  height: 250px;
  border-bottom: 4px solid #000;
}

.product-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s;
}

.product-card:hover .product-image {
  transform: scale(1.05);
}

.product-info {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  flex-grow: 1;
}

.product-name {
  font-size: 1.5rem;
  font-weight: 900;
  margin: 0;
  line-height: 1.2;
  text-transform: uppercase;
  letter-spacing: -1px;
}

.product-price {
  font-size: 1.3rem;
  font-weight: 800;
  margin: 0.5rem 0;
  color: #333;
}

.product-category {
  display: inline-block;
  background-color: #333;
  color: white;
  padding: 0.25rem 0.5rem;
  font-size: 0.8rem;
  text-transform: uppercase;
  font-weight: 700;
  letter-spacing: 1px;
  transform: skew(-5deg);
  margin: 0.5rem 0;
}

.product-actions {
  padding: 0 1.5rem 1.5rem;
}

.add-to-cart-btn {
  background-color: #000;
  color: white;
  border: 0;
  padding: 0.8rem 1rem;
  font-weight: 800;
  text-transform: uppercase;
  cursor: pointer;
  letter-spacing: 1px;
  width: 100%;
  transition: all 0.2s;
  font-size: 1rem;
}

.add-to-cart-btn:hover {
  background-color: #333;
  transform: translateY(-2px);
}

.add-to-cart-btn:active {
  transform: translateY(2px);
}
</style>
