import { CONFIG } from '../config/constants.js';
import { getOppositeDirection, isPositionOccupied } from '../utils/helpers.js';

export class Snake {
    constructor() {
        this.reset();
    }

    reset() {
        const startX = Math.floor(CONFIG.CANVAS_WIDTH / CONFIG.GRID_SIZE / 2);
        const startY = Math.floor(CONFIG.CANVAS_HEIGHT / CONFIG.GRID_SIZE / 2);

        this.body = [];
        for (let i = 0; i < CONFIG.INITIAL_SNAKE_LENGTH; i++) {
            this.body.push({ x: startX - i, y: startY });
        }

        this.currentDirection = 'RIGHT';
        this.nextDirection = 'RIGHT';
        this.growPending = false;
    }

    setDirection(newDirection) {
        if (getOppositeDirection(this.currentDirection) !== newDirection) {
            this.nextDirection = newDirection;
        }
    }

    move() {
        this.currentDirection = this.nextDirection;
        const direction = CONFIG.DIRECTIONS[this.currentDirection];
        const head = this.body[0];
        const newHead = {
            x: head.x + direction.x,
            y: head.y + direction.y
        };

        this.body.unshift(newHead);

        if (this.growPending) {
            this.growPending = false;
        } else {
            this.body.pop();
        }
    }

    grow() {
        this.growPending = true;
    }

    getHead() {
        return this.body[0];
    }

    checkWallCollision() {
        const head = this.getHead();
        const maxX = Math.floor(CONFIG.CANVAS_WIDTH / CONFIG.GRID_SIZE) - 1;
        const maxY = Math.floor(CONFIG.CANVAS_HEIGHT / CONFIG.GRID_SIZE) - 1;

        return head.x < 0 || head.x > maxX || head.y < 0 || head.y > maxY;
    }

    checkSelfCollision() {
        const head = this.getHead();
        const bodyWithoutHead = this.body.slice(1);
        return isPositionOccupied(head, bodyWithoutHead);
    }

    checkCollision() {
        return this.checkWallCollision() || this.checkSelfCollision();
    }

    getBodyPositions() {
        return [...this.body];
    }

    getLength() {
        return this.body.length;
    }
}
