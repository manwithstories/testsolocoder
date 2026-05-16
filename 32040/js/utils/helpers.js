export function randomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

export function getRandomPosition(gridSize, canvasWidth, canvasHeight) {
    const maxX = Math.floor(canvasWidth / gridSize) - 1;
    const maxY = Math.floor(canvasHeight / gridSize) - 1;
    return {
        x: randomInt(0, maxX),
        y: randomInt(0, maxY)
    };
}

export function isPositionOccupied(position, occupiedPositions) {
    return occupiedPositions.some(
        pos => pos.x === position.x && pos.y === position.y
    );
}

export function getOppositeDirection(direction) {
    const opposites = {
        'UP': 'DOWN',
        'DOWN': 'UP',
        'LEFT': 'RIGHT',
        'RIGHT': 'LEFT'
    };
    return opposites[direction] || direction;
}

export function saveToStorage(key, value) {
    try {
        localStorage.setItem(key, JSON.stringify(value));
    } catch (e) {
        console.warn('无法保存到本地存储:', e);
    }
}

export function loadFromStorage(key, defaultValue = null) {
    try {
        const item = localStorage.getItem(key);
        return item ? JSON.parse(item) : defaultValue;
    } catch (e) {
        console.warn('无法从本地存储读取:', e);
        return defaultValue;
    }
}

export function throttle(func, limit) {
    let inThrottle;
    return function(...args) {
        if (!inThrottle) {
            func.apply(this, args);
            inThrottle = true;
            setTimeout(() => inThrottle = false, limit);
        }
    };
}
