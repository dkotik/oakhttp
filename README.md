# Oak Access Control System

## Features

1. Safety: static types, immutability, and proper defaults.
2. Minimalism: tracks the least amount of information possible without compromising safety.
3. Consistency: assumes single fully synchronized source of truth. Revocations are instant.
4. Flexibility: simple, independent, and configurable models that support multiple back-ends.

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
