<template>
  <div
    class="morpion-cell"
    :class="{
      'winning-cell': isWinningCell,
      'x-cell': value === 'X',
      'o-cell': value === 'O',
    }"
    @click="onClick"
  >
    {{ value }}
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import type { Player } from "@/stores/morpionStore";

const props = defineProps<{
  value: Player;
  row: number;
  col: number;
  isWinningCell: boolean;
}>();

const emit = defineEmits<{
  (e: "cell-click", row: number, col: number): void;
}>();

const onClick = () => {
  if (!props.value) {
    emit("cell-click", props.row, props.col);
  }
};
</script>

<style scoped>
.morpion-cell {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100px;
  height: 100px;
  font-size: 60px;
  font-weight: bold;
  cursor: pointer;
  border: 3px solid #000;
  background-color: #fff;
  transition: all 0.2s;
  user-select: none;
  position: relative;
  box-shadow: inset 0 0 0 3px rgba(0, 0, 0, 0.1);
}

.morpion-cell:hover:not(.x-cell):not(.o-cell) {
  background-color: #f0f0f0;
  transform: translate(-2px, -2px);
  box-shadow: 4px 4px 0 #000;
}

.x-cell {
  color: #e74c3c;
  background-color: rgba(231, 76, 60, 0.1);
  font-size: 65px;
  text-shadow: 2px 2px 0 rgba(0, 0, 0, 0.2);
}

.o-cell {
  color: #3498db;
  background-color: rgba(52, 152, 219, 0.1);
  font-size: 65px;
  text-shadow: 2px 2px 0 rgba(0, 0, 0, 0.2);
}

.winning-cell {
  background-color: #a3f7bf;
  transform: translate(-3px, -3px);
  box-shadow: 5px 5px 0 #000;
  z-index: 1;
  border-width: 4px;
}

/* Applique différents angles de rotation aux cellules pour un effet néobrutalism */
.morpion-cell:nth-child(odd) {
  transform: rotate(0.5deg);
}

.morpion-cell:nth-child(even) {
  transform: rotate(-0.5deg);
}

.winning-cell:nth-child(odd),
.winning-cell:nth-child(even) {
  transform: translate(-3px, -3px) rotate(0deg);
}
</style>
