/**
 * Extracts dominant colors from an image for background gradients
 */

interface RGB {
  r: number;
  g: number;
  b: number;
}

/**
 * Load an image and extract its dominant colors
 */
export async function extractColorsFromImage(imageUrl: string): Promise<RGB[]> {
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.crossOrigin = 'anonymous';

    img.onload = () => {
      try {
        const colors = getColorPalette(img);
        resolve(colors);
      } catch (err) {
        reject(err);
      }
    };

    img.onerror = () => reject(new Error('Failed to load image'));
    img.src = imageUrl;
  });
}

/**
 * Calculate color saturation (0-1)
 */
function getSaturation(rgb: RGB): number {
  const r = rgb.r / 255;
  const g = rgb.g / 255;
  const b = rgb.b / 255;

  const max = Math.max(r, g, b);
  const min = Math.min(r, g, b);
  const delta = max - min;

  if (max === 0) return 0;
  return delta / max; // HSV saturation
}

/**
 * Extract color palette from an image - prioritizes vibrant, saturated colors
 */
function getColorPalette(img: HTMLImageElement): RGB[] {
  const canvas = document.createElement('canvas');
  const ctx = canvas.getContext('2d');

  if (!ctx) {
    throw new Error('Could not get canvas context');
  }

  // Resize for performance
  const size = 100;
  canvas.width = size;
  canvas.height = size;

  ctx.drawImage(img, 0, 0, size, size);
  const imageData = ctx.getImageData(0, 0, size, size);
  const pixels = imageData.data;

  // Sample colors with saturation data
  const colorCounts = new Map<string, { color: RGB; count: number; saturation: number }>();

  for (let i = 0; i < pixels.length; i += 4 * 10) { // Sample every 10th pixel
    const r = pixels[i];
    const g = pixels[i + 1];
    const b = pixels[i + 2];
    const a = pixels[i + 3];

    // Skip transparent pixels
    if (a < 128) continue;

    // Quantize colors to reduce variations
    const qr = Math.round(r / 30) * 30;
    const qg = Math.round(g / 30) * 30;
    const qb = Math.round(b / 30) * 30;

    const color = { r: qr, g: qg, b: qb };
    const saturation = getSaturation(color);

    // Skip very desaturated colors (grays, near-whites, near-blacks)
    if (saturation < 0.2) continue;

    const key = `${qr},${qg},${qb}`;
    const existing = colorCounts.get(key);

    if (existing) {
      existing.count++;
    } else {
      colorCounts.set(key, { color, count: 1, saturation });
    }
  }

  // Score colors by both frequency AND saturation (favor vibrant colors)
  const sorted = Array.from(colorCounts.values())
    .map(item => ({
      ...item,
      score: item.count * (1 + item.saturation * 2) // Saturation weighted heavily
    }))
    .sort((a, b) => b.score - a.score)
    .slice(0, 3)
    .map(item => item.color);

  return sorted.length > 0 ? sorted : [
    { r: 255, g: 214, b: 214 },
    { r: 232, g: 213, b: 255 },
    { r: 213, g: 232, b: 255 }
  ];
}

 /**
 * Lighten and soften a color to harmonize with the UI
 */
function harmonizeColor(color: RGB): RGB {
  // Mix with white (tint) for a lighter appearance - reduced for more vibrant colors
  const tintAmount = 0.35; // 35% white mix - more color saturation
  const r = Math.round(color.r + (255 - color.r) * tintAmount);
  const g = Math.round(color.g + (255 - color.g) * tintAmount);
  const b = Math.round(color.b + (255 - color.b) * tintAmount);

  // Ensure minimum brightness - reduced to allow more vibrant colors
  const minBrightness = 180;
  const brightness = (r + g + b) / 3;

  if (brightness < minBrightness) {
    const boost = minBrightness - brightness;
    return {
      r: Math.min(255, r + boost),
      g: Math.min(255, g + boost),
      b: Math.min(255, b + boost)
    };
  }

  return { r, g, b };
}

// Store current colors for smooth interpolation
let currentColors: RGB[] = [
  { r: 255, g: 214, b: 214 },
  { r: 232, g: 213, b: 255 },
  { r: 213, g: 232, b: 255 }
];

let animationFrameId: number | null = null;

/**
 * Interpolate between two colors
 */
function lerpColor(from: RGB, to: RGB, t: number): RGB {
  return {
    r: Math.round(from.r + (to.r - from.r) * t),
    g: Math.round(from.g + (to.g - from.g) * t),
    b: Math.round(from.b + (to.b - from.b) * t)
  };
}

/**
 * Easing function for smooth animation
 */
function easeInOutCubic(t: number): number {
  return t < 0.5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2;
}

/**
 * Apply gradient to document with smooth animation
 */
export function applyGradient(colors: RGB[]) {
  if (colors.length < 3) return;

  const targetColors = [
    harmonizeColor(colors[0]),
    harmonizeColor(colors[1]),
    harmonizeColor(colors[2])
  ];

  // Cancel any ongoing animation
  if (animationFrameId !== null) {
    cancelAnimationFrame(animationFrameId);
  }

  const startColors = [...currentColors];
  const duration = 1500; // 1.5 seconds
  const startTime = performance.now();

  function animate(currentTime: number) {
    const elapsed = currentTime - startTime;
    const progress = Math.min(elapsed / duration, 1);
    const easedProgress = easeInOutCubic(progress);

    // Interpolate all three colors
    const interpolated = [
      lerpColor(startColors[0], targetColors[0], easedProgress),
      lerpColor(startColors[1], targetColors[1], easedProgress),
      lerpColor(startColors[2], targetColors[2], easedProgress)
    ];

    // Update CSS variables
    const root = document.documentElement;
    root.style.setProperty('--gradient-color-1', `${interpolated[0].r}, ${interpolated[0].g}, ${interpolated[0].b}`);
    root.style.setProperty('--gradient-color-2', `${interpolated[1].r}, ${interpolated[1].g}, ${interpolated[1].b}`);
    root.style.setProperty('--gradient-color-3', `${interpolated[2].r}, ${interpolated[2].g}, ${interpolated[2].b}`);

    if (progress < 1) {
      animationFrameId = requestAnimationFrame(animate);
    } else {
      // Animation complete - store final colors
      currentColors = targetColors;
      animationFrameId = null;
    }
  }

  animationFrameId = requestAnimationFrame(animate);
}

/**
 * Reset to default gradient
 */
export function resetGradient() {
  const defaultColors = [
    { r: 255, g: 214, b: 214 },
    { r: 232, g: 213, b: 255 },
    { r: 213, g: 232, b: 255 }
  ];

  applyGradient(defaultColors);
}