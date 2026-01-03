/**
 * Player theme constants for colors, animations, and visual parameters
 */

import type { RGB } from '../utils/colorAnimation';

/**
 * Default gradient colors used when no album artwork is available
 * or when color extraction fails
 */
export const DEFAULT_COLORS: [RGB, RGB, RGB] = [
  { r: 255, g: 214, b: 214 }, // Soft pink
  { r: 232, g: 213, b: 255 }, // Soft purple
  { r: 213, g: 232, b: 255 }  // Soft blue
];

/**
 * Duration of color transition animation in milliseconds
 */
export const COLOR_ANIMATION_DURATION = 1200;

/**
 * Amount of white tint to apply to extracted colors (0-1)
 * Lower values = more vibrant, higher values = more pastel
 */
export const TINT_AMOUNT = 0.15;

/**
 * Minimum brightness threshold for colors (0-255)
 * Colors darker than this will be boosted
 */
export const MIN_BRIGHTNESS = 100;

/**
 * Aurora animation durations for background gradient layers
 */
export const AURORA_ANIMATION_DURATIONS = {
  layer1: '35s',
  layer2: '50s'
} as const;