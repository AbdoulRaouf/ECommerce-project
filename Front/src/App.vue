<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import { RouterLink, RouterView } from "vue-router";
import NotificationToast from "./components/NotificationToast.vue";
import { useAuthStore } from "./stores/authStore";
import { useProductStore } from "./stores/productStore";

const productStore = useProductStore();
const authStore = useAuthStore();

// État des notifications
const toast = ref({
  visible: false,
  message: "",
  type: "success",
});

// État pour le menu déroulant utilisateur
const userMenuOpen = ref(false);

// Observer les changements dans le panier pour afficher des notifications
const cartSize = ref(0);
watch(
  () => productStore.cartItemsCount,
  (newCount, oldCount) => {
    if (newCount > oldCount) {
      showToast("Produit ajouté au panier avec succès !", "success");
    } else if (newCount < oldCount) {
      showToast("Produit retiré du panier", "info");
    }
    cartSize.value = newCount;
  }
);

// Observer les changements d'authentification pour rafraîchir le panier
watch(
  () => authStore.isLoggedIn,
  (isLoggedIn) => {
    if (isLoggedIn) {
      productStore.fetchCart();
    } else {
      productStore.clearCart();
    }
  }
);

function showToast(message: string, type = "success") {
  toast.value = {
    visible: true,
    message,
    type,
  };
}

function closeToast() {
  toast.value.visible = false;
}

function toggleUserMenu() {
  userMenuOpen.value = !userMenuOpen.value;
}

function handleLogout() {
  authStore.logout();
  productStore.clearCart();
  showToast("Vous avez été déconnecté", "info");
  userMenuOpen.value = false;
}

// Fermer le menu utilisateur lorsque l'utilisateur clique ailleurs
function handleOutsideClick(event: MouseEvent) {
  const userMenu = document.querySelector(".user-menu");
  const userButton = document.querySelector(".user-button");
  if (
    userMenu &&
    userButton &&
    !userMenu.contains(event.target as Node) &&
    !userButton.contains(event.target as Node)
  ) {
    userMenuOpen.value = false;
  }
}

onMounted(() => {
  document.addEventListener("click", handleOutsideClick);

  // Charger le profil utilisateur si l'utilisateur est connecté
  if (authStore.isLoggedIn) {
    authStore.fetchUserProfile();
  }
});
</script>

