# background

A tiny library to make it easier to run multiple top-level goroutines in the same application:

- [runner](runner.go) provides an abstraction that handles shutdown signal and wait group (so that goroutines don't need to care about it)
- [job](job/job.go) represents logic that is run inside top-level goroutine. There are 2 jobs that I use most often and therefore provided by default:
  - [recurring job](job/recurring.go)
    - run [work](job/work.go) periodically (using Go ticker)
    - an example can be found [here](examples/recurring_job/main.go). This example spins up 2 goroutines that run ticker periodically. 2nd job has some custom clean up logic
  - [blocking job](job/blocking.go)
    - run [work](job/work.go) and block waiting for it to complete (.e.g. a http server)
    - an example can be found [here](examples/blocking_job/main.go). This example runs a recurring job in 1 goroutine and a http server in the other. The http server has custom clean up logic
  - when a `work` implements `CleanUp` interface from [job](job/work.go), both recurring and blocking jobs above will execute clean up logic when they receive shutdown signal
