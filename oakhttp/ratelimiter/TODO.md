# TODO

- [ ] Split map[IP]*rate.Limiter out as its own type, so that Middleware can be written with it without duplicating the logic.
- [ ] Add a request-examiner for rate limiting
- [ ] Add a stackable request-examiner as exemplified [here](https://larasec.substack.com/p/security-tip-multiple-rate-limits?utm_source=twitter&utm_campaign=auto_share&r=54brl):

    ```php
    RateLimiter::for('login', function (Request $request) {
        return [
            Limit::perMinute(500),
            Limit::perMinute(5)->by($request->ip()),
            Limit::perMinute(5)->by($request->input('email')),
        ];
    });
    ```

## Ideas

- [ ] SPA router generator with sha256 hashes for security that can be used for sign-in, recovery workflows!
