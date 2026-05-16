import pygame
from config import *


class HUD:
    def __init__(self):
        pass

    def draw(self, screen, player, score, wave, enemy_manager):
        score_text = FONT_MEDIUM.render(f'分数: {score}', True, WHITE)
        screen.blit(score_text, (10, 10))

        wave_text = FONT_MEDIUM.render(f'波次: {wave}', True, WHITE)
        screen.blit(wave_text, (10, 50))

        lives_text = FONT_MEDIUM.render('生命: ', True, WHITE)
        screen.blit(lives_text, (10, 90))
        for i in range(player.max_lives):
            color = GREEN if i < player.lives else DARK_GRAY
            pygame.draw.circle(screen, color, (110 + i * 30, 105), 10)

        powerup_y = 140
        if player.shield_active:
            shield_text = FONT_SMALL.render('护盾激活中', True, CYAN)
            screen.blit(shield_text, (10, powerup_y))
            powerup_y += 25
        if player.rapid_fire:
            rapid_text = FONT_SMALL.render('射速提升中', True, YELLOW)
            screen.blit(rapid_text, (10, powerup_y))
            powerup_y += 25
        if player.multi_shot:
            multi_text = FONT_SMALL.render('散射激活中', True, MAGENTA)
            screen.blit(multi_text, (10, powerup_y))

        if enemy_manager.boss:
            boss = enemy_manager.boss
            bar_width = SCREEN_WIDTH - 200
            bar_height = 20
            bar_x = 100
            bar_y = SCREEN_HEIGHT - 40
            hp_ratio = boss.hp / boss.max_hp

            pygame.draw.rect(screen, DARK_GRAY, (bar_x, bar_y, bar_width, bar_height))
            pygame.draw.rect(screen, RED, (bar_x, bar_y, int(bar_width * hp_ratio), bar_height))
            pygame.draw.rect(screen, WHITE, (bar_x, bar_y, bar_width, bar_height), 2)

            boss_text = FONT_SMALL.render(f'BOSS - 第 {boss.wave} 波', True, WHITE)
            boss_text_rect = boss_text.get_rect(center=(SCREEN_WIDTH // 2, bar_y - 15))
            screen.blit(boss_text, boss_text_rect)


class MainMenu:
    def __init__(self):
        self.selected = 0
        self.options = ['开始游戏', '退出']

    def handle_event(self, event):
        if event.type == pygame.KEYDOWN:
            if event.key == pygame.K_UP:
                self.selected = (self.selected - 1) % len(self.options)
            elif event.key == pygame.K_DOWN:
                self.selected = (self.selected + 1) % len(self.options)
            elif event.key == pygame.K_RETURN or event.key == pygame.K_SPACE:
                return self.options[self.selected]
        return None

    def draw(self, screen):
        screen.fill(BLACK)

        title = FONT_TITLE.render('太空射击', True, CYAN)
        title_rect = title.get_rect(center=(SCREEN_WIDTH // 2, SCREEN_HEIGHT // 3))
        screen.blit(title, title_rect)

        subtitle = FONT_MEDIUM.render('SPACE SHOOTER', True, WHITE)
        subtitle_rect = subtitle.get_rect(center=(SCREEN_WIDTH // 2, SCREEN_HEIGHT // 3 + 60))
        screen.blit(subtitle, subtitle_rect)

        for i, option in enumerate(self.options):
            color = YELLOW if i == self.selected else WHITE
            text = FONT_MEDIUM.render(option, True, color)
            text_rect = text.get_rect(center=(SCREEN_WIDTH // 2, SCREEN_HEIGHT // 2 + i * 50))
            screen.blit(text, text_rect)

        controls = FONT_SMALL.render('方向键/A D移动 | 空格射击 | 回车选择', True, GRAY)
        controls_rect = controls.get_rect(center=(SCREEN_WIDTH // 2, SCREEN_HEIGHT - 50))
        screen.blit(controls, controls_rect)


class WaveBreakScreen:
    def __init__(self):
        self.start_time = 0
        self.wave = 1

    def start(self, wave):
        self.wave = wave
        self.start_time = pygame.time.get_ticks()

    def is_done(self):
        return pygame.time.get_ticks() - self.start_time >= WAVE_BREAK_TIME

    def draw(self, screen):
        overlay = pygame.Surface((SCREEN_WIDTH, SCREEN_HEIGHT))
        overlay.set_alpha(200)
        overlay.fill(BLACK)
        screen.blit(overlay, (0, 0))

        wave_text = FONT_LARGE.render(f'第 {self.wave} 波', True, CYAN)
        wave_rect = wave_text.get_rect(center=(SCREEN_WIDTH // 2, SCREEN_HEIGHT // 2 - 30))
        screen.blit(wave_text, wave_rect)

        ready_text = FONT_MEDIUM.render('准备战斗!', True, YELLOW)
        ready_rect = ready_text.get_rect(center=(SCREEN_WIDTH // 2, SCREEN_HEIGHT // 2 + 30))
        screen.blit(ready_text, ready_rect)

        elapsed = pygame.time.get_ticks() - self.start_time
        remaining = max(0, (WAVE_BREAK_TIME - elapsed) // 1000 + 1)
        count_text = FONT_MEDIUM.render(str(remaining), True, WHITE)
        count_rect = count_text.get_rect(center=(SCREEN_WIDTH // 2, SCREEN_HEIGHT // 2 + 80))
        screen.blit(count_text, count_rect)


class Shop:
    def __init__(self):
        self.selected = 0
        self.options = [
            ('子弹威力', SHOP_UPGRADE_DAMAGE_COST, 'damage'),
            ('移动速度', SHOP_UPGRADE_SPEED_COST, 'speed'),
            ('射击速度', SHOP_UPGRADE_FIRE_RATE_COST, 'fire_rate'),
            ('恢复生命', SHOP_BUY_LIFE_COST, 'heal'),
            ('继续游戏', 0, 'continue')
        ]
        self.player = None
        self.score = 0

    def open(self, player, score):
        self.player = player
        self.score = score
        self.selected = 0

    def handle_event(self, event):
        if event.type == pygame.KEYDOWN:
            if event.key == pygame.K_UP:
                self.selected = (self.selected - 1) % len(self.options)
            elif event.key == pygame.K_DOWN:
                self.selected = (self.selected + 1) % len(self.options)
            elif event.key == pygame.K_RETURN or event.key == pygame.K_SPACE:
                return self.buy()
        return None

    def buy(self):
        option = self.options[self.selected]
        name, cost, action = option

        if action == 'continue':
            return 'continue'

        if self.score < cost:
            return None

        if action == 'damage':
            self.player.upgrade_damage()
        elif action == 'speed':
            self.player.upgrade_speed()
        elif action == 'fire_rate':
            self.player.upgrade_fire_rate()
        elif action == 'heal':
            if not self.player.heal():
                return None

        self.score -= cost
        return 'purchased'

    def get_updated_score(self):
        return self.score

    def draw(self, screen):
        overlay = pygame.Surface((SCREEN_WIDTH, SCREEN_HEIGHT))
        overlay.set_alpha(220)
        overlay.fill(BLACK)
        screen.blit(overlay, (0, 0))

        title = FONT_LARGE.render('商店', True, YELLOW)
        title_rect = title.get_rect(center=(SCREEN_WIDTH // 2, 80))
        screen.blit(title, title_rect)

        score_text = FONT_MEDIUM.render(f'积分: {self.score}', True, WHITE)
        score_rect = score_text.get_rect(center=(SCREEN_WIDTH // 2, 130))
        screen.blit(score_text, score_rect)

        stats_y = 180
        stat_texts = [
            f'子弹威力 Lv.{self.player.damage_level}',
            f'移动速度 Lv.{self.player.speed_level}',
            f'射击速度 Lv.{self.player.fire_rate_level}',
            f'生命: {self.player.lives}/{self.player.max_lives}'
        ]
        for stat in stat_texts:
            text = FONT_SMALL.render(stat, True, GRAY)
            text_rect = text.get_rect(center=(SCREEN_WIDTH // 2, stats_y))
            screen.blit(text, text_rect)
            stats_y += 22

        for i, (name, cost, action) in enumerate(self.options):
            color = YELLOW if i == self.selected else WHITE
            if cost > self.score and action != 'continue':
                color = DARK_GRAY

            if action == 'continue':
                display_text = name
            else:
                display_text = f'{name} - {cost} 积分'

            text = FONT_MEDIUM.render(display_text, True, color)
            text_rect = text.get_rect(center=(SCREEN_WIDTH // 2, 300 + i * 45))
            screen.blit(text, text_rect)

        controls = FONT_SMALL.render('方向键选择 | 回车购买 | 选择"继续游戏"退出商店', True, GRAY)
        controls_rect = controls.get_rect(center=(SCREEN_WIDTH // 2, SCREEN_HEIGHT - 30))
        screen.blit(controls, controls_rect)


class GameOverScreen:
    def __init__(self):
        self.score = 0
        self.wave = 1
        self.selected = 0
        self.options = ['重新开始', '返回主菜单']

    def show(self, score, wave):
        self.score = score
        self.wave = wave
        self.selected = 0

    def handle_event(self, event):
        if event.type == pygame.KEYDOWN:
            if event.key == pygame.K_UP:
                self.selected = (self.selected - 1) % len(self.options)
            elif event.key == pygame.K_DOWN:
                self.selected = (self.selected + 1) % len(self.options)
            elif event.key == pygame.K_RETURN or event.key == pygame.K_SPACE:
                return self.options[self.selected]
        return None

    def draw(self, screen):
        overlay = pygame.Surface((SCREEN_WIDTH, SCREEN_HEIGHT))
        overlay.set_alpha(200)
        overlay.fill(BLACK)
        screen.blit(overlay, (0, 0))

        title = FONT_TITLE.render('游戏结束', True, RED)
        title_rect = title.get_rect(center=(SCREEN_WIDTH // 2, SCREEN_HEIGHT // 3))
        screen.blit(title, title_rect)

        score_text = FONT_LARGE.render(f'最终分数: {self.score}', True, WHITE)
        score_rect = score_text.get_rect(center=(SCREEN_WIDTH // 2, SCREEN_HEIGHT // 2 - 20))
        screen.blit(score_text, score_rect)

        wave_text = FONT_MEDIUM.render(f'坚持到第 {self.wave} 波', True, CYAN)
        wave_rect = wave_text.get_rect(center=(SCREEN_WIDTH // 2, SCREEN_HEIGHT // 2 + 30))
        screen.blit(wave_text, wave_rect)

        for i, option in enumerate(self.options):
            color = YELLOW if i == self.selected else WHITE
            text = FONT_MEDIUM.render(option, True, color)
            text_rect = text.get_rect(center=(SCREEN_WIDTH // 2, SCREEN_HEIGHT // 2 + 90 + i * 45))
            screen.blit(text, text_rect)
