<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';
  import ResultSection from '../lib/components/ResultSection.svelte';
  import type { AnalyzeResult, ScreenshotResponse } from '../lib/types';
  import { isValidUrl, formatUrl } from '../lib/utils';
  import { replaceState } from '$app/navigation';


  let input = '';
  let loading = false;
  let error: string | null = null;
  let data: AnalyzeResult | null = null;
  let screenshotData: ScreenshotResponse | null = null;
  let screenshotUrl: string | null = null;

  type Verdict = 'Safe' | 'Risky' | 'Unclear' | 'Unknown';
  const ACCENTS: Record<Verdict, { ring: string; glow: string; badge: string }> = {
    Safe: { ring: 'focus:ring-emerald-600', glow: 'from-emerald-600/20', badge: 'bg-emerald-600/20 text-emerald-300 border-emerald-700' },
    Risky: { ring: 'focus:ring-red-600', glow: 'from-red-600/20', badge: 'bg-red-600/20 text-red-300 border-red-700' },
    Unclear: { ring: 'focus:ring-yellow-600', glow: 'from-yellow-500/20', badge: 'bg-yellow-600/20 text-yellow-300 border-yellow-700' },
    Unknown: { ring: 'focus:ring-yellow-600', glow: 'from-yellow-500/20', badge: 'bg-yellow-600/20 text-yellow-300 border-yellow-700' }
  };

  function normalizeVerdict(v: string | null | undefined): Verdict {
    switch (v) {
      case 'Safe':
      case 'Risky':
      case 'Unclear':
      case 'Unknown':
        return v;
      default:
        return 'Unknown';
    }
  }

  $: verdict = normalizeVerdict(data?.result?.verdict);
  $: accent = ACCENTS[verdict];
  $: isLanding = !data && !loading && !error;


  function buildScreenshotUrl(targetUrl: string): string | null {
    // If backend exposes screenshot, define pattern here later; placeholder for now
    // e.g., `${PUBLIC_BASE_URL.replace(/\/api\/v1$/, '')}/api/v1/screenshot?url=...`
    return null;
  }

  async function runAnalyze(q: string) {
  const url = formatUrl(q);
  if (!isValidUrl(url)) {
    error = 'Please enter a valid URL';
    return;
  }

  loading = true;
  error = null;
  data = null;
  screenshotData = null;

  try {
    // kick off screenshot but don't block on it
    api.screenshot(url)
      .then((res) => {
        screenshotData = res as ScreenshotResponse;
      })
      .catch(() => {
        console.warn('Screenshot request failed');
      });

    // await analyze so loading reflects this call only
    const res = await api.analyze(url);

    if (res.error) {
      error = res.error;
    } else {
      data = res.data as AnalyzeResult;

      const share = new URL(window.location.href);
      share.searchParams.set('q', url);
      replaceState(share.toString(), {});
    }
  } catch (err) {
    error = 'Analyze request failed';
  } finally {
    loading = false;
  }
}





  function onSubmit(e: Event) {
    e.preventDefault();
    runAnalyze(input);
  }

  onMount(() => {
    const params = new URLSearchParams(window.location.search);
    const q = params.get('q');
    if (q) {
      input = q;
      runAnalyze(q);
    }

    const onKey = (e: KeyboardEvent) => {
      if (e.key === '/' && !(e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement)) {
        e.preventDefault();
        const el = document.getElementById('url-input') as HTMLInputElement | null;
        el?.focus();
      }
      if (e.key === 'Escape') {
        input = '';
      }
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  });
</script>

<section>
  <title>SafeSurf</title>
  <div class={`max-w-4xl mx-auto px-6 ${isLanding ? 'min-h-[70vh] flex flex-col justify-center' : 'py-12'}`}>
    <header class="mb-10">
      <h1 class="text-3xl md:text-5xl font-semibold tracking-tight text-white">SafeSurf</h1>
      <p class="mt-3 text-gray-400 text-base">Surf safe with SafeSurf.</p>
    </header>

    <form class="bg-gray-950 rounded-xl border border-gray-800 p-4 md:p-5" on:submit|preventDefault={onSubmit}>
      <label for="url-input" class="sr-only">URL to analyze</label>
      <div class="flex flex-col md:flex-row gap-2">
        <input
          id="url-input"
          type="text"
          class={`flex-1 rounded-lg bg-gray-900 border border-gray-800 px-4 py-3 text-sm placeholder-gray-500 text-gray-200 focus:outline-none focus:ring-2 ${accent.ring}`}
          placeholder="Enter a URL (e.g. google.com)"
          bind:value={input}
          autocomplete="url"
          inputmode="url"
          required
        />

        

        <button
          type="submit"
          class="inline-flex items-center justify-center gap-2 px-5 py-3 rounded-lg bg-blue-600 hover:bg-blue-500 text-white text-sm font-medium disabled:opacity-50 focus:outline-none focus:ring-2 focus:ring-offset-0 focus:ring-blue-600"
          disabled={loading}
          aria-busy={loading}
        >{loading ? 'Analyzing…' : 'Analyze'}</button>
      </div>
      <!-- <p class="mt-2 text-xs text-gray-500">Include http:// or https:// in the URL.</p> -->
    </form>

    <!-- {#if data}
      <div class="mt-8 flex items-center justify-between bg-gray-950 border border-gray-800 rounded-lg p-4">
        <div>
          <div class="text-[11px] text-gray-400">Analyzed</div>
          <div class="text-sm text-gray-200 truncate max-w-[60vw]">{data.url}</div>
        </div>
        <div class="flex items-center gap-2">
          <span class={`text-[10px] uppercase tracking-wide px-2 py-0.5 rounded-full border ${accent.badge}`}>{verdict || '—'}</span>
          <span class="text-xs text-gray-400">Score</span>
          <div class="w-24 h-2 bg-gray-900 rounded overflow-hidden"><div class="h-2 bg-gray-200" style={`width:${Math.min(100, data.result?.final_score ?? 0)}%`}></div></div>
        </div>
      </div>
    {/if} -->

    <div class="mt-8" aria-live="polite">
      <ResultSection {data} {loading} {error} {screenshotData} />
    </div>
  </div>
</section>

