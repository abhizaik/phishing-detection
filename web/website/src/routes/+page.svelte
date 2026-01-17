<script lang="ts">
  import { replaceState } from "$app/navigation";
  import { onMount } from "svelte";
  import { api } from "../lib/api";
  import ResultSection from "../lib/components/ResultSection.svelte";
  import type { AnalyzeResult } from "../lib/types";
  import { formatUrl, isValidUrl } from "../lib/utils";

  let input = "";
  let loading = false;
  let error: string | null = null;
  let data: AnalyzeResult | null = null;
  let screenshotUrl: string | null = null;

  type Verdict = "Safe" | "Risky" | "Unclear" | "Suspicious";
  const ACCENTS: Record<Verdict, { ring: string; glow: string; badge: string }> = {
    Safe: {
      ring: "focus:ring-emerald-600",
      glow: "from-emerald-600/20",
      badge: "bg-emerald-600/20 text-emerald-300 border-emerald-700",
    },
    Risky: {
      ring: "focus:ring-red-600",
      glow: "from-red-600/20",
      badge: "bg-red-600/20 text-red-300 border-red-700",
    },
    Unclear: {
      ring: "focus:ring-yellow-600",
      glow: "from-yellow-500/20",
      badge: "bg-yellow-600/20 text-yellow-300 border-yellow-700",
    },
    Suspicious: {
      ring: "focus:ring-yellow-600",
      glow: "from-yellow-500/20",
      badge: "bg-yellow-600/20 text-yellow-300 border-yellow-700",
    },
  };

  function normalizeVerdict(v: string | null | undefined): Verdict {
    switch (v) {
      case "Safe":
      case "Risky":
      case "Unclear":
      case "Suspicious":
        return v;
      default:
        return "Unclear";
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
      error = "Please enter a valid URL";
      return;
    }

    loading = true;
    error = null;
    data = null;
    // Clean up previous screenshot blob URL if it exists
    if (screenshotUrl) {
      URL.revokeObjectURL(screenshotUrl);
      screenshotUrl = null;
    }

    try {
      // kick off screenshot but don't block on it
      api
        .screenshot(url)
        .then((res) => {
          if (res.data) {
            screenshotUrl = res.data as string;
          }
        })
        .catch(() => {
          console.warn("Screenshot request failed");
        });

      // await analyze so loading reflects this call only
      const res = await api.analyze(url);

      if (res.error) {
        error = res.error;
      } else {
        data = res.data as AnalyzeResult;

        const share = new URL(window.location.href);
        share.searchParams.set("q", url);
        replaceState(share.toString(), {});
      }
    } catch (err) {
      error = "Analyze request failed";
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
    const q = params.get("q");
    if (q) {
      input = q;
      runAnalyze(q);
    }

    const onKey = (e: KeyboardEvent) => {
      if (
        e.key === "/" &&
        !(e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement)
      ) {
        e.preventDefault();
        const el = document.getElementById("url-input") as HTMLInputElement | null;
        el?.focus();
      }
      if (e.key === "Escape") {
        input = "";
      }
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  });
</script>

<section>
  <title>SafeSurf</title>
  <div
    class={`max-w-4xl mx-auto px-6 ${isLanding ? "min-h-[70vh] flex flex-col justify-center" : "py-12"}`}
  >
    <header
      class="relative mb-14 flex flex-col items-center md:items-start text-center md:text-left"
    >
      <!-- Background accent -->
      <div
        class="absolute -top-10 -left-10 w-40 h-40 bg-blue-600/30 rounded-full blur-3xl animate-blob z-0"
      ></div>
      <div
        class="absolute top-0 right-0 w-32 h-32 bg-emerald-500/20 rounded-full blur-3xl animate-blob animation-delay-2000 z-0"
      ></div>

      <!-- Heading -->
      <h1 class="relative text-6xl md:text-6xl font-extrabold tracking-tight text-white z-10">
        <a
          href="/"
          on:click={() => (location.href = "/")}
          class="bg-clip-text text-transparent bg-gradient-to-r from-blue-400 via-indigo-500 to-purple-500 hover:from-purple-500 hover:to-pink-500 transition-all"
        >
          SafeSurf
        </a>
      </h1>

      <!-- Subheading -->
      <p
        class="relative mt-4 text-gray-300 md:text-lg text-center md:text-left max-w-xl z-10 animate-fadeIn"
      >
        Check if a link is safe.
      </p>
    </header>

    <style>
      /* Blob animation */
      @keyframes blob {
        0%,
        100% {
          transform: translate(0px, 0px) scale(1);
        }
        33% {
          transform: translate(20px, -10px) scale(1.1);
        }
        66% {
          transform: translate(-15px, 15px) scale(0.95);
        }
      }
      .animate-blob {
        animation: blob 8s infinite;
      }
      .animation-delay-2000 {
        animation-delay: 2s;
      }

      /* Fade-in animation for subheading */
      @keyframes fadeIn {
        from {
          opacity: 0;
          transform: translateY(10px);
        }
        to {
          opacity: 1;
          transform: translateY(0);
        }
      }
      .animate-fadeIn {
        animation: fadeIn 1s ease-out forwards;
      }

      /* Autofill fix for dark input added by browsers */
      input:-webkit-autofill,
      input:-webkit-autofill:hover,
      input:-webkit-autofill:focus,
      input:-webkit-autofill:active {
        -webkit-text-fill-color: #e5e7eb; /* Tailwind text-gray-200 */
        transition: background-color 5000s ease-in-out 0s; /* prevent yellow flash */
        box-shadow: 0 0 0px 1000px #1f2937 inset; /* Tailwind bg-gray-900 */
        -webkit-box-shadow: 0 0 0px 1000px #1f2937 inset;
      }
    </style>

    <form
      class="relative bg-gray-950 rounded-xl border border-gray-800 p-6 md:p-8 overflow-hidden"
      on:submit|preventDefault={onSubmit}
    >
      <!-- Background accent -->
      <div
        class="absolute -top-10 -left-10 w-32 h-32 bg-blue-600/20 rounded-full blur-3xl animate-blob z-0"
      ></div>
      <div
        class="absolute -bottom-10 -right-10 w-28 h-28 bg-purple-500/20 rounded-full blur-3xl animate-blob animation-delay-3000 z-0"
      ></div>

      <label for="url-input" class="sr-only">URL to analyze</label>
      <div class="relative flex flex-col md:flex-row gap-3 z-10">
        <!-- Input -->
        <input
          id="url-input"
          type="text"
          class={`flex-1 rounded-lg bg-gray-900 border border-gray-800 px-4 py-3 text-sm placeholder-gray-500 text-gray-200 focus:outline-none focus:ring-2 focus:ring-offset-0 transition-all duration-200 ${accent.ring} focus:shadow-lg`}
          placeholder="Enter a URL (e.g. example.com)"
          bind:value={input}
          autocomplete="url"
          inputmode="url"
          required
        />

        <!-- Button -->
        <button
          type="submit"
          class="inline-flex items-center justify-center gap-2 px-5 py-3 rounded-lg bg-blue-600 hover:bg-blue-500 text-white text-sm font-medium disabled:opacity-50 focus:outline-none focus:ring-2 focus:ring-offset-0 focus:ring-blue-600"
          disabled={loading}
          aria-busy={loading}
        >
          {#if loading}
            <svg
              class="w-4 h-4 animate-spin text-white"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M12 4v4m0 8v4m8-8h-4M4 12H0"
              />
            </svg>
            Scanning ..
          {:else}
            Scan Now
          {/if}
        </button>
      </div>
    </form>

    <style>
      /* Blob animation */
      @keyframes blob {
        0%,
        100% {
          transform: translate(0px, 0px) scale(1);
        }
        33% {
          transform: translate(20px, -10px) scale(1.1);
        }
        66% {
          transform: translate(-15px, 15px) scale(0.95);
        }
      }
      .animate-blob {
        animation: blob 8s infinite;
      }
      .animation-delay-3000 {
        animation-delay: 3s;
      }
    </style>

    <div class="mt-8" aria-live="polite">
      <ResultSection {data} {loading} {error} {screenshotUrl} />
    </div>
  </div>
</section>
