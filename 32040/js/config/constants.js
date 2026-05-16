export const CONFIG = {
    CANVAS_WIDTH: 400,
    CANVAS_HEIGHT: 400,
    GRID_SIZE: 20,
    INITIAL_SPEED: 200,
    SPEED_INCREMENT: 15,
    MIN_SPEED: 60,
    SPEED_UP_THRESHOLD: 5,
    INITIAL_SNAKE_LENGTH: 3,
    DIRECTIONS: {
        UP: { x: 0, y: -1 },
        DOWN: { x: 0, y: 1 },
        LEFT: { x: -1, y: 0 },
        RIGHT: { x: 1, y: 0 }
    },
    COLORS: {
        SNAKE_HEAD: '#4ade80',
        SNAKE_BODY: '#22c55e',
        FOOD: '#f87171',
        GRID: '#1e293b',
        BACKGROUND: '#0f172a'
    },
    STORAGE_KEY: 'snake_high_score'
};

export const KEY_CODES = {
    ARROW_UP: 'ArrowUp',
    ARROW_DOWN: 'ArrowDown',
    ARROW_LEFT: 'ArrowLeft',
    ARROW_RIGHT: 'ArrowRight',
    SPACE: ' ',
    W: 'KeyW',
    A: 'KeyA',
    S: 'KeyS',
    D: 'KeyD'
};

export const GAME_STATE = {
    MENU: 'menu',
    PLAYING: 'playing',
    PAUSED: 'paused',
    GAME_OVER: 'gameOver'
};
