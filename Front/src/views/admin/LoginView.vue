<template>
  <div class="login-page">
    <div class="login-container">
      <h1 class="login-title">ADMIN LOGIN</h1>

      <div class="login-form-container">
        <form @submit.prevent="login" class="login-form">
          <!-- Error alert -->
          <div v-if="error" class="error-message">
            {{ error }}
          </div>

          <div class="form-group">
            <label for="username">Nom d'utilisateur</label>
            <input
              type="text"
              id="username"
              v-model="username"
              required
              placeholder="Entrez votre nom d'utilisateur"
            />
          </div>

          <div class="form-group">
            <label for="password">Mot de passe</label>
            <input
              type="password"
              id="password"
              v-model="password"
              required
              placeholder="Entrez votre mot de passe"
            />
          </div>

          <button type="submit" class="login-btn" :disabled="isLoading">
            {{ isLoading ? "CONNEXION EN COURS..." : "SE CONNECTER" }}
          </button>
        </form>
      </div>

      <div class="back-to-site">
        <router-link to="/" class="back-link"> ← RETOUR AU SITE </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAdminStore } from "@/stores/adminStore";
import { ref } from "vue";
import { useRouter } from "vue-router";

// État local
const username = ref("");
const password = ref("");
const error = ref<string | null>(null);

// Store et router
const adminStore = useAdminStore();
const router = useRouter();
const isLoading = ref(false);

// Fonction de connexion
const login = async () => {
  error.value = null;
  isLoading.value = true;

  try {
    const success = await adminStore.login(username.value, password.value);

    if (success) {
      router.push("/admin/dashboard");
    } else {
      error.value = "Identifiants incorrects. Veuillez réessayer.";
    }
  } catch (err: any) {
    error.value = err.message || "Une erreur est survenue. Veuillez réessayer.";
  } finally {
    isLoading.value = false;
  }
};
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  position: relative;
  background-color: var(--color-background);
  z-index: 0;
}

.login-page::before {
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
  z-index: -1;
}

.login-container {
  width: 100%;
  max-width: 500px;
  text-align: center;
}

.login-title {
  font-family: var(--font-heading);
  font-size: 3.5rem;
  font-weight: 900;
  text-transform: uppercase;
  letter-spacing: -2px;
  margin-bottom: 2.5rem;
  color: var(--color-text);
  position: relative;
  display: inline-block;
  transform: skew(-3deg);
}

.login-title::after {
  content: "";
  position: absolute;
  bottom: -10px;
  left: 0;
  width: 100%;
  height: 6px;
  background-color: #000;
}

.login-form-container {
  background-color: white;
  padding: 2.5rem;
  border: 4px solid #000;
  box-shadow: 12px 12px 0 rgba(0, 0, 0, 0.9);
  transform: rotate(-1deg);
  margin-bottom: 2rem;
}

.login-form {
  transform: rotate(1deg);
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  text-align: left;
}

label {
  font-family: var(--font-heading);
  font-weight: 700;
  text-transform: uppercase;
  font-size: 1rem;
}

input {
  padding: 1rem;
  border: 3px solid #000;
  font-family: var(--font-body);
  font-size: 1rem;
  transform: skew(-1deg);
  background-color: #f5f5f5;
  transition: all 0.2s;
}

input:focus {
  outline: none;
  border-color: #000;
  box-shadow: 4px 4px 0 rgba(0, 0, 0, 0.2);
  transform: skew(-1deg) translateY(-3px);
}

.login-btn {
  background-color: #000;
  color: white;
  border: 3px solid #000;
  padding: 1rem;
  font-family: var(--font-heading);
  font-weight: 700;
  font-size: 1.2rem;
  text-transform: uppercase;
  cursor: pointer;
  transition: all 0.2s;
  transform: skew(-3deg);
  margin-top: 1rem;
}

.login-btn:hover:not(:disabled) {
  background-color: #222;
  transform: skew(-3deg) translateY(-3px);
  box-shadow: 6px 6px 0 rgba(0, 0, 0, 0.3);
}

.login-btn:disabled {
  background-color: #999;
  cursor: not-allowed;
  opacity: 0.7;
  border-color: #999;
}

.error-message {
  background-color: #fff5f5;
  border-left: 6px solid #ff3b30;
  color: #e74c3c;
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  transform: skew(-2deg);
  font-family: var(--font-body);
}

.back-to-site {
  margin-top: 2rem;
  font-family: var(--font-heading);
}

.back-link {
  display: inline-block;
  font-weight: 700;
  color: #000;
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-bottom: 3px solid #000;
  transition: all 0.2s;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.back-link:hover {
  background-color: #f5f5f5;
  transform: skew(-3deg) translateY(-2px);
  box-shadow: 4px 4px 0 rgba(0, 0, 0, 0.2);
}

@media (max-width: 600px) {
  .login-title {
    font-size: 2.5rem;
  }

  .login-form-container {
    padding: 1.5rem;
    box-shadow: 8px 8px 0 rgba(0, 0, 0, 0.9);
  }
}
</style>
