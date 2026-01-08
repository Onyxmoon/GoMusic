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
