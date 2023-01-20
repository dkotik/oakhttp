# Oak Access Control System

## TODO

- [ ] Explain how ReBAC is implemented using oakrbac.Predicate to address the distinction explained here: https://dev.to/egeaytin/rbac-vs-rebac-when-to-use-them-47c4
- [ ] Figure out the best way to mitigate CSRF:
    - [ ] https://dev.to/justlorain/use-session-middleware-to-implement-distributed-session-solution-based-on-redis-5b65 (example is buried)

## Features

1. Safety: static types, immutability, and proper defaults.
2. Minimalism: tracks the least amount of information possible without compromising safety.
3. Consistency: assumes single fully synchronized source of truth. Revocations are instant.
4. Flexibility: simple, independent, and configurable models that support multiple back-ends.

## Packages

- [oakrbac](https://pkg.go.dev/github.com/dkotik/oakacs/oakrbac): role-based access control

    [![GoDoc](https://godoc.org/github.com/dkotik/oakacs/oakrbac?status.png)](https://pkg.go.dev/github.com/dkotik/oakacs/oakrbac?tab=doc)
    [![Go report card](https://goreportcard.com/badge/github.com/dkotik/oakacs/oakrbac)](https://goreportcard.com/report/github.com/dkotik/oakacs/oakrbac)
    [![Test coverage](http://gocover.io/_badge/github.com/dkotik/oakacs/oakrbac)](https://gocover.io/github.com/dkotik/oakacs/oakrbac)

## Functions

1. Humanity Recognition
2. Throttling
   - Prevent password-reuse?
3. Timing modulation
4. Registration
5. Authentication
   - Password policy
   - Revocation
   - Kill switch
   - Recovery
6. Authorization
7. Observability
   - Logging

## Model

- Identity: provides authentication.
- Group: enumerates roles which are available for identities.
- Role: provides authorization by granular permissions.
  - Permission
    - Service
    - Domain
    - Resource
    - Action
- Session: the result of pairing identity to a role.
- Token: one-time utility codes.

## Security

The library is created in ways that anticipate misconfiguration by aiming at simplicity.

1. All roles and policies deny access by default.
2. Policies and predicates must return explicit sentinel value `Allow`.
3. Comes with a code generation tool that helps build tight access control policies and **test cases**.

## Logging

Logging can be approached in several different ways:

1. By writing a Policy wrapper. Use the function `WithLogger` for an example.
2. Inside the policies themselves.
3. At a higher level with request logs.

## Interesting Access Control Projects

### Authorization

- [goRBAC](https://github.com/mikespook/gorbac)
- [authzed SpiceDB](https://github.com/authzed)
- [Keto](https://github.com/ory/keto)

### Authentication

- [OpenFGA, Zanzibar-based](https://github.com/openfga/openfga)
- [Authelia](https://github.com/authelia/authelia)
- [Ballerine](https://github.com/ballerine-io/ballerine)

### Observability

- [Ntfy.sh](https://ntfy.sh/docs/publish/)
