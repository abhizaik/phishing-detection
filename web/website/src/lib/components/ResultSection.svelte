<script lang="ts">
  import type { AnalyzeResult, ScreenshotResponse } from '../types';
  import { browser } from '$app/environment';
  export let data: AnalyzeResult | null = null;
  export let screenshotData: ScreenshotResponse | null = null;
  export let loading = false;
  export let error: string | null = null;

  let showAdvanced = false;
  $: primary = data?.result;
  $: reasons = primary?.reasons;

  let copied = false;
  let showModal = false;

  async function copyShareLink() { // Function for first copy button, commented out as we added new Share Button
    try {
      const url = new URL(window.location.href);
      if (data?.url) {
        url.searchParams.set('q', data.url);
        await navigator.clipboard.writeText(url.toString());
        copied = true;
        setTimeout(() => (copied = false), 1200);
      }
    } catch {
      error = "Could not copy link";
    }
  }


  let shareCopied = false;
  let shareUrl = '';
  if (browser) {
    shareUrl = window.location.href; // to only run on client-side
  }

  async function shareLink() {
    const shareText = `Check out this SafeSurf result for ${data?.domain}`;
    if (browser && navigator.share) {
      try {
        await navigator.share({
          title: "SafeSurf",
          text: shareText,
          url: shareUrl,
        });
      } catch (err) {
        console.error("Share failed:", err);
      }
    } else if (browser) {
      try {
        await navigator.clipboard.writeText(`${shareText} \n${shareUrl}`);
        shareCopied = true;
        setTimeout(() => (shareCopied = false), 1200);
      } catch (err) {
        console.error("Clipboard copy failed:", err);
      }
    }
  }

  function toggleAdvanced() {
    showAdvanced = !showAdvanced;
  }
</script>

