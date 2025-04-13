<template>
  <div class="cart-page container">
    <h1 class="page-title">Mon Panier</h1>
    <ShoppingCart />
  </div>
</template>

<script setup lang="ts">
import ShoppingCart from "@/components/ShoppingCart.vue";
</script>

<style scoped>
.cart-page {
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
</style>
