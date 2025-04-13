<template>
  <div class="auth-container">
    <div class="container">
      <div class="auth-card neo-border">
        <h1 class="auth-title">Inscription</h1>

        <div v-if="authStore.error" class="error-message">
          {{ authStore.error }}
        </div>

        <form @submit.prevent="handleRegister" class="auth-form">
          <div class="form-group">
            <label for="name">Nom complet</label>
            <input
              type="text"
              id="name"
              v-model="userData.name"
              required
              class="form-input"
              placeholder="Votre nom"
            />
          </div>

          <div class="form-group">
            <label for="email">Email</label>
            <input
              type="email"
              id="email"
              v-model="userData.email"
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
              v-model="userData.password"
              required
              minlength="6"
              class="form-input"
              placeholder="Minimum 6 caractères"
            />
          </div>

          <div class="form-group">
            <label for="address">Adresse (optionnel)</label>
            <textarea
              id="address"
              v-model="userData.address"
              class="form-input textarea"
              placeholder="Votre adresse de livraison"
            ></textarea>
          </div>

          <div class="form-group">
            <label for="phone">Téléphone (optionnel)</label>
            <input
              type="tel"
              id="phone"
              v-model="userData.phone"
              class="form-input"
              placeholder="Votre numéro de téléphone"
            />
          </div>

          <button
            type="submit"
            class="submit-button"
            :disabled="authStore.isLoading"
          >
            <span v-if="authStore.isLoading">Inscription en cours...</span>
            <span v-else>Créer mon compte</span>
          </button>
        </form>

        <div class="auth-links">
          <p>
            Déjà un compte ?
            <router-link to="/login" class="auth-link"
              >Se connecter</router-link
            >
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "../../stores/authStore";
import { useProductStore } from "../../stores/productStore";

const authStore = useAuthStore();
const productStore = useProductStore();
const router = useRouter();

const userData = reactive({
  name: "",
  email: "",
  password: "",
  address: "",
  phone: "",
});

async function handleRegister() {
  const success = await authStore.register(userData);
  if (success) {
    // Après l'inscription, on initialise le panier de l'utilisateur
    await productStore.fetchCart();
    router.push("/");
  }
}
</script>

<style scoped>
.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 3rem 1rem;
}

.auth-card {
  background-color: white;
  border: 4px solid #000;
  box-shadow: 8px 8px 0 rgba(0, 0, 0, 0.8);
  padding: 2.5rem;
  width: 100%;
  max-width: 500px;
  transform: rotate(-1deg);
  margin: 2rem 0;
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

.textarea {
  min-height: 100px;
  resize: vertical;
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
