import pygame
import sys
from config import *
from player import Player
from bullets import BulletManager
from enemies import EnemyManager
from powerups import PowerUpManager
from ui import HUD, MainMenu, WaveBreakScreen, Shop, GameOverScreen
from audio import AudioManager


class Game:
    def __init__(self):
        pygame.init()
        self.screen = pygame.display.set_mode((SCREEN_WIDTH, SCREEN_HEIGHT))
        pygame.display.set_caption('太空射击 - Space Shooter')
        self.clock = pygame.time.Clock()
        self.running = True

        self.audio = AudioManager()
        self.player = Player()
        self.bullet_manager = BulletManager()
        self.enemy_manager = EnemyManager()
        self.powerup_manager = PowerUpManager()
        self.hud = HUD()
        self.main_menu = MainMenu()
        self.wave_break = WaveBreakScreen()
        self.shop = Shop()
        self.game_over = GameOverScreen()

        self.state = 'menu'
        self.score = 0
        self.wave = 1

        self.stars = []
        for _ in range(100):
            self.stars.append({
                'x': pygame.time.get_ticks() % SCREEN_WIDTH,
                'y': pygame.time.get_ticks() % SCREEN_HEIGHT,
                'speed': 1 + pygame.time.get_ticks() % 3
            })
            pygame.time.wait(1)

    def handle_events(self):
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                self.running = False

            if self.state == 'menu':
                result = self.main_menu.handle_event(event)
                if result == '开始游戏':
                    self.start_game()
                elif result == '退出':
                    self.running = False

            elif self.state == 'playing':
                self.player.handle_event(event)

            elif self.state == 'wave_break':
                pass

            elif self.state == 'shop':
                result = self.shop.handle_event(event)
                if result == 'continue':
                    self.exit_shop()
                elif result == 'purchased':
                    self.score = self.shop.get_updated_score()
                    self.audio.play_buy()

            elif self.state == 'game_over':
                result = self.game_over.handle_event(event)
                if result == '重新开始':
                    self.start_game()
                elif result == '返回主菜单':
                    self.state = 'menu'

    def start_game(self):
        self.player.reset()
        self.bullet_manager.clear_all()
        self.enemy_manager.clear_all()
        self.powerup_manager.clear_all()
        self.score = 0
        self.wave = 1
        self.state = 'wave_break'
        self.wave_break.start(self.wave)
        self.audio.play_wave_complete()

    def start_next_wave(self):
        self.wave += 1

        if self.wave % SHOP_WAVE_INTERVAL == 1 and self.wave > 1:
            self.state = 'shop'
            self.shop.open(self.player, self.score)
        else:
            self.state = 'wave_break'
            self.wave_break.start(self.wave)
            self.audio.play_wave_complete()

        if self.wave % BOSS_SPAWN_WAVE == 0:
            self.audio.play_boss_appear()

    def exit_shop(self):
        self.score = self.shop.get_updated_score()
        self.state = 'wave_break'
        self.wave_break.start(self.wave)
        self.audio.play_wave_complete()

    def update(self):
        if self.state == 'menu':
            self.update_stars()

        elif self.state == 'playing':
            self.update_stars()
            self.player.update(self.bullet_manager)
            self.enemy_manager.update(self.bullet_manager)
            self.bullet_manager.update()
            self.powerup_manager.update()

            self.check_collisions()
            self.powerup_manager.check_collisions(self.player)

            if self.enemy_manager.check_wave_complete():
                self.start_next_wave()

            if self.player.lives <= 0:
                self.state = 'game_over'
                self.game_over.show(self.score, self.wave)

        elif self.state == 'wave_break':
            self.update_stars()
            if self.wave_break.is_done():
                self.enemy_manager.start_wave(self.wave)
                self.state = 'playing'

        elif self.state == 'shop':
            pass

        elif self.state == 'game_over':
            self.update_stars()

    def update_stars(self):
        for star in self.stars:
            star['y'] += star['speed']
            if star['y'] > SCREEN_HEIGHT:
                star['y'] = 0
                star['x'] = pygame.time.get_ticks() % SCREEN_WIDTH

    def check_collisions(self):
        player_rect = self.player.get_rect()

        for bullet in self.bullet_manager.player_bullets:
            if not bullet.active:
                continue

            bullet_rect = bullet.get_rect()
            for enemy in self.enemy_manager.get_all_enemies():
                if enemy.get_rect().colliderect(bullet_rect):
                    bullet.active = False
                    if enemy.take_damage(bullet.damage):
                        self.score += enemy.points
                        self.audio.play_explosion()
                        self.powerup_manager.spawn(
                            enemy.x + enemy.width // 2 - POWERUP_SIZE // 2,
                            enemy.y + enemy.height // 2
                        )
                    else:
                        self.audio.play_hit()
                    break

        for bullet in self.bullet_manager.enemy_bullets:
            if not bullet.active:
                continue

            if bullet.get_rect().colliderect(player_rect):
                bullet.active = False
                if self.player.take_damage():
                    self.audio.play_player_hit()

        for enemy in self.enemy_manager.get_all_enemies():
            if enemy.get_rect().colliderect(player_rect):
                if self.player.take_damage():
                    self.audio.play_player_hit()
                if hasattr(enemy, 'entering'):
                    if not enemy.entering:
                        enemy.take_damage(5)
                else:
                    enemy.active = False

    def draw(self):
        self.screen.fill(BLACK)
        self.draw_stars()

        if self.state == 'menu':
            self.main_menu.draw(self.screen)

        elif self.state == 'playing':
            self.powerup_manager.draw(self.screen)
            self.bullet_manager.draw(self.screen)
            self.enemy_manager.draw(self.screen)
            self.player.draw(self.screen)
            self.hud.draw(self.screen, self.player, self.score, self.wave, self.enemy_manager)

        elif self.state == 'wave_break':
            self.powerup_manager.draw(self.screen)
            self.bullet_manager.draw(self.screen)
            self.enemy_manager.draw(self.screen)
            self.player.draw(self.screen)
            self.hud.draw(self.screen, self.player, self.score, self.wave, self.enemy_manager)
            self.wave_break.draw(self.screen)

        elif self.state == 'shop':
            self.shop.draw(self.screen)

        elif self.state == 'game_over':
            self.powerup_manager.draw(self.screen)
            self.bullet_manager.draw(self.screen)
            self.enemy_manager.draw(self.screen)
            self.player.draw(self.screen)
            self.hud.draw(self.screen, self.player, self.score, self.wave, self.enemy_manager)
            self.game_over.draw(self.screen)

        pygame.display.flip()

    def draw_stars(self):
        for star in self.stars:
            size = 1 if star['speed'] < 2 else 2
            pygame.draw.circle(self.screen, WHITE, (int(star['x']), int(star['y'])), size)

    def run(self):
        while self.running:
            self.handle_events()
            self.update()
            self.draw()
            self.clock.tick(FPS)

        pygame.quit()
        sys.exit()


if __name__ == '__main__':
    game = Game()
    game.run()
