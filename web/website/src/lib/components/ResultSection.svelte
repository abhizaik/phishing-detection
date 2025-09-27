<script lang="ts">
  import type { AnalyzeResult, ScreenshotResponse } from '../types';
  export let data: AnalyzeResult | null = null;
  export let screenshotData: ScreenshotResponse | null = null;
  export let loading = false;
  export let error: string | null = null;
  let showAdvanced = false;
  $: primary = data?.result;
  $: reasons = primary?.reasons;
  // export let screenshotUrl: string | null = null;
  function toggleAdvanced() {
    showAdvanced = !showAdvanced;
  }
  let copied = false;
  async function copyShareLink() {
  try {
    const url = new URL(window.location.href);
    if (data?.url) {
      url.searchParams.set('q', data.url);
      await navigator.clipboard.writeText(url.toString());
      copied = true;
      setTimeout(() => (copied = false), 2000);
    }
  } catch {
    error = "Could not copy link";
  }
}

</script>

{#if error}
  <div class="max-w-3xl mx-auto p-4 bg-red-900/30 border border-red-700 text-red-200 rounded-md">{error}</div>
{:else if loading}
  <div class="max-w-3xl mx-auto space-y-4">
    <div class="animate-pulse rounded-lg border border-gray-800 bg-gray-950/60 p-4">
      <div class="h-4 w-24 bg-gray-800 rounded mb-3"></div>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
        <div class="h-2 bg-gray-800 rounded"></div>
        <div class="h-2 bg-gray-800 rounded"></div>
        <div class="h-2 bg-gray-800 rounded"></div>
      </div>
    </div>
    <div class="animate-pulse rounded-md border border-gray-800 bg-gray-950/60 p-4 h-24"></div>
  </div>
{:else if data}
  <section class="max-w-3xl mx-auto space-y-6">
    <div class="flex items-center justify-between gap-3">
      <h2 class="text-xl font-semibold text-gray-100">Result</h2>
      <div class="flex items-center gap-2">
        <!-- <button class="px-3 py-1.5 text-xs rounded bg-gray-800 hover:bg-gray-700 border border-gray-700" on:click={copyShareLink}>Copy link</button> -->
        <button
            class="inline-flex items-center justify-center gap-2 px-5 py-3 rounded-lg bg-gray-800 hover:bg-gray-700 text-white text-sm font-medium focus:outline-none"
            on:click={copyShareLink}
          >
            {#if copied}
              <!-- Tickmark icon -->
              <svg class="w-4 h-4 text-emerald-300" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
              </svg>
              <span class="text-emerald-300">Copied!</span>
            {:else}
              <!-- Clipboard icon -->
              <svg class="w-4 h-4 text-gray-300" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
                <path d="M8 2a2 2 0 00-2 2v2h2V4h6v6h-2v2h2a2 2 0 002-2V4a2 2 0 00-2-2H8zM4 8a2 2 0 00-2 2v6a2 2 0 002 2h6a2 2 0 002-2v-6a2 2 0 00-2-2H4zm0 2h6v6H4v-6z"/>
              </svg>
              <span>Copy Result Link</span>
            {/if}
          </button>



        <!-- {#if copied}
          <span class="text-[11px] text-emerald-300">Copied</span>
        {/if} -->
      </div>
    </div>

    <div class="rounded-xl border border-gray-800 bg-gray-950 p-5">
  <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-3">
    <div>
      <div class="text-xs text-gray-400">Verdict :</div>
      <div class="inline-flex items-center gap-2 mt-0.5">
        <span class="text-lg font-semibold">{primary?.verdict ?? '-'}</span>

      </div>
    </div>

    <div class="flex items-center justify-between text-lg font-bold text-gray-200 w-full md:w-auto">
      <!-- <span>Trust Score : &nbsp; </span> -->
      <div class="text-xs text-gray-400">
        Trust Score : &nbsp;
      </div>
      <div class="inline-flex items-center gap-2 mt-0.5">
      <span class="text-blue-400">{primary?.final_score ?? '-'} / 100</span>
    </div>
  </div>
</div>


    {#if reasons?.bad_reasons?.length || reasons?.good_reasons?.length}
  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    {#if reasons?.bad_reasons?.length}
      <div class="rounded-lg border border-red-900 bg-red-950/40 p-4 shadow-sm">
        <div class="text-sm font-semibold text-red-300 mb-2 flex items-center gap-2">
          <svg class="w-4 h-4 text-red-400" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.721-1.36 3.486 0l6.518 11.59c.75 1.335-.213 3.011-1.743 3.011H3.482c-1.53 0-2.493-1.676-1.743-3.01L8.257 3.1zM11 14a1 1 0 10-2 0 1 1 0 002 0zm-1-2a.75.75 0 01-.75-.75V8a.75.75 0 011.5 0v3.25A.75.75 0 0110 12z" clip-rule="evenodd" />
          </svg>
          Red flags
        </div>
        <ul class="text-sm text-red-200 space-y-1">
          {#each reasons.bad_reasons as r}
            <li class="flex items-start gap-2">
              <span class="mt-1 inline-block h-1.5 w-1.5 rounded-full bg-red-400"></span>
              <span>{r}</span>
            </li>
          {/each}
        </ul>
      </div>
    {/if}

    {#if reasons?.good_reasons?.length}
      <div class="rounded-lg border border-emerald-900 bg-emerald-950/40 p-4 shadow-sm">
        <div class="text-sm font-semibold text-emerald-300 mb-2 flex items-center gap-2">
          <svg class="w-4 h-4 text-emerald-400" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M16.704 5.29a1 1 0 010 1.42l-7.388 7.388a1 1 0 01-1.42 0L3.296 9.498a1 1 0 111.408-1.42L8.5 11.874l6.796-6.795a1 1 0 011.408 0z" clip-rule="evenodd" />
          </svg>
          Green flags
        </div>
        <ul class="text-sm text-emerald-200 space-y-1">
          {#each reasons.good_reasons as r}
            <li class="flex items-start gap-2">
              <span class="mt-1 inline-block h-1.5 w-1.5 rounded-full bg-emerald-400"></span>
              <span>{r}</span>
            </li>
          {/each}
        </ul>
      </div>
    {/if}
  </div>
{/if}

  {#if screenshotData}
          <div class="mt-3">
            <h4 class="text-xs font-semibold text-gray-400 mb-1">Screenshot</h4>
            <img
              src={screenshotData.file}
              alt="Website screenshot"
              class="w-full rounded border border-gray-800"
              loading="lazy"
            />
          </div>
        {/if}


    <div>
      <button
  class="inline-flex items-center justify-center gap-2 px-5 py-3 rounded-lg bg-gray-800 hover:bg-gray-700 text-white text-sm font-medium focus:outline-none focus-visible:ring-2 focus-visible:ring-emerald-400 focus-visible:ring-offset-0 transition-colors duration-150"
  on:click={toggleAdvanced}
  aria-expanded={showAdvanced}
  aria-controls="advanced-panel"
>
  {#if showAdvanced}
    <!-- Up chevron -->
    <svg class="w-4 h-4" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
      <path fill-rule="evenodd" d="M14.77 12.79a.75.75 0 01-1.06-.02L10 8.812l-3.71 3.958a.75.75 0 11-1.08-1.04l4.25-4.53a.75.75 0 011.08 0l4.25 4.53a.75.75 0 01-.02 1.06z" clip-rule="evenodd" />
    </svg>
  {:else}
    <!-- Down chevron -->
    <svg class="w-4 h-4" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
      <path fill-rule="evenodd" d="M5.23 7.21a.75.75 0 011.06.02L10 11.188l3.71-3.958a.75.75 0 111.08 1.04l-4.25 4.53a.75.75 0 01-1.08 0l-4.25-4.53a.75.75 0 01.02-1.06z" clip-rule="evenodd" />
    </svg>
  {/if}

  {showAdvanced ? 'Hide info for nerds' : 'Show info for nerds'}
</button>


    </div>

    {#if showAdvanced}
      <div id="advanced-panel" class="rounded-md border border-gray-800 bg-gray-950 p-4 space-y-6 transition-all duration-200">
        {#if data?.features}
          <section class="bg-gray-950 border border-gray-800 rounded-lg p-4">
            <div class="flex items-center justify-between mb-3">
              <h3 class="text-sm font-semibold text-gray-200">Features</h3>
              <span class="text-[10px] text-gray-500 uppercase tracking-wide">URL + TLD</span>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-3 text-xs text-gray-300">
              {#if data.features.rank !== undefined}
                <div>Rank <span class="text-gray-100">{data.features.rank}</span></div>
              {/if}
              {#if data.features.tld}
                <div>TLD <span class="text-gray-100">trusted={String(data.features.tld.is_trusted_tld)} risky={String(data.features.tld.is_risky_tld)} icann={String(data.features.tld.is_icann)}</span></div>
              {/if}
              {#if data.features.url}
                <div>URL shortener <span class="text-gray-100">{String(data.features.url.url_shortener)}</span></div>
                <div>Uses IP <span class="text-gray-100">{String(data.features.url.uses_ip)}</span></div>
                <div>Punycode <span class="text-gray-100">{String(data.features.url.contains_punycode)}</span></div>
                <div>Too long <span class="text-gray-100">{String(data.features.url.too_long)}</span></div>
                <div>Too deep <span class="text-gray-100">{String(data.features.url.too_deep)}</span></div>
                <div>Homoglyph <span class="text-gray-100">{String(data.features.url.has_homoglyph)}</span></div>
                <div>Subdomain count <span class="text-gray-100">{data.features.url.subdomain_count}</span></div>
              {/if}
            </div>
          </section>
        {/if}

        {#if data?.infrastructure}
          <section class="bg-gray-950 border border-gray-800 rounded-lg p-4">
            <div class="flex items-center justify-between mb-3">
              <h3 class="text-sm font-semibold text-gray-200">Infrastructure</h3>
              <span class="text-[10px] text-gray-500 uppercase tracking-wide">Network</span>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-3 text-xs text-gray-300">
              <div class="md:col-span-2">
                {#if data.infrastructure.ip_addresses?.length}
                  <div class="text-xs text-gray-300 mb-1">IP addresses</div>
                  <ul class="text-xs text-gray-100 list-disc list-inside">
                    {#each data.infrastructure.ip_addresses as ip}
                      <li>{ip}</li>
                    {/each}
                  </ul>
                {/if}
              </div>
              <div>Nameservers valid <span class="text-gray-100">{String(data.infrastructure.nameservers_valid)}</span></div>
              <div>MX records valid <span class="text-gray-100">{String(data.infrastructure.mx_records_valid)}</span></div>
            </div>
          </section>
        {/if}

        {#if data?.domain_info}
          <section class="bg-gray-950 border border-gray-800 rounded-lg p-4">
            <div class="flex items-center justify-between mb-3">
              <h3 class="text-sm font-semibold text-gray-200">Domain info</h3>
              <span class="text-[10px] text-gray-500 uppercase tracking-wide">WHOIS</span>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-3 text-xs text-gray-300">
              <div>Domain <span class="text-gray-100">{data.domain_info.domain}</span></div>
              <div>Registrar <span class="text-gray-100">{data.domain_info.registrar || '-'}</span></div>
              <div>DNSSEC <span class="text-gray-100">{String(data.domain_info.dnssec)}</span></div>
              <div>Created <span class="text-gray-100">{data.domain_info.created}</span></div>
              <div>Updated <span class="text-gray-100">{data.domain_info.updated}</span></div>
              <div>Expiry <span class="text-gray-100">{data.domain_info.expiry}</span></div>
              <div>Age <span class="text-gray-100">{data.domain_info.age_human}</span></div>
              <!-- <div>Age (days) <span class="text-gray-100">{data.domain_info.age_days}</span></div> -->
            </div>
            {#if data.domain_info.nameservers?.length}
              <div class="text-xs text-gray-300 mt-3">Nameservers</div>
              <ul class="text-xs text-gray-100 list-disc list-inside">
                {#each data.domain_info.nameservers as ns}
                  <li>{ns}</li>
                {/each}
              </ul>
            {/if}
            {#if data.domain_info.status?.length}
              <div class="text-xs text-gray-300 mt-3">Status</div>
              <ul class="text-xs text-gray-100 list-disc list-inside">
                {#each data.domain_info.status as st}
                  <li>{st}</li>
                {/each}
              </ul>
            {/if}
          </section>
        {/if}

        {#if data?.analysis}
          <section class="bg-gray-950 border border-gray-800 rounded-lg p-4">
            <div class="flex items-center justify-between mb-3">
              <h3 class="text-sm font-semibold text-gray-200">Analysis</h3>
              <span class="text-[10px] text-gray-500 uppercase tracking-wide">HTTP/Redirects</span>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-3 text-xs text-gray-300">
              {#if data.analysis.redirection_result}
                <div>Redirected <span class="text-gray-100">{String(data.analysis.redirection_result.is_redirected)}</span></div>
                <div>Chain length <span class="text-gray-100">{data.analysis.redirection_result.chain_length}</span></div>
                {#if data.analysis.redirection_result.final_url}
                  <div class="md:col-span-2">Final URL <span class="text-gray-100 break-all">{data.analysis.redirection_result.final_url}</span></div>
                {/if}
                {#if data.analysis.redirection_result.chain?.length}
                  <div class="md:col-span-2">
                    <div class="text-xs text-gray-300 mb-1">Chain</div>
                    <ul class="text-xs text-gray-100 list-disc list-inside">
                      {#each data.analysis.redirection_result.chain as c}
                        <li class="break-all">{c}</li>
                      {/each}
                    </ul>
                  </div>
                {/if}
              {/if}
              {#if data.analysis.http_status}
                <div>HTTP status <span class="text-gray-100">{data.analysis.http_status.code} {data.analysis.http_status.text}</span></div>
                <div>Success <span class="text-gray-100">{String(data.analysis.http_status.success)}</span></div>
              {/if}
              {#if data.analysis.is_hsts_supported !== undefined}
                <div>HSTS supported <span class="text-gray-100">{String(data.analysis.is_hsts_supported)}</span></div>
              {/if}
            </div>
          </section>
        {/if}

        {#if data?.performance}
          <section class="bg-gray-950 border border-gray-800 rounded-lg p-4">
            <div class="flex items-center justify-between mb-3">
              <h3 class="text-sm font-semibold text-gray-200">Performance</h3>
              <span class="text-[10px] text-gray-500 uppercase tracking-wide">Timings</span>
            </div>
            <div class="text-xs text-gray-300 mb-2">Total time <span class="text-gray-100">{data.performance.total_time}</span></div>
            {#if data.performance.timings}
              <ul class="text-[11px] text-gray-100 grid grid-cols-1 md:grid-cols-2 gap-2">
                {#each Object.entries(data.performance.timings) as [k, v]}
                  <li class="flex justify-between gap-2 border-b border-gray-900 pb-1"><span class="text-gray-400">{k}</span><span>{v as string}</span></li>
                {/each}
              </ul>
            {/if}
          </section>
        {/if}

        <!-- {#if screenshotUrl}
          <section class="bg-gray-950 border border-gray-800 rounded-lg p-4">
            <h3 class="text-sm font-semibold text-gray-200 mb-2">Screenshot</h3>
            <img src={screenshotUrl} alt="Website screenshot" class="w-full rounded border border-gray-800" loading="lazy" />
          </section>
        {/if} -->
      </div>
    {/if}
  </section>
{/if}

<style>
</style>

