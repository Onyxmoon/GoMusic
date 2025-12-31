<script lang="ts">
  export let variant: 'primary' | 'secondary' | 'ghost' | 'danger' = 'primary';
  export let size: 'sm' | 'md' | 'lg' = 'md';
  export let disabled: boolean = false;
  export let fullWidth: boolean = false;
</script>

<button
  class="btn btn-{variant} btn-{size}"
  class:btn-full={fullWidth}
  {disabled}
  on:click
>
  <slot />
</button>

<style>
  .btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border: none;
    border-radius: 6px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    font-family: inherit;
  }

  .btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  /* Sizes */
  .btn-sm {
    padding: 6px 12px;
    font-size: 13px;
  }

  .btn-md {
    padding: 8px 16px;
    font-size: 14px;
  }

  .btn-lg {
    padding: 12px 24px;
    font-size: 16px;
  }

  /* Variants */
  .btn-primary {
    background: #5b8cff;
    color: white;
  }

  .btn-primary:hover:not(:disabled) {
    background: #4a7de8;
  }

  .btn-secondary {
    background: #3a3f4b;
    color: white;
  }

  .btn-secondary:hover:not(:disabled) {
    background: #4a5060;
  }

  .btn-ghost {
    background: transparent;
    color: #c4c7d0;
  }

  .btn-ghost:hover:not(:disabled) {
    background: rgba(255, 255, 255, 0.05);
  }

  .btn-danger {
    background: #dc2626;
    color: white;
  }

  .btn-danger:hover:not(:disabled) {
    background: #b91c1c;
  }

  .btn-full {
    width: 100%;
  }
</style>