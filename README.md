# An asynchronous hook for [Logrus](https://github.com/sirupsen/logrus)

[![Build][Build-Status-Image]][Build-Status-Url] [![Codecov][codecov-image]][codecov-url] [![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]

## Quick Start

### Download and install

```bash
$ go get -u -v github.com/LyricTian/logrus-async-hook
```

### Usage

```go
package main

import (
    "github.com/LyricTian/logrus-async-hook"
    "github.com/sirupsen/logrus"
)

func main() {
    hook := asynchook.New(...)
    defer hook.Flush()

    log := logrus.New()
    log.AddHook(hook)
}
```

## MIT License

    Copyright (c) 2018 Lyric

[Build-Status-Url]: https://travis-ci.org/LyricTian/logrus-async-hook
[Build-Status-Image]: https://travis-ci.org/LyricTian/logrus-async-hook.svg?branch=master
[codecov-url]: https://codecov.io/gh/LyricTian/logrus-async-hook
[codecov-image]: https://codecov.io/gh/LyricTian/logrus-async-hook/branch/master/graph/badge.svg
[reportcard-url]: https://goreportcard.com/report/github.com/LyricTian/logrus-async-hook
[reportcard-image]: https://goreportcard.com/badge/github.com/LyricTian/logrus-async-hook
[godoc-url]: https://godoc.org/github.com/LyricTian/logrus-async-hook
[godoc-image]: https://godoc.org/github.com/LyricTian/logrus-async-hook?status.svg
[license-url]: http://opensource.org/licenses/MIT
[license-image]: https://img.shields.io/npm/l/express.svg
