<template>
  <nav class="admin-navbar">
    <div class="logo">
      <router-link to="/admin/dashboard"> ADMIN PANEL </router-link>
    </div>

    <div class="nav-links">
      <router-link to="/admin/dashboard" class="nav-link" active-class="active">
        <div class="icon">📊</div>
        <span>Dashboard</span>
      </router-link>

      <router-link to="/admin/products" class="nav-link" active-class="active">
        <div class="icon">📦</div>
        <span>Produits</span>
      </router-link>

      <router-link to="/admin/orders" class="nav-link" active-class="active">
        <div class="icon">🧾</div>
        <span>Commandes</span>
      </router-link>
    </div>

    <div class="user-actions">
      <button @click="logout" class="logout-btn">
        <div class="icon">🚪</div>
        <span>Déconnexion</span>
      </button>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { useAdminStore } from "@/stores/adminStore";
import { useRouter } from "vue-router";

const adminStore = useAdminStore();
const router = useRouter();

const logout = () => {
  adminStore.logout();
  router.push("/admin/login");
};
</script>

<style scoped>
.admin-navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #111;
  padding: 0.8rem 1.5rem;
  color: white;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
  font-family: "Space Grotesk", sans-serif;
}

.logo a {
  color: white;
  text-decoration: none;
  font-size: 1.4rem;
  font-weight: 800;
  letter-spacing: 1px;
}

.nav-links {
  display: flex;
  gap: 1rem;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #ccc;
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.nav-link:hover {
  color: white;
  background-color: #333;
}

.nav-link.active {
  color: white;
  background-color: #333;
  box-shadow: 0 0 0 2px #555;
}

.user-actions {
  display: flex;
  align-items: center;
}

.logout-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background-color: #333;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  font-family: inherit;
  font-size: 0.9rem;
  transition: all 0.2s ease;
}

.logout-btn:hover {
  background-color: #d32f2f;
}

.icon {
  font-size: 1.2rem;
}

@media (max-width: 768px) {
  .admin-navbar {
    flex-direction: column;
    gap: 1rem;
    padding: 1rem;
  }

  .nav-links {
    width: 100%;
    justify-content: space-around;
  }

  .nav-link {
    flex-direction: column;
    font-size: 0.8rem;
    padding: 0.5rem;
  }

  .user-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>
