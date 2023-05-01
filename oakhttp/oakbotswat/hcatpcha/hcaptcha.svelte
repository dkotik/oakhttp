<script>
    import { onMount } from "svelte";

    export let siteKey;
    export let value;

    let visible = false;
    let mounted = false;
    var resolvePromise = null;

    // https://docs.hcaptcha.com/invisible
    // https://github.com/mrmikardo/svelte-hcaptcha/blob/master/src/HCaptcha.svelte

    /*
    // const cookieName = "HCaptchaProof";
    // const expiration = 30 * 60 * 1000; // 30 minutes
curl https://hcaptcha.com/siteverify \
  -X POST \
  -d "response=CLIENT-RESPONSE&secret=YOUR-SECRET"
*/

    export const getToken = () => {
        visible = true;

        if (value) return Promise.resolve(value);

        return new Promise((resolve, reject) => {
            // if (true) {
            //     visible = false;
            //     resolve("passthrough");
            //     return;
            // }
            // if (resolvePromise) return; // already waiting
            resolvePromise = resolve;
        });
    };

    onMount(() => {
        mounted = true;
        window.HCaptchaReset = () => {
            value = null;
            visible = true;
        };
        window.HCaptchaPass = (token) => {
            value = token;
            // const d = new Date();
            // d.setTime(d.getTime() + expiration);
            // let expires = "expires=" + d.toUTCString();
            // document.cookie =
            //     cookieName + "=" + token + ";" + expires + ";path=/";

            visible = false;
            // dispatch("pass", {
            //     token: token,
            // });
            if (resolvePromise) {
                resolvePromise(token);
                resolvePromise = null;
            }
            // console.log("got token:", token);

            // request humanity refresh when token is about to time out
            // window.setTimeout(() => {
            //     visible = true;
            // }, expiration * 0.9);
        };

        // if (getCookie(cookieName) === "") visible = true;
    });

    // function setCookie(cname, cvalue, minutes) {}
    //
    // function getCookie(cname) {
    //     let name = cname + "=";
    //     let ca = document.cookie.split(";");
    //     for (let i = 0; i < ca.length; i++) {
    //         let c = ca[i];
    //         while (c.charAt(0) == " ") {
    //             c = c.substring(1);
    //         }
    //         if (c.indexOf(name) == 0) {
    //             return c.substring(name.length, c.length);
    //         }
    //     }
    //     return "";
    // }
</script>

<div class="curtain" class:is-hidden={!visible}>
    <fieldset>
        {#if mounted && window?.HCaptchaPass}
            <div
                class="h-captcha"
                data-size="compact"
                data-sitekey={siteKey}
                data-callback="HCaptchaPass"
                data-expired-callback="HCaptchaReset"
                data-chalexpired-callback="HCaptchaReset"
            />
            <script
                src="https://js.hcaptcha.com/1/api.js?recaptchacompat=off"
                async
                defer></script>
        {/if}
        <input type="hidden" name="hcaptcha" {value} />
        <slot>
            <p class="help has-text-centered">
                Please verify that you are human. Verification requires
                Javascript and Cookies enabled.
            </p>
        </slot>
    </fieldset>
</div>

<style>
    .curtain {
        display: block;
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-color: rgba(0, 0, 0, 0.9);
        z-index: 9999999;
    }

    .curtain.is-hidden {
        display: none;
    }

    fieldset {
        margin: 0 auto;
        margin-top: 10vh;
        max-width: 25em;
        background-color: white;
        border: 1px solid black;
        border-radius: 0.6em;
    }

    .h-captcha {
        width: fit-content;
        margin: 0 auto;
        margin-top: 3em;
        margin-bottom: 3em;
    }

    p {
        margin: 1.5em 2em;
    }
</style>
