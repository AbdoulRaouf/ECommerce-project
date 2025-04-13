import { defineStore } from "pinia";
import { computed, ref } from "vue";

export type Player = "X" | "O" | null;
export type Board = Player[][];

export const useMorpionStore = defineStore("morpion", () => {
  // État du jeu
  const board = ref<Board>([
    [null, null, null],
    [null, null, null],
    [null, null, null],
  ]);

  const currentPlayer = ref<"X" | "O">("X");
  const winner = ref<Player>(null);
  const isGameOver = ref<boolean>(false);
  const winningCells = ref<number[][]>([]);

  // Getters
  const isBoardFull = computed(() => {
    return board.value.every((row) => row.every((cell) => cell !== null));
  });

  // Actions
  function makeMove(row: number, col: number): boolean {
    // Vérifier si la cellule est déjà occupée ou si le jeu est terminé
    if (board.value[row][col] !== null || isGameOver.value) {
      return false;
    }

    // Placer le symbole du joueur actuel
    board.value[row][col] = currentPlayer.value;

    // Vérifier s'il y a un gagnant
    checkWinner();

    // Si le jeu n'est pas terminé, passer au joueur suivant
    if (!isGameOver.value) {
      currentPlayer.value = currentPlayer.value === "X" ? "O" : "X";
    }

    return true;
  }

  function checkWinner(): void {
    const winPatterns = [
      // Lignes horizontales
      [
        [0, 0],
        [0, 1],
        [0, 2],
      ],
      [
        [1, 0],
        [1, 1],
        [1, 2],
      ],
      [
        [2, 0],
        [2, 1],
        [2, 2],
      ],
      // Lignes verticales
      [
        [0, 0],
        [1, 0],
        [2, 0],
      ],
      [
        [0, 1],
        [1, 1],
        [2, 1],
      ],
      [
        [0, 2],
        [1, 2],
        [2, 2],
      ],
      // Diagonales
      [
        [0, 0],
        [1, 1],
        [2, 2],
      ],
      [
        [0, 2],
        [1, 1],
        [2, 0],
      ],
    ];

    for (const pattern of winPatterns) {
      const [a, b, c] = pattern;
      const [rowA, colA] = a;
      const [rowB, colB] = b;
      const [rowC, colC] = c;

      if (
        board.value[rowA][colA] !== null &&
        board.value[rowA][colA] === board.value[rowB][colB] &&
        board.value[rowA][colA] === board.value[rowC][colC]
      ) {
        winner.value = board.value[rowA][colA];
        isGameOver.value = true;
        winningCells.value = pattern;
        return;
      }
    }

    // Vérifier si le jeu est nul (match nul)
    if (isBoardFull.value) {
      isGameOver.value = true;
    }
  }

  function resetGame(): void {
    board.value = [
      [null, null, null],
      [null, null, null],
      [null, null, null],
    ];
    currentPlayer.value = "X";
    winner.value = null;
    isGameOver.value = false;
    winningCells.value = [];
  }

  return {
    board,
    currentPlayer,
    winner,
    isGameOver,
    winningCells,
    isBoardFull,
    makeMove,
    resetGame,
  };
});
