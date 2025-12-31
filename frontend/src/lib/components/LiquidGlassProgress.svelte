<script lang="ts">
  export let progress: number = 0; // 0-100
  export let onClick: ((e: MouseEvent) => void) | undefined = undefined;
  export let animated: boolean = true;

  let isHovering = false;
  let isDragging = false;

  function handleMouseDown(e: MouseEvent) {
    if (!onClick) return;

    isDragging = true;
    onClick(e);

    function onMouseMove(event: MouseEvent) {
      if (onClick) onClick(event);
    }

    function onMouseUp() {
      isDragging = false;
      window.removeEventListener('mousemove', onMouseMove);
      window.removeEventListener('mouseup', onMouseUp);
    }

    window.addEventListener('mousemove', onMouseMove);
    window.addEventListener('mouseup', onMouseUp);
  }
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
  class="liquid-glass-progress"
  class:clickable={onClick !== undefined}
  class:hovering={isHovering}
  class:dragging={isDragging}
  on:mouseenter={() => isHovering = true}
  on:mouseleave={() => isHovering = false}
>
  <div
    class="glass-track"
    on:mousedown={handleMouseDown}
  >
    <div
      class="liquid-fill"
      class:animated={animated && !isDragging}
      style="width: {Math.max(0, Math.min(100, progress))}%"
    >
      <div class="shimmer"></div>
    </div>
  </div>
</div>

<style>
  .liquid-glass-progress {
    width: 100%;
    height: 12px;
    position: relative;
    padding: 2px 0;
  }

  .liquid-glass-progress.clickable {
    cursor: pointer;
  }

  .liquid-glass-progress.dragging {
    cursor: grabbing;
    user-select: none;
  }

  .glass-track {
    position: relative;
    width: 100%;
    height: 100%;
    background: rgba(255, 255, 255, 0.15);
    backdrop-filter: blur(20px) saturate(200%) brightness(1.1);
    -webkit-backdrop-filter: blur(20px) saturate(200%) brightness(1.1);
    border-radius: 8px;
    overflow: hidden;
    border: 0.5px solid rgba(255, 255, 255, 0.25);
    box-shadow:
      inset 0 1px 2px rgba(0, 0, 0, 0.05),
      0 1px 4px rgba(0, 0, 0, 0.04);
    transition: all 0.3s ease;
  }

  .liquid-glass-progress.hovering .glass-track,
  .liquid-glass-progress.dragging .glass-track {
    background: rgba(255, 255, 255, 0.22);
    backdrop-filter: blur(20px) saturate(220%) brightness(1.15);
    -webkit-backdrop-filter: blur(20px) saturate(220%) brightness(1.15);
    border-color: rgba(255, 255, 255, 0.35);
  }

  .liquid-glass-progress.dragging .glass-track {
    transform: scaleY(1.15);
  }

  .liquid-fill {
    position: absolute;
    left: 0;
    top: 0;
    height: 100%;
    background: rgba(129, 129, 129, 0.39);
    backdrop-filter: blur(10px) saturate(250%) brightness(0.5);
    -webkit-backdrop-filter: blur(10px) saturate(10050%) brightness(0.5);
    border-radius: 7px;
    transition: width 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    overflow: hidden;
    box-shadow:
      0 0 8px rgba(255, 255, 255, 0.2),
      inset 0 1px 1px rgba(255, 255, 255, 0.4);
  }

  .liquid-fill.animated {
    transition: width 0.3s cubic-bezier(0.34, 1.2, 0.64, 1);
  }

  .shimmer {
    position: absolute;
    inset: 0;
    background: linear-gradient(
      90deg,
      transparent 0%,
      rgba(255, 255, 255, 0.3) 50%,
      transparent 100%
    );
    animation: shimmerFlow 4s ease-in-out infinite;
    opacity: 0.5;
  }

  .liquid-glass-progress.hovering .shimmer {
    animation-duration: 2.5s;
    opacity: 0.7;
  }

  @keyframes shimmerFlow {
    0% {
      transform: translateX(-100%);
      opacity: 0;
    }
    30% {
      opacity: 0.5;
    }
    70% {
      opacity: 0.5;
    }
    100% {
      transform: translateX(200%);
      opacity: 0;
    }
  }
</style>