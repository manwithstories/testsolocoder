import pygame
from config import *


class Bullet:
    def __init__(self, x, y, is_enemy=False, damage=BULLET_DAMAGE, speed=None):
        self.x = x
        self.y = y
        self.is_enemy = is_enemy
        self.damage = damage
        self.speed = speed if speed else BULLET_SPEED
        self.width = BULLET_WIDTH
        self.height = BULLET_HEIGHT
        self.active = True

        if is_enemy:
            self.color = RED
        else:
            self.color = YELLOW

    def update(self):
        if self.is_enemy:
            self.y += self.speed
        else:
            self.y -= self.speed

        if self.y < -self.height or self.y > SCREEN_HEIGHT:
            self.active = False

    def draw(self, screen):
        pygame.draw.rect(screen, self.color, (self.x, self.y, self.width, self.height))

    def get_rect(self):
        return pygame.Rect(self.x, self.y, self.width, self.height)


class BulletManager:
    def __init__(self):
        self.player_bullets = []
        self.enemy_bullets = []

    def add_player_bullet(self, x, y, damage=BULLET_DAMAGE, speed=None):
        bullet = Bullet(x, y, False, damage, speed)
        self.player_bullets.append(bullet)

    def add_enemy_bullet(self, x, y, damage=1, speed=None):
        bullet = Bullet(x, y, True, damage, speed)
        self.enemy_bullets.append(bullet)

    def update(self):
        for bullet in self.player_bullets:
            bullet.update()
        for bullet in self.enemy_bullets:
            bullet.update()

        self.player_bullets = [b for b in self.player_bullets if b.active]
        self.enemy_bullets = [b for b in self.enemy_bullets if b.active]

    def draw(self, screen):
        for bullet in self.player_bullets:
            bullet.draw(screen)
        for bullet in self.enemy_bullets:
            bullet.draw(screen)

    def clear_all(self):
        self.player_bullets.clear()
        self.enemy_bullets.clear()