<template>
  <header class="main-header">
    <div class="container header-container">
      <div class="logo">
        <RouterLink to="/">BRUT DESIGN</RouterLink>
      </div>

      <nav class="main-nav">
        <RouterLink to="/" class="nav-link">Accueil</RouterLink>
        <RouterLink to="/about" class="nav-link">À propos</RouterLink>
      </nav>

      <div class="user-actions">
        <!-- Afficher le bouton de connexion si l'utilisateur n'est pas connecté -->
        <div v-if="!authStore.isLoggedIn" class="auth-buttons">
          <RouterLink to="/login" class="login-button">Connexion</RouterLink>
          <RouterLink to="/register" class="register-button"
            >Inscription</RouterLink
          >
        </div>

        <!-- Afficher le menu utilisateur si l'utilisateur est connecté -->
        <div v-else class="user-profile">
          <button @click="toggleUserMenu" class="user-button">
            {{ authStore.username }}
            <span class="user-icon">👤</span>
          </button>
          <div v-if="userMenuOpen" class="user-menu neo-border">
            <RouterLink to="/profile" class="menu-item">Mon profil</RouterLink>
            <RouterLink to="/orders" class="menu-item"
              >Mes commandes</RouterLink
            >
            <button @click="handleLogout" class="menu-item logout-button">
              Déconnexion
            </button>
          </div>
        </div>

        <div class="cart-icon">
          <RouterLink to="/cart" class="cart-button">
            <span class="cart-icon-text">Panier</span>
            <span class="cart-count" v-if="productStore.cartItemsCount > 0">{{
              productStore.cartItemsCount
            }}</span>
          </RouterLink>
        </div>
      </div>
    </div>
  </header>

  <RouterView />

  <footer class="main-footer">
    <div class="container footer-content">
      <div class="footer-logo">
        <RouterLink to="/">BRUT DESIGN</RouterLink>
      </div>

      <div class="footer-links">
        <div class="footer-column">
          <h3>Navigation</h3>
          <RouterLink to="/">Accueil</RouterLink>
          <RouterLink to="/about">À propos</RouterLink>
          <RouterLink to="/cart">Panier</RouterLink>
        </div>

        <div class="footer-column">
          <h3>Compte</h3>
          <RouterLink v-if="!authStore.isLoggedIn" to="/login"
            >Connexion</RouterLink
          >
          <RouterLink v-if="!authStore.isLoggedIn" to="/register"
            >Inscription</RouterLink
          >
          <RouterLink v-if="authStore.isLoggedIn" to="/profile"
            >Mon profil</RouterLink
          >
          <RouterLink v-if="authStore.isLoggedIn" to="/orders"
            >Mes commandes</RouterLink
          >
        </div>

        <div class="footer-column">
          <h3>Contact</h3>
          <p>contact@brutdesign.fr</p>
          <p>+33 1 23 45 67 89</p>
        </div>

        <div class="footer-column">
          <h3>Mentions légales</h3>
          <a href="#">Conditions générales</a>
          <a href="#">Politique de confidentialité</a>
        </div>
      </div>

      <div class="footer-bottom">
        <p>© 2025 BRUT DESIGN. Tous droits réservés.</p>
      </div>
    </div>
  </footer>

  <!-- Système de notifications -->
  <NotificationToast
    :visible="toast.visible"
    :message="toast.message"
    :type="toast.type"
    @close="closeToast"
  />
</template>

<style>
/* Importer les polices nécessaires */
@import url("https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@400;500;700&family=Inter:wght@400;700;900&display=swap");

/* Variables globales de couleur (palette sobre et néobrutaliste) */
:root {
  --color-background: #f9f9f9;
  --color-text: #222222;
  --color-black: #000000;
  --color-white: #ffffff;
  --color-accent: #e5e5e5;
  --color-border: #000000;

  --font-heading: "Space Grotesk", sans-serif;
  --font-body: "Inter", sans-serif;
}

/* Reset et styles de base */
*,
*::before,
*::after {
  box-sizing: border-box;
  margin: 0;
  font-weight: normal;
}

