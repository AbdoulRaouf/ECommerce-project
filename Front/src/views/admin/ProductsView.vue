<template>
  <div class="products-container">
    <AdminNavbar />

    <div class="products-content">
      <div class="page-header">
        <h1 class="page-title">Gestion des Produits</h1>
        <button @click="openAddProductModal" class="add-btn">
          Ajouter un produit
        </button>
      </div>

      <!-- Loading state -->
      <div v-if="adminStore.isLoading" class="loading">
        <div class="loading-spinner"></div>
        <p>Chargement des produits...</p>
      </div>

      <!-- Error state -->
      <div v-else-if="adminStore.error" class="error-message">
        <h3>Erreur de chargement</h3>
        <p>{{ adminStore.error }}</p>
        <button @click="loadProducts" class="retry-btn">Réessayer</button>
      </div>

      <!-- Products List -->
      <div
        v-else-if="adminStore.products.length > 0"
        class="products-table-container neo-border"
      >
        <table class="data-table products-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Image</th>
              <th>Nom</th>
              <th>Prix</th>
              <th>Stock</th>
              <th>Catégorie</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="product in adminStore.products" :key="product.id">
              <td>{{ product.id }}</td>
              <td class="product-image">
                <img
                  v-if="product.image"
                  :src="product.image"
                  alt="Product image"
                  class="product-thumbnail"
                />
                <div v-else class="no-image">Pas d'image</div>
              </td>
              <td>{{ product.name }}</td>
              <td>{{ formatPrice(product.price) }}</td>
              <td>
                <span
                  class="stock-count"
                  :class="{ critical: product.stock < 5 }"
                >
                  {{ product.stock }}
                </span>
              </td>
              <td>{{ product.category }}</td>
              <td class="table-actions">
                <button
                  @click="openEditProductModal(product)"
                  class="action-btn edit-btn"
                >
                  Modifier
                </button>
                <button
                  @click="confirmDeleteProduct(product)"
                  class="action-btn delete-btn"
                >
                  Supprimer
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Empty state -->
      <div v-else class="empty-products neo-border">
        <div class="empty-icon">🏷️</div>
        <h2>Aucun produit trouvé</h2>
        <p>Commencez par ajouter des produits à votre catalogue</p>
        <button @click="openAddProductModal" class="add-product-btn">
          Ajouter un produit
        </button>
      </div>

      <!-- Product Form Modal -->
      <product-form-modal
        :is-open="isModalOpen"
        :product="selectedProduct"
        @close="closeModal"
        @submit="handleProductSubmit"
      />

      <!-- Confirmation Modal -->
      <div class="modal" :class="{ 'is-active': isConfirmDialogOpen }">
        <div
          class="modal-background"
          @click="isConfirmDialogOpen = false"
        ></div>
        <div class="modal-content confirm-dialog">
          <h3>Confirmer la suppression</h3>
          <p>
            Êtes-vous sûr de vouloir supprimer ce produit ? Cette action est
            irréversible.
          </p>
          <div class="confirm-actions">
            <button @click="isConfirmDialogOpen = false" class="btn-cancel">
              Annuler
            </button>
            <button @click="deleteProduct" class="btn-confirm-delete">
              Supprimer
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import AdminNavbar from "@/components/admin/AdminNavbar.vue";
import ProductFormModal from "@/components/admin/ProductFormModal.vue";
import { useAdminStore, type Product } from "@/stores/adminStore";
import { onMounted, ref } from "vue";

// Store
const adminStore = useAdminStore();

// State
const isModalOpen = ref(false);
const selectedProduct = ref<Product | null>(null);
const isConfirmDialogOpen = ref(false);
const productToDelete = ref<Product | null>(null);

// Load products
const loadProducts = async () => {
  await adminStore.fetchProducts();
};

// Formatter les prix
const formatPrice = (price: number): string => {
  return new Intl.NumberFormat("fr-FR", {
    style: "currency",
    currency: "EUR",
  }).format(price);
};

// Methods
const openAddProductModal = () => {
  selectedProduct.value = {
    name: "",
    price: 0,
    description: "",
    category: "",
    image: "",
    stock: 0,
  };
  isModalOpen.value = true;
};

const openEditProductModal = (product: Product) => {
  selectedProduct.value = { ...product };
  isModalOpen.value = true;
};

const closeModal = () => {
  isModalOpen.value = false;
  selectedProduct.value = null;
};

const handleProductSubmit = async (productData: Product) => {
  let success = false;

  if (productData.id) {
    // Update existing product
    success = await adminStore.updateProduct(productData.id, productData);
  } else {
    // Create new product
    success = await adminStore.createProduct(productData);
  }

  if (success) {
    closeModal();
  }
};

const confirmDeleteProduct = (product: Product) => {
  productToDelete.value = product;
  isConfirmDialogOpen.value = true;
};

const deleteProduct = async () => {
  if (productToDelete.value?.id) {
    const success = await adminStore.deleteProduct(productToDelete.value.id);

    if (success) {
      isConfirmDialogOpen.value = false;
      productToDelete.value = null;
    }
  }
};

// Lifecycle hooks
onMounted(loadProducts);
</script>

<style scoped>
.products-container {
  min-height: 100vh;
  background-color: #f9f9f9;
}

.products-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2.5rem;
}

.page-title {
  font-family: var(--font-heading);
  font-size: 3rem;
  font-weight: 900;
  text-transform: uppercase;
  letter-spacing: -1px;
  position: relative;
  display: inline-block;
  margin: 0;
}

.page-title::after {
  content: "";
  position: absolute;
  bottom: -10px;
  left: 0;
  width: 100px;
  height: 6px;
  background-color: #000;
}

