package server

// TODO: add base context! so that requests are .Deadline-aware!!! Otherwise http package uses context.Background by default. Make this into WithBaseContext() option.
// httpServer := &http.Server{
//     Addr: ":8000",
//     Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         ctx := r.Context()
//
//         for {
//             select {
//             case <-ctx.Done():
//                 fmt.Println("Graceful handler exit")
//                 w.WriteHeader(http.StatusOK)
//                 return
//             case <-time.After(1 * time.Second):
//                 fmt.Println("Hello in a loop")
//             }
//         }
//     }),
//     BaseContext: func(_ net.Listener) context.Context {
//         return mainCtx
//     },
// }