body {
  min-height: 100vh;
  color: var(--color-text);
  background: var(--color-background);
  line-height: 1.6;
  font-family: var(--font-body);
  font-size: 16px;
  text-rendering: optimizeLegibility;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

h1,
h2,
h3,
h4,
h5 {
  font-family: var(--font-heading);
  font-weight: 900;
}

/* Container pour aligner le contenu */
.container {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
}

/* Header */
.main-header {
  background-color: var(--color-white);
  border-bottom: 5px solid var(--color-black);
  padding: 1.5rem 0;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo a {
  font-family: var(--font-heading);
  font-size: 1.8rem;
  font-weight: 900;
  text-transform: uppercase;
  text-decoration: none;
  color: var(--color-black);
  letter-spacing: -1px;
}

.main-nav {
  display: flex;
  gap: 2rem;
}

.nav-link {
  text-decoration: none;
  color: var(--color-text);
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1px;
  font-size: 0.9rem;
  position: relative;
  padding: 0.5rem 0;
}

.nav-link:hover::after {
  content: "";
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 3px;
  background-color: var(--color-black);
}

.nav-link.router-link-active::after {
  content: "";
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 3px;
  background-color: var(--color-black);
}

.user-actions {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.auth-buttons {
  display: flex;
  gap: 1rem;
}

.login-button,
.register-button {
  text-decoration: none;
  font-weight: 700;
  padding: 0.5rem 1rem;
  transition: all 0.2s;
  text-transform: uppercase;
  font-size: 0.9rem;
}

.login-button {
  color: var(--color-black);
  border-bottom: 2px solid var(--color-black);
}

.register-button {
  background-color: var(--color-black);
  color: var(--color-white);
  transform: skew(-5deg);
}

.register-button:hover {
  transform: skew(-5deg) translateY(-3px);
  box-shadow: 3px 3px 0 rgba(0, 0, 0, 0.2);
}

.user-profile {
  position: relative;
}

.user-button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: none;
  border: none;
  font-family: var(--font-body);
  font-size: 1rem;
  font-weight: 700;
  cursor: pointer;
  padding: 0.5rem;
}

.user-icon {
  font-size: 1.2rem;
}

.user-menu {
  position: absolute;
  top: 100%;
  right: 0;
  width: 200px;
  background-color: white;
  border: 3px solid var(--color-black);
  box-shadow: 5px 5px 0 rgba(0, 0, 0, 0.8);
  margin-top: 0.5rem;
  z-index: 200;
  transform: rotate(-1deg);
}

.menu-item {
  display: block;
  padding: 0.75rem 1rem;
  text-decoration: none;
  color: var(--color-text);
  font-weight: 700;
  transition: background-color 0.2s;
  text-align: left;
  width: 100%;
  border: none;
  background: none;
  font-family: var(--font-body);
  font-size: 1rem;
  cursor: pointer;
}

.menu-item:not(:last-child) {
  border-bottom: 1px solid #eee;
}

.menu-item:hover {
  background-color: #f0f0f0;
}

.logout-button {
  color: #b91c1c;
}

.cart-button {
  display: flex;
  align-items: center;
  background-color: var(--color-black);
  color: var(--color-white);
  padding: 0.6rem 1.2rem;
  text-decoration: none;
  font-weight: 700;
  text-transform: uppercase;
  transform: skew(-5deg);
  transition: transform 0.2s;
  position: relative;
}

.cart-button:hover {
  transform: skew(-5deg) translateY(-3px);
}

.cart-count {
  position: absolute;
  top: -8px;
  right: -8px;
  background-color: var(--color-text);
  color: var(--color-white);
  width: 22px;
  height: 22px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  font-weight: 700;
}

/* Footer */
.main-footer {
  background-color: var(--color-black);
  color: var(--color-white);
  padding: 4rem 0 2rem;
  margin-top: 4rem;
}

.footer-content {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.footer-logo a {
  font-family: var(--font-heading);
  font-size: 2rem;
  font-weight: 900;
  text-transform: uppercase;
  text-decoration: none;
  color: var(--color-white);
  letter-spacing: -1px;
}

.footer-links {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 2rem;
}

.footer-column h3 {
  font-size: 1.2rem;
  margin-bottom: 1rem;
  text-transform: uppercase;
  letter-spacing: 1px;
  font-weight: 700;
}

.footer-column a,
.footer-column p {
  display: block;
  color: var(--color-white);
  text-decoration: none;
  margin-bottom: 0.5rem;
  opacity: 0.8;
  transition: opacity 0.2s;
}

.footer-column a:hover {
  opacity: 1;
}

.footer-bottom {
  border-top: 1px solid rgba(255, 255, 255, 0.2);
  padding-top: 1.5rem;
  text-align: center;
  font-size: 0.9rem;
  opacity: 0.8;
}

/* Style responsive */
@media (max-width: 768px) {
  .header-container {
    flex-wrap: wrap;
  }

  .logo {
    margin-bottom: 1rem;
    width: 100%;
    text-align: center;
  }

  .main-nav {
    flex: 1;
    justify-content: center;
  }

  .user-actions {
    width: 100%;
    justify-content: space-between;
    margin-top: 1rem;
  }

  .footer-links {
    grid-template-columns: 1fr;
    text-align: center;
  }
}
</style>
