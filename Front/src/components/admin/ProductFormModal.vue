<template>
  <div class="modal" :class="{ 'is-active': isOpen }">
    <div class="modal-background" @click="closeModal"></div>
    <div class="modal-content">
      <div class="form-container">
        <h2 class="form-title">
          {{ product?.id ? "Modifier le produit" : "Ajouter un produit" }}
        </h2>

        <form @submit.prevent="submitForm" class="product-form">
          <div class="form-group">
            <label for="name">Nom du produit</label>
            <input
              id="name"
              type="text"
              v-model="productData.name"
              required
              placeholder="Nom du produit"
            />
          </div>

          <div class="form-group">
            <label for="description">Description</label>
            <textarea
              id="description"
              v-model="productData.description"
              required
              placeholder="Description du produit"
              rows="3"
            ></textarea>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="price">Prix (€)</label>
              <input
                id="price"
                type="number"
                v-model.number="productData.price"
                step="0.01"
                min="0"
                required
                placeholder="0.00"
              />
            </div>

            <div class="form-group">
              <label for="stock">Stock</label>
              <input
                id="stock"
                type="number"
                v-model.number="productData.stock"
                min="0"
                required
                placeholder="0"
              />
            </div>
          </div>

          <div class="form-group">
            <label for="category">Catégorie</label>
            <select id="category" v-model="productData.category" required>
              <option value="">-- Sélectionner une catégorie --</option>
              <option
                v-for="category in categories"
                :key="category"
                :value="category"
              >
                {{ category }}
              </option>
            </select>
          </div>

          <div class="form-group">
            <label for="image">URL de l'image</label>
            <input
              id="image"
              type="text"
              v-model="productData.image"
              required
              placeholder="https://example.com/image.jpg"
            />
          </div>

          <div v-if="productData.image" class="image-preview">
            <img :src="productData.image" :alt="productData.name" />
          </div>

          <div class="form-actions">
            <button type="button" class="btn btn-secondary" @click="closeModal">
              Annuler
            </button>
            <button type="submit" class="btn btn-primary" :disabled="isLoading">
              {{
                isLoading
                  ? "Chargement..."
                  : product?.id
                  ? "Mettre à jour"
                  : "Ajouter"
              }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ProductService } from "@/services/api";
import { type Product } from "@/stores/adminStore";
import { defineEmits, defineProps, onMounted, reactive, ref, watch } from "vue";

const props = defineProps<{
  isOpen: boolean;
  product: Product | null;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "submit", product: Product): void;
}>();

const isLoading = ref(false);
const categories = ref<string[]>([]);
const productData = reactive<Product>({
  name: "",
  description: "",
  price: 0,
  stock: 0,
  category: "",
  image: "",
});

// Récupérer les catégories disponibles
const fetchCategories = async () => {
  try {
    const response = await ProductService.getCategories();
    categories.value = response.data;
  } catch (error) {
    console.error("Erreur lors du chargement des catégories:", error);
  }
};

// Réinitialiser le formulaire avec les données du produit lorsqu'elles changent
watch(
  () => props.product,
  (newProduct) => {
    if (newProduct && newProduct.id) {
      // Édition d'un produit existant
      Object.assign(productData, { ...newProduct });
    } else {
      // Réinitialiser pour un nouveau produit
      Object.assign(productData, {
        name: "",
        description: "",
        price: 0,
        stock: 0,
        category: "",
        image: "",
      });
    }
  },
  { immediate: true, deep: true }
);

// Soumettre le formulaire
const submitForm = () => {
  emit("submit", { ...productData });
};

// Fermer la modal
const closeModal = () => {
  emit("close");
};

onMounted(fetchCategories);
</script>

<style scoped>
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
}

.modal-content {
  position: relative;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  z-index: 101;
}

.form-container {
  background-color: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.2);
}

.form-title {
  margin-top: 0;
  margin-bottom: 1.5rem;
  font-size: 1.5rem;
  font-weight: 600;
  color: #333;
}

.product-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-row {
  display: flex;
  gap: 1rem;
}

.form-row .form-group {
  flex: 1;
}

label {
  font-weight: 600;
  color: #444;
  font-size: 0.9rem;
}

input,
select,
textarea {
  padding: 0.7rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-family: inherit;
  font-size: 1rem;
}

input:focus,
select:focus,
textarea:focus {
  border-color: #4a7bff;
  outline: none;
  box-shadow: 0 0 0 2px rgba(74, 123, 255, 0.2);
}

.image-preview {
  margin-top: 1rem;
  text-align: center;
}

.image-preview img {
  max-width: 100%;
  max-height: 200px;
  object-fit: contain;
  border-radius: 4px;
  border: 1px solid #ddd;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 1.5rem;
}

.btn {
  padding: 0.7rem 1.2rem;
  border: none;
  border-radius: 4px;
  font-family: inherit;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-primary {
  background-color: #4a7bff;
  color: white;
}

.btn-primary:hover {
  background-color: #3a6ae6;
}

.btn-primary:disabled {
  background-color: #a0b5ec;
  cursor: not-allowed;
}

.btn-secondary {
  background-color: #f0f0f0;
  color: #444;
}

.btn-secondary:hover {
  background-color: #e0e0e0;
}
</style>
