import numpy as np
from tqdm import tqdm
import matplotlib.pyplot as plt

# Constants
BOARD_SIZE = 10
SHIP_SIZES = [5, 4, 3, 3, 2]  # Carrier, Battleship, Cruiser, Submarine, Destroyer

class Board:
    def __init__(self, size=BOARD_SIZE, ship_sizes=SHIP_SIZES):
        self.size = size
        self.ship_sizes = ship_sizes
        self.reset()

    def reset(self):
        self.grid = np.zeros((self.size, self.size), dtype=int)
        self._place_ships()
        self.hits = set()

    def _place_ships(self):
        self.ship_coords = []
        for length in self.ship_sizes:
            placed = False
            while not placed:
                horizontal = np.random.rand() < 0.5
                if horizontal:
                    row = np.random.randint(0, self.size)
                    col = np.random.randint(0, self.size - length + 1)
                    coords = [(row, c) for c in range(col, col + length)]
                else:
                    row = np.random.randint(0, self.size - length + 1)
                    col = np.random.randint(0, self.size)
                    coords = [(r, col) for r in range(row, row + length)]
                if all(self.grid[r, c] == 0 for r, c in coords):
                    for r, c in coords:
                        self.grid[r, c] = 1
                    self.ship_coords.append(set(coords))
                    placed = True

    def shoot(self, pos):
        r, c = pos
        if pos in self.hits:
            return 'repeat'
        self.hits.add(pos)
        if self.grid[r, c] == 1:
            for ship in self.ship_coords:
                if pos in ship:
                    ship.remove(pos)
                    if not ship:
                        return 'sunk'
                    return 'hit'
        return 'miss'

    def all_sunk(self):
        return all(len(ship) == 0 for ship in self.ship_coords)
    

class BaseAgent:
    def __init__(self):
        self.reset()

    def reset(self):
        self.shots = []  # record shots

    def next_shot(self):
        raise NotImplementedError

    def update(self, pos, result):
        pass

class RandomAgent(BaseAgent):
    def next_shot(self):
        candidates = [(r, c) for r in range(BOARD_SIZE) for c in range(BOARD_SIZE) if (r, c) not in self.shots]
        choice = candidates[np.random.randint(len(candidates))]
        self.shots.append(choice)
        return choice

# TODO: Implement GridAgent, PDFAgent, GPAAgent, MCTSAgent, NNAgent skeletons

# Simulation functions
def play_game(agent, board=None):
    if board is None:
        board = Board()
    board.reset()
    agent.reset()
    turns = 0
    while not board.all_sunk():
        pos = agent.next_shot()
        result = board.shoot(pos)
        agent.update(pos, result)
        turns += 1
    return turns



class GridAgent(BaseAgent):
    def reset(self):
        super().reset()
        self.mode = 'hunt'
        # Only shoot on cells where (row+col) is odd
        self.hunt_cells = [(r, c) for r in range(BOARD_SIZE) for c in range(BOARD_SIZE) if (r + c) % 2 == 1]
        np.random.shuffle(self.hunt_cells)
        self.origin = None
        self.direction = None
        self.tried_dirs = []
        self.next_target = None

    def next_shot(self):
        # HUNT MODE
        if self.mode == 'hunt':
            if self.hunt_cells:
                shot = self.hunt_cells.pop()
            else:
                # Fallback: shoot any remaining cell
                candidates = [(r, c) for r in range(BOARD_SIZE) for c in range(BOARD_SIZE) if (r, c) not in self.shots]
                shot = candidates[np.random.randint(len(candidates))]
            self.shots.append(shot)
            return shot

        # TARGET MODE
        # Try directions around origin until ship is sunk
        while self.mode == 'target':
            # pick new direction if needed
            if self.next_target is None:
                # if all dirs tried, abort back to hunt
                all_dirs = [(1, 0), (-1, 0), (0, 1), (0, -1)]
                if set(self.tried_dirs) >= set(all_dirs):
                    # no more directions, reset
                    self.mode = 'hunt'
                    return self.next_shot()
                # select next untried dir
                for dx, dy in all_dirs:
                    if (dx, dy) not in self.tried_dirs:
                        self.direction = (dx, dy)
                        self.next_target = (self.origin[0] + dx, self.origin[1] + dy)
                        break
            r, c = self.next_target
            # check valid
            if not (0 <= r < BOARD_SIZE and 0 <= c < BOARD_SIZE) or self.next_target in self.shots:
                # mark miss on this direction and pick next
                self.tried_dirs.append(self.direction)
                self.direction = None
                self.next_target = None
                continue
            # valid shot
            shot = self.next_target
            self.shots.append(shot)
            return shot

    def update(self, pos, result):
        if self.mode == 'hunt':
            if result == 'hit':
                self.mode = 'target'
                self.origin = pos
                self.tried_dirs = []
                self.direction = None
                self.next_target = None
            return

        # TARGET MODE UPDATE
        if result == 'hit':
            dx, dy = self.direction
            # continue in same direction
            self.next_target = (pos[0] + dx, pos[1] + dy)
        elif result == 'miss':
            # mark this dir tried and reset next_target
            self.tried_dirs.append(self.direction)
            self.direction = None
            self.next_target = None
        elif result == 'sunk':
            # ship down, go back to hunt
            self.mode = 'hunt'
            self.origin = None
            self.direction = None
            self.tried_dirs = []
            self.next_target = None


