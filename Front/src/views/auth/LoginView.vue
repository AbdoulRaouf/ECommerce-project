<template>
  <div class="auth-container">
    <div class="container">
      <div class="auth-card neo-border">
        <h1 class="auth-title">Connexion</h1>

        <div v-if="authStore.error" class="error-message">
          {{ authStore.error }}
        </div>

        <form @submit.prevent="handleLogin" class="auth-form">
          <div class="form-group">
            <label for="email">Email</label>
            <input
              type="email"
              id="email"
              v-model="email"
              required
              class="form-input"
              placeholder="votre@email.com"
            />
          </div>

          <div class="form-group">
            <label for="password">Mot de passe</label>
            <input
              type="password"
              id="password"
              v-model="password"
              required
              class="form-input"
              placeholder="Votre mot de passe"
            />
          </div>

          <button
            type="submit"
            class="submit-button"
            :disabled="authStore.isLoading"
          >
            <span v-if="authStore.isLoading">Connexion en cours...</span>
            <span v-else>Se connecter</span>
          </button>
        </form>

        <div class="auth-links">
          <p>
            Pas encore de compte ?
            <router-link to="/register" class="auth-link"
              >Créer un compte</router-link
            >
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "../../stores/authStore";
import { useProductStore } from "../../stores/productStore";

const authStore = useAuthStore();
const productStore = useProductStore();
const router = useRouter();
const route = useRoute();

const email = ref("");
const password = ref("");

async function handleLogin() {
  const success = await authStore.login(email.value, password.value);
  if (success) {
    // Après la connexion, on charge le panier de l'utilisateur
    await productStore.fetchCart();

    // Redirection vers la page demandée si elle existe, sinon vers la page d'accueil
    const redirectPath = (route.query.redirect as string) || "/";
    router.push(redirectPath);
  }
}
</script>

<style scoped>
.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 3rem 1rem;
  min-height: calc(100vh - 260px);
}

.auth-card {
  background-color: white;
  border: 4px solid #000;
  box-shadow: 8px 8px 0 rgba(0, 0, 0, 0.8);
  padding: 2.5rem;
  width: 100%;
  max-width: 500px;
  transform: rotate(-1deg);
}

.auth-title {
  font-family: var(--font-heading);
  font-size: 2.5rem;
  font-weight: 900;
  margin-bottom: 2rem;
  text-transform: uppercase;
  letter-spacing: -1px;
  position: relative;
  display: inline-block;
}

.auth-title::after {
  content: "";
  position: absolute;
  bottom: -10px;
  left: 0;
  width: 60px;
  height: 4px;
  background-color: #000;
}

.auth-form {
  margin-top: 2rem;
  transform: rotate(1deg);
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 700;
  text-transform: uppercase;
  font-size: 0.9rem;
  letter-spacing: 0.5px;
}

.form-input {
  width: 100%;
  padding: 0.75rem 1rem;
  font-size: 1rem;
  border: 3px solid #000;
  background-color: #f9f9f9;
  transition: all 0.2s;
}

.form-input:focus {
  outline: none;
  box-shadow: 4px 4px 0 rgba(0, 0, 0, 0.2);
  transform: translateY(-2px);
}

.submit-button {
  width: 100%;
  background-color: var(--color-black);
  color: var(--color-white);
  border: none;
  padding: 1rem;
  font-family: var(--font-heading);
  font-size: 1.1rem;
  font-weight: 700;
  text-transform: uppercase;
  cursor: pointer;
  transform: skew(-5deg);
  transition: all 0.2s;
  margin-top: 1rem;
}

.submit-button:hover:not(:disabled) {
  transform: skew(-5deg) translateY(-3px);
  box-shadow: 5px 5px 0 rgba(0, 0, 0, 0.2);
}

.submit-button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.auth-links {
  margin-top: 2rem;
  text-align: center;
  font-size: 0.9rem;
}

.auth-link {
  color: var(--color-black);
  font-weight: 700;
  text-decoration: underline;
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
