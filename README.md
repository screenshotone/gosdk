# paddle

[![Build](https://github.com/screenshotone/gosdk/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/screenshotone/gosdk/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/screenshotone/gosdk/branch/main/graph/badge.svg?token=rh8BDdHc2v)](https://codecov.io/gh/screenshotone/gosdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/screenshotone/gosdk)](https://goreportcard.com/report/github.com/screenshotone/gosdk)
[![GoDoc](https://godoc.org/https://godoc.org/github.com/screenshotone/gosdk?status.svg)](https://godoc.org/github.com/screenshotone/gosdk)

An official [Screenshot API (screenshotone.com)](https://screenshotone.com/) client for Go. 

## Installation

```shell
go get github.com/krasun/paddle
```

## Usage

Import the library: 
```go
import "github.com/krasun/paddle"
```

Generate a screenshot URL without executing request: 
```go
client, err := screenshots.NewClient("IVmt2ghj9TG_jQ", "Sxt94yAj9aQSgg")
if err != nil {
    // ...
}

options := screenshots.NewTakeOptions("https://scalabledeveloper.com").
    Format("png").
    FullPage(true).
    DeviceScaleFactor(2).
    BlockAds(true).
    BlockTrackers(true)

u, err := client.GenerateTakeURL(options)
if err != nil {
    // ...
}

fmt.Println(u.String())
// Output: https://api.screenshotone.com/take?access_key=IVmt2ghj9TG_jQ&block_ads=true&block_trackers=true&device_scale_factor=2&format=png&full_page=true&url=https%3A%2F%2Fscalabledeveloper.com&signature=85aabf7ac251563ec6158ef6839dd019bb79ce222cc85288a2e8cea0291a824e
```

Take a screenshot and save the image in the file: 
```go 
client, err := screenshots.NewClient("IVmt2ghj9TG_jQ", "Sxt94yAj9aQSgg")
if err != nil {
    // ...
}

options := screenshots.NewTakeOptions("https://example.com").
    Format("png").
    FullPage(true).
    DeviceScaleFactor(2).
    BlockAds(true).
    BlockTrackers(true)

image, err := client.Take(context.TODO(), options)
if err != nil {
    // ...
}

defer image.Close()
out, err := os.Create("example.png")
if err != nil {
    // ...
}

defer out.Close()
io.Copy(out, image)
```

## Tests 

To run tests, just execute: 
```
$ go test . 
```

## License 

`screenshotone/gosdk` is released under [the MIT license](LICENSE).