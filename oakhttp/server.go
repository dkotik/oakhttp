// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
// go func() {
//     <-sigs
//     log.Println("got interruption signal")
//     if err := srv.Shutdown(context.TODO()); err != nil {
//         log.Printf("server shutdown returned an err: %v\n", err)
//     }
//     close(done)
// }()
//
// if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
//     log.Fatalf("listen and serve returned err: %v", err)
// }
// <-done
// userService.Stop()

// Event better:
// func main() {
// 	// Make a signal-based context. The stop function, when called, unregisters
// 	// the signals and restores the default signal behaviour.
// 	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
//
// 	// Wait for the signal.
// 	<-ctx.Done()
// 	stop() // After calling stop, another SIGINT will terminate the program.
//
// 	fmt.Println("Interrupted. Exiting.")
//
// 	// Long clean-up code goes here.
// 	time.Sleep(5 * time.Second)
// }
// more: https://dev.to/mokiat/proper-http-shutdown-in-go-3fji
