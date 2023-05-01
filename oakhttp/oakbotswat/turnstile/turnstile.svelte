<script>
  import { onMount } from 'svelte';
  import { createEventDispatcher } from "svelte";
  const dispatch = createEventDispatcher();
  // https://github.com/ghostdevv/svelte-turnstile/blob/main/src/lib/Turnstile.svelte

  export let siteKey;
  export let siteAction = "view";
  export let locale = "en";
  let mounted = false;
  let widget;
  let cleanUp;

  const throwError = (err) => {
    dispatch("error", new Error(err));
    dispatch("token", null);
  }

  const loadCallback = () => {
      cleanUp = turnstile.render(widget, {
          language: locale,
          theme: "light",
          size: "normal",
          action: siteAction,
          sitekey: siteKey,
          callback: (token) => {
              dispatch("token", token);
              dispatch("error", null);
          },
          'timeout-callback': () => throwError(new Error("humanity check timed out")),
          'expired-callback': () => throwError(new Error("humanity check expired, please refresh")),
          'unsupported-callback': () => throwError(new Error("your browser does not support Turnstile humanity checks")),
          'error-callback': throwError,
      });
  }

  onMount(() => {
    mounted = true;

    return () => {
        mounted = false;
        // console.log("cleaning up", cleanUp);
        if (cleanUp && window.turnstile) {
          // console.log("cleaned up");
          window.turnstile.remove(cleanUp);
        }
    };
  });

</script>

<svelte:head>
  {#if mounted}
    <script
        src="https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit"
        on:load={loadCallback}
        async defer></script>
  {/if}
</svelte:head>

{#if mounted}
  <div bind:this={widget} class="turnstileWidget" />
{/if}

<style>
.turnstileWidget {
  width: fit-content;
  margin: 2em auto;
}
</style>
