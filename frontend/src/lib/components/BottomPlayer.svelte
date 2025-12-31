<script lang="ts">
  import { player } from '../stores/player.svelte';
  import { SkipBack, Play, Pause, SkipForward, Heart, Shuffle, Volume2, Maximize2 } from 'lucide-svelte';
  import { extractColorsFromImage } from '../utils/colorExtractor';
  import { formatTime } from '../utils/timeFormat';
  import { lerpColor, easeInOutCubic, type RGB } from '../utils/colorAnimation';
  import LiquidGlassProgress from './LiquidGlassProgress.svelte';

  let playerElement: HTMLDivElement;

  // Current colors for smooth interpolation
  let currentColor1: RGB = { r: 255, g: 214, b: 214 };
  let currentColor2: RGB = { r: 232, g: 213, b: 255 };
  let currentColor3: RGB = { r: 213, g: 232, b: 255 };

  let color1 = $state('255, 214, 214');
  let color2 = $state('232, 213, 255');
  let color3 = $state('213, 232, 255');

  let animationFrameId: number | null = null;

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

  function animateColors(target1: RGB, target2: RGB, target3: RGB) {
    if (animationFrameId !== null) {
      cancelAnimationFrame(animationFrameId);
    }

    const start1 = { ...currentColor1 };
    const start2 = { ...currentColor2 };
    const start3 = { ...currentColor3 };

    const duration = 1200; // 1.2 seconds
    const startTime = performance.now();

    function animate(currentTime: number) {
      const elapsed = currentTime - startTime;
      const progress = Math.min(elapsed / duration, 1);
      const easedProgress = easeInOutCubic(progress);

      const interpolated1 = lerpColor(start1, target1, easedProgress);
      const interpolated2 = lerpColor(start2, target2, easedProgress);
      const interpolated3 = lerpColor(start3, target3, easedProgress);

      color1 = `${interpolated1.r}, ${interpolated1.g}, ${interpolated1.b}`;
      color2 = `${interpolated2.r}, ${interpolated2.g}, ${interpolated2.b}`;
      color3 = `${interpolated3.r}, ${interpolated3.g}, ${interpolated3.b}`;

      if (progress < 1) {
        animationFrameId = requestAnimationFrame(animate);
      } else {
        currentColor1 = target1;
        currentColor2 = target2;
        currentColor3 = target3;
        animationFrameId = null;
      }
    }

    animationFrameId = requestAnimationFrame(animate);
  }

  // Extract colors from album artwork and apply to player background
  $effect(() => {
    if (player.currentTrack?.hasArtwork) {
      const artworkUrl = `/artwork/stream?id=${encodeURIComponent(player.currentTrack.id)}`;
      extractColorsFromImage(artworkUrl)
        .then(colors => {
          const harmonized = colors.map(c => {
            const tintAmount = 0.15; // Reduced from 0.4 for more color intensity
            const r = Math.round(c.r + (255 - c.r) * tintAmount);
            const g = Math.round(c.g + (255 - c.g) * tintAmount);
            const b = Math.round(c.b + (255 - c.b) * tintAmount);

            const brightness = (r + g + b) / 3;
            const minBrightness = 100;

            if (brightness < minBrightness) {
              const boost = minBrightness - brightness;
              return {
                r: Math.min(255, r + boost),
                g: Math.min(255, g + boost),
                b: Math.min(255, b + boost)
              };
            }
            return { r, g, b };
          });

          animateColors(harmonized[0], harmonized[1], harmonized[2]);
        })
        .catch(() => {
          animateColors(
            { r: 255, g: 214, b: 214 },
            { r: 232, g: 213, b: 255 },
            { r: 213, g: 232, b: 255 }
          );
        });
    } else {
      animateColors(
        { r: 255, g: 214, b: 214 },
        { r: 232, g: 213, b: 255 },
        { r: 213, g: 232, b: 255 }
      );
    }
  });
</script>

