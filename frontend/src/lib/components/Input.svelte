<script lang="ts">
  export let value: string = '';
  export let placeholder: string = '';
  export let disabled: boolean = false;
</script>

<input
  bind:value
  type="text"
  {placeholder}
  {disabled}
  class="input"
  on:input
  on:focus
  on:blur
/>

<style>
  .input {
    width: 100%;
    padding: 10px 14px;
    background: #2a2e39;
    border: 1px solid #3a3f4b;
    border-radius: 6px;
    color: #e4e6eb;
    font-size: 14px;
    font-family: inherit;
    outline: none;
    transition: all 0.2s ease;
  }

  .input::placeholder {
    color: #6b7280;
  }

  .input:focus {
    border-color: #5b8cff;
    background: #2f3441;
  }

  .input:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>