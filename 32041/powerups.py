import pygame
import random
from config import *


class PowerUp:
    def __init__(self, x, y, power_type=None):
        self.x = x
        self.y = y
        self.size = POWERUP_SIZE
        self.speed = POWERUP_SPEED
        self.active = True

        if power_type is None:
            self.type = random.choice(POWERUP_TYPES)
        else:
            self.type = power_type

        self.color_map = {
            'shield': (0, 255, 255),
            'rapid_fire': (255, 255, 0),
            'heal': (0, 255, 0),
            'multi_shot': (255, 0, 255)
        }
        self.color = self.color_map.get(self.type, WHITE)

    def update(self):
        self.y += self.speed
        if self.y > SCREEN_HEIGHT:
            self.active = False

    def draw(self, screen):
        pygame.draw.circle(
            screen,
            self.color,
            (self.x + self.size // 2, self.y + self.size // 2),
            self.size // 2
        )
        pygame.draw.circle(
            screen,
            WHITE,
            (self.x + self.size // 2, self.y + self.size // 2),
            self.size // 2,
            2
        )

        icon_text = ''
        if self.type == 'shield':
            icon_text = 'S'
        elif self.type == 'rapid_fire':
            icon_text = 'F'
        elif self.type == 'heal':
            icon_text = '+'
        elif self.type == 'multi_shot':
            icon_text = 'M'

        text = FONT_SMALL.render(icon_text, True, WHITE)
        text_rect = text.get_rect(center=(self.x + self.size // 2, self.y + self.size // 2))
        screen.blit(text, text_rect)

    def apply(self, player):
        if self.type == 'shield':
            player.activate_shield()
        elif self.type == 'rapid_fire':
            player.activate_rapid_fire()
        elif self.type == 'heal':
            player.heal()
        elif self.type == 'multi_shot':
            player.activate_multi_shot()

    def get_rect(self):
        return pygame.Rect(self.x, self.y, self.size, self.size)


class PowerUpManager:
    def __init__(self):
        self.powerups = []

    def spawn(self, x, y):
        if random.random() < POWERUP_DROP_CHANCE:
            powerup = PowerUp(x, y)
            self.powerups.append(powerup)

    def update(self):
        for powerup in self.powerups:
            powerup.update()
        self.powerups = [p for p in self.powerups if p.active]

    def draw(self, screen):
        for powerup in self.powerups:
            powerup.draw(screen)

    def check_collisions(self, player):
        player_rect = player.get_rect()
        collected = []

        for powerup in self.powerups:
            if player_rect.colliderect(powerup.get_rect()):
                powerup.apply(player)
                collected.append(powerup)
                powerup.active = False

        self.powerups = [p for p in self.powerups if p.active]
        return collected

    def clear_all(self):
        self.powerups.clear()
