<script lang="ts">
    import type {dto} from '../../../wailsjs/go/models';
    import {MoreVertical, Play} from 'lucide-svelte';
    import NowPlayingIndicator from './NowPlayingIndicator.svelte';
    import {player} from '../stores/player.svelte';

    interface Props {
        tracks: dto.TrackDTO[];
        onTrackClick?: (track: dto.TrackDTO) => void;
    }

    let {
        tracks = [], onTrackClick = () => {
        }
    }: Props = $props();

    let selectedIndex = $state(0);

    function keyboardNav(node: HTMLElement) {
        node.tabIndex = 0;

        const handleKeyDown = (e: KeyboardEvent) => {
            if (!tracks.length) return;

            switch (e.key) {
                case 'ArrowDown':
                    e.preventDefault();
                    selectedIndex = Math.min(selectedIndex + 1, tracks.length - 1);
                    break;
                case 'ArrowUp':
                    e.preventDefault();
                    selectedIndex = Math.max(selectedIndex - 1, 0);
                    break;
                case 'Enter':
                    e.preventDefault();
                    onTrackClick(tracks[selectedIndex]);
                    break;
            }
        };

        node.addEventListener('keydown', handleKeyDown);
        return {destroy: () => node.removeEventListener('keydown', handleKeyDown)};
    }

    const formatDuration = (seconds: number): string => {
        const mins = Math.floor(seconds / 60);
        const secs = Math.floor(seconds % 60);
        return `${mins}:${secs.toString().padStart(2, '0')}`;
    };

    const handleClick = (track: dto.TrackDTO, index: number) => {
        selectedIndex = index;
        onTrackClick(track);
    };
</script>

<div class="track-table">
    <div class="table-header">
        <div class="col col-number">#</div>
        <div class="col col-title">Title</div>
        <div class="col col-album">Album</div>
        <div class="col col-time">Time</div>
    </div>

    <div class="table-body" use:keyboardNav>
        {#each tracks as track, index}
            {@const {isCurrentTrack, isPlaying} = player.isTrackPlaying(track.id)}
            {@const isPaused = isCurrentTrack && !isPlaying}
            <div
                    class="track-row"
                    class:selected={selectedIndex === index}
                    onclick={() => handleClick(track, index)}
                    role="button"
                    tabindex="-1"
            >
                <div class="col col-number" class:playing={isPlaying} class:paused={isPaused}>
                    <NowPlayingIndicator trackId={track.id} size={14}/>
                    <span class="track-number" class:hide={isCurrentTrack}>{index + 1}</span>
                    <span class="play-icon"><Play size={14}/></span>
                </div>

                <div class="col col-title">
                    <div class="track-cover">
                        {#if track.hasArtwork}
                            <img src={`/artwork/stream?id=${encodeURIComponent(track.id)}`} alt={track.title}/>
                        {:else}
                            <div class="cover-placeholder"></div>
                        {/if}
                    </div>
                    <div class="title-info">
                        <div class="title">{track.title}</div>
                        <div class="artist">{track.artist}</div>
                    </div>
                </div>

                <div class="col col-album">{track.album}</div>
                <div class="col col-time">
                    {formatDuration(track.duration)}
                    <button class="more-btn" onclick={(e) => e.stopPropagation()}>
                        <MoreVertical size={16}/>
                    </button>
                </div>
            </div>
        {/each}
    </div>
</div>

<style>
    /* Layout */
    .track-table {
        width: 100%;
    }

    .table-header {
        display: grid;
        grid-template-columns: 50px minmax(200px, 2fr) minmax(120px, 1fr) 100px;
        gap: 16px;
        padding: 12px 16px;
        border-bottom: 1px solid rgba(0, 0, 0, 0.06);
        font-size: 12px;
        font-weight: 500;
        color: #6b7280;
    }

    .table-body {
        display: flex;
        flex-direction: column;
        outline: none;
    }

    /* Track Row */
    .track-row {
        display: grid;
        grid-template-columns: 50px minmax(200px, 2fr) minmax(120px, 1fr) 100px;
        gap: 16px;
        padding: 10px 16px;
        margin-top: 5px;
        background: transparent;
        border: 1px transparent solid;
        border-radius: 8px;
        cursor: pointer;
        text-align: left;
        transition: all 0.15s ease;
        outline: none;
    }

    .track-row:hover {
        background: rgba(255, 255, 255, 0.5);
    }

    .track-row.selected {
        background: rgba(255, 255, 255, 0.3);
        border: 1px rgba(60, 63, 71, 0.15) solid;
    }

    .track-row.selected:hover {
        background: rgba(255, 255, 255, 0.6);
    }

    /* Columns Base */
    .col {
        display: flex;
        align-items: center;
        color: #2d2d2d;
        font-size: 14px;
    }

    /* Column: Number */
    .col-number {
        position: relative;
        justify-content: center;
        color: #6b7280;
        font-weight: 500;
    }

    .track-number {
        transition: opacity 0.15s ease;
    }

    .track-number.hide,
    .track-row:hover .track-number {
        opacity: 0;
    }

    .play-icon {
        position: absolute;
        opacity: 0;
        transition: opacity 0.15s ease;
        color: #2d2d2d;
    }

    .track-row:hover .play-icon {
        opacity: 1;
    }

    .track-row:hover .col-number.playing .play-icon {
        opacity: 0;
    }

    .track-row:hover .col-number.paused :global(.now-playing-indicator) {
        opacity: 0;
    }

    /* Column: Title */
    .col-title {
        gap: 12px;
    }

    .track-cover {
        width: 40px;
        height: 40px;
        border-radius: 6px;
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
        font-size: 18px;
    }

    .title-info {
        flex: 1;
        min-width: 0;
    }

    .title {
        font-weight: 500;
        color: #2d2d2d;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .artist {
        font-size: 12px;
        color: #6b7280;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    /* Column: Album */
    .col-album {
        color: #6b7280;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        min-width: 0;
    }

    /* Column: Time */
    .col-time {
        justify-content: space-between;
        color: #6b7280;
    }

    .more-btn {
        width: 28px;
        height: 28px;
        border: none;
        background: transparent;
        border-radius: 4px;
        color: #6b7280;
        cursor: pointer;
        opacity: 0;
        transition: all 0.15s ease;
    }

    .track-row:hover .more-btn {
        opacity: 1;
    }

    .more-btn:hover {
        background: rgba(0, 0, 0, 0.05);
        color: #2d2d2d;
    }
</style>