<div class="bottom-player" class:playing={player.isPlaying} bind:this={playerElement} style="--c1: {color1}; --c2: {color2}; --c3: {color3};">
  {#if player.currentTrack}
    <div class="player-left">
      <div class="track-cover">
        {#if player.currentTrack.hasArtwork}
          <img src={`/artwork/stream?id=${encodeURIComponent(player.currentTrack.id)}`} alt={player.currentTrack.title} />
        {:else}
          <div class="cover-placeholder"></div>
        {/if}
      </div>
      <div class="track-info">
        <div class="track-title">{player.currentTrack.title}</div>
        <div class="track-artist">{player.currentTrack.artist}</div>
      </div>
      <button class="heart-btn"><Heart size={18} /></button>
    </div>

    <div class="player-center">
      <div class="controls">
        <button class="control-btn" on:click={player.previous}><SkipBack size={18} /></button>
        <button class="control-btn play-btn" on:click={togglePlayPause}>
          {#if player.isPlaying}
            <Pause size={18} />
          {:else}
            <Play size={18} />
          {/if}
        </button>
        <button class="control-btn" on:click={player.next}><SkipForward size={18} /></button>
      </div>

      <div class="progress-section">
        <span class="time">{formatTime(player.currentTime)}</span>
        <LiquidGlassProgress progress={player.progress} onClick={handleProgressClick} animated={player.isPlaying} />
        <span class="time">{formatTime(player.duration)}</span>
      </div>
    </div>

    <div class="player-right">
      <button class="icon-btn"><Shuffle size={16} /></button>
      <button class="icon-btn"><Volume2 size={16} /></button>
      <button class="icon-btn"><Maximize2 size={16} /></button>
    </div>
  {:else}
    <div class="player-empty">
      No track playing
    </div>
  {/if}
</div>

<style>
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

  @keyframes auroraFlow3 {
    0% {
      transform: translate(-15%, 30%) rotate(-2deg);
      opacity: 0.3;
    }
    40% {
      transform: translate(25%, -25%) rotate(5deg);
      opacity: 0.5;
    }
    80% {
      transform: translate(-35%, 15%) rotate(-3deg);
      opacity: 0.4;
    }
    100% {
      transform: translate(-15%, 30%) rotate(-2deg);
      opacity: 0.3;
    }
  }

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

  /* Aurora Layer 1 - Balanced flowing gradient */
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
  }

  /* Aurora Layer 2 - Balanced counter flow */
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
  }

  .bottom-player::before {
    animation: auroraFlow1 35s ease-in-out infinite;
  }

  .bottom-player::after {
    animation: auroraFlow2 50s ease-in-out infinite;
  }

  .player-left {
    display: flex;
    align-items: center;
    gap: 12px;
    flex: 1;
    min-width: 0;
    position: relative;
    z-index: 1;
  }

  .track-cover {
    width: 56px;
    height: 56px;
    border-radius: 8px;
    overflow: hidden;
    flex-shrink: 0;
  }

  .track-cover img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .cover-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #8a65ff, #ff6b9d);
    font-size: 24px;
  }

  .track-info {
    flex: 1;
    min-width: 0;
  }

  .track-title {
    font-size: 14px;
    font-weight: 500;
    color: #2d2d2d;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .track-artist {
    font-size: 12px;
    color: #31343a;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .heart-btn {
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: #31343a;
    font-size: 18px;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .heart-btn:hover {
    color: #ff5686;
  }

  .player-center {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8px;
    align-items: center;
    position: relative;
    z-index: 1;
  }

  .controls {
    display: flex;
    gap: 12px;
    align-items: center;
  }

  .control-btn {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    border: none;
    background: transparent;
    color: #2d2d2d;
    font-size: 16px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.15s ease;
  }

  .control-btn:hover {
    background: rgba(0, 0, 0, 0.05);
  }

  .play-btn {
    width: 40px;
    height: 40px;
    background: #2d2d2d;
    color: white;
  }

  .play-btn:hover {
    background: #1a1a1a;
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

  .player-empty {
    width: 100%;
    text-align: center;
    color: #9ca3af;
  }
</style>