# TODO: Implement PDFAgent, GPAAgent, MCTSAgent, NNAgent skeletons

# Simulation functions
def play_game(agent, board=None):
    if board is None:
        board = Board()
    board.reset()
    agent.reset()
    turns = 0
    while not board.all_sunk():
        pos = agent.next_shot()
        result = board.shoot(pos)
        agent.update(pos, result)
        turns += 1
    return turns

if __name__ == '__main__':
    MODELS = [
        ('Random', RandomAgent()),
        ('Grid', GridAgent()),
        # ('PDF', PDFAgent()),
        # ('GP', GPAAgent()),
        # ('MCTS', MCTSAgent()),
        # ('NN', NNAgent()),
    ]
    results = {}
    total_turns = []
    for name, agent in MODELS:
        turns_list = []
        for _ in tqdm(range(1000), desc=f"Simulating {name}"):
            turns = play_game(agent)
            turns_list.append(turns)
        results[name] = np.mean(turns_list)
        total_turns.append(turns_list)

    # Plotting
    names = list(results.keys())
    values = [results[n] for n in names]
    plt.figure()
    plt.bar(names, values)
    plt.ylabel('Average Turns to Sink All')
    plt.title('Battleship Agent Performance')
    plt.show()


# Extra Plots
plt.plot(np.arange(1000), total_turns[0], total_turns[1])
plt.show() 

# Plotting helpers
def plot_board(grid, shot_results, title, show_ships):
    ax = plt.gca()
    ax.clear()
    ax.set_xticks(np.arange(-.5, BOARD_SIZE, 1))
    ax.set_yticks(np.arange(-.5, BOARD_SIZE, 1))
    ax.set_xticklabels([])
    ax.set_yticklabels([])
    ax.grid(True)
    if show_ships:
        for (r, c), val in np.ndenumerate(grid):
            if val == 1:
                ax.add_patch(plt.Rectangle((c - .5, r - .5), 1, 1, fill=True,
                                            edgecolor='black', facecolor='lightgray'))
    for (r, c), res in shot_results.items():
        color = 'red' if res in ('hit', 'sunk') else 'blue'
        ax.text(c, r, 'X', ha='center', va='center', color=color, fontsize=12)
    ax.set_title(title)
    plt.gca().invert_yaxis()
    plt.draw()
    plt.pause(0.001)

# Simulation and animation

def simulate_and_plot(agent, show_ships=False):
    board = Board()
    board.reset()
    agent.reset()
    shot_results = {}
    while not board.all_sunk():
        pos = agent.next_shot()
        result = board.shoot(pos)
        agent.update(pos, result)
        shot_results[pos] = result
    plt.figure(figsize=(6,6))
    plot_board(board.grid, shot_results,
               f"{agent.__class__.__name__} Simulation ({len(shot_results)} shots)",
               show_ships)
    plt.show()


def simulate_with_steps(agent, show_ships=False, pause=0.5):
    board = Board()
    board.reset()
    agent.reset()
    shot_results = {}
    plt.ion()
    fig = plt.figure(figsize=(6,6))
    for turn in range(BOARD_SIZE * BOARD_SIZE):
        if board.all_sunk():
            break
        pos = agent.next_shot()
        result = board.shoot(pos)
        agent.update(pos, result)
        shot_results[pos] = result
        plot_board(board.grid, shot_results,
                   f"{agent.__class__.__name__} Turn {turn+1}",
                   show_ships)
        plt.pause(pause)
    plt.ioff()
    plt.show()


simulate_with_steps(GridAgent(), show_ships=True, pause=0.1)