.add-btn {
  display: inline-block;
  padding: 0.75rem 1.5rem;
  background-color: #000;
  color: white;
  border: 3px solid #000;
  font-family: var(--font-heading);
  font-weight: 700;
  text-transform: uppercase;
  cursor: pointer;
  transform: skew(-3deg);
  transition: all 0.2s;
  font-size: 0.9rem;
}

.add-btn:hover {
  transform: skew(-3deg) translateY(-3px);
  box-shadow: 5px 5px 0 rgba(0, 0, 0, 0.3);
}

.products-table-container {
  background-color: white;
  border: 4px solid #000;
  box-shadow: 10px 10px 0 rgba(0, 0, 0, 0.8);
  margin-bottom: 2rem;
  overflow: hidden;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th {
  border-bottom: 2px solid #000;
  padding: 1rem;
  text-align: left;
  font-family: var(--font-heading);
  text-transform: uppercase;
  font-weight: 700;
  font-size: 0.9rem;
}

.data-table td {
  padding: 1rem;
  border-bottom: 1px solid #eee;
}

.data-table tr:last-child td {
  border-bottom: none;
}

.data-table tr {
  transition: background-color 0.2s;
}

.data-table tr:hover {
  background-color: #f9f9f9;
}

.product-image {
  width: 60px;
}

.product-thumbnail {
  width: 50px;
  height: 50px;
  object-fit: cover;
  border: 2px solid #000;
}

.no-image {
  width: 50px;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f0f0f0;
  border: 2px solid #000;
  font-size: 0.7rem;
  text-align: center;
}

.stock-count {
  font-weight: 700;
}

.critical {
  color: #dc2626;
}

.table-actions {
  display: flex;
  gap: 0.5rem;
}

.action-btn {
  padding: 0.4rem 0.8rem;
  border: none;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  transform: skew(-3deg);
  font-size: 0.9rem;
}

.edit-btn {
  background-color: #0369a1;
  color: white;
}

.edit-btn:hover {
  transform: skew(-3deg) translateY(-2px);
  box-shadow: 2px 2px 0 rgba(3, 105, 161, 0.3);
}

.delete-btn {
  background-color: #b91c1c;
  color: white;
}

.delete-btn:hover {
  transform: skew(-3deg) translateY(-2px);
  box-shadow: 2px 2px 0 rgba(185, 28, 28, 0.3);
}

.empty-products {
  background-color: white;
  border: 4px solid #000;
  box-shadow: 10px 10px 0 rgba(0, 0, 0, 0.8);
  padding: 3rem;
  text-align: center;
  margin-bottom: 2rem;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
  opacity: 0.3;
}

.empty-products h2 {
  font-family: var(--font-heading);
  font-size: 1.8rem;
  margin-bottom: 1rem;
}

.empty-products p {
  color: #666;
  margin-bottom: 1.5rem;
}

.add-product-btn {
  display: inline-block;
  padding: 0.75rem 1.5rem;
  background-color: #000;
  color: white;
  border: 3px solid #000;
  font-family: var(--font-heading);
  font-weight: 700;
  text-transform: uppercase;
  cursor: pointer;
  transform: skew(-3deg);
  transition: all 0.2s;
}

.add-product-btn:hover {
  transform: skew(-3deg) translateY(-3px);
  box-shadow: 5px 5px 0 rgba(0, 0, 0, 0.3);
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
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

.error-message {
  background-color: #fee2e2;
  border-left: 6px solid #ef4444;
  padding: 1.5rem;
  margin: 2rem 0;
}

.error-message h3 {
  margin-top: 0;
  font-weight: 700;
  color: #b91c1c;
}

.retry-btn {
  display: inline-block;
  margin-top: 1rem;
  background-color: #000;
  color: white;
  border: 3px solid #000;
  padding: 0.75rem 1.5rem;
  font-family: var(--font-heading);
  font-weight: 700;
  text-transform: uppercase;
  cursor: pointer;
  transform: skew(-3deg);
  transition: all 0.2s;
}

.retry-btn:hover {
  transform: skew(-3deg) translateY(-3px);
  box-shadow: 5px 5px 0 rgba(0, 0, 0, 0.3);
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 100;
  visibility: hidden;
  opacity: 0;
  transition: all 0.3s ease;
}

.modal.is-active {
  visibility: visible;
  opacity: 1;
}

.modal-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.7);
  z-index: 101;
}

.modal-content {
  position: relative;
  z-index: 102;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  background-color: white;
}

.confirm-dialog {
  border: 4px solid #000;
  box-shadow: 10px 10px 0 rgba(0, 0, 0, 0.8);
  padding: 1.5rem;
  transform: rotate(-1deg);
}

.confirm-dialog h3 {
  font-family: var(--font-heading);
  font-size: 1.5rem;
  font-weight: 700;
  margin-top: 0;
  margin-bottom: 1rem;
}

.confirm-dialog p {
  margin-bottom: 1.5rem;
}

.confirm-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
}

.btn-cancel {
  padding: 0.75rem 1.25rem;
  background-color: #e5e5e5;
  border: 3px solid #000;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  transform: skew(-3deg);
}

.btn-cancel:hover {
  background-color: #d5d5d5;
  transform: skew(-3deg) translateY(-2px);
  box-shadow: 3px 3px 0 rgba(0, 0, 0, 0.2);
}

.btn-confirm-delete {
  padding: 0.75rem 1.25rem;
  background-color: #b91c1c;
  border: 3px solid #b91c1c;
  color: white;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  transform: skew(-3deg);
}

.btn-confirm-delete:hover {
  background-color: #991b1b;
  transform: skew(-3deg) translateY(-2px);
  box-shadow: 3px 3px 0 rgba(185, 28, 28, 0.4);
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .page-title {
    font-size: 2rem;
  }

  .data-table {
    display: block;
    overflow-x: auto;
  }
}
</style>
