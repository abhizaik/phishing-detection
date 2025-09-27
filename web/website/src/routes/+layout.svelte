<script>
  import '../app.css';
  const year = new Date().getFullYear();
</script>

<div class="min-h-screen flex flex-col">
  <main class="flex-1">
    <slot />
  </main>
  <footer class="border-t border-gray-800 text-gray-400 text-xs py-6">
    <div class="max-w-4xl mx-auto px-6 flex flex-col md:flex-row items-center justify-between gap-2">
      <a href="https://github.com/abhizaik/SafeSurf" target="_blank" rel="noopener noreferrer" class="hover:text-gray-300">
        View source at github.com/abhizaik/SafeSurf
      </a>
      <div class="text-center md:text-right">
        SafeSurf is licensed under MIT — © abhizaik {year}
      </div>
    </div>
  </footer>
</div>