
Your services should have a `/server` package that includes all the database plumbing, adapters with error handling, and the router that takes in domain methods from the root package. There should be no HTTP primitives in the domain logic whatsoever. This package is a scaffold of a minimal security gearbox that prepares your server for production.

- [ ] timeouts
- [ ] signals
- [ ] panic catcher
- [ ] global rate limit
- [ ] session-based rate limiter? with a fallback to IP-based rate limiter
- [ ] alerts for %-based failures on requests, say log critical when 10% of requests are failing.


## Signals

```
1 	SIGHUP 	Terminate 	Hang up controlling terminal or process. Sometimes used as a signal to reread configuration file for the program.
2 	SIGINT 	Terminate 	Interrupt from keyboard, Ctrl + C.
3 	SIGQUIT 	Dump 	Quit from keyboard, Ctrl + \.
9 	SIGKILL 	Terminate 	Forced termination.
15 	SIGTERM 	Terminate 	Graceful termination.
17 	SIGCHLD 	Ignore 	Child process exited.
18 	SIGCONT 	Continue 	Resume process execution.
19 	SIGSTOP 	Stop 	Stop process execution, Ctrl + Z.
```
