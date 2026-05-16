import pygame
import random
from config import *


class Enemy:
    def __init__(self, wave=1):
        self.width = ENEMY_WIDTH
        self.height = ENEMY_HEIGHT
        self.x = random.randint(0, SCREEN_WIDTH - self.width)
        self.y = -self.height
        self.speed = ENEMY_BASE_SPEED + (wave - 1) * WAVE_SPEED_MULTIPLIER
        self.hp = ENEMY_BASE_HP + (wave - 1) // 3
        self.max_hp = self.hp
        self.points = ENEMY_POINTS * wave
        self.active = True
        self.wave = wave

        self.move_pattern = random.choice(['straight', 'zigzag', 'sine'])
        self.zigzag_direction = random.choice([-1, 1])
        self.sine_offset = random.uniform(0, 3.14)
        self.sine_timer = 0

    def update(self):
        self.y += self.speed

        if self.move_pattern == 'zigzag':
            self.x += self.zigzag_direction * 2
            if self.x <= 0 or self.x >= SCREEN_WIDTH - self.width:
                self.zigzag_direction *= -1
        elif self.move_pattern == 'sine':
            self.sine_timer += 0.05
            self.x += int(5 * (self.sine_timer + self.sine_offset))

        self.x = max(0, min(SCREEN_WIDTH - self.width, self.x))

        if self.y > SCREEN_HEIGHT:
            self.active = False

    def draw(self, screen):
        enemy_points = [
            (self.x + self.width // 2, self.y + self.height),
            (self.x, self.y),
            (self.x + self.width // 4, self.y + self.height * 0.3),
            (self.x + self.width // 2, self.y + self.height * 0.15),
            (self.x + self.width * 3 // 4, self.y + self.height * 0.3),
            (self.x + self.width, self.y)
        ]
        pygame.draw.polygon(screen, RED, enemy_points)
        pygame.draw.polygon(screen, WHITE, enemy_points, 2)

        if self.max_hp > 1:
            bar_width = self.width
            bar_height = 4
            bar_x = self.x
            bar_y = self.y - 8
            hp_ratio = self.hp / self.max_hp
            pygame.draw.rect(screen, DARK_GRAY, (bar_x, bar_y, bar_width, bar_height))
            pygame.draw.rect(screen, GREEN, (bar_x, bar_y, int(bar_width * hp_ratio), bar_height))

    def take_damage(self, damage):
        self.hp -= damage
        if self.hp <= 0:
            self.active = False
            return True
        return False

    def get_rect(self):
        return pygame.Rect(self.x, self.y, self.width, self.height)


class Boss:
    def __init__(self, wave=5):
        self.width = BOSS_WIDTH
        self.height = BOSS_HEIGHT
        self.x = (SCREEN_WIDTH - self.width) // 2
        self.y = -self.height
        self.target_y = 60
        self.speed = BOSS_BASE_SPEED
        self.hp = BOSS_BASE_HP + (wave // BOSS_SPAWN_WAVE - 1) * 25
        self.max_hp = self.hp
        self.points = BOSS_POINTS * (wave // BOSS_SPAWN_WAVE)
        self.active = True
        self.wave = wave

        self.move_direction = 1
        self.last_shot_time = 0
        self.attack_pattern = 0
        self.attack_timer = 0
        self.entering = True

    def update(self, bullet_manager):
        if self.entering:
            self.y += 2
            if self.y >= self.target_y:
                self.y = self.target_y
                self.entering = False
            return

        self.x += self.speed * self.move_direction
        if self.x <= 0 or self.x >= SCREEN_WIDTH - self.width:
            self.move_direction *= -1

        current_time = pygame.time.get_ticks()
        self.attack_timer += 1

        if self.attack_timer % 300 == 0:
            self.attack_pattern = (self.attack_pattern + 1) % 3

        if current_time - self.last_shot_time >= BOSS_FIRE_RATE:
            self.shoot(bullet_manager)
            self.last_shot_time = current_time

    def shoot(self, bullet_manager):
        if self.entering:
            return

        center_x = self.x + self.width // 2 - BULLET_WIDTH // 2
        bottom_y = self.y + self.height

        if self.attack_pattern == 0:
            for i in range(-2, 3):
                bullet_manager.add_enemy_bullet(
                    center_x + i * 25, bottom_y, speed=5
                )
        elif self.attack_pattern == 1:
            for i in range(4):
                offset_x = i * 30 - 45
                bullet_manager.add_enemy_bullet(
                    center_x + offset_x, bottom_y + 10, speed=6
                )
        else:
            bullet_manager.add_enemy_bullet(center_x, bottom_y, speed=7)
            bullet_manager.add_enemy_bullet(center_x - 30, bottom_y + 5, speed=6)
            bullet_manager.add_enemy_bullet(center_x + 30, bottom_y + 5, speed=6)

    def draw(self, screen):
        boss_color = (180, 0, 180)

        main_body = [
            (self.x + self.width // 2, self.y),
            (self.x, self.y + self.height * 0.4),
            (self.x + self.width * 0.1, self.y + self.height),
            (self.x + self.width * 0.9, self.y + self.height),
            (self.x + self.width, self.y + self.height * 0.4)
        ]
        pygame.draw.polygon(screen, boss_color, main_body)
        pygame.draw.polygon(screen, WHITE, main_body, 3)

        pygame.draw.circle(
            screen,
            RED,
            (self.x + self.width // 2, self.y + self.height * 0.4),
            15
        )
        pygame.draw.circle(
            screen,
            YELLOW,
            (self.x + self.width // 2, self.y + self.height * 0.4),
            8
        )

        wing_left = [
            (self.x, self.y + self.height * 0.5),
            (self.x - 20, self.y + self.height * 0.8),
            (self.x + self.width * 0.15, self.y + self.height * 0.9)
        ]
        wing_right = [
            (self.x + self.width, self.y + self.height * 0.5),
            (self.x + self.width + 20, self.y + self.height * 0.8),
            (self.x + self.width * 0.85, self.y + self.height * 0.9)
        ]
        pygame.draw.polygon(screen, boss_color, wing_left)
        pygame.draw.polygon(screen, WHITE, wing_left, 2)
        pygame.draw.polygon(screen, boss_color, wing_right)
        pygame.draw.polygon(screen, WHITE, wing_right, 2)

    def take_damage(self, damage):
        self.hp -= damage
        if self.hp <= 0:
            self.active = False
            return True
        return False

    def get_rect(self):
        return pygame.Rect(self.x, self.y, self.width, self.height)


class EnemyManager:
    def __init__(self):
        self.enemies = []
        self.boss = None
        self.last_spawn_time = 0
        self.spawned_count = 0
        self.total_enemies_per_wave = 5
        self.current_wave = 1

    def start_wave(self, wave):
        self.current_wave = wave
        self.enemies.clear()
        self.boss = None
        self.spawned_count = 0
        self.last_spawn_time = 0
        self.total_enemies_per_wave = 5 + (wave - 1) * WAVE_ENEMY_INCREMENT

        if wave % BOSS_SPAWN_WAVE == 0:
            self.boss = Boss(wave)
            self.total_enemies_per_wave = 0

    def update(self, bullet_manager):
        current_time = pygame.time.get_ticks()

        if self.boss:
            self.boss.update(bullet_manager)
            if not self.boss.active:
                self.boss = None
            return

        spawn_rate = max(500, ENEMY_SPAWN_RATE - (self.current_wave - 1) * 50)
        if (self.spawned_count < self.total_enemies_per_wave and
                current_time - self.last_spawn_time >= spawn_rate):
            enemy = Enemy(self.current_wave)
            self.enemies.append(enemy)
            self.spawned_count += 1
            self.last_spawn_time = current_time

        for enemy in self.enemies:
            enemy.update()

        self.enemies = [e for e in self.enemies if e.active]

    def draw(self, screen):
        for enemy in self.enemies:
            enemy.draw(screen)
        if self.boss:
            self.boss.draw(screen)

    def check_wave_complete(self):
        if self.boss:
            return not self.boss.active
        return (self.spawned_count >= self.total_enemies_per_wave and
                len(self.enemies) == 0)

    def get_all_enemies(self):
        all_enemies = list(self.enemies)
        if self.boss:
            all_enemies.append(self.boss)
        return all_enemies

    def clear_all(self):
        self.enemies.clear()
        self.boss = None
