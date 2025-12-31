<script lang="ts">
  export let clickable: boolean = false;
</script>

<div
  class="card"
  class:clickable
  on:click
  role={clickable ? 'button' : undefined}
  tabindex={clickable ? 0 : undefined}
  on:keydown={(e) => clickable && e.key === 'Enter' && e.currentTarget.click()}
>
  <slot />
</div>

<style>
  .card {
    background: #2a2e39;
    border-radius: 8px;
    padding: 16px;
    border: 1px solid #3a3f4b;
  }

  .card.clickable {
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .card.clickable:hover {
    background: #2f3441;
    border-color: #4a5060;
    transform: translateY(-2px);
  }
</style>