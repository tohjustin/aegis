# badger

SVG badge generation service, icons powered by [Font Awesome](https://fontawesome.com/https://fontawesome.com/).

[![Go Report Card](https://goreportcard.com/badge/github.com/tohjustin/badger)](https://goreportcard.com/report/github.com/tohjustin/badger)
[![CircleCI](https://circleci.com/gh/tohjustin/badger/tree/master.svg?style=shield&circle-token=fbdca44ece1ce1c6e2a907a530476138578946a2)](https://circleci.com/gh/tohjustin/badger/tree/master)
[![CodeCov](https://codecov.io/gh/tohjustin/badger/branch/master/graph/badge.svg?token=HRJhI2JVS0)](https://codecov.io/gh/tohjustin/badger)
[![GoDoc](https://godoc.org/github.com/tohjustin/badger/pkg/badge?status.svg)](http://godoc.org/github.com/tohjustin/badger/pkg/badge/)
[![FOSSA License Scan](https://app.fossa.io/api/projects/git%2Bgithub.com%2Ftohjustin%2Fbadger.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Ftohjustin%2Fbadger?ref=badge_shield)
[![Font Awesome](https://badger.now.sh/static?icon=brands/font-awesome&subject=Font%20Awesome&status=v5.6.3&color=318FE0)](https://fontawesome.com/)
[![License](https://badger.now.sh/static?subject=license&status=MIT)](https://opensource.org/licenses/MIT)

## Usage

### Query Parameters

Use the following list of available query parameters to configure your badge appearance:

| Query Parameter | Description                                  | Input Format                                                                                             | Example                                      |
| --------------- | -------------------------------------------- | -------------------------------------------------------------------------------------------------------- | --------------------------------------------- |
| color           | Configures or overwrites badge primary color | RGB Hex Values, [CSS Color Keywords](https://developer.mozilla.org/en-US/docs/Web/CSS/color_value)       | "fff", "1BACBF", "mediumturquoise"            |
| icon            | Adds font-awesome icon                       | Any one of the available [Font Awesome Icons](https://fontawesome.com/icons): `<ICON_STYLE>/<ICON_NAME>` | "brands/github", "regular/star", "solid/star" |
| status          | Configures or overwrites badge status text   | -                                                                                                        | "license"                                     |
| style           | Configures badge style                       | Any one of the 4 available badge styles (classic, flat, plastic, semaphore)                              | "classic", "flat", "plastic", "semaphore"     |
| subject         | Configures or overwrites badge subject text  | -                                                                                                        | "MIT"                                         |

### Static Badge Service

| Path    | Description  | Example                                                                                                           |
| ------- |------------- | ----------------------------------------------------------------------------------------------------------------- |
| /static | Static badge | ![static](https://badger.now.sh/static?icon=brands/font-awesome&subject=Font%20Awesome&status=v5.6.3&color=318FE0) |

### Bitbucket Badge Service

[![Bitbucket Cloud REST API](https://badger.now.sh/static?icon=brands/bitbucket&subject=Bitbucket%20Cloud%20REST%20API&status=v2.0)](https://developer.atlassian.com/bitbucket/api/2/reference/)

| Path                                           | Description        | Example                                                                             |
| ---------------------------------------------- |------------------- | ----------------------------------------------------------------------------------- |
| /bitbucket/`<USERNAME>`/`<REPO_SLUG>`/forks         | Forks count         | ![bitbucket/forks](https://badger.now.sh/bitbucket/tutorials/tutorials.bitbucket.org/forks)                 |
| /bitbucket/`<USERNAME>`/`<REPO_SLUG>`/issues<br>/bitbucket/`<USERNAME>`/`<REPO_SLUG>`/issues?state=new<br>/bitbucket/`<USERNAME>`/`<REPO_SLUG>`/issues?state=open<br>/bitbucket/`<USERNAME>`/`<REPO_SLUG>`/issues?state=resolved<br>/bitbucket/`<USERNAME>`/`<REPO_SLUG>`/issues?state=on-hold<br>/bitbucket/`<USERNAME>`/`<REPO_SLUG>`/issues?state=invalid<br>/bitbucket/`<USERNAME>`/`<REPO_SLUG>`/issues?state=duplicate<br>/bitbucket/`<USERNAME>`/`<REPO_SLUG>`/issues?state=wontfix<br>/bitbucket/`<USERNAME>`/`<REPO_SLUG>`/issues?state=closed<br> | Issue count        | ![bitbucket/issues](https://badger.now.sh/bitbucket/tutorials/tutorials.bitbucket.org/issues)<br>![bitbucket/new-issues](https://badger.now.sh/bitbucket/tutorials/tutorials.bitbucket.org/issues?state=new)<br>![bitbucket/open-issues](https://badger.now.sh/bitbucket/tutorials/tutorials.bitbucket.org/issues?state=open)<br>![bitbucket/resolved-issues](https://badger.now.sh/bitbucket/tutorials/tutorials.bitbucket.org/issues?state=resolved)<br>![bitbucket/on-hold-issues](https://badger.now.sh/bitbucket/tutorials/tutorials.bitbucket.org/issues?state=on-hold)<br>![bitbucket/invalid-issues](https://badger.now.sh/bitbucket/tutorials/tutorials.bitbucket.org/issues?state=invalid)<br>![bitbucket/duplicate-issues](https://badger.now.sh/bitbucket/tutorials/tutorials.bitbucket.org/issues?state=duplicate)<br>![bitbucket/wontfix-issues](https://badger.now.sh/bitbucket/tutorials/tutorials.bitbucket.org/issues?state=wontfix)<br>![bitbucket/closed-issues](https://badger.now.sh/bitbucket/tutorials/tutorials.bitbucket.org/issues?state=closed)<br>  |
| /bitbucket/`<USERNAME>`/`<REPO_SLUG>`/pull-requests<br>/bitbucket/`<USERNAME>`/`<REPO_SLUG>`/pull-requests?state=open<br>/bitbucket/`<USERNAME>`/`<REPO_SLUG>`/pull-requests?state=declined<br>/bitbucket/`<USERNAME>`/`<REPO_SLUG>`/pull-requests?state=merged<br>/bitbucket/`<USERNAME>`/`<REPO_SLUG>`/pull-requests?state=superseded<br> | Pull Request count | ![bitbucket/pull-requests](https://badger.now.sh/bitbucket/tutorials/tutorials.bitbucket.org/pull-requests)<br>![bitbucket/open-pull-requests](https://badger.now.sh/bitbucket/tutorials/tutorials.bitbucket.org/pull-requests?state=open)<br>![bitbucket/declined-pull-requests](https://badger.now.sh/bitbucket/tutorials/tutorials.bitbucket.org/pull-requests?state=declined)<br>![bitbucket/merged-pull-requests](https://badger.now.sh/bitbucket/tutorials/tutorials.bitbucket.org/pull-requests?state=merged)<br>![bitbucket/superseded-pull-requests](https://badger.now.sh/bitbucket/tutorials/tutorials.bitbucket.org/pull-requests?state=superseded) |

### Github Badge Service

[![Github GraphQL API](https://badger.now.sh/static?icon=brands/github&subject=Github%20GraphQL%20API&status=v4)](https://developer.github.com/v4/)

| Path                                           | Description        | Example                                                                             |
| ---------------------------------------------- |------------------- | ----------------------------------------------------------------------------------- |
| /github/`<OWNER>`/`<REPOSITORY>`/forks         | Forks count         | ![github/forks](https://badger.now.sh/github/google/gopacket/forks)                 |
| /github/`<OWNER>`/`<REPOSITORY>`/issues<br>/github/`<OWNER>`/`<REPOSITORY>`/issues?state=open<br>/github/`<OWNER>`/`<REPOSITORY>`/issues?state=closed<br> | Issue count        | ![github/issues](https://badger.now.sh/github/google/gopacket/issues)<br>![github/open-issues](https://badger.now.sh/github/google/gopacket/issues?state=open)<br>![github/closed-issues](https://badger.now.sh/github/google/gopacket/issues?state=closed)  |
| /github/`<OWNER>`/`<REPOSITORY>`/pull-requests<br>/github/`<OWNER>`/`<REPOSITORY>`/pull-requests?state=open<br>/github/`<OWNER>`/`<REPOSITORY>`/pull-requests?state=closed<br>/github/`<OWNER>`/`<REPOSITORY>`/pull-requests?state=merged<br> | Pull Request count | ![github/pull-requests](https://badger.now.sh/github/google/gopacket/pull-requests)<br>![github/open-pull-requests](https://badger.now.sh/github/google/gopacket/pull-requests?state=open)<br>![github/closed-pull-requests](https://badger.now.sh/github/google/gopacket/pull-requests?state=closed)<br>![github/merged-pull-requests](https://badger.now.sh/github/google/gopacket/pull-requests?state=merged) |
| /github/`<OWNER>`/`<REPOSITORY>`/stars         | Star count         | ![github/stars](https://badger.now.sh/github/google/gopacket/stars)                 |

## Getting Started

This project includes a [Makefile](Makefile) for testing and building the project. To see all available options:

```
❯ make help
all                            Runs a clean, build, fmt, lint, test, staticcheck, vet and install
build                          Builds a dynamic executable or package
bump-version                   Bump the version in the version file. Set BUMP to [ patch | major | minor ]
clean                          Cleanup any build binaries or packages
cover                          Runs go test with coverage
cross                          Builds the cross-compiled binaries, creating a clean directory structure (eg. GOOS/GOARCH/binary)
fmt                            Verifies all files have been `gofmt`ed
install                        Installs the executable or package
lint                           Verifies `golint` passes
release                        Builds the cross-compiled binaries, naming them in such a way for release (eg. binary-GOOS-GOARCH)
static                         Builds a static executable
staticcheck                    Verifies `staticcheck` passes
tag                            Create a new git tag to prepare to build a release
test                           Runs the go tests
vet                            Verifies `go vet` passes
```

To run the badger server locally, make sure to run `make all` or `make build` to build the binary & execute it:

```
❯ ./badger
2019/01/11 23:39:10 HTTP service listening on port 8080...
```

## Package Usage

Sample program for using the [badge](https://godoc.org/github.com/tohjustin/badger/pkg/badge) package:

```go
package main

import (
	"fmt"

	"github.com/tohjustin/badger/pkg/badge"
)

func main() {
	badgeOptions := badge.Options{
		Color: "318FE0",
		Icon:  "brands/font-awesome",
		Style: badge.ClassicStyle,
	}
	generatedBadge, _ := badge.Create("Font Awesome", "v5.6.3", &badgeOptions)
}
```
