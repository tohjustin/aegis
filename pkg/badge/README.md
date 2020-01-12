# badge

A package for generating SVG badges.

## Usage

Sample program:

```go
package main

import "github.com/tohjustin/badger/pkg/badge"

func main() {
  generatedBadge, _ := badge.Create(&badge.Params{
    Subject: "Font Awesome",
    Status:  "v5.12.0",
    Color:   "#318FE0",
    Icon:    "brands/font-awesome",
    Style:   badge.ClassicStyle,
  })
}
```
