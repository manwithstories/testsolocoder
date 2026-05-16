import pygame
import random
from config import *


class AudioManager:
    def __init__(self):
        self.enabled = True
        self.sounds = {}
        self.music_playing = False

        try:
            pygame.mixer.init()
            self._generate_sounds()
        except pygame.error:
            self.enabled = False

    def _generate_sounds(self):
        self.sounds['shoot'] = self._create_tone(880, 0.05, 0.3)
        self.sounds['hit'] = self._create_tone(220, 0.1, 0.3)
        self.sounds['explosion'] = self._create_noise(0.2, 0.4)
        self.sounds['powerup'] = self._create_tone(660, 0.15, 0.3)
        self.sounds['player_hit'] = self._create_tone(150, 0.3, 0.5)
        self.sounds['boss_appear'] = self._create_tone(100, 0.5, 0.6)
        self.sounds['buy'] = self._create_tone(523, 0.1, 0.3)
        self.sounds['wave_complete'] = self._create_tone_sequence([523, 659, 784], 0.15, 0.3)

    def _create_tone(self, frequency, duration, volume=0.3):
        sample_rate = 44100
        n_samples = int(sample_rate * duration)
        sound = pygame.mixer.Sound(
            buffer=bytes(
                [int(volume * 127 * (0.5 + 0.5 * (i * frequency * 2 * 3.14159 / sample_rate) % 1 * 2 - 1))
                 for i in range(n_samples)]
            )
        )
        return sound

    def _create_noise(self, duration, volume=0.3):
        sample_rate = 44100
        n_samples = int(sample_rate * duration)
        sound = pygame.mixer.Sound(
            buffer=bytes(
                [int(volume * 127 * (random.random() * 2 - 1)) for _ in range(n_samples)]
            )
        )
        return sound

    def _create_tone_sequence(self, frequencies, duration, volume=0.3):
        sample_rate = 44100
        total_samples = int(sample_rate * duration * len(frequencies))
        samples = []

        for freq in frequencies:
            n_samples = int(sample_rate * duration)
            for i in range(n_samples):
                sample = volume * 127 * (0.5 + 0.5 * (i * freq * 2 * 3.14159 / sample_rate) % 1 * 2 - 1)
                samples.append(int(sample))

        return pygame.mixer.Sound(buffer=bytes(samples))

    def play(self, sound_name):
        if self.enabled and sound_name in self.sounds:
            try:
                self.sounds[sound_name].play()
            except:
                pass

    def play_shoot(self):
        self.play('shoot')

    def play_hit(self):
        self.play('hit')

    def play_explosion(self):
        self.play('explosion')

    def play_powerup(self):
        self.play('powerup')

    def play_player_hit(self):
        self.play('player_hit')

    def play_boss_appear(self):
        self.play('boss_appear')

    def play_buy(self):
        self.play('buy')

    def play_wave_complete(self):
        self.play('wave_complete')
