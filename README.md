
<p align="center">
  <h3 align="center">Scaffold</h3>
  <p align="center">Setup new projects to use with <a href="https://github.com/buildtool/build-tools">build-tools</a> 
</p>

---

<p align="center">
  <a href="https://github.com/buildtool/scaffold/actions"><img alt="GitHub Actions" src="https://github.com/buildtool/scaffold/workflows/Go/badge.svg"></a>
  <a href="https://github.com/buildtool/scaffold/releases"><img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/buildtool/scaffold"></a>
  <a href="pulls"><img alt="GitHub pull requests" src="https://img.shields.io/github/issues-pr/buildtool/scaffold"></a>
  <a href="https://github.com/buildtool/scaffold/releases"><img alt="GitHub All Releases" src="https://img.shields.io/github/downloads/buildtool/scaffold/total"></a>
</p>

<p align="center">
  <a href="https://github.com/buildtool/scaffold/blob/master/LICENSE"><img alt="LICENSE" src="https://img.shields.io/badge/license-MIT-blue.svg?maxAge=43200"></a>
  <a href="https://codecov.io/github/buildtool/scaffold"><img alt="Coverage Status" src="https://codecov.io/gh/buildtool/scaffold/branch/master/graph/badge.svg"></a>
  <a href="https://codebeat.co/projects/github-com-buildtool-scaffold-master"><img alt="codebeat badge" src="https://codebeat.co/badges/bec62abc-78b9-4a9f-8048-1e29117c512b" /></a>  <a href="https://goreportcard.com/report/github.com/buildtool/scaffold"><img alt="goreportcard badge" src="https://goreportcard.com/badge/github.com/buildtool/scaffold" /></a>
  <a href="https://libraries.io/github/buildtool/scaffold"><img alt="" src="https://img.shields.io/librariesio/github/buildtool/scaffold"></a>
</p>

# Setup
You can install the pre-compiled binary (in several different ways), use Docker or compile from source.

## Installation pre-built binaries
**Homebrew tap**

```sh 
$ brew install buildtool/taps/scaffold
```

**Shell script**
```sh
$ curl -sfL https://raw.githubusercontent.com/buildtool/scaffold/master/install.sh | sh
```
**Manually**

Download the pre-compiled binaries from the [releases](https://github.com/buildtool/scaffold/releases) page and copy to the desired location.

      
## Compiling from source
```sh

    # clone it outside GOPATH
    $ git clone https://github.com/buildtool/scaffold
    $ cd build-tools
    
    # get dependencies using go modules (needs go 1.11+)
    $ go get ./...
    
    # build
    $ go build ./cmd
    
    # check it works
    ./scaffold -version
```



This project adheres to the Contributor Covenant [code of conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.
We appreciate your contribution. Please refer to our [contributing guidelines](CONTRIBUTING.md) for further information.


