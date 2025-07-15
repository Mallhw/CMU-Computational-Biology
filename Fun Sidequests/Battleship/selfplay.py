import random
import numpy as np
import torch
import torch.nn as nn
import torch.optim as optim
from collections import deque, namedtuple
from tqdm import tqdm

# --- Board Environment (from baseline) ---
class Board:
    def __init__(self, size=10, ship_sizes=[5,4,3,3,2]):
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

# --- Environment: two-player Battleship self-play ---
class BattleshipEnv:
    def __init__(self, board_size=10, ship_sizes=[5,4,3,3,2]):
        self.size = board_size
        self.ship_sizes = ship_sizes
        self.reset()

    def reset(self):
        # create two independent boards
        self.boards = [Board(self.size, self.ship_sizes), Board(self.size, self.ship_sizes)]
        for b in self.boards:
            b.reset()
        self.current_player = 0
        # both players see only their shot history
        return self._get_state(self.current_player)

    def _get_state(self, player):
        # encode UNKNOWN=0, MISS=1, HIT=2 in grid
        view = np.zeros((self.size, self.size), dtype=int)
        hits = self.boards[1-player].hits
        grid = self.boards[1-player].grid
        for (r, c) in hits:
            if grid[r,c] == 1:
                view[r,c] = 2
            else:
                view[r,c] = 1
        return view.flatten().astype(np.float32)

    def step(self, action):
        # action is idx 0..size*size-1
        r, c = divmod(action, self.size)
        board = self.boards[1-self.current_player]
        result = board.shoot((r,c))
        # assign reward
        if result == 'hit': reward = 1
        elif result == 'sunk': reward = 5
        elif result == 'miss': reward = -1
        else: reward = 0  # repeat, unlikely
        done = board.all_sunk()
        if done:
            # winning bonus
            reward += 20
        # switch player
        self.current_player = 1 - self.current_player
        next_state = self._get_state(self.current_player)
        return next_state, reward, done, {}

# --- Replay Buffer ---
Transition = namedtuple('Transition', ('state','action','reward','next_state','done'))
class ReplayBuffer:
    def __init__(self, capacity=100000):
        self.buffer = deque(maxlen=capacity)
    def push(self, *args):
        self.buffer.append(Transition(*args))
    def sample(self, batch_size):
        batch = random.sample(self.buffer, batch_size)
        return Transition(*zip(*batch))
    def __len__(self): return len(self.buffer)

# --- Q-Network ---
class DQN(nn.Module):
    def __init__(self, input_dim, output_dim):
        super(DQN, self).__init__()
        self.net = nn.Sequential(
            nn.Linear(input_dim, 256), nn.ReLU(),
            nn.Linear(256, 256), nn.ReLU(),
            nn.Linear(256, output_dim)
        )
    def forward(self, x): return self.net(x)

# --- Agent ---
class DQNAgent:
    def __init__(self, state_dim, action_dim, lr=1e-3, gamma=0.99, epsilon_start=1.0,
                 epsilon_final=0.01, epsilon_decay=10000):
        self.device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
        self.model = DQN(state_dim, action_dim).to(self.device)
        self.target = DQN(state_dim, action_dim).to(self.device)
        self.target.load_state_dict(self.model.state_dict())
        self.optimizer = optim.Adam(self.model.parameters(), lr=lr)
        self.gamma = gamma
        self.epsilon = epsilon_start
        self.epsilon_final = epsilon_final
        self.epsilon_decay = epsilon_decay
        self.steps_done = 0

    def select_action(self, state):
        self.steps_done += 1
        # decaying epsilon
        self.epsilon = self.epsilon_final + (self.epsilon - self.epsilon_final) * \
                       np.exp(-1. * self.steps_done / self.epsilon_decay)
        if random.random() < self.epsilon:
            return random.randrange(state.shape[0])
        else:
            with torch.no_grad():
                state_v = torch.tensor(state).unsqueeze(0).to(self.device)
                q = self.model(state_v)
                return int(q.argmax().item())

    def update(self, replay_buffer, batch_size):
        if len(replay_buffer) < batch_size:
            return
        transitions = replay_buffer.sample(batch_size)
        states = torch.tensor(np.stack(transitions.state)).to(self.device)
        actions = torch.tensor(transitions.action).unsqueeze(1).to(self.device)
        rewards = torch.tensor(transitions.reward).unsqueeze(1).to(self.device)
        next_states = torch.tensor(np.stack(transitions.next_state)).to(self.device)
        dones = torch.tensor(transitions.done).unsqueeze(1).to(self.device)

        q_values = self.model(states).gather(1, actions)
        with torch.no_grad():
            q_next = self.target(next_states).max(1)[0].unsqueeze(1)
            q_target = rewards + self.gamma * q_next * (1 - dones)

        loss = nn.MSELoss()(q_values, q_target)
        self.optimizer.zero_grad()
        loss.backward()
        self.optimizer.step()

    def sync_target(self):
        self.target.load_state_dict(self.model.state_dict())

# --- Training Loop ---
def train_selfplay(num_steps=10000, batch_size=64, target_update=1000):
    env = BattleshipEnv()
    state_dim = env.size * env.size
    action_dim = env.size * env.size
    agent = DQNAgent(state_dim, action_dim)
    buffer = ReplayBuffer()

    state = env.reset()
    total_reward = 0
    for step in range(1, num_steps+1):
        action = agent.select_action(state)
        next_state, reward, done, _ = env.step(action)
        buffer.push(state, action, reward, next_state, float(done))
        agent.update(buffer, batch_size)
        state = next_state
        total_reward += reward

        if done:
            state = env.reset()
        if step % target_update == 0:
            agent.sync_target()
            print(f"Step {step}: avg reward {total_reward/target_update:.2f}, epsilon {agent.epsilon:.3f}")
            total_reward = 0

    return agent

if __name__ == '__main__':
    trained_agent = train_selfplay()
    torch.save(trained_agent.model.state_dict(), 'dqn_battleship.pth')
