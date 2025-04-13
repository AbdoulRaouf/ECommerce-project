import { createRouter, createWebHistory } from "vue-router";
import { useAdminStore } from "../stores/adminStore"; // Import admin store for guard
import { useAuthStore } from "../stores/authStore"; // Import auth store for user routes
import HomeView from "../views/HomeView.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: HomeView,
    },
    {
      path: "/about",
      name: "about",
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import("../views/AboutView.vue"),
    },
    {
      path: "/cart",
      name: "cart",
      component: () => import("../views/CartView.vue"),
    },
    {
      path: "/product/:id",
      name: "product-detail",
      component: () => import("../views/ProductDetailView.vue"),
    },
    // User Authentication Routes
    {
      path: "/login",
      name: "login",
      component: () => import("../views/auth/LoginView.vue"),
      meta: { guestOnly: true },
    },
    {
      path: "/register",
      name: "register",
      component: () => import("../views/auth/RegisterView.vue"),
      meta: { guestOnly: true },
    },
    {
      path: "/profile",
      name: "profile",
      component: () => import("../views/auth/ProfileView.vue"),
      meta: { requiresUserAuth: true },
    },
    {
      path: "/orders",
      name: "orders",
      component: () => import("../views/auth/OrdersView.vue"),
      meta: { requiresUserAuth: true },
    },
    {
      path: "/wishlist",
      name: "wishlist",
      component: () => import("../views/auth/WishlistView.vue"),
      meta: { requiresUserAuth: true },
    },
    // Admin Routes
    {
      path: "/admin/login",
      name: "admin-login",
      component: () => import("../views/admin/LoginView.vue"),
      meta: { requiresAuth: false },
    },
    {
      path: "/admin",
      redirect: "/admin/dashboard",
      meta: { requiresAuth: true },
    },
    {
      path: "/admin/dashboard",
      name: "admin-dashboard",
      component: () => import("../views/admin/DashboardView.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/admin/products",
      name: "admin-products",
      component: () => import("../views/admin/ProductsView.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/admin/orders",
      name: "admin-orders",
      component: () => import("../views/admin/OrdersView.vue"),
      meta: { requiresAuth: true },
    },
  ],
});

// Navigation guard for admin routes and user authenticated routes
router.beforeEach((to, from, next) => {
  const adminStore = useAdminStore();
  const authStore = useAuthStore();

  // Check if the route requires admin authentication
  if (to.matched.some((record) => record.meta.requiresAuth === true)) {
    if (!adminStore.isAuthenticated) {
      // Redirect to admin login page if not authenticated
      next({ name: "admin-login" });
    } else {
      // Allow access if authenticated
      next();
    }
  }
  // Check if the route requires user authentication
  else if (to.matched.some((record) => record.meta.requiresUserAuth === true)) {
    if (!authStore.isLoggedIn) {
      // Redirect to user login page if not authenticated
      next({ name: "login", query: { redirect: to.fullPath } });
    } else {
      // Allow access if authenticated
      next();
    }
  }
  // Check if the route is for guests only (login/register)
  else if (to.matched.some((record) => record.meta.guestOnly === true)) {
    if (authStore.isLoggedIn) {
      // Redirect to home if already logged in
      next({ name: "home" });
    } else {
      // Allow access for guests
      next();
    }
  } else {
    // Allow access to non-protected routes
    next();
  }
});

export default router;
