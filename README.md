# retry [![Build Status](https://travis-ci.org/flowchartsman/retry.svg?branch=master)](https://travis-ci.org/flowchartsman/v8) [![Go Report Card](https://goreportcard.com/badge/github.com/flowchartsman/retry)](https://goreportcard.com/report/github.com/flowchartsman/retry) [![GoDoc](https://godoc.org/github.com/flowchartsman/retry?status.svg)](https://godoc.org/github.com/flowchartsman/retry)

retry is a simple retrier for golang with exponential backoff[1][2] and context support.

[1]: https://en.wikipedia.org/wiki/Exponential_backoff
[2]: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/