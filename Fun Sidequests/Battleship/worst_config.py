from typing import Dict, Any, Tuple, List
from .base_placement_ai import PlacementAI

class WorstCasePlacementAI(PlacementAI):
    def generate_placement(self) -> List[Dict[str, Any]]:
        return [
            {'row': 6, 'col': 0, 'length': 5, 'direction': 'horizontal'},
            {'row': 0, 'col': 9, 'length': 4, 'direction': 'vertical'},
            {'row': 0, 'col': 0, 'length': 3, 'direction': 'horizontal'},
            {'row': 2, 'col': 5, 'length': 2, 'direction': 'vertical'},
            {'row': 6, 'col': 7, 'length': 2, 'direction': 'horizontal'},
        ]
