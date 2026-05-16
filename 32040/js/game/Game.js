import { CONFIG, KEY_CODES, GAME_STATE } from '../config/constants.js';
import { Snake } from './Snake.js';
import { Food } from './Food.js';
import { CanvasRenderer } from '../renderer/CanvasRenderer.js';
import { UIManager } from '../ui/UIManager.js';

export class Game {
    constructor() {
        this.snake = new Snake();
        this.food = new Food();
        this.renderer = new CanvasRenderer('game-canvas');
        this.ui = new UIManager();

        this.score = 0;
        this.foodEaten = 0;
        this.currentSpeed = CONFIG.INITIAL_SPEED;
        this.speedLevel = 1;
        this.gameState = GAME_STATE.MENU;
        this.gameLoopId = null;

        this.init();
    }

    init() {
        this.setupEventListeners();
        this.render();
    }

    setupEventListeners() {
        document.addEventListener('keydown', (e) => this.handleKeyDown(e));
        this.ui.onStartClick(() => this.start());
        this.ui.onRestartClick(() => this.restart());
    }

    handleKeyDown(e) {
        if (this.gameState === GAME_STATE.MENU) {
            if (e.code === 'Space' || e.code === 'Enter') {
                e.preventDefault();
                this.start();
            }
            return;
        }

        if (this.gameState === GAME_STATE.GAME_OVER) {
            if (e.code === 'Space' || e.code === 'Enter') {
                e.preventDefault();
                this.restart();
            }
            return;
        }

        if (e.code === KEY_CODES.SPACE || e.code === 'Space') {
            e.preventDefault();
            this.togglePause();
            return;
        }

        if (this.gameState !== GAME_STATE.PLAYING) return;

        const keyMap = {
            [KEY_CODES.ARROW_UP]: 'UP',
            [KEY_CODES.ARROW_DOWN]: 'DOWN',
            [KEY_CODES.ARROW_LEFT]: 'LEFT',
            [KEY_CODES.ARROW_RIGHT]: 'RIGHT',
            [KEY_CODES.W]: 'UP',
            [KEY_CODES.S]: 'DOWN',
            [KEY_CODES.A]: 'LEFT',
            [KEY_CODES.D]: 'RIGHT'
        };

        const direction = keyMap[e.code];
        if (direction) {
            e.preventDefault();
            this.snake.setDirection(direction);
        }
    }

    start() {
        this.ui.hideStartScreen();
        this.gameState = GAME_STATE.PLAYING;
        this.food.spawn(this.snake.getBodyPositions());
        this.startGameLoop();
    }

    restart() {
        this.stopGameLoop();
        this.snake.reset();
        this.score = 0;
        this.foodEaten = 0;
        this.currentSpeed = CONFIG.INITIAL_SPEED;
        this.speedLevel = 1;
        this.ui.resetDisplay();
        this.food.spawn(this.snake.getBodyPositions());
        this.gameState = GAME_STATE.PLAYING;
        this.startGameLoop();
    }

    togglePause() {
        if (this.gameState === GAME_STATE.PLAYING) {
            this.gameState = GAME_STATE.PAUSED;
            this.stopGameLoop();
        } else if (this.gameState === GAME_STATE.PAUSED) {
            this.gameState = GAME_STATE.PLAYING;
            this.startGameLoop();
        }
    }

    startGameLoop() {
        this.gameLoop();
    }

    stopGameLoop() {
        if (this.gameLoopId) {
            clearTimeout(this.gameLoopId);
            this.gameLoopId = null;
        }
    }

    gameLoop() {
        if (this.gameState !== GAME_STATE.PLAYING) return;

        this.update();
        this.render();

        this.gameLoopId = setTimeout(() => this.gameLoop(), this.currentSpeed);
    }

    update() {
        this.snake.move();

        if (this.snake.checkCollision()) {
            this.gameOver();
            return;
        }

        if (this.food.isEaten(this.snake.getHead())) {
            this.eatFood();
        }
    }

    eatFood() {
        this.snake.grow();
        this.score += 10;
        this.foodEaten++;
        this.ui.updateScore(this.score);

        if (this.foodEaten % CONFIG.SPEED_UP_THRESHOLD === 0) {
            this.increaseSpeed();
        }

        this.food.spawn(this.snake.getBodyPositions());
    }

    increaseSpeed() {
        if (this.currentSpeed > CONFIG.MIN_SPEED) {
            this.currentSpeed = Math.max(CONFIG.MIN_SPEED, this.currentSpeed - CONFIG.SPEED_INCREMENT);
            this.speedLevel++;
            this.ui.updateSpeed(this.speedLevel);
        }
    }

    gameOver() {
        this.gameState = GAME_STATE.GAME_OVER;
        this.stopGameLoop();
        this.ui.showGameOver(this.score);
    }

    render() {
        this.renderer.render(
            this.snake.getBodyPositions(),
            this.food.getPosition()
        );
    }
}
