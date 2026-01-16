<script lang="ts">
  import { browser } from "$app/environment";
  import type { AnalyzeResult, ScreenshotResponse } from "../types";
  import TooltipIcon from "./TooltipIcon.svelte";
  export let data: AnalyzeResult | null = null;
  export let screenshotData: ScreenshotResponse | null = null;
  export let loading = false;
  export let error: string | null = null;

  let showAdvanced = false;
  $: primary = data?.result;
  $: reasons = primary?.reasons;

  // Available tabs based on data
  $: availableTabs = [
    { id: "features", label: "URL Features", icon: "⚙️", condition: data?.features },
    { id: "infrastructure", label: "Infra Details", icon: "🌐", condition: data?.infrastructure },
    { id: "domain", label: "Domain Info", icon: "📋", condition: data?.domain_info },
    { id: "analysis", label: "Redirections", icon: "🔍", condition: data?.analysis },
    // { id: "performance", label: "Performance", icon: "⚡", condition: data?.performance },
  ].filter((tab) => tab.condition);

  function scrollToSection(sectionId: string) {
    const element = document.getElementById(`section-${sectionId}`);
    if (element) {
      const offset = 80;
      const top = element.getBoundingClientRect().top + window.scrollY - offset;
      window.scrollTo({ top, behavior: "smooth" });
    }
  }

  $: gridColumns = availableTabs.length > 0 ? `repeat(${availableTabs.length}, 1fr)` : "1fr";

  let copied = false;
  let showModal = false;

  async function copyShareLink() {
    // Function for first copy button, commented out as we added new Share Button
    try {
      const url = new URL(window.location.href);
      if (data?.url) {
        url.searchParams.set("q", data.url);
        await navigator.clipboard.writeText(url.toString());
        copied = true;
        setTimeout(() => (copied = false), 1200);
      }
    } catch {
      error = "Could not copy link";
    }
  }

  let shareCopied = false;

  async function shareLink() {
    if (!browser) return;

    // Get current URL from address bar
    const currentUrl = window.location.href;
    const shareText = `Check out this SafeSurf result for ${data?.domain}`;

    if (navigator.share) {
      try {
        await navigator.share({
          title: "SafeSurf",
          text: shareText,
          url: currentUrl,
        });
      } catch (err: unknown) {
        // User cancelled or share failed, fall back to clipboard
        if (err instanceof Error && err.name !== "AbortError") {
          console.error("Share failed:", err);
        }
        // Fall through to clipboard copy
        try {
          await navigator.clipboard.writeText(`${shareText}\n${currentUrl}`);
          shareCopied = true;
          setTimeout(() => (shareCopied = false), 2000);
        } catch (clipboardErr) {
          console.error("Clipboard copy failed:", clipboardErr);
        }
      }
    } else {
      // Fallback to clipboard copy
      try {
        await navigator.clipboard.writeText(`${shareText}\n${currentUrl}`);
        shareCopied = true;
        setTimeout(() => (shareCopied = false), 2000);
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
  <div class="max-w-3xl mx-auto p-4 bg-red-900/30 border border-red-700 text-red-200 rounded-md">
    {error}
  </div>
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
        <h2 class="text-2xl font-semibold text-white" id="analysis-summary">Analysis Summary</h2>
        <p class="text-gray-400 text-sm mt-1">
          Here’s the security profile for {data?.domain}
        </p>
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
        class="group inline-flex items-center justify-center gap-2.5 px-6 py-3.5 rounded-lg bg-gradient-to-r from-blue-600 to-blue-500 hover:from-blue-500 hover:to-blue-400 text-white text-sm font-semibold shadow-lg hover:shadow-xl transition-all duration-200 hover:scale-105 active:scale-95 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:ring-offset-2 focus:ring-offset-gray-950 disabled:opacity-50 {shareCopied
          ? 'from-emerald-600 to-emerald-500 hover:from-emerald-500 hover:to-emerald-400 animate-pulse'
          : ''}"
        on:click={shareLink}
        disabled={shareCopied}
      >
        {#if shareCopied}
          <svg
            class="w-5 h-5 text-white animate-in fade-in"
            fill="none"
            stroke="currentColor"
            stroke-width="2.5"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
          </svg>
          <span class="text-white font-medium">Copied!</span>
        {:else}
          <svg
            class="w-5 h-5 text-white group-hover:rotate-12 transition-transform duration-200"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M7.217 10.907a2.25 2.25 0 100 2.186m0-2.186c.18.324.283.696.283 1.093s-.103.77-.283 1.093m0-2.186l9.566-5.314m-9.566 7.5l9.566 5.314m0 0a2.25 2.25 0 103.935 2.186 2.25 2.25 0 00-3.935-2.186zm0-12.814a2.25 2.25 0 103.935-2.186 2.25 2.25 0 00-3.935 2.186z"
            />
          </svg>
          <span class="font-medium">Share</span>
        {/if}
      </button>
    </div>

    <!-- Verdict & Trust Score -->
    <div
      class="flex flex-col md:flex-row gap-6 p-6 bg-gray-900/80 rounded-xl shadow-md hover:shadow-lg transition-transform hover:scale-[1.01]"
    >
      <div class="flex-1">
        <div class="text-sm font-medium text-gray-300 uppercase tracking-wide mb-1">Verdict</div>
        <div class="flex items-center gap-3">
          <span class="text-2xl font-bold text-white">{primary?.verdict ?? "-"}</span>
          {#if primary?.verdict === "Safe"}
            <span
              class="px-3 py-1 rounded-full bg-green-700 text-white font-medium text-xs uppercase tracking-wide"
            >
              Trusted
            </span>
          {:else if primary?.verdict === "Suspicious"}
            <span
              class="px-3 py-1 rounded-full bg-yellow-500 text-black font-medium text-xs uppercase tracking-wide"
            >
              Be Cautious
            </span>
          {:else if primary?.verdict === "Risky"}
            <span
              class="px-3 py-1 rounded-full bg-red-700 text-white font-medium text-xs uppercase tracking-wide"
            >
              High Risk
            </span>
          {:else if primary?.verdict === "Unclear"}
            <span
              class="px-3 py-1 rounded-full bg-gray-500 text-white font-medium text-xs uppercase tracking-wide"
            >
              Not Enough Data
            </span>
          {:else}
            <span
              class="px-3 py-1 rounded-full bg-red-600 text-white font-medium text-xs uppercase tracking-wide"
            >
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
        <div class="text-sm font-medium text-gray-300 uppercase tracking-wide mb-1">
          Trust Score
        </div>
        <span class="text-3xl font-extrabold text-white-50"
          >{primary?.final_score ?? "-"} / 100</span
        >
      </div>
    </div>

    <!-- Flags -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      {#if reasons}
        <!-- Red Flags -->
        <div
          class="rounded-xl border border-red-700 bg-red-900/20 p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all"
        >
          <div class="flex items-center gap-3 mb-3">
            <svg class="w-5 h-5 text-red-500 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
              <path
                fill-rule="evenodd"
                d="M8.257 3.099c.765-1.36 2.721-1.36 3.486 0l6.518 11.59c.75 1.335-.213 3.011-1.743 3.011H3.482c-1.53 0-2.493-1.676-1.743-3.01L8.257 3.1zM11 14a1 1 0 10-2 0 1 1 0 002 0zm-1-2a.75.75 0 01-.75-.75V8a.75.75 0 011.5 0v3.25A.75.75 0 0110 12z"
                clip-rule="evenodd"
              />
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
        <div
          class="rounded-xl border border-emerald-700 bg-emerald-900/20 p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all"
        >
          <div class="flex items-center gap-3 mb-3">
            <svg
              class="w-5 h-5 text-emerald-400 flex-shrink-0"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fill-rule="evenodd"
                d="M16.704 5.29a1 1 0 010 1.42l-7.388 7.388a1 1 0 01-1.42 0L3.296 9.498a1 1 0 111.408-1.42L8.5 11.874l6.796-6.795a1 1 0 011.408 0z"
                clip-rule="evenodd"
              />
            </svg>
            <h3 class="text-sm font-semibold text-emerald-400 uppercase tracking-wide">
              Green Flags
            </h3>
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
    {#if screenshotData?.status === "success" && screenshotData?.file && !screenshotData.file.startsWith("server")}
      <div
        class="mt-6 rounded-xl border border-gray-800 bg-gray-900/70 p-4 shadow-md hover:shadow-lg transition-all"
      >
        <h4 class="text-sm font-semibold text-gray-300 mb-2">Website Screenshot</h4>
        <!-- Image for testing -->
        <!-- src="screenshot-google-com.png" -->
        <img
          src={screenshotData.file}
          alt="Website screenshot"
          class="w-full rounded-lg border border-gray-800 cursor-pointer hover:opacity-90"
          loading="lazy"
          on:click={() => (showModal = true)}
        />
      </div>
    {/if}

    <!-- Screenshot Modal -->
    {#if showModal}
      <div class="fixed inset-0 bg-black/80 flex items-center justify-center z-50">
        <button
          class="absolute top-4 right-4 text-gray-300 hover:text-white text-2xl"
          on:click={() => (showModal = false)}>×</button
        >
        <img
          src={screenshotData?.file}
          alt="Full screenshot"
          class="max-h-[90vh] max-w-[90vw] rounded-lg shadow-lg"
        />
        <!-- <img src=screenshot-google-com.png alt="Full screenshot" class="max-h-[90vh] max-w-[90vw] rounded-lg shadow-lg" /> -->
      </div>
    {/if}

    <!-- Advanced Panel Toggle -->
    <div class="mt-6 flex justify-center">
      <button
        id="full-report-button"
        class="inline-flex items-center justify-center gap-2 px-5 py-3 rounded-lg {showAdvanced
          ? 'bg-gray-800 hover:bg-gray-700 text-white'
          : 'bg-blue-600 hover:bg-blue-500 text-white'} text-sm font-medium disabled:opacity-50 focus:outline-none focus:ring-2 focus:ring-offset-0 focus:ring-blue-600 transition-all duration-200"
        on:click={toggleAdvanced}
        aria-expanded={showAdvanced}
        aria-controls="advanced-panel"
      >
        <svg class="w-4 h-4" viewBox="0 0 20 20" fill="currentColor">
          {#if showAdvanced}
            <path
              fill-rule="evenodd"
              d="M14.77 12.79a.75.75 0 01-1.06-.02L10 8.812l-3.71 3.958a.75.75 0 11-1.08-1.04l4.25-4.53a.75.75 0 011.08 0l4.25 4.53a.75.75 0 01-.02 1.06z"
              clip-rule="evenodd"
            />
          {:else}
            <path
              fill-rule="evenodd"
              d="M5.23 7.21a.75.75 0 011.06.02L10 11.188l3.71-3.958a.75.75 0 111.08 1.04l-4.25 4.53a.75.75 0 01-1.08 0l-4.25-4.53a.75.75 0 01.02-1.06z"
              clip-rule="evenodd"
            />
          {/if}
        </svg>
        <span>{showAdvanced ? "Hide Full Report" : "View Full Report"}</span>
      </button>
    </div>

    <!-- Advanced Panel -->
    <div
      id="advanced-panel"
      class="transition-all duration-500 ease-in-out {showAdvanced
        ? 'max-h-[5000px] opacity-100 mt-4 overflow-visible'
        : 'max-h-0 opacity-0 overflow-hidden'}"
    >
      <div class="rounded-xl border border-gray-800 bg-gray-950 shadow-md overflow-hidden">
        <!-- Tabs Navigation -->
        {#if availableTabs.length > 0}
          <div class="border-b border-gray-800 bg-gray-900/50">
            <div class="grid" style="grid-template-columns: {gridColumns};">
              {#each availableTabs as tab}
                <button
                  class="flex items-center justify-center gap-2 px-4 py-4 text-sm font-medium transition-all duration-200 whitespace-nowrap border-b-2 border-transparent text-gray-400 hover:text-gray-300 hover:bg-gray-900/40 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-emerald-400"
                  on:click={() => scrollToSection(tab.id)}
                  role="tab"
                >
                  <span class="text-base">{tab.icon}</span>
                  <span>{tab.label}</span>
                </button>
              {/each}
            </div>
          </div>
        {/if}

        <!-- All Content - Continuous Scroll -->
        <div class="p-6 space-y-6">
          <!-- Features Section -->
          {#if data?.features}
            <section
              id="section-features"
              class="bg-gray-900/80 border border-gray-800 rounded-lg p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all scroll-mt-20"
            >
              <div class="flex items-center justify-between mb-4">
                <h3 class="text-base font-semibold text-white">Website Features</h3>
                <span
                  class="text-[10px] text-gray-400 uppercase tracking-wide px-2 py-0.5 bg-gray-800 rounded"
                  >URL + Domain</span
                >
              </div>

              <div
                class="space-y-0 divide-y divide-gray-800 text-sm text-gray-200 max-w-4xl w-full mx-auto"
              >
                {#if data.features.rank !== undefined}
                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>Global Traffic Rank:</span>
                      <TooltipIcon
                        text="A rough estimate of the website’s global popularity, lower numbers mean more visitors. Derived from traffic and engagement data."
                      />
                    </div>

                    <span class="font-medium text-white">
                      {data.features.rank === 0 ? "Unranked" : data.features.rank}
                    </span>
                  </div>
                {/if}

                {#if data.features.tld}
                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>Domain Ending (TLD):</span>
                      <TooltipIcon
                        text="The last part of a domain name (like .com, .org, .io). It can hint at the site’s trust level or purpose."
                      />
                    </div>
                    <span class="font-medium text-white">.{data.features.tld.tld}</span>
                  </div>

                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>High Trust Domain Ending:</span>
                      <TooltipIcon
                        text="Indicates whether this domain ending (TLD) is widely recognized and commonly used by highly trusted entities like government and other institutions."
                      />
                    </div>
                    {#if data.features.tld.is_trusted_tld}
                      <span class="text-green-400 font-medium flex items-center gap-1">✅ Yes</span>
                    {:else}
                      <span class="text-red-400 font-medium flex items-center gap-1">❌ No</span>
                    {/if}
                  </div>

                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>Is Risky TLD:</span>
                      <TooltipIcon
                        text="Some TLDs are frequently abused by scammers or malicious sites. 'Yes' suggests caution."
                      />
                    </div>
                    {#if data.features.tld.is_risky_tld}
                      <span class="text-red-400 font-medium flex items-center gap-1">❌ Yes</span>
                    {:else}
                      <span class="text-green-400 font-medium flex items-center gap-1">✅ No</span>
                    {/if}
                  </div>

                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>TLD Recognized by ICANN:</span>
                      <TooltipIcon
                        text="ICANN oversees global domain names. Recognition means this TLD is officially managed and monitored."
                      />
                    </div>
                    {#if data.features.tld.is_icann}
                      <span class="text-green-400 font-medium flex items-center gap-1">✅ Yes</span>
                    {:else}
                      <span class="text-red-400 font-medium flex items-center gap-1">❌ No</span>
                    {/if}
                  </div>
                {/if}

                {#if data.features.url}
                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>Uses Link Shortener:</span>
                      <TooltipIcon
                        text="Shows if the URL uses a shortening service (like bit.ly). Shortened links can hide a site’s real destination."
                      />
                    </div>
                    {#if data.features.url.url_shortener}
                      <span class="text-red-400 font-medium flex items-center gap-1">❌ Yes</span>
                    {:else}
                      <span class="text-green-400 font-medium flex items-center gap-1">✅ No</span>
                    {/if}
                  </div>

                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>Uses Direct IP Address:</span>
                      <TooltipIcon
                        text="Some malicious sites use IP addresses instead of domain names to avoid detection. Legitimate websites rarely do this."
                      />
                    </div>
                    {#if data.features.url.uses_ip}
                      <span class="text-red-400 font-medium flex items-center gap-1">❌ Yes</span>
                    {:else}
                      <span class="text-green-400 font-medium flex items-center gap-1">✅ No</span>
                    {/if}
                  </div>

                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>Contains Punycode Characters:</span>
                      <TooltipIcon
                        text="Punycode allows special or non-Latin characters in domains (like xn--example). Scammers sometimes exploit this for lookalike attacks."
                      />
                    </div>
                    {#if data.features.url.contains_punycode}
                      <span class="text-red-400 font-medium flex items-center gap-1">❌ Yes</span>
                    {:else}
                      <span class="text-green-400 font-medium flex items-center gap-1">✅ No</span>
                    {/if}
                  </div>

                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>URL Too Long:</span>
                      <TooltipIcon
                        text="Very long URLs can be used to obscure malicious content or trick users into trusting a fake site."
                      />
                    </div>
                    {#if data.features.url.too_long}
                      <span class="text-red-400 font-medium flex items-center gap-1">❌ Yes</span>
                    {:else}
                      <span class="text-green-400 font-medium flex items-center gap-1">✅ No</span>
                    {/if}
                  </div>

                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>URL Too Deep (Many Slashes):</span>
                      <TooltipIcon
                        text="A URL with many nested paths may indicate redirections or hidden content, often seen in phishing attempts."
                      />
                    </div>
                    {#if data.features.url.too_deep}
                      <span class="text-red-400 font-medium flex items-center gap-1">❌ Yes</span>
                    {:else}
                      <span class="text-green-400 font-medium flex items-center gap-1">✅ No</span>
                    {/if}
                  </div>

                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>Has Look-Alike Letters (Homoglyph):</span>
                      <TooltipIcon
                        text="Detects if the domain name includes characters that look like others (e.g., 'go0gle.com'). Often used for impersonation."
                      />
                    </div>
                    {#if data.features.url.has_homoglyph}
                      <span class="text-red-400 font-medium flex items-center gap-1">❌ Yes</span>
                    {:else}
                      <span class="text-green-400 font-medium flex items-center gap-1">✅ No</span>
                    {/if}
                  </div>

                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>Subdomain Count:</span>
                      <TooltipIcon
                        text="Shows how many subdomains (like shop.example.com) are used. Too many can hint at suspicious or temporary setups."
                      />
                    </div>
                    <span class="font-medium text-white">{data.features.url.subdomain_count}</span>
                  </div>
                {/if}
              </div>
            </section>
          {/if}

          <!-- Infrastructure Section -->
          {#if data?.infrastructure}
            <section
              id="section-infrastructure"
              class="bg-gray-900/80 border border-gray-800 rounded-lg p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all scroll-mt-20"
            >
              <div class="flex items-center justify-between mb-4">
                <h3 class="text-base font-semibold text-white">Website Infrastructure</h3>
                <span
                  class="text-[10px] text-gray-400 uppercase tracking-wide px-2 py-0.5 bg-gray-800 rounded"
                  >Network Details</span
                >
              </div>

              <div
                class="space-y-0 divide-y divide-gray-800 text-sm text-gray-200 max-w-4xl w-full mx-auto"
              >
                {#if data.infrastructure.ip_addresses?.length}
                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span
                        >Server IP Address{data.infrastructure.ip_addresses.length > 1
                          ? "es"
                          : ""}:</span
                      >
                      <TooltipIcon
                        text="The actual network address where the website is hosted. Each IP points to a specific physical or cloud server."
                      />
                    </div>
                    <div class="flex flex-wrap gap-2">
                      {#each data.infrastructure.ip_addresses as ip}
                        <span class="px-2 py-1 bg-gray-800 text-white rounded text-xs">{ip}</span>
                      {/each}
                    </div>
                  </div>
                {/if}

                <div
                  class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                >
                  <div class="flex items-center gap-1 text-gray-400">
                    <span>DNS Configuration:</span>
                    <TooltipIcon
                      text="Nameservers control where your domain points. They act like the internet’s ‘address book’, linking your domain name to the right hosting provider."
                    />
                  </div>
                  {#if data.infrastructure.nameservers_valid}
                    <span class="text-green-400 font-medium flex items-center gap-1"
                      >✅ Detected</span
                    >
                  {:else}
                    <span class="text-red-400 font-medium flex items-center gap-1"
                      >❌ Not Detected</span
                    >
                  {/if}
                </div>

                {#if data.infrastructure.ns_hosts?.length > 0}
                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>Nameserver Hosts:</span>
                      <TooltipIcon
                        text="The servers responsible for managing your domain’s DNS settings. These typically belong to your registrar or hosting provider."
                      />
                    </div>
                    <div class="flex flex-wrap gap-2">
                      {#each data.infrastructure.ns_hosts as ns_host}
                        <span class="px-2 py-1 bg-gray-800 text-white rounded text-xs"
                          >{ns_host}</span
                        >
                      {/each}
                    </div>
                  </div>
                {/if}

                <div
                  class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                >
                  <div class="flex items-center gap-1 text-gray-400">
                    <span>Email Server Records:</span>
                    <TooltipIcon
                      text="MX records define where emails for this domain are delivered, essential for sending and receiving mail securely."
                    />
                  </div>
                  {#if data.infrastructure.mx_records_valid}
                    <span class="text-green-400 font-medium flex items-center gap-1"
                      >✅ Detected</span
                    >
                  {:else}
                    <span class="text-red-400 font-medium flex items-center gap-1"
                      >❌ Not Detected</span
                    >
                  {/if}
                </div>

                {#if data.infrastructure.mx_hosts?.length > 0}
                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>MX Hosts:</span>
                      <TooltipIcon
                        text="The mail servers responsible for handling your domain’s email traffic."
                      />
                    </div>
                    <div class="flex flex-wrap gap-2">
                      {#each data.infrastructure.mx_hosts as mx_host}
                        <span class="px-2 py-1 bg-gray-800 text-white rounded text-xs"
                          >{mx_host}</span
                        >
                      {/each}
                    </div>
                  </div>
                {/if}
              </div>
            </section>
          {/if}

          <!-- Domain Info Section -->
          {#if data?.domain_info}
            <section
              id="section-domain"
              class="bg-gray-900/80 border border-gray-800 rounded-lg p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all scroll-mt-20"
            >
              <div class="flex items-center justify-between mb-4">
                <h3 class="text-base font-semibold text-white">Domain Information</h3>
                <span
                  class="text-[10px] text-gray-400 uppercase tracking-wide px-2 py-0.5 bg-gray-800 rounded"
                  >{data.domain_info.source}</span
                >
              </div>

              <div
                class="space-y-0 divide-y divide-gray-800 text-sm text-gray-200 max-w-4xl w-full mx-auto"
              >
                <div
                  class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                >
                  <div class="flex items-center gap-1 text-gray-400">
                    <span>Domain:</span>
                    <TooltipIcon
                      text="The registered name of the website — what users type in the browser to visit it."
                    />
                  </div>
                  <span class="font-medium text-white">{data.domain_info.domain}</span>
                </div>

                <div
                  class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                >
                  <div class="flex items-center gap-1 text-gray-400">
                    <span>Registrar:</span>
                    <TooltipIcon
                      text="The company or organization that manages the registration of this domain (e.g., GoDaddy, Namecheap, Google Domains)."
                    />
                  </div>
                  <span class="font-medium text-white">{data.domain_info.registrar || "-"}</span>
                </div>

                <div
                  class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                >
                  <div class="flex items-center gap-1 text-gray-400">
                    <span>Domain Age:</span>
                    <TooltipIcon
                      text="How long ago the domain was first registered. Older domains often suggest more established or legitimate websites."
                    />
                  </div>
                  <span class="font-medium text-white">{data.domain_info.age_human}</span>
                </div>

                <div
                  class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                >
                  <div class="flex items-center gap-1 text-gray-400">
                    <span>DNSSEC Enabled:</span>
                    <TooltipIcon
                      text="A security feature that helps protect against DNS tampering and redirection attacks by digitally signing DNS data."
                    />
                  </div>
                  {#if data.domain_info.dnssec}
                    <span class="text-green-400 font-medium flex items-center gap-1">✅ Yes</span>
                  {:else}
                    <span class="text-red-400 font-medium flex items-center gap-1">❌ No</span>
                  {/if}
                </div>

                <div
                  class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                >
                  <div class="flex items-center gap-1 text-gray-400">
                    <span>Created:</span>
                    <TooltipIcon
                      text="The date when this domain was first registered and became active on the internet."
                    />
                  </div>
                  <span class="font-medium text-white">{data.domain_info.created}</span>
                </div>

                <div
                  class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                >
                  <div class="flex items-center gap-1 text-gray-400">
                    <span>Updated:</span>
                    <TooltipIcon
                      text="The last date the domain registration information was modified (e.g., contact change or nameserver update)."
                    />
                  </div>
                  <span class="font-medium text-white">{data.domain_info.updated}</span>
                </div>

                <div
                  class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                >
                  <div class="flex items-center gap-1 text-gray-400">
                    <span>Expiry:</span>
                    <TooltipIcon
                      text="The date when this domain’s registration will expire unless renewed by the owner."
                    />
                  </div>
                  <span class="font-medium text-white">{data.domain_info.expiry}</span>
                </div>

                {#if data.domain_info.nameservers?.length}
                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>Nameservers:</span>
                      <TooltipIcon
                        text="Servers responsible for directing internet traffic to the correct web host. They link your domain to its hosting service."
                      />
                    </div>
                    <div class="flex flex-wrap gap-2">
                      {#each data.domain_info.nameservers as ns}
                        <span class="px-2 py-1 bg-gray-800 text-white rounded text-xs">{ns}</span>
                      {/each}
                    </div>
                  </div>
                {/if}

                {#if data.domain_info.status?.length}
                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>Status:</span>
                      <TooltipIcon
                        text="Domain lifecycle or control states — such as 'active', 'clientTransferProhibited', or 'pendingDelete'. They indicate administrative or operational restrictions."
                      />
                    </div>
                    <div class="flex flex-wrap gap-2">
                      {#each data.domain_info.status as st}
                        <span class="px-2 py-1 bg-gray-800 text-white rounded text-xs">{st}</span>
                      {/each}
                    </div>
                  </div>
                {/if}
              </div>
            </section>
          {/if}

          <!-- Analysis Section -->
          {#if data?.analysis}
            <section
              id="section-analysis"
              class="bg-gray-900/80 border border-gray-800 rounded-lg p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all scroll-mt-20"
            >
              <div class="flex items-center justify-between mb-4">
                <h3 class="text-base font-semibold text-white">Website Analysis</h3>
                <span
                  class="text-[10px] text-gray-400 uppercase tracking-wide px-2 py-0.5 bg-gray-800 rounded"
                  >HTTP / Redirects</span
                >
              </div>

              <div
                class="space-y-0 divide-y divide-gray-800 text-sm text-gray-200 max-w-4xl w-full mx-auto"
              >
                {#if data.analysis.redirection_result}
                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>Redirected:</span>
                      <TooltipIcon
                        text="Indicates whether visiting given URL automatically forwards you to another URL."
                      />
                    </div>
                    {#if data.analysis.redirection_result.is_redirected}
                      <span class="text-gray-400 font-medium flex items-center gap-1"> Yes</span>
                    {:else}
                      <span class="text-gray-400 font-medium flex items-center gap-1"> No</span>
                    {/if}
                  </div>

                  {#if data.analysis.redirection_result.final_url}
                    <div
                      class="flex flex-col md:grid md:grid-cols-[350px,1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                    >
                      <div class="flex items-center gap-1 text-gray-400">
                        <span>Final URL Domain:</span>
                        <TooltipIcon
                          text="The domain where the visitor finally lands after any redirects. Useful to detect domain changes or phishing redirects."
                        />
                      </div>
                      <span class="font-medium text-white break-all"
                        >{data.analysis.redirection_result.final_url_domain}</span
                      >
                    </div>
                    <div
                      class="flex flex-col md:grid md:grid-cols-[350px,1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                    >
                      <div class="flex items-center gap-1 text-gray-400">
                        <span>Final URL:</span>
                        <TooltipIcon
                          text="The complete URL where the user ends up after all redirections."
                        />
                      </div>
                      <span class="font-medium text-white break-all"
                        >{data.analysis.redirection_result.final_url}</span
                      >
                    </div>
                  {/if}

                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>Domain Jumped to Another Domain:</span>
                      <TooltipIcon
                        text="Checks if the website redirects to a completely different domain, which can indicate phishing or tracking."
                      />
                    </div>
                    {#if data.analysis.redirection_result.has_domain_jump}
                      <span class="text-red-400 font-medium flex items-center gap-1">❌ Yes</span>
                    {:else}
                      <span class="text-green-400 font-medium flex items-center gap-1">✅ No</span>
                    {/if}
                  </div>

                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>Redirection Chain Length:</span>
                      <TooltipIcon
                        text="Shows how many redirect steps the website takes before reaching the final destination."
                      />
                    </div>
                    <span class="font-medium text-white"
                      >{data.analysis.redirection_result.chain_length}</span
                    >
                  </div>

                  {#if data.analysis.redirection_result.chain?.length}
                    <div
                      class="flex flex-col md:grid md:grid-cols-[350px,1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                    >
                      <div class="flex items-center gap-1 text-gray-400">
                        <span>Redirection Chain:</span>
                        <TooltipIcon
                          text="A step-by-step list of all URLs in the redirection path. Warning icons highlight jumps to unexpected domains."
                        />
                      </div>
                      <ul class="text-sm text-gray-100 list-none">
                        {#each data.analysis.redirection_result.chain as url, index}
                          <li class="break-all flex items-center gap-2 mb-1">
                            <span class="text-gray-400">{index + 1}.</span>
                            <span class="font-medium text-white">{url}</span>
                            {#if url.includes(data.domain) === false}
                              <span class="text-red-400 text-xs">⚠️</span>
                            {/if}
                          </li>
                        {/each}
                      </ul>
                    </div>
                  {/if}
                {/if}

                {#if data.analysis.http_status}
                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>HTTP Status Code:</span>
                      <TooltipIcon
                        text="The server response code returned when accessing the URL (e.g., 200 = OK, 404 = Not Found)."
                      />
                    </div>
                    <span class="font-medium text-white"
                      >{data.analysis.http_status.code} {data.analysis.http_status.text}</span
                    >
                  </div>

                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>Redirection Status Code (3xx):</span>
                      <TooltipIcon
                        text="Indicates whether the HTTP status is a redirection (3xx) code, which automatically sends visitors to another URL."
                      />
                    </div>
                    {#if data.analysis.http_status.is_redirect}
                      <span class="text-red-400 font-medium flex items-center gap-1">❌ Yes</span>
                    {:else}
                      <span class="text-green-400 font-medium flex items-center gap-1">✅ No</span>
                    {/if}
                  </div>
                {/if}

                {#if data.analysis.is_hsts_supported !== undefined}
                  <div
                    class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                  >
                    <div class="flex items-center gap-1 text-gray-400">
                      <span>HSTS Supported (HTTPS Only):</span>
                      <TooltipIcon
                        text="Shows if the website enforces HTTPS connections automatically to improve security and prevent attacks."
                      />
                    </div>
                    {#if data.analysis.is_hsts_supported}
                      <span class="text-green-400 font-medium flex items-center gap-1">✅ Yes</span>
                    {:else}
                      <span class="text-red-400 font-medium flex items-center gap-1">❌ No</span>
                    {/if}
                  </div>
                {/if}
              </div>
            </section>
          {/if}

          <!-- Performance Section -->
          {#if data?.performance}
            <section
              id="section-performance"
              class="bg-gray-900/80 border border-gray-800 rounded-lg p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all scroll-mt-20"
            >
              <div class="flex items-center justify-between mb-4">
                <h3 class="text-base font-semibold text-white">Performance</h3>
                <span
                  class="text-[10px] text-gray-400 uppercase tracking-wide px-2 py-0.5 bg-gray-800 rounded"
                  >Timings</span
                >
              </div>

              <div
                class="space-y-0 divide-y divide-gray-800 text-sm text-gray-200 max-w-4xl w-full mx-auto"
              >
                <div
                  class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                >
                  <span class="text-gray-400">Total time:</span>
                  <span class="font-medium text-white">{data.performance.total_time}</span>
                </div>

                {#if data.performance.timings && data.performance.timings.length > 0}
                  <ul class="text-sm text-gray-100 space-y-0 divide-y divide-gray-800">
                    {#each data.performance.timings as timing}
                      <li
                        class="flex flex-col md:grid md:grid-cols-[350px,1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
                      >
                        <span class="text-gray-400">{timing.task}</span>
                        <span>{timing.time}</span>
                      </li>
                    {/each}
                  </ul>
                {/if}
              </div>
            </section>
          {/if}

          <!-- Scroll Back to Top Button -->
          <div class="mt-6 flex justify-center">
            <button
              class="inline-flex items-center justify-center gap-2 px-5 py-3 rounded-full bg-gray-800 hover:bg-gray-700 text-white text-sm font-medium focus:outline-none focus-visible:ring-2 focus-visible:ring-emerald-400 focus-visible:ring-offset-0 transition-colors duration-150"
              on:click={() => {
                const target = document.getElementById("analysis-summary");
                if (target) {
                  // Scroll so the button is 50px from top
                  const offset = 50;
                  const top = target.getBoundingClientRect().top + window.scrollY - offset;
                  window.scrollTo({ top, behavior: "smooth" });
                }
              }}
            >
              <svg class="w-4 h-4" viewBox="0 0 20 20" fill="currentColor">
                <path
                  fill-rule="evenodd"
                  d="M5.23 12.79a.75.75 0 001.06.02L10 8.812l3.71 3.998a.75.75 0 101.08-1.04l-4.25-4.53a.75.75 0 00-1.08 0l-4.25 4.53a.75.75 0 00.02 1.06z"
                  clip-rule="evenodd"
                />
              </svg>
              <!-- Scroll to Top -->
            </button>
          </div>
        </div>
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
