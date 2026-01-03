<script lang="ts">
  // Stores
  import { player } from '../stores/player.svelte';

  // Icons
  import { Shuffle, Maximize2 } from 'lucide-svelte';

  // Utils & Composables
  import { extractColorsFromImage } from '../utils/colorExtractor';
  import { formatTime } from '../utils/timeFormat';
  import { harmonizeColors } from '../utils/colorAnimation';
  import { useColorAnimation } from '../composables/useColorAnimation.svelte';

  // Constants
  import { DEFAULT_COLORS, TINT_AMOUNT, MIN_BRIGHTNESS } from '../constants/playerTheme';

  // Components
  import LiquidGlassProgress from './BottomPlayer/LiquidGlassProgress.svelte';
  import TrackInfo from './BottomPlayer/TrackInfo.svelte';
  import PlaybackControls from './BottomPlayer/PlaybackControls.svelte';
  import VolumeControl from './BottomPlayer/VolumeControl.svelte';


  //// Color Animation

  const colors = useColorAnimation();

  // Extract colors from album artwork and apply to player background
  $effect(() => {
    if (player.currentTrack?.hasArtwork) {
      const artworkUrl = `/artwork/stream?id=${encodeURIComponent(player.currentTrack.id)}`;
      extractColorsFromImage(artworkUrl)
        .then(extractedColors => {
          const harmonized = harmonizeColors(extractedColors, TINT_AMOUNT, MIN_BRIGHTNESS);
          colors.animateColors(harmonized[0], harmonized[1], harmonized[2]);
        })
        .catch(() => {
          colors.animateColors(DEFAULT_COLORS[0], DEFAULT_COLORS[1], DEFAULT_COLORS[2]);
        });
    } else {
      colors.animateColors(DEFAULT_COLORS[0], DEFAULT_COLORS[1], DEFAULT_COLORS[2]);
    }
  });

  // ============================================================================
  // Event Handlers
  // ============================================================================

  function handleProgressClick(e: MouseEvent) {
    const target = e.currentTarget as HTMLElement;
    const rect = target.getBoundingClientRect();
    const percent = (e.clientX - rect.left) / rect.width;
    const newTime = percent * player.duration;
    player.seek(newTime);
  }

  function togglePlayPause() {
    if (player.isPlaying) {
      player.pause();
    } else {
      player.resume();
    }
  }

  function handleVolumeChange(volume: number) {
    player.setVolume(volume);
  }

  function handlePrevious() {
    player.previous();
  }

  function handleNext() {
    player.next();
  }

  function handleToggleMute() {
    player.toggleMute();
  }
</script>


<div class="bottom-player" class:playing={player.isPlaying} style="--c1: {colors.color1}; --c2: {colors.color2}; --c3: {colors.color3};">
  {#if player.currentTrack}
    <TrackInfo track={player.currentTrack} />

    <div class="player-center">
      <PlaybackControls
        isPlaying={player.isPlaying}
        onPlayPause={togglePlayPause}
        onPrevious={handlePrevious}
        onNext={handleNext}
      />

      <div class="progress-section">
        <span class="time">{formatTime(player.currentTime)}</span>
        <LiquidGlassProgress progress={player.progress} onClick={handleProgressClick} animated={player.isPlaying} />
        <span class="time">{formatTime(player.duration)}</span>
      </div>
    </div>

    <div class="player-right">
      <button class="icon-btn" aria-label="Toggle shuffle"><Shuffle size={16} /></button>
      <VolumeControl
        volume={player.volume}
        isMuted={player.isMuted}
        onVolumeChange={handleVolumeChange}
        onToggleMute={handleToggleMute}
      />
      <button class="icon-btn" aria-label="Toggle fullscreen"><Maximize2 size={16} /></button>
    </div>
  {:else}
    <div class="player-empty">
      No track playing
    </div>
  {/if}
</div>


<style>
  /* Aurora Animations */
  @keyframes auroraFlow1 {
    0% {
      transform: translate(-30%, -20%) rotate(-5deg);
      opacity: 0.4;
    }
    25% {
      transform: translate(20%, 30%) rotate(3deg);
      opacity: 0.6;
    }
    50% {
      transform: translate(-10%, 40%) rotate(-2deg);
      opacity: 0.5;
    }
    75% {
      transform: translate(30%, -10%) rotate(4deg);
      opacity: 0.55;
    }
    100% {
      transform: translate(-30%, -20%) rotate(-5deg);
      opacity: 0.4;
    }
  }

  @keyframes auroraFlow2 {
    0% {
      transform: translate(40%, 20%) rotate(3deg);
      opacity: 0.35;
    }
    33% {
      transform: translate(-20%, -30%) rotate(-4deg);
      opacity: 0.5;
    }
    66% {
      transform: translate(10%, 35%) rotate(2deg);
      opacity: 0.45;
    }
    100% {
      transform: translate(40%, 20%) rotate(3deg);
      opacity: 0.35;
    }
  }

  /* Main Player Container */
  .bottom-player {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    height: 80px;
    background: rgba(255, 255, 255, 0.3);
    backdrop-filter: blur(30px) saturate(180%);
    -webkit-backdrop-filter: blur(30px) saturate(180%);
    border-top: 1px solid rgba(255, 255, 255, 0.5);
    box-shadow: 0 -4px 24px rgba(0, 0, 0, 0.08);
    display: flex;
    align-items: center;
    padding: 0 20px;
    gap: 24px;
    overflow: hidden;
  }

  /* Aurora Background Layers */
  .bottom-player::before {
    content: '';
    position: absolute;
    top: -50%;
    left: -50%;
    width: 200%;
    height: 200%;
    background: radial-gradient(
      ellipse at center,
      rgba(var(--c1), 0.6) 0%,
      rgba(var(--c2), 0.7) 35%,
      rgba(var(--c3), 0.55) 70%,
      transparent 100%
    );
    filter: blur(50px);
    pointer-events: none;
    z-index: 0;
    animation: auroraFlow1 35s ease-in-out infinite;
  }

  .bottom-player::after {
    content: '';
    position: absolute;
    top: -50%;
    left: -50%;
    width: 200%;
    height: 200%;
    background: radial-gradient(
      ellipse at center,
      rgba(var(--c3), 0.5) 0%,
      rgba(var(--c1), 0.65) 40%,
      rgba(var(--c2), 0.5) 75%,
      transparent 100%
    );
    filter: blur(55px);
    pointer-events: none;
    z-index: 0;
    animation: auroraFlow2 50s ease-in-out infinite;
  }

  /* Player Sections */
  .player-center {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8px;
    align-items: center;
    position: relative;
    z-index: 1;
  }

  .progress-section {
    display: flex;
    align-items: center;
    gap: 12px;
    width: 100%;
    max-width: 600px;
  }

  .time {
    font-size: 11px;
    color: #31343a;
    min-width: 40px;
    text-align: center;
  }

  .player-right {
    display: flex;
    gap: 8px;
    align-items: center;
    flex: 1;
    justify-content: flex-end;
    position: relative;
    z-index: 1;
  }

  /* Buttons */
  .icon-btn {
    width: 32px;
    height: 32px;
    border-radius: 6px;
    border: none;
    background: transparent;
    color: #40434c;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .icon-btn:hover {
    background: rgba(0, 0, 0, 0.05);
    color: #2d2d2d;
  }

  /* Empty State */
  .player-empty {
    width: 100%;
    text-align: center;
    color: #9ca3af;
  }
</style>