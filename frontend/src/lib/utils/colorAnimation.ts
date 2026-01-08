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

/**
 * Harmonizes colors for player backgrounds by applying tinting and brightness adjustments
 * @param colors - Array of RGB colors to harmonize
 * @param tintAmount - Amount of white tint to apply (0-1). Higher = more pastel
 * @param minBrightness - Minimum brightness threshold (0-255). Colors darker than this will be boosted
 * @returns Array of harmonized RGB colors
 */
export function harmonizeColors(colors: RGB[], tintAmount: number, minBrightness: number): RGB[] {
  return colors.map(c => {
    // Apply tint to make colors lighter/more pastel
    const r = Math.round(c.r + (255 - c.r) * tintAmount);
    const g = Math.round(c.g + (255 - c.g) * tintAmount);
    const b = Math.round(c.b + (255 - c.b) * tintAmount);

    // Calculate brightness (average of RGB components)
    const brightness = (r + g + b) / 3;

    // Boost dark colors to meet minimum brightness threshold
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
}