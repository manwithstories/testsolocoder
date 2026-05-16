import { CONFIG } from '../config/constants.js';
import { getRandomPosition, isPositionOccupied } from '../utils/helpers.js';

export class Food {
    constructor() {
        this.position = { x: 0, y: 0 };
    }

    spawn(snakeBodyPositions) {
        let newPosition;
        let attempts = 0;
        const maxAttempts = 100;

        do {
            newPosition = getRandomPosition(
                CONFIG.GRID_SIZE,
                CONFIG.CANVAS_WIDTH,
                CONFIG.CANVAS_HEIGHT
            );
            attempts++;
        } while (
            isPositionOccupied(newPosition, snakeBodyPositions) &&
            attempts < maxAttempts
        );

        this.position = newPosition;
    }

    getPosition() {
        return { ...this.position };
    }

    isEaten(snakeHead) {
        return snakeHead.x === this.position.x && snakeHead.y === this.position.y;
    }
}
