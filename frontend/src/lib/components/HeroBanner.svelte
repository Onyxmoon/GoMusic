<script lang="ts">
  import { MoreVertical, Heart, Play } from 'lucide-svelte';

  export let title: string = 'Chill Vibes';
  export let subtitle: string = 'Playlist • 67 songs • 2h 31 min';
  export let coverUrl: string = 'https://via.placeholder.com/400x200';
  export let onPlay: () => void = () => {};
</script>

<div class="hero-banner">
  <img src={coverUrl} alt={title} class="hero-image" />

  <div class="hero-overlay">
    <div class="hero-content">
      <h1 class="hero-title">{title}</h1>
      <p class="hero-subtitle">{subtitle}</p>

      <div class="hero-actions">
        <button class="more-btn"><MoreVertical size={16} /></button>
        <button class="favorite-btn"><Heart size={16} /></button>
      </div>
    </div>

    <button class="play-button" on:click={onPlay}>
      <Play size={24} />
    </button>
  </div>
</div>

<style>
  .hero-banner {
    position: relative;
    height: 200px;
    border-radius: 16px;
    overflow: hidden;
    margin-bottom: 24px;
  }

  .hero-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .hero-overlay {
    position: absolute;
    inset: 0;
    background: linear-gradient(
      to bottom,
      rgba(0, 0, 0, 0.1) 0%,
      rgba(0, 0, 0, 0.4) 100%
    );
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    padding: 24px;
  }

  .hero-content {
    flex: 1;
  }

  .hero-title {
    font-size: 48px;
    font-weight: 700;
    color: white;
    margin: 0 0 8px 0;
    text-shadow: 0 2px 20px rgba(0, 0, 0, 0.3);
  }

  .hero-subtitle {
    font-size: 14px;
    color: rgba(255, 255, 255, 0.9);
    margin: 0 0 16px 0;
    text-shadow: 0 1px 10px rgba(0, 0, 0, 0.3);
  }

  .hero-actions {
    display: flex;
    gap: 12px;
  }

  .more-btn,
  .favorite-btn {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    border: none;
    background: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
    color: white;
    font-size: 16px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
  }

  .more-btn:hover,
  .favorite-btn:hover {
    background: rgba(255, 255, 255, 0.3);
    transform: scale(1.05);
  }

  .play-button {
    width: 64px;
    height: 64px;
    border-radius: 50%;
    border: none;
    background: linear-gradient(135deg, #8a65ff, #ff6b9d);
    color: white;
    font-size: 24px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 8px 24px rgba(138, 101, 255, 0.4);
    transition: all 0.2s ease;
  }

  .play-button:hover {
    transform: scale(1.05);
    box-shadow: 0 12px 32px rgba(138, 101, 255, 0.5);
  }

  .play-button:active {
    transform: scale(0.98);
  }

  .play-icon {
    margin-left: 3px;
  }
</style>
