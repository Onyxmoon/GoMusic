import { lerpColor, easeInOutCubic, rgbToString, type RGB } from '../utils/colorAnimation';
import { DEFAULT_COLORS, COLOR_ANIMATION_DURATION } from '../constants/playerTheme';

/**
 * Class for managing smooth color transitions with animation
 * Uses Svelte 5's $state for reactive color strings
 */
class ColorAnimation {
  // Current colors for smooth interpolation (not reactive)
  currentColor1: RGB = DEFAULT_COLORS[0];
  currentColor2: RGB = DEFAULT_COLORS[1];
  currentColor3: RGB = DEFAULT_COLORS[2];

  // Reactive CSS color strings
  color1 = $state(rgbToString(DEFAULT_COLORS[0]));
  color2 = $state(rgbToString(DEFAULT_COLORS[1]));
  color3 = $state(rgbToString(DEFAULT_COLORS[2]));

  animationFrameId: number | null = null;

  /**
   * Animates from current colors to target colors over COLOR_ANIMATION_DURATION
   */
  animateColors(target1: RGB, target2: RGB, target3: RGB) {
    if (this.animationFrameId !== null) {
      cancelAnimationFrame(this.animationFrameId);
    }

    const start1 = { ...this.currentColor1 };
    const start2 = { ...this.currentColor2 };
    const start3 = { ...this.currentColor3 };

    const startTime = performance.now();

    const animate = (currentTime: number) => {
      const elapsed = currentTime - startTime;
      const progress = Math.min(elapsed / COLOR_ANIMATION_DURATION, 1);
      const easedProgress = easeInOutCubic(progress);

      const interpolated1 = lerpColor(start1, target1, easedProgress);
      const interpolated2 = lerpColor(start2, target2, easedProgress);
      const interpolated3 = lerpColor(start3, target3, easedProgress);

      this.color1 = rgbToString(interpolated1);
      this.color2 = rgbToString(interpolated2);
      this.color3 = rgbToString(interpolated3);

      if (progress < 1) {
        this.animationFrameId = requestAnimationFrame(animate);
      } else {
        this.currentColor1 = target1;
        this.currentColor2 = target2;
        this.currentColor3 = target3;
        this.animationFrameId = null;
      }
    };

    this.animationFrameId = requestAnimationFrame(animate);
  }
}

/**
 * Composable for managing smooth color transitions with animation
 * Returns a reactive ColorAnimation instance
 */
export function useColorAnimation() {
  return new ColorAnimation();
}