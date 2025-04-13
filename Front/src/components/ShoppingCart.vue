<template>
  <div class="shopping-cart">
    <h2 class="cart-title">VOTRE PANIER</h2>

    <!-- État de chargement -->
    <div v-if="productStore.isLoading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>Chargement du panier...</p>
    </div>

    <!-- Gestion des erreurs -->
    <div v-else-if="productStore.error" class="error-message">
      <p>{{ productStore.error }}</p>
      <button @click="productStore.fetchCart" class="retry-btn">
        Réessayer
      </button>
    </div>

    <div v-else-if="productStore.cart.length === 0" class="empty-cart">
      <p>Votre panier est vide</p>
      <router-link to="/" class="continue-shopping-btn"
        >CONTINUER VOS ACHATS</router-link
      >
    </div>

    <div v-else class="cart-contents">
      <div class="cart-items">
        <div v-for="item in productStore.cart" :key="item.id" class="cart-item">
          <div class="item-image">
            <img :src="item.product?.image" :alt="item.product?.name" />
          </div>

          <div class="item-details">
            <h3 class="item-name">{{ item.product?.name }}</h3>
            <p class="item-price">
              {{ formatPrice(item.product?.price || 0) }}
            </p>
          </div>

          <div class="item-quantity">
            <button
              @click="updateQuantity(item.id || 0, item.quantity - 1)"
              class="quantity-btn"
              :disabled="productStore.isLoading"
            >
              -
            </button>
            <span class="quantity">{{ item.quantity }}</span>
            <button
              @click="updateQuantity(item.id || 0, item.quantity + 1)"
              class="quantity-btn"
              :disabled="
                productStore.isLoading ||
                (item.product && item.quantity >= item.product.stock)
              "
            >
              +
            </button>
          </div>

          <div class="item-total">
            <p>{{ formatPrice((item.product?.price || 0) * item.quantity) }}</p>
          </div>

          <button
            @click="removeItem(item.id || 0)"
            class="remove-btn"
            :disabled="productStore.isLoading"
          >
            ×
          </button>
        </div>
      </div>

      <div class="cart-summary">
        <div class="summary-row">
          <span>Total</span>
          <span class="cart-total">{{
            formatPrice(productStore.cartTotal)
          }}</span>
        </div>

        <button
          @click="checkout"
          class="checkout-btn"
          :disabled="productStore.isLoading || checkoutInProgress"
        >
          {{ checkoutInProgress ? "TRAITEMENT..." : "PASSER LA COMMANDE" }}
        </button>

        <button
          @click="productStore.clearCart"
          class="clear-cart-btn"
          :disabled="productStore.isLoading"
        >
          VIDER LE PANIER
        </button>
      </div>
    </div>

    <!-- Message de confirmation de commande -->
    <div v-if="orderSuccessful" class="order-success">
      <div class="success-content">
        <h3>Commande confirmée !</h3>
        <p>Votre commande a été passée avec succès.</p>
        <router-link to="/" class="continue-btn"
          >Continuer mes achats</router-link
        >
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useProductStore } from "@/stores/productStore";
import { OrderService } from "@/services/api";

const productStore = useProductStore();
const checkoutInProgress = ref(false);
const orderSuccessful = ref(false);

onMounted(() => {
  productStore.fetchCart();
});

function formatPrice(price: number): string {
  return new Intl.NumberFormat("fr-FR", {
    style: "currency",
    currency: "EUR",
  }).format(price);
}

function updateQuantity(itemId: number, quantity: number) {
  if (quantity < 1) quantity = 1;
  productStore.updateCartItemQuantity(itemId, quantity);
}

function removeItem(itemId: number) {
  productStore.removeFromCart(itemId);
}

async function checkout() {
  if (productStore.isLoading || checkoutInProgress.value) return;

  checkoutInProgress.value = true;

  try {
    await OrderService.createOrder();
    orderSuccessful.value = true;

    // Rafraîchir le panier après la commande
    await productStore.fetchCart();

    // Masquer le message après quelques secondes
    setTimeout(() => {
      orderSuccessful.value = false;
    }, 5000);
  } catch (error) {
    console.error("Erreur lors du passage de la commande:", error);
    alert(
      "Une erreur est survenue lors du traitement de votre commande. Veuillez réessayer ultérieurement."
    );
  } finally {
    checkoutInProgress.value = false;
  }
}
</script>

<style scoped>
.shopping-cart {
  max-width: 1000px;
  margin: 0 auto;
  padding: 2rem;
  position: relative;
}

.cart-title {
  font-size: 2.5rem;
  margin-bottom: 2rem;
  font-weight: 900;
  letter-spacing: -1px;
  text-transform: uppercase;
  transform: skew(-5deg);
  display: inline-block;
  background-color: #000;
  color: white;
  padding: 0.5rem 1.5rem;
}

.empty-cart {
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

.cart-contents {
  border: 4px solid #000;
  background-color: #f5f5f5;
}

.cart-items {
  border-bottom: 4px solid #000;
}

.cart-item {
  display: grid;
  grid-template-columns: 100px 1fr auto auto auto;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
  border-bottom: 2px solid #ddd;
  position: relative;
}

.cart-item:last-child {
  border-bottom: none;
}

.item-image {
  height: 80px;
  width: 80px;
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
}

.item-price {
  margin: 0.5rem 0 0;
  color: #333;
}

.item-quantity {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  border: 3px solid #000;
  padding: 0.3rem;
}

.quantity-btn {
  background-color: #e5e5e5;
  border: 2px solid #000;
  width: 30px;
  height: 30px;
  font-weight: bold;
  cursor: pointer;
  font-size: 1.2rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.quantity-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.quantity {
  font-weight: bold;
  width: 30px;
  text-align: center;
}

.item-total {
  font-weight: 800;
  font-size: 1.1rem;
}

.remove-btn {
  position: absolute;
  top: 1rem;
  right: 1rem;
  background-color: #000;
  color: white;
  border: none;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  font-size: 1.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: transform 0.2s;
}

.remove-btn:hover:not(:disabled) {
  transform: scale(1.1);
}

.remove-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.cart-summary {
  padding: 1.5rem;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 0;
  font-size: 1.2rem;
  font-weight: bold;
  text-transform: uppercase;
}

.cart-total {
  font-size: 1.5rem;
  font-weight: 900;
}

.checkout-btn,
.clear-cart-btn {
  display: block;
  width: 100%;
  padding: 1rem;
  margin-top: 1rem;
  font-weight: 800;
  text-transform: uppercase;
  cursor: pointer;
  font-size: 1rem;
  border: 3px solid #000;
  transition: all 0.2s;
  transform: skew(-3deg);
}

.checkout-btn {
  background-color: #000;
  color: white;
}

.checkout-btn:hover:not(:disabled) {
  background-color: #222;
  transform: skew(-3deg) translateY(-3px);
  box-shadow: 5px 5px 0 rgba(0, 0, 0, 0.3);
}

.clear-cart-btn {
  background-color: #e5e5e5;
  color: #000;
  margin-top: 0.5rem;
}

.clear-cart-btn:hover:not(:disabled) {
  background-color: #ddd;
}

.checkout-btn:disabled,
.clear-cart-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

/* État de chargement */
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

/* Message d'erreur */
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

/* Message de succès de commande */
.order-success {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.success-content {
  background-color: white;
  border: 6px solid #000;
  padding: 2rem;
  text-align: center;
  max-width: 500px;
  transform: rotate(-1deg);
  box-shadow: 10px 10px 0 rgba(0, 0, 0, 0.8);
}

.success-content h3 {
  font-size: 1.8rem;
  margin-bottom: 1rem;
}

.continue-btn {
  display: inline-block;
  margin-top: 1rem;
  background-color: #000;
  color: white;
  text-decoration: none;
  padding: 0.8rem 1.5rem;
  font-weight: 700;
  text-transform: uppercase;
}

@media (max-width: 768px) {
  .cart-item {
    grid-template-columns: 80px 1fr;
    grid-template-rows: auto auto auto;
    gap: 0.5rem;
  }

  .item-image {
    grid-row: 1 / span 2;
  }

  .item-quantity {
    grid-column: 1 / span 2;
    justify-self: start;
  }

  .item-total {
    justify-self: end;
  }
}
</style>
