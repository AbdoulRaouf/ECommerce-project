<script setup lang="ts">
import ProductList from "@/components/ProductList.vue";
</script>

<template>
  <main>
    <div class="hero-section">
      <div class="hero-content">
        <h1 class="hero-title">DESIGN BRUTALISTE</h1>
        <p class="hero-subtitle">
          Des meubles et accessoires au style néobrutaliste
        </p>
        <div class="hero-cta">
          <a href="#products" class="cta-button">VOIR NOS PRODUITS</a>
        </div>
      </div>
    </div>

    <div id="products" class="product-section container">
      <ProductList />
    </div>
  </main>
</template>

<style scoped>
.hero-section {
  min-height: 70vh;
  display: flex;
  align-items: center;
  background-color: #f5f5f5;
  padding: 2rem;
  border-bottom: 6px solid #000;
  position: relative;
  overflow: hidden;
}

.hero-section::after {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: repeating-linear-gradient(
    -45deg,
    transparent,
    transparent 20px,
    rgba(0, 0, 0, 0.03) 20px,
    rgba(0, 0, 0, 0.03) 40px
  );
  z-index: 1;
}

.hero-content {
  max-width: 800px;
  padding: 3rem;
  background-color: white;
  border: 6px solid #000;
  box-shadow: 12px 12px 0 rgba(0, 0, 0, 1);
  position: relative;
  z-index: 2;
  transform: rotate(-2deg);
}

.hero-title {
  font-size: 4rem;
  font-weight: 900;
  margin-bottom: 1rem;
  text-transform: uppercase;
  line-height: 1;
  letter-spacing: -2px;
}

.hero-subtitle {
  font-size: 1.5rem;
  margin-bottom: 2rem;
  font-weight: 500;
}

.cta-button {
  display: inline-block;
  background-color: #000;
  color: white;
  text-decoration: none;
  padding: 1rem 2rem;
  font-size: 1.2rem;
  font-weight: 800;
  text-transform: uppercase;
  transform: skew(-5deg);
  transition: all 0.3s;
  border: 3px solid #000;
}

.cta-button:hover {
  transform: skew(-5deg) translateY(-5px);
  box-shadow: 5px 5px 0 rgba(0, 0, 0, 0.3);
}

.product-section {
  padding: 3rem 1rem;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
}

@media (max-width: 768px) {
  .hero-title {
    font-size: 2.5rem;
  }

  .hero-content {
    padding: 2rem;
  }
}
</style>
