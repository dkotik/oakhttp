package oakpolicy

// type powerlessAuthority struct {
// 	cause error
// }
//
// func (p *powerlessAuthority) IsAuthorized(context.Context, Action, Resource) error {
// 	return p.cause
// }
//
// func (p *powerlessAuthority) WithContext(context.Context) Authority {
// 	return p
// }
//
// // NewPowerlessAuthority returns an [Authority] that always fails with the same error. If provided error is <nil>, [ErrInvalidContext] is assumed to be the cause.
// func NewPowerlessAuthority(fromCause error) Authority {
// 	if fromCause == nil {
// 		fromCause = ErrInvalidContext
// 	}
// 	return &helplessness{cause: fromCause}
// }
//
// // type Authorization interface {
// // 	AuthorizedTo(context.Context, Action, Resource) error
// // }
//
// // type Capabilities interface {
// //   IsAuthorized (Action, Resource) error
// // }
//
// // type AuthorizedContext interface {
// // 	To(Action, Resource) error
// // }
//
// type Authority interface {
// 	// IsAuthorized returns <nil> if [context.Context] is empowered to perform a given [Action] on a [Resource]. Returns [AuthorizationError] in case of a failure.
// 	IsAuthorized(context.Context, Action, Resource) error
//
// 	// WithContext binds [Authority] to a given context to allow multiple [Authority.Authorize] calls without having to analyze [context.Context] more than once. Use for batching [Authority.IsAuthorized] calls.
// 	WithContext(context.Context) Authority
// }
