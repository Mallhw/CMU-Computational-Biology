from abc import ABC, abstractmethod
from typing import Tuple, Dict, Any, List
import numpy as np
from enum import Enum
from base_ai import BattleshipAI


class WellState(Enum):
    UNKNOWN = 0
    MISS = 1
    HIT = 2

class HeatmapBattleshipAI(BattleshipAI):
    def __init__(self, player_id: str, board_shape: Tuple[int, int], ship_schema: Dict[str, Any]):
        super().__init__(player_id, board_shape, ship_schema)
        self.rows, self.cols = board_shape

    def select_next_move(self) -> Tuple[int, int]:
        board_with_hits = (self.board_state == WellState.HIT).astype(int)
        board_with_misses = (self.board_state == WellState.MISS).astype(int) * 2

        prob_matrix = self.generate_probabilities_for_all_ships(board_with_hits, board_with_misses)
        move = np.unravel_index(np.argmax(prob_matrix), prob_matrix.shape)
        return move

    def generate_probabilities_for_all_ships(self, board_with_hits, board_with_misses):
        final = np.zeros((self.rows, self.cols))
        ship_lengths = []
        for ship, data in self.ship_schema.items():
            ship_lengths.extend([data['length']] * data['count'])
        for length in ship_lengths:
            final += self.possible_locations_probability(board_with_hits, board_with_misses, length)
        return final

    def possible_locations_probability(self, board_with_hits, board_with_misses, ship_length):
        list_of_probabilities = []

        # Horizontal
        for row in range(self.rows):
            for col in range(self.cols - ship_length + 1):
                segment = range(col, col + ship_length)
                positions_with_hits = []
                empty_slots = 0
                for c in segment:
                    if board_with_misses[row, c] == 0:
                        if board_with_hits[row, c] == 1:
                            positions_with_hits.append(c)
                        empty_slots += 1
                if empty_slots == ship_length:
                    prob_matrix = np.zeros((self.rows, self.cols))
                    multiplier = 4 * len(positions_with_hits) if positions_with_hits else 1
                    for c in segment:
                        if c not in positions_with_hits:
                            prob_matrix[row, c] = ship_length * multiplier
                    list_of_probabilities.append(prob_matrix)

        # Vertical
        for col in range(self.cols):
            for row in range(self.rows - ship_length + 1):
                segment = range(row, row + ship_length)
                positions_with_hits = []
                empty_slots = 0
                for r in segment:
                    if board_with_misses[r, col] == 0:
                        if board_with_hits[r, col] == 1:
                            positions_with_hits.append(r)
                        empty_slots += 1
                if empty_slots == ship_length:
                    prob_matrix = np.zeros((self.rows, self.cols))
                    multiplier = 4 * len(positions_with_hits) if positions_with_hits else 1
                    for r in segment:
                        if r not in positions_with_hits:
                            prob_matrix[r, col] = ship_length * multiplier
                    list_of_probabilities.append(prob_matrix)

        final_matrix = np.zeros((self.rows, self.cols))
        for m in list_of_probabilities:
            final_matrix += m
        return final_matrix