{#if error}
  <div class="max-w-3xl mx-auto p-4 bg-red-900/30 border border-red-700 text-red-200 rounded-md">{error}</div>
{:else if loading}
  <div class="max-w-3xl mx-auto space-y-4">
    <div class="animate-pulse rounded-xl border border-gray-800 bg-gray-950/60 p-6 h-32"></div>
    <div class="animate-pulse rounded-xl border border-gray-800 bg-gray-950/60 p-6 h-24"></div>
  </div>
{:else if data}
<section class="max-w-4xl mx-auto space-y-8 px-4">
  <!-- Header & Copy Button -->
  <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-3">
  <!-- Title + Paragraph -->
  <div class="flex flex-col">
    <h2 class="text-2xl font-semibold text-white">Analysis Summary</h2>
    <p class="text-gray-400 text-sm mt-1">Check the verdict, score, and flags for {data?.domain}</p>
  </div>

  <!-- Copy Button -->
  <!-- <button
    class="inline-flex items-center gap-2 px-5 py-3 rounded-full bg-gray-800 hover:bg-gray-700 text-white text-sm font-medium transition-all {copied ? 'animate-pulse bg-emerald-700' : ''}"
    on:click={copyShareLink}
  >
    {#if copied}
      <svg class="w-4 h-4 text-emerald-300" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
      </svg>
      <span class="text-emerald-300">Copied!</span>
    {:else}
      <svg class="w-4 h-4 text-gray-300" fill="currentColor" viewBox="0 0 20 20">
        <path d="M8 2a2 2 0 00-2 2v2h2V4h6v6h-2v2h2a2 2 0 002-2V4a2 2 0 00-2-2H8zM4 8a2 2 0 00-2 2v6a2 2 0 002 2h6a2 2 0 002-2v-6a2 2 0 00-2-2H4zm0 2h6v6H4v-6z"/>
      </svg>
      <span>Copy Result</span>
    {/if}
  </button> -->


<!-- Share Button -->
    <button
      class="inline-flex items-center gap-2 px-5 py-3 rounded-full bg-gray-800 hover:bg-gray-700 text-white text-sm font-medium transition-all {shareCopied ? 'animate-pulse bg-emerald-700' : ''}"
      on:click={shareLink}
    >
      {#if shareCopied}
        <svg class="w-4 h-4 text-emerald-300" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
        </svg>
        <span class="text-emerald-300">Copied to clipboard!</span>
      {:else}
        <svg class="w-4 h-4 text-gray-300" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M4 12v7a2 2 0 002 2h12a2 2 0 002-2v-7M16 6l-4-4-4 4m4-4v14" />
        </svg>

        <span>Share Result</span>
      {/if}
    </button>

</div>


  <!-- Verdict & Trust Score -->
  <div class="flex flex-col md:flex-row gap-6 p-6 bg-gray-900/80 rounded-xl shadow-md hover:shadow-lg transition-transform hover:scale-[1.01]">
    <div class="flex-1">
      <div class="text-sm font-medium text-gray-300 uppercase tracking-wide mb-1">Verdict</div>
      <div class="flex items-center gap-3">
        <span class="text-2xl font-bold text-white">{primary?.verdict ?? '-'}</span>
         {#if primary?.verdict === 'Safe'}
    <span class="px-3 py-1 rounded-full bg-green-700 text-white font-medium text-xs uppercase tracking-wide">
      Trusted
    </span>
  {:else if primary?.verdict === 'Suspicious'}
    <span class="px-3 py-1 rounded-full bg-yellow-500 text-black font-medium text-xs uppercase tracking-wide">
      Be Cautious
    </span>
  {:else if primary?.verdict === 'Risky'}
    <span class="px-3 py-1 rounded-full bg-red-700 text-white font-medium text-xs uppercase tracking-wide">
      High Risk
    </span>
  {:else if primary?.verdict === 'Unclear'}
    <span class="px-3 py-1 rounded-full bg-gray-500 text-white font-medium text-xs uppercase tracking-wide">
      Not Enough Data
    </span>
  {:else}
    <span class="px-3 py-1 rounded-full bg-red-600 text-white font-medium text-xs uppercase tracking-wide">
      Dangerous
    </span>
  {/if}
      </div>
      <!-- Trust Score Percentage Bar -->
      <!-- <div class="mt-3 h-2 w-full bg-gray-800 rounded-full overflow-hidden">
        <div
          class="h-2 bg-blue-500 rounded-full transition-all duration-700 ease-out"
          style="width:{primary?.final_score ?? 0}%"
        ></div>
      </div> -->
    </div>

    <div class="flex-1 md:text-right flex flex-col justify-center">
      <div class="text-sm font-medium text-gray-300 uppercase tracking-wide mb-1">Trust Score</div>
      <span class="text-3xl font-extrabold text-white-50">{primary?.final_score ?? '-'} / 100</span>
    </div>
  </div>

  <!-- Flags -->
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    {#if reasons}
      <!-- Red Flags -->
      <div class="rounded-xl border border-red-700 bg-red-900/20 p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all">
        <div class="flex items-center gap-3 mb-3">
          <svg class="w-5 h-5 text-red-500 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.721-1.36 3.486 0l6.518 11.59c.75 1.335-.213 3.011-1.743 3.011H3.482c-1.53 0-2.493-1.676-1.743-3.01L8.257 3.1zM11 14a1 1 0 10-2 0 1 1 0 002 0zm-1-2a.75.75 0 01-.75-.75V8a.75.75 0 011.5 0v3.25A.75.75 0 0110 12z" clip-rule="evenodd" />
          </svg>
          <h3 class="text-sm font-semibold text-red-400 uppercase tracking-wide">Red Flags</h3>
        </div>
        {#if reasons.bad_reasons?.length}
          <ul class="space-y-2 text-red-200 text-sm">
            {#each reasons.bad_reasons as r}
              <li class="flex items-start gap-2" title="Potential risk">
                <span class="mt-1 h-2 w-2 rounded-full bg-red-500 flex-shrink-0"></span>
                <span class="break">{r}</span>
              </li>
            {/each}
          </ul>
        {:else}
          <p class="text-gray-400 text-sm">No red flags found.</p>
        {/if}
      </div>

      <!-- Green Flags -->
      <div class="rounded-xl border border-emerald-700 bg-emerald-900/20 p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all">
        <div class="flex items-center gap-3 mb-3">
          <svg class="w-5 h-5 text-emerald-400 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M16.704 5.29a1 1 0 010 1.42l-7.388 7.388a1 1 0 01-1.42 0L3.296 9.498a1 1 0 111.408-1.42L8.5 11.874l6.796-6.795a1 1 0 011.408 0z" clip-rule="evenodd" />
          </svg>
          <h3 class="text-sm font-semibold text-emerald-400 uppercase tracking-wide">Green Flags</h3>
        </div>
        {#if reasons.good_reasons?.length}
          <ul class="space-y-2 text-emerald-200 text-sm">
            {#each reasons.good_reasons as r}
              <li class="flex items-start gap-2" title="Positive sign">
                <span class="mt-1 h-2 w-2 rounded-full bg-emerald-400 flex-shrink-0"></span>
                <span class="break">{r}</span>
              </li>
            {/each}
          </ul>
        {:else}
          <p class="text-gray-400 text-sm">No green flags found.</p>
        {/if}
      </div>
    {/if}
  </div>

  <!-- Screenshot -->
  {#if screenshotData?.status === 'success' && screenshotData?.file && !screenshotData.file.startsWith('server')}
    <div class="mt-6 rounded-xl border border-gray-800 bg-gray-900/70 p-4 shadow-md hover:shadow-lg transition-all">
      <h4 class="text-sm font-semibold text-gray-300 mb-2">Website Screenshot</h4>
      <!-- Image for testing -->
      <!-- src="screenshot-google-com.png" -->
      <img
      src={screenshotData.file}
        alt="Website screenshot"
        class="w-full rounded-lg border border-gray-800 cursor-pointer hover:opacity-90"
        loading="lazy"
        on:click={() => showModal = true}
      />
    </div>
  {/if}

  <!-- Screenshot Modal -->
  {#if showModal}
    <div class="fixed inset-0 bg-black/80 flex items-center justify-center z-50">
      <button class="absolute top-4 right-4 text-gray-300 hover:text-white text-2xl" on:click={() => showModal = false}>×</button>
      <img src={screenshotData?.file} alt="Full screenshot" class="max-h-[90vh] max-w-[90vw] rounded-lg shadow-lg" />
      <!-- <img src=screenshot-google-com.png alt="Full screenshot" class="max-h-[90vh] max-w-[90vw] rounded-lg shadow-lg" /> -->
    </div>
  {/if}

  <!-- Advanced Panel Toggle -->
  <div class="mt-6 flex justify-center">
    <button
      id="full-report-button"
      class="inline-flex items-center justify-center gap-2 px-5 py-3 rounded-full bg-gray-800 hover:bg-gray-700 text-white text-sm font-medium focus:outline-none focus-visible:ring-2 focus-visible:ring-emerald-400 focus-visible:ring-offset-0 transition-colors duration-150"
      on:click={toggleAdvanced}
      aria-expanded={showAdvanced}
      aria-controls="advanced-panel"
    >
      <svg class="w-4 h-4" viewBox="0 0 20 20" fill="currentColor">
        {#if showAdvanced}
          <path fill-rule="evenodd" d="M14.77 12.79a.75.75 0 01-1.06-.02L10 8.812l-3.71 3.958a.75.75 0 11-1.08-1.04l4.25-4.53a.75.75 0 011.08 0l4.25 4.53a.75.75 0 01-.02 1.06z" clip-rule="evenodd"/>
        {:else}
          <path fill-rule="evenodd" d="M5.23 7.21a.75.75 0 011.06.02L10 11.188l3.71-3.958a.75.75 0 111.08 1.04l-4.25 4.53a.75.75 0 01-1.08 0l-4.25-4.53a.75.75 0 01.02-1.06z" clip-rule="evenodd"/>
        {/if}
      </svg>
      {showAdvanced ? 'Hide Full Report' : 'View Full Report'}
    </button>
  </div>

  <!-- Advanced Panel -->
  <div id="advanced-panel" class="transition-all duration-500 ease-in-out overflow-hidden {showAdvanced ? 'max-h-[5000px] opacity-100 mt-4' : 'max-h-0 opacity-0'}">
    <div class="rounded-xl border border-gray-800 bg-gray-950 p-6 space-y-6 shadow-md">
      
      <!-- Features -->
      {#if data?.features}
        <section class="bg-gray-900/80 border border-gray-800 rounded-lg p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-base font-semibold text-white">Features</h3>
            <span class="text-[10px] text-gray-400 uppercase tracking-wide px-2 py-0.5 bg-gray-800 rounded">URL + TLD</span>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm text-gray-200">
            {#if data.features.rank !== undefined}
              <div>Rank: <span class="font-medium text-white">{data.features.rank}</span></div>
            {/if}

            {#if data.features.tld}
              <div>TLD: 
                <div class="flex flex-wrap gap-1 mt-1">
                  <span class="px-2 py-0.5 bg-gray-800 text-white rounded text-[11px]">Trusted={String(data.features.tld.is_trusted_tld)}</span>
                  <span class="px-2 py-0.5 bg-gray-800 text-white rounded text-[11px]">Risky={String(data.features.tld.is_risky_tld)}</span>
                  <span class="px-2 py-0.5 bg-gray-800 text-white rounded text-[11px]">ICANN={String(data.features.tld.is_icann)}</span>
                </div>
              </div>
            {/if}

            {#if data.features.url}
              <div>Shortener: <span class="font-medium text-white">{String(data.features.url.url_shortener)}</span></div>
              <div>IP used: <span class="font-medium text-white">{String(data.features.url.uses_ip)}</span></div>
              <div>Punycode: <span class="font-medium text-white">{String(data.features.url.contains_punycode)}</span></div>
              <div>Too long: <span class="font-medium text-white">{String(data.features.url.too_long)}</span></div>
              <div>Too deep: <span class="font-medium text-white">{String(data.features.url.too_deep)}</span></div>
              <div>Homoglyph: <span class="font-medium text-white">{String(data.features.url.has_homoglyph)}</span></div>
              <div>Subdomains: <span class="font-medium text-white">{data.features.url.subdomain_count}</span></div>
            {/if}
          </div>
        </section>
      {/if}

      <!-- Infrastructure -->
      {#if data?.infrastructure}
        <section class="bg-gray-900/80 border border-gray-800 rounded-lg p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-base font-semibold text-white">Infrastructure</h3>
            <span class="text-[10px] text-gray-400 uppercase tracking-wide px-2 py-0.5 bg-gray-800 rounded">Network</span>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm text-gray-200">
            {#if data.infrastructure.ip_addresses?.length}
              <div class="md:col-span-2">
                <div class="text-sm text-gray-300 mb-1">IP addresses</div>
                <div class="flex flex-wrap gap-2">
                  {#each data.infrastructure.ip_addresses as ip}
                    <span class="px-2 py-1 bg-gray-700 text-white rounded text-xs">{ip}</span>
                  {/each}
                </div>
              </div>
            {/if}
            <div>Nameservers valid: <span class="font-medium text-white">{String(data.infrastructure.nameservers_valid)}</span></div>
            <div>MX records valid: <span class="font-medium text-white">{String(data.infrastructure.mx_records_valid)}</span></div>
          </div>
        </section>
      {/if}

      <!-- Domain Info -->
      {#if data?.domain_info}
        <section class="bg-gray-900/80 border border-gray-800 rounded-lg p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-base font-semibold text-white">Domain info</h3>
            <span class="text-[10px] text-gray-400 uppercase tracking-wide px-2 py-0.5 bg-gray-800 rounded">WHOIS</span>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-3 gap-3 text-sm text-gray-200">
            <div>Domain: <span class="font-medium text-white">{data.domain_info.domain}</span></div>
            <div>Registrar: <span class="font-medium text-white">{data.domain_info.registrar || '-'}</span></div>
            <div>DNSSEC: <span class="font-medium text-white">{String(data.domain_info.dnssec)}</span></div>
            <div>Created: <span class="font-medium text-white">{data.domain_info.created}</span></div>
            <div>Updated: <span class="font-medium text-white">{data.domain_info.updated}</span></div>
            <div>Expiry: <span class="font-medium text-white">{data.domain_info.expiry}</span></div>
            <div>Age: <span class="font-medium text-white">{data.domain_info.age_human}</span></div>
          </div>

          {#if data.domain_info.nameservers?.length}
            <div class="text-sm text-gray-300 mt-3">Nameservers</div>
            <div class="flex flex-wrap gap-2 mt-1">
              {#each data.domain_info.nameservers as ns}
                <span class="px-2 py-1 bg-gray-700 text-white rounded text-xs">{ns}</span>
              {/each}
            </div>
          {/if}

          {#if data.domain_info.status?.length}
            <div class="text-sm text-gray-300 mt-3">Status</div>
            <div class="flex flex-wrap gap-2 mt-1">
              {#each data.domain_info.status as st}
                <span class="px-2 py-1 bg-gray-700 text-white rounded text-xs">{st}</span>
              {/each}
            </div>
          {/if}
        </section>
      {/if}

      <!-- Analysis -->
      {#if data?.analysis}
        <section class="bg-gray-900/80 border border-gray-800 rounded-lg p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-base font-semibold text-white">Analysis</h3>
            <span class="text-[10px] text-gray-400 uppercase tracking-wide px-2 py-0.5 bg-gray-800 rounded">HTTP/Redirects</span>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm text-gray-200">
            {#if data.analysis.redirection_result}
              <div>Redirected: <span class="font-medium text-white">{String(data.analysis.redirection_result.is_redirected)}</span></div>
              <div>Chain length: <span class="font-medium text-white">{data.analysis.redirection_result.chain_length}</span></div>
              {#if data.analysis.redirection_result.final_url}
                <div class="md:col-span-2">Final URL: <span class="font-medium text-white break-all">{data.analysis.redirection_result.final_url}</span></div>
              {/if}
              {#if data.analysis.redirection_result.chain?.length}
                <div class="md:col-span-2">
                  <div class="text-sm text-gray-300 mb-1">Chain</div>
                  <ul class="text-sm text-gray-100 list-disc list-inside">
                    {#each data.analysis.redirection_result.chain as c}
                      <li class="break-all">{c}</li>
                    {/each}
                  </ul>
                </div>
              {/if}
            {/if}

            {#if data.analysis.http_status}
              <div>HTTP status: <span class="font-medium text-white">{data.analysis.http_status.code} {data.analysis.http_status.text}</span></div>
              <div>Success: <span class="font-medium text-white">{String(data.analysis.http_status.success)}</span></div>
            {/if}

            {#if data.analysis.is_hsts_supported !== undefined}
              <div>HSTS supported: <span class="font-medium text-white">{String(data.analysis.is_hsts_supported)}</span></div>
            {/if}
          </div>
        </section>
      {/if}

      <!-- Performance -->
      {#if data?.performance}
        <section class="bg-gray-900/80 border border-gray-800 rounded-lg p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-base font-semibold text-white">Performance</h3>
            <span class="text-[10px] text-gray-400 uppercase tracking-wide px-2 py-0.5 bg-gray-800 rounded">Timings</span>
          </div>

          <div class="text-sm text-gray-200 mb-2">Total time: <span class="font-medium text-white">{data.performance.total_time}</span></div>

          {#if data.performance.timings}
            <ul class="text-sm text-gray-100 grid grid-cols-1 md:grid-cols-2 gap-2">
              {#each Object.entries(data.performance.timings) as [k, v]}
                <li class="flex justify-between gap-2 border-b border-gray-800 pb-1">
                  <span class="text-gray-400">{k}</span><span>{v as string}</span>
                </li>
              {/each}
            </ul>
          {/if}
        </section>
      {/if}

      <!-- Scroll Back to Top Button -->
      <div class="mt-6 flex justify-center">
        <button
          class="inline-flex items-center justify-center gap-2 px-5 py-3 rounded-full bg-gray-800 hover:bg-gray-700 text-white text-sm font-medium focus:outline-none focus-visible:ring-2 focus-visible:ring-emerald-400 focus-visible:ring-offset-0 transition-colors duration-150"
          on:click={() => {
              const target = document.getElementById('full-report-button'); 
              if (target) {
                // Scroll so the button is 100px from top
                const offset = 100;
                const top = target.getBoundingClientRect().top + window.scrollY - offset;
                window.scrollTo({ top, behavior: 'smooth' });
              }
            }}
          >
          <svg class="w-4 h-4" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M5.23 12.79a.75.75 0 001.06.02L10 8.812l3.71 3.998a.75.75 0 101.08-1.04l-4.25-4.53a.75.75 0 00-1.08 0l-4.25 4.53a.75.75 0 00.02 1.06z" clip-rule="evenodd"/>
          </svg>
          <!-- Scroll to Top -->
        </button>
      </div>

    <!-- </div>
  </div> -->

    </div>
  </div>

</section>
{/if}

<style>
  img {
    transition: transform 0.2s ease-in-out;
  }
  img:hover {
    transform: scale(1.02);
  }
</style>
