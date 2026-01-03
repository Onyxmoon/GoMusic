# GoMusic

> **Work in Progress** - This project is still in active development

A prototype desktop music player built with Go and Svelte. Just me experimenting

## Stack

- **Backend**: Go 1.23 + Wails v2 (hopefully v3 sppm)
- **Frontend**: Svelte 5 + TypeScript
- **Build**: Vite

## Current Features

- Local music library scanning
- Audio playback with HTML5 Audio API
- Track metadata extraction (MP3, FLAC, M4A, OGG)
- Album artwork display
- Keyboard shortcuts for playback control
- Preview mode for quick track browsing
- Logarithmic volume curve
- Track search and filtering

## Development

```bash
# Install dependencies
npm install

# Run in dev mode
wails dev

# Build for production
wails build
```

## License

Mine