<template>
  <Transition name="toast">
    <div v-if="visible" class="toast-notification" :class="type">
      <div class="toast-content">
        <span class="toast-message">{{ message }}</span>
        <button class="toast-close" @click="close">×</button>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";

const props = defineProps({
  message: {
    type: String,
    default: "",
  },
  type: {
    type: String,
    default: "success",
    validator: (value: string) =>
      ["success", "error", "warning", "info"].includes(value),
  },
  duration: {
    type: Number,
    default: 3000,
  },
  visible: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["close"]);

let timeout: number | null = null;

watch(
  () => props.visible,
  (newValue) => {
    if (newValue) {
      if (timeout) clearTimeout(timeout);
      timeout = setTimeout(() => {
        close();
      }, props.duration) as unknown as number;
    }
  }
);

function close() {
  if (timeout) clearTimeout(timeout);
  emit("close");
}
</script>

<style scoped>
.toast-notification {
  position: fixed;
  bottom: 20px;
  right: 20px;
  z-index: 9999;
  max-width: 350px;
  border: 4px solid #000;
  box-shadow: 6px 6px 0 #000;
  transform: rotate(-1deg);
  font-family: var(--font-body);
  font-weight: 700;
  padding: 0;
  overflow: hidden;
}

.toast-notification.success {
  background-color: #4cd964;
}

.toast-notification.error {
  background-color: #ff3b30;
}

.toast-notification.warning {
  background-color: #ffcc00;
}

.toast-notification.info {
  background-color: #34aadc;
}

.toast-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 15px 20px;
}

.toast-message {
  margin-right: 30px;
  font-size: 1rem;
  color: #000;
}

.toast-close {
  background: #000;
  color: #fff;
  border: none;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  cursor: pointer;
  padding: 0;
  line-height: 1;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  transform: translateX(100%) rotate(-1deg);
  opacity: 0;
}

.toast-leave-to {
  transform: translateX(100%) rotate(-1deg);
  opacity: 0;
}
</style>
