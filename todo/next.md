## Access control

- Oauth as credential type
- Re-authenticate users prior to performing critical operations.
- Use MFA for highly sensitive or high value transactional accounts.
- Implement monitoring to identify attacks against multiple user accounts using the same password. This attack is used to bypass standard lockouts, when user IDs can be harvested or guessed.
- Change all vendor-supplied default passwords and user IDs or disable the associated accounts.
- Enforce account disabling after an established number of invalid login attempts. Five attempts is common. The account must be disabled for a period of time sufficient to discourage brute force guessing of credentials, but not so long as to allow for a denial-of-service attack to be performed.

## Notes from paper

```go
NewSession()
NewHumanSession()?
ACS.NewSessionFromToken()?
session.Authorize(ctx, action string) error
ACS.Authorize(ctx, sessionUUID, action) error
```

## Need to add full Secrets Manager Interface

## 3-Phase Login

- Client requests salt after form is filled out
- Server sends salt (signed with a time-limit), determenistically generated if client ID is missing, real salt if not, time-equalized both
- Client hashes using the given salt, sends just the hash, never shares their password

## Timing attack mitigation

- Implement time-constant comparison logic where possible
- Modulate with noise?
- Add CAPTCHAs to forms with user interactions
- Test for timing differences in critical operations
