<template>
  <div class="product-detail-view">
    <div class="container">
      <ProductDetail :productId="Number(productId)" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from "vue-router";
import ProductDetail from "@/components/ProductDetail.vue";

const route = useRoute();
const productId = route.params.id as string;
</script>

<style scoped>
.product-detail-view {
  padding: 2rem 0;
}
</style>
