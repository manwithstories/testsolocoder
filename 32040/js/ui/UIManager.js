import { CONFIG } from '../config/constants.js';
import { saveToStorage, loadFromStorage } from '../utils/helpers.js';

export class UIManager {
    constructor() {
        this.scoreElement = document.getElementById('score');
        this.highScoreElement = document.getElementById('high-score');
        this.speedElement = document.getElementById('speed');
        this.finalScoreElement = document.getElementById('final-score');
        this.startScreen = document.getElementById('start-screen');
        this.gameOverScreen = document.getElementById('game-over-screen');
        this.newRecordElement = document.getElementById('new-record');
        this.startBtn = document.getElementById('start-btn');
        this.restartBtn = document.getElementById('restart-btn');

        this.highScore = loadFromStorage(CONFIG.STORAGE_KEY, 0);
        this.updateHighScoreDisplay();
    }

    updateScore(score) {
        this.scoreElement.textContent = score;
    }

    updateSpeed(speedLevel) {
        this.speedElement.textContent = speedLevel;
    }

    updateHighScoreDisplay() {
        this.highScoreElement.textContent = this.highScore;
    }

    checkNewRecord(score) {
        if (score > this.highScore) {
            this.highScore = score;
            saveToStorage(CONFIG.STORAGE_KEY, this.highScore);
            this.updateHighScoreDisplay();
            return true;
        }
        return false;
    }

    showGameOver(score) {
        this.finalScoreElement.textContent = score;
        const isNewRecord = this.checkNewRecord(score);
        
        if (isNewRecord) {
            this.newRecordElement.classList.remove('hidden');
        } else {
            this.newRecordElement.classList.add('hidden');
        }
        
        this.gameOverScreen.classList.remove('hidden');
    }

    hideGameOver() {
        this.gameOverScreen.classList.add('hidden');
    }

    showStartScreen() {
        this.startScreen.classList.remove('hidden');
    }

    hideStartScreen() {
        this.startScreen.classList.add('hidden');
    }

    resetDisplay() {
        this.updateScore(0);
        this.updateSpeed(1);
        this.hideGameOver();
    }

    onStartClick(callback) {
        this.startBtn.addEventListener('click', callback);
    }

    onRestartClick(callback) {
        this.restartBtn.addEventListener('click', callback);
    }
}
