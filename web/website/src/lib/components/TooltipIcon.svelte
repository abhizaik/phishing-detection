<script lang="ts">
  export let text: string;
</script>

<div class="relative inline-flex items-center group">
  <!-- Info icon -->
  <div
    class="w-4 h-4 flex items-center justify-center rounded-full
           bg-gray-700 text-gray-200 text-[10px] font-bold cursor-pointer
           hover:bg-gray-600 transition-colors duration-150 select-none"
  >
    i
  </div>

  <!-- Tooltip bubble -->
  <div
    class="absolute left-1/2 bottom-full mb-2 w-max max-w-xs -translate-x-1/2
           opacity-0 group-hover:opacity-100 translate-y-1 group-hover:translate-y-0
           bg-gray-800 text-gray-100 text-xs px-3 py-1.5 rounded-lg shadow-lg
           border border-gray-700 transition-all duration-200 ease-out
           pointer-events-none z-50 whitespace-normal"
  >
    {text}
    <div
      class="absolute left-1/2 top-full -translate-x-1/2 w-2 h-2
             bg-gray-800 border-b border-r border-gray-700 rotate-45"
    ></div>
  </div>
</div>
