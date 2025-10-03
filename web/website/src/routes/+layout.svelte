<script>
  import '../app.css';
  const year = new Date().getFullYear();
</script>

<div class="min-h-screen flex flex-col">
  <main class="flex-1">
    <slot />
  </main>

  <footer class="bg-gray-900 text-gray-300 border-t border-gray-800 py-6">
    <div class="max-w-4xl mx-auto px-6 flex flex-col md:flex-row md:justify-between items-center gap-3 text-sm md:text-base">

      <!-- GitHub link -->
      <p class="text-gray-400 flex items-center gap-2">
        <svg class="w-5 h-5 text-gray-400" fill="currentColor" viewBox="0 0 24 24">
          <path d="M12 .5C5.648.5.5 5.648.5 12c0 5.084 3.292 9.387 7.872 10.905.576.106.784-.25.784-.556 0-.274-.01-1-.016-1.974-3.2.694-3.876-1.544-3.876-1.544-.522-1.33-1.276-1.683-1.276-1.683-1.042-.714.08-.699.08-.699 1.152.08 1.756 1.183 1.756 1.183 1.024 1.754 2.688 1.246 3.344.952.104-.742.4-1.246.728-1.532-2.554-.292-5.238-1.276-5.238-5.674 0-1.252.448-2.274 1.184-3.076-.12-.292-.512-1.468.112-3.064 0 0 .964-.308 3.16 1.176.916-.254 1.9-.382 2.876-.388.976.006 1.96.134 2.876.388 2.192-1.484 3.156-1.176 3.156-1.176.624 1.596.232 2.772.112 3.064.736.802 1.184 1.824 1.184 3.076 0 4.408-2.69 5.38-5.25 5.664.412.354.78 1.05.78 2.118 0 1.532-.014 2.768-.014 3.144 0 .308.2.672.788.556C20.708 21.385 24 17.084 24 12c0-6.352-5.148-11.5-12-11.5z"/>
        </svg>
        <a href="https://github.com/abhizaik/SafeSurf" target="_blank" class="text-gray-300 hover:text-white transition-colors">
          View source on GitHub
        </a>
      </p>

      <!-- License and author -->
      <p class="text-gray-400 text-center md:text-right">
        SafeSurf is licensed under 
        <a href="https://github.com/abhizaik/SafeSurf/blob/main/LICENSE" target="_blank" class="text-gray-300 hover:text-white transition-colors">
          MIT
        </a> — © 
        <a href="https://abhizaik.com" target="_blank" class="text-gray-300 hover:text-white transition-colors">
          abhizaik
        </a> {year}
      </p>
    </div>
  </footer>
</div>
