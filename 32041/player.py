import pygame
import time
from config import *


class Player:
    def __init__(self):
        self.width = PLAYER_WIDTH
        self.height = PLAYER_HEIGHT
        self.x = (SCREEN_WIDTH - self.width) // 2
        self.y = SCREEN_HEIGHT - self.height - 20
        self.speed = PLAYER_SPEED
        self.lives = PLAYER_MAX_LIVES
        self.max_lives = PLAYER_MAX_LIVES
        self.bullet_damage = BULLET_DAMAGE
        self.fire_rate = FIRE_RATE

        self.last_shot_time = 0
        self.invincible = False
        self.invincible_end_time = 0
        self.shield_active = False
        self.shield_end_time = 0
        self.rapid_fire = False
        self.rapid_fire_end_time = 0
        self.multi_shot = False
        self.multi_shot_end_time = 0

        self.move_left = False
        self.move_right = False
        self.shooting = False

        self.damage_level = 1
        self.speed_level = 1
        self.fire_rate_level = 1

    def handle_event(self, event):
        if event.type == pygame.KEYDOWN:
            if event.key == pygame.K_LEFT or event.key == pygame.K_a:
                self.move_left = True
            if event.key == pygame.K_RIGHT or event.key == pygame.K_d:
                self.move_right = True
            if event.key == pygame.K_SPACE:
                self.shooting = True
        elif event.type == pygame.KEYUP:
            if event.key == pygame.K_LEFT or event.key == pygame.K_a:
                self.move_left = False
            if event.key == pygame.K_RIGHT or event.key == pygame.K_d:
                self.move_right = False
            if event.key == pygame.K_SPACE:
                self.shooting = False

    def update(self, bullet_manager):
        current_time = pygame.time.get_ticks()

        if self.move_left:
            self.x -= self.speed
        if self.move_right:
            self.x += self.speed

        self.x = max(0, min(SCREEN_WIDTH - self.width, self.x))

        if self.invincible and current_time >= self.invincible_end_time:
            self.invincible = False

        if self.shield_active and current_time >= self.shield_end_time:
            self.shield_active = False

        if self.rapid_fire and current_time >= self.rapid_fire_end_time:
            self.rapid_fire = False

        if self.multi_shot and current_time >= self.multi_shot_end_time:
            self.multi_shot = False

        if self.shooting:
            current_fire_rate = self.fire_rate // 2 if self.rapid_fire else self.fire_rate
            if current_time - self.last_shot_time >= current_fire_rate:
                self.shoot(bullet_manager)
                self.last_shot_time = current_time

    def shoot(self, bullet_manager):
        bullet_x = self.x + self.width // 2 - BULLET_WIDTH // 2
        bullet_y = self.y

        if self.multi_shot:
            bullet_manager.add_player_bullet(bullet_x - 15, bullet_y, self.bullet_damage)
            bullet_manager.add_player_bullet(bullet_x, bullet_y, self.bullet_damage)
            bullet_manager.add_player_bullet(bullet_x + 15, bullet_y, self.bullet_damage)
        else:
            bullet_manager.add_player_bullet(bullet_x, bullet_y, self.bullet_damage)

    def draw(self, screen):
        if self.invincible:
            current_time = pygame.time.get_ticks()
            if (current_time // 100) % 2 == 0:
                return

        ship_points = [
            (self.x + self.width // 2, self.y),
            (self.x, self.y + self.height),
            (self.x + self.width // 4, self.y + self.height * 0.7),
            (self.x + self.width // 2, self.y + self.height * 0.85),
            (self.x + self.width * 3 // 4, self.y + self.height * 0.7),
            (self.x + self.width, self.y + self.height)
        ]
        pygame.draw.polygon(screen, CYAN, ship_points)
        pygame.draw.polygon(screen, WHITE, ship_points, 2)

        if self.shield_active:
            shield_surface = pygame.Surface((self.width + 20, self.height + 20), pygame.SRCALPHA)
            pygame.draw.circle(
                shield_surface,
                (0, 255, 255, 80),
                ((self.width + 20) // 2, (self.height + 20) // 2),
                (self.width + 20) // 2
            )
            screen.blit(shield_surface, (self.x - 10, self.y - 10))
            pygame.draw.circle(
                screen,
                CYAN,
                (self.x + self.width // 2, self.y + self.height // 2),
                self.width // 2 + 10,
                2
            )

    def get_rect(self):
        return pygame.Rect(self.x, self.y, self.width, self.height)

    def take_damage(self):
        if self.invincible or self.shield_active:
            return False

        self.lives -= 1
        self.invincible = True
        self.invincible_end_time = pygame.time.get_ticks() + PLAYER_INVINCIBLE_TIME
        return True

    def heal(self):
        if self.lives < self.max_lives:
            self.lives += 1
            return True
        return False

    def activate_shield(self, duration=POWERUP_DURATION):
        self.shield_active = True
        self.shield_end_time = pygame.time.get_ticks() + duration

    def activate_rapid_fire(self, duration=POWERUP_DURATION):
        self.rapid_fire = True
        self.rapid_fire_end_time = pygame.time.get_ticks() + duration

    def activate_multi_shot(self, duration=POWERUP_DURATION):
        self.multi_shot = True
        self.multi_shot_end_time = pygame.time.get_ticks() + duration

    def upgrade_damage(self):
        self.damage_level += 1
        self.bullet_damage += 1

    def upgrade_speed(self):
        self.speed_level += 1
        self.speed += 1

    def upgrade_fire_rate(self):
        self.fire_rate_level += 1
        self.fire_rate = max(100, self.fire_rate - 40)

    def reset(self):
        self.x = (SCREEN_WIDTH - self.width) // 2
        self.y = SCREEN_HEIGHT - self.height - 20
        self.lives = self.max_lives
        self.bullet_damage = BULLET_DAMAGE
        self.fire_rate = FIRE_RATE
        self.speed = PLAYER_SPEED
        self.invincible = False
        self.shield_active = False
        self.rapid_fire = False
        self.multi_shot = False
        self.damage_level = 1
        self.speed_level = 1
        self.fire_rate_level = 1
