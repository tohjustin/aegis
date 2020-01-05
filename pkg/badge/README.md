# badge

A package for generating badges.

## Usage

Sample program for using the package:

```go
package main

import "github.com/tohjustin/badger/pkg/badge"

func main() {
  badgeOptions := badge.Options{
    Color: "318FE0",
    Icon:  "brands/font-awesome",
    Style: badge.ClassicStyle,
  }
  generatedBadge, _ := badge.Create("Font Awesome", "v5.6.3", &badgeOptions)
}
```
