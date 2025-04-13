<template>
  <div class="morpion-container">
    <div class="game-info">
      <h1>Jeu de Morpion</h1>
      <div v-if="!morpionStore.isGameOver" class="current-player">
        Tour du joueur:
        <span
          :class="{
            'player-x': morpionStore.currentPlayer === 'X',
            'player-o': morpionStore.currentPlayer === 'O',
          }"
          >{{ morpionStore.currentPlayer }}</span
        >
      </div>
      <div v-else-if="morpionStore.winner" class="winner-info">
        Le joueur
        <span
          :class="{
            'player-x': morpionStore.winner === 'X',
            'player-o': morpionStore.winner === 'O',
          }"
          >{{ morpionStore.winner }}</span
        >
        a gagné!
      </div>
      <div v-else class="draw-info">Match nul!</div>
      <button @click="morpionStore.resetGame" class="reset-button">
        Nouvelle partie
      </button>
    </div>

    <div class="board">
      <div
        v-for="(row, rowIndex) in morpionStore.board"
        :key="rowIndex"
        class="board-row"
      >
        <MorpionCell
          v-for="(cell, colIndex) in row"
          :key="colIndex"
          :value="cell"
          :row="rowIndex"
          :col="colIndex"
          :is-winning-cell="isWinningCell(rowIndex, colIndex)"
          @cell-click="onCellClick"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useMorpionStore } from "@/stores/morpionStore";
import MorpionCell from "./MorpionCell.vue";
import confetti from "canvas-confetti";
import { watch } from "vue";

const morpionStore = useMorpionStore();

const onCellClick = (row: number, col: number) => {
  morpionStore.makeMove(row, col);
};

const isWinningCell = (row: number, col: number): boolean => {
  return morpionStore.winningCells.some(
    (cell) => cell[0] === row && cell[1] === col
  );
};

// Lancer des confettis en cas de victoire
watch(
  () => morpionStore.winner,
  (newValue) => {
    if (newValue) {
      // Lancer des confettis avec les couleurs du joueur gagnant
      const colors =
        newValue === "X" ? ["#e74c3c", "#c0392b"] : ["#3498db", "#2980b9"];

      // Animation de confettis
      confetti({
        particleCount: 200,
        spread: 100,
        origin: { y: 0.6 },
        colors: colors,
      });

      // Animation supplémentaire après un court délai
      setTimeout(() => {
        confetti({
          particleCount: 100,
          angle: 60,
          spread: 80,
          origin: { x: 0 },
          colors: colors,
        });

        confetti({
          particleCount: 100,
          angle: 120,
          spread: 80,
          origin: { x: 1 },
          colors: colors,
        });
      }, 500);
    }
  }
);
</script>

<style scoped>
.morpion-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 40px;
  gap: 20px;
  font-family: "Arial", sans-serif;
}

.game-info {
  text-align: center;
  margin-bottom: 20px;
  background-color: #ffde59;
  padding: 20px;
  border: 4px solid #000;
  box-shadow: 8px 8px 0 #000;
  transform: rotate(-1deg);
  width: 100%;
  max-width: 400px;
}

h1 {
  font-size: 2.5rem;
  text-transform: uppercase;
  letter-spacing: -1px;
  margin: 0 0 15px;
  transform: skew(-2deg);
}

.current-player,
.winner-info,
.draw-info {
  font-size: 24px;
  margin: 15px 0;
  font-weight: bold;
  text-transform: uppercase;
}

.player-x {
  color: #e74c3c;
  font-weight: bold;
  background-color: rgba(231, 76, 60, 0.2);
  padding: 2px 8px;
  border: 2px solid #e74c3c;
}

.player-o {
  color: #3498db;
  font-weight: bold;
  background-color: rgba(52, 152, 219, 0.2);
  padding: 2px 8px;
  border: 2px solid #3498db;
}

.board {
  display: flex;
  flex-direction: column;
  background-color: #f8f9fa;
  border: 6px solid #000;
  box-shadow: 12px 12px 0 #000;
  transform: rotate(1deg);
}

.board-row {
  display: flex;
}

.reset-button {
  margin-top: 20px;
  padding: 12px 25px;
  font-size: 18px;
  font-weight: bold;
  text-transform: uppercase;
  background-color: #ff6b6b;
  color: white;
  border: 4px solid #000;
  box-shadow: 6px 6px 0 #000;
  cursor: pointer;
  transition: all 0.2s;
  transform: skew(-3deg);
}

.reset-button:hover {
  background-color: #ff8e8e;
  transform: skew(-3deg) translate(-2px, -2px);
  box-shadow: 8px 8px 0 #000;
}

.reset-button:active {
  transform: skew(-3deg) translate(2px, 2px);
  box-shadow: 2px 2px 0 #000;
}
</style>
