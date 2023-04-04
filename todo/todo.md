Your services should have a `/server` package that includes all the database plumbing, adapters with error handling, and the router that takes in domain methods from the root package. There should be no HTTP primitives in the domain logic whatsoever. This package is a scaffold of a minimal security gearbox that prepares your server for production.

- [ ] global rate limit
- [ ] session-based rate limiter? with a fallback to IP-based rate limiter
- [ ] alerts for %-based failures on requests, say log critical when 10% of requests are failing.

# Primary

- [ ] remove alike oakwords
- [ ] is `go get cuelang.org/go/pkg/uuid@v0.4.0` better UUID package?
- [ ] consider simpler Permissions matching to map[string]string
  - allow domain:glob=\* id:any forum:regex=wowza

# Secondary

- [ ] Study https://github.com/TwiN/g8
- [ ] Study Zanzibar and Authzed https://authzed.com/
- [ ] Study casbin: https://casbin.org/docs/en/how-it-works
- [ ] Study OIDC: https://github.com/XenitAB/go-oidc-middleware
- [ ] Study session offerings:
  - [ ] <https://github.com/alexedwards/scs>
  - [ ] https://github.com/Acebond/session/blob/main/session.go
- [ ] Study passwordless:
    - [ ] https://github.com/teamhanko/hanko
