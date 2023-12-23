# background

A tiny library to make it easier to run multiple top-level goroutines in the same application

Running a goroutine with ticker is my most common usage so a [recurring job](job/recurring.go) is provided by default.
[recurring_job example](examples/recurring_job/main.go) gives following output:

```
{"time":"2023-11-15T21:27:00.071747+11:00","level":"INFO","msg":"job started","job":"job-2"}
{"time":"2023-11-15T21:27:00.07175+11:00","level":"INFO","msg":"job started","job":"job-1"}
{"time":"2023-11-15T21:27:00.071933+11:00","level":"INFO","msg":"processing","job":"job-2"}
{"time":"2023-11-15T21:27:00.071934+11:00","level":"INFO","msg":"processing","job":"job-1"}
{"time":"2023-11-15T21:27:03.072711+11:00","level":"INFO","msg":"processing","job":"job-1"}
^Ctask: Signal received: "interrupt"
{"time":"2023-11-15T21:27:04.402977+11:00","level":"INFO","msg":"received interrupt signal"}
{"time":"2023-11-15T21:27:05.071878+11:00","level":"INFO","msg":"ticker stopped","job":"job-2"}
{"time":"2023-11-15T21:27:05.071898+11:00","level":"INFO","msg":"job stopped","job":"job-2"}
{"time":"2023-11-15T21:27:06.072015+11:00","level":"INFO","msg":"ticker stopped","job":"job-1"}
{"time":"2023-11-15T21:27:06.072043+11:00","level":"INFO","msg":"job stopped","job":"job-1"}
{"time":"2023-11-15T21:27:06.072068+11:00","level":"INFO","msg":"exit"}
```

[blocking_job example](examples/blocking_job/main.go) gives an example of a job that runs and blocks (http server)
