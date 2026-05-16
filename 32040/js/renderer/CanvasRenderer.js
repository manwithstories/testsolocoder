import { CONFIG } from '../config/constants.js';

export class CanvasRenderer {
    constructor(canvasId) {
        this.canvas = document.getElementById(canvasId);
        this.ctx = this.canvas.getContext('2d');
        this.gridSize = CONFIG.GRID_SIZE;
    }

    clear() {
        this.ctx.fillStyle = CONFIG.COLORS.BACKGROUND;
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
    }

    drawGrid() {
        this.ctx.strokeStyle = CONFIG.COLORS.GRID;
        this.ctx.lineWidth = 0.5;

        for (let x = 0; x <= this.canvas.width; x += this.gridSize) {
            this.ctx.beginPath();
            this.ctx.moveTo(x, 0);
            this.ctx.lineTo(x, this.canvas.height);
            this.ctx.stroke();
        }

        for (let y = 0; y <= this.canvas.height; y += this.gridSize) {
            this.ctx.beginPath();
            this.ctx.moveTo(0, y);
            this.ctx.lineTo(this.canvas.width, y);
            this.ctx.stroke();
        }
    }

    drawSnake(snakeBody) {
        snakeBody.forEach((segment, index) => {
            const x = segment.x * this.gridSize;
            const y = segment.y * this.gridSize;
            const isHead = index === 0;

            this.ctx.fillStyle = isHead ? CONFIG.COLORS.SNAKE_HEAD : CONFIG.COLORS.SNAKE_BODY;
            this.ctx.beginPath();
            this.ctx.roundRect(x + 1, y + 1, this.gridSize - 2, this.gridSize - 2, 4);
            this.ctx.fill();

            if (isHead) {
                this.drawEyes(x, y);
            }
        });
    }

    drawEyes(headX, headY) {
        const eyeSize = 3;
        const offset = 5;

        this.ctx.fillStyle = '#fff';
        this.ctx.beginPath();
        this.ctx.arc(headX + offset, headY + offset, eyeSize, 0, Math.PI * 2);
        this.ctx.arc(headX + this.gridSize - offset, headY + offset, eyeSize, 0, Math.PI * 2);
        this.ctx.fill();

        this.ctx.fillStyle = '#000';
        this.ctx.beginPath();
        this.ctx.arc(headX + offset, headY + offset, eyeSize / 2, 0, Math.PI * 2);
        this.ctx.arc(headX + this.gridSize - offset, headY + offset, eyeSize / 2, 0, Math.PI * 2);
        this.ctx.fill();
    }

    drawFood(foodPosition) {
        const x = foodPosition.x * this.gridSize;
        const y = foodPosition.y * this.gridSize;
        const centerX = x + this.gridSize / 2;
        const centerY = y + this.gridSize / 2;
        const radius = this.gridSize / 2 - 2;

        const gradient = this.ctx.createRadialGradient(centerX - 2, centerY - 2, 0, centerX, centerY, radius);
        gradient.addColorStop(0, '#fca5a5');
        gradient.addColorStop(1, CONFIG.COLORS.FOOD);

        this.ctx.fillStyle = gradient;
        this.ctx.beginPath();
        this.ctx.arc(centerX, centerY, radius, 0, Math.PI * 2);
        this.ctx.fill();

        this.ctx.fillStyle = 'rgba(255, 255, 255, 0.3)';
        this.ctx.beginPath();
        this.ctx.arc(centerX - 3, centerY - 3, radius / 3, 0, Math.PI * 2);
        this.ctx.fill();
    }

    render(snakeBody, foodPosition) {
        this.clear();
        this.drawGrid();
        this.drawSnake(snakeBody);
        this.drawFood(foodPosition);
    }
}
