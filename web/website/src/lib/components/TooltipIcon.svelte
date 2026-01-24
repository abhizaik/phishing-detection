<script lang="ts">
  export let text: string;
  let bubble: HTMLDivElement;
  let arrow: HTMLDivElement;

  function adjustPosition() {
    if (!bubble || !arrow) return;

    // Reset to calculate natural position
    bubble.style.transform = "translateX(-50%)";
    arrow.style.left = "50%";

    const rect = bubble.getBoundingClientRect();
    const vw = window.innerWidth;
    const padding = 10;

    let shift = 0;
    if (rect.left < padding) {
      shift = padding - rect.left;
    } else if (rect.right > vw - padding) {
      shift = vw - padding - rect.right;
    }

    if (shift !== 0) {
      bubble.style.transform = `translateX(calc(-50% + ${shift}px))`;
      // Keep arrow over the icon but within bubble bounds
      const arrowShift = -shift;
      const halfWidth = rect.width / 2;
      // 12px padding from edges of bubble
      const safeArrowShift = Math.max(-halfWidth + 12, Math.min(halfWidth - 12, arrowShift));
      arrow.style.left = `calc(50% + ${safeArrowShift}px)`;
    }
  }
</script>

<div
  class="relative inline-flex items-center group z-[99999]"
  on:mouseenter={adjustPosition}
>
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
    bind:this={bubble}
    class="absolute left-1/2 bottom-full mb-2 w-max max-w-[85vw] md:max-w-xs -translate-x-1/2
           opacity-0 group-hover:opacity-100 translate-y-1 group-hover:translate-y-0
           bg-gray-800 text-gray-100 text-xs px-3 py-1.5 rounded-lg shadow-lg
           border border-gray-700 transition-all duration-200 ease-out
           pointer-events-none z-[99999] whitespace-normal"
  >
    {text}
    <div
      bind:this={arrow}
      class="absolute left-1/2 top-full -translate-x-1/2 w-2 h-2
             bg-gray-800 border-b border-r border-gray-700 rotate-45"
    ></div>
  </div>
</div>
