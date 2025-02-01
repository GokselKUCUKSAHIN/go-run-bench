# go-run-bench

[![][workflow-badge]][workflow-actions] [![Release][release-badge]][release] [![][license-badge]][license]

Go Run Bench is a command-line tool for benchmarking Go tests in specified directories. It allows you to execute benchmarks with customizable parameters such as cooldown periods, benchmark durations, and output formats. The tool is designed to easily integrate with your workflow for efficient benchmarking and performance analysis.

## Features:
- Cooldown: Configure the cooldown period between benchmark runs.
- Benchmark Duration: Set the duration for each benchmark run.
- Memory Benchmarking: Optionally include memory usage metrics in the benchmark results.
- Flexible Output: Save benchmark results in JSON or CSV format.
- Custom Working Directory: Specify a custom directory to run benchmarks on.
- This tool is ideal for developers who need a lightweight, configurable solution for running and recording benchmark results in Go projects.

# Install

To use this tool, ensure that Go is installed on your system. Then, download the repository and follow the build instructions provided in the project.
```sh
go install github.com/GokselKUCUKSAHIN/go-run-bench
```

# Examples
Running with default settings (20 seconds cooldown, 5 seconds benchmark duration, saving results in JSON format, using the current directory):
```sh
go-run-bench
```

Setting a 10-second cooldown, 5-second benchmark duration, saving results in CSV format, and specifying a working directory:
```sh
go-run-bench -cooldown=10 -benchtime=5 -save=csv -wd=/absolute/path/to/directory
```

Disabling the cooldown (0 seconds), setting a 10-second benchmark duration, saving results in JSON format:
```sh
go-run-bench -cooldown=0 -benchtime=10 -save=json
```

Setting a 5-second cooldown, a 30-second benchmark duration, saving results in CSV format, and using the current directory:
```sh
go-run-bench -cooldown=5 -benchtime=30 -save=csv
```

Disabling the cooldown (0 seconds), setting a 3-second benchmark duration, saving results in CSV format, and specifying a different working directory:
```sh
go-run-bench -cooldown=5 -benchtime=30 -save=csv
```

Running with default settings, including memory benchmarking (20 seconds cooldown, 5 seconds benchmark duration, saving results in JSON format):
```sh
go-run-bench -benchmem=true
```

Displaying the help message for go-run-bench:
```sh
go-run-bench -help
go-run-bench --help
go-run-bench help
```

# License

MIT - Please check the [LICENSE](./LICENSE) file for full text.

[workflow-actions]: https://github.com/GokselKUCUKSAHIN/go-run-bench/actions
[workflow-badge]: https://github.com/GokselKUCUKSAHIN/go-run-bench/workflows/build/badge.svg
[license]:https://github.com/GokselKUCUKSAHIN/go-run-bench/blob/main/LICENSE
[license-badge]:https://img.shields.io/badge/License-MIT-blue.svg
[release]: https://github.com/GokselKUCUKSAHIN/go-run-bench/releases
[release-badge]: https://img.shields.io/github/v/release/GokselKUCUKSAHIN/go-run-bench.svg