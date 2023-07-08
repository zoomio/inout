# Contributing

## Testing

- benchmarks are `func BenchmarkXxx(b *testing.B) {` use `go test -v -run=^$ -bench=<BenchmarkXxx> -cpuprofile=prof.cpu ./<path_to_package>`, then profiling `go tool pprof <package>.test prof.cpu`

## Guidelines for pull requests

- Write tests for any changes (use `./_bin/test.sh` to trigger tests locally).
- Separate unrelated changes into multiple pull requests.
- For bigger changes, make sure you start a discussion first by creating an issue and explaining the intended change.
- Ensure the build is green before you open your PR.

## Build

* [Go](https://golang.org/dl/)

## Release

* All notable changes comming with the new version should be documented in [CHANGELOG.md](https://raw.githubusercontent.com/zoomio/inout/master/CHANGELOG.md).
* Use `./_bin/tag.sh <x.y.z>` to tag, push and trigger new release. 
