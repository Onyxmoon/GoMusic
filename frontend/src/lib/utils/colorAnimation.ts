/**
 * Color animation utilities for smooth color transitions
 */

export interface RGB {
  r: number;
  g: number;
  b: number;
}

/**
 * Linear interpolation between two colors
 * @param from - Starting color
 * @param to - Target color
 * @param t - Progress (0-1)
 * @returns Interpolated color
 */
export function lerpColor(from: RGB, to: RGB, t: number): RGB {
  return {
    r: Math.round(from.r + (to.r - from.r) * t),
    g: Math.round(from.g + (to.g - from.g) * t),
    b: Math.round(from.b + (to.b - from.b) * t)
  };
}

/**
 * Cubic ease-in-out function for smooth animations
 * @param t - Progress (0-1)
 * @returns Eased progress (0-1)
 */
export function easeInOutCubic(t: number): number {
  return t < 0.5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2;
}

/**
 * Converts RGB color to CSS string format
 * @param color - RGB color object
 * @returns CSS color string (e.g., "255, 128, 0")
 */
export function rgbToString(color: RGB): string {
  return `${color.r}, ${color.g}, ${color.b}`;
}