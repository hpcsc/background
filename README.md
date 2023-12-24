# background

A tiny library to make it easier to run multiple top-level goroutines in the same application:

- [runner](runner.go) provides an abstraction that handles shutdown signal and wait group (so that goroutines don't need to care about it)
- [job](job/job.go) represents logic that is run inside goroutine. There are 2 jobs that I use most often and therefore provided by default:
  - [recurring job](job/recurring.go)
    - run some logic periodically (using Go ticker)
    - an example can be found [here](examples/recurring_job/main.go). This example spins up 2 goroutines that run ticker periodically. 2nd job has some custom clean up logic
  - [blocking job](job/blocking.go)
    - run some logic and block waiting for something to happen (.e.g. a http server)
    - an example can be found [here](examples/blocking_job/main.go). This example runs a recurring job in 1 goroutine and a http server in the other. The http server has custom clean up logic
