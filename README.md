# repos-stat
Simple tool to check the status of multiple repo at once.

This utility run `git status` on all directory in the filesystem that contains a `.git` folder. It then display the status of both the repo **index** and **working tree** (just a wrapper for `git status -s`). 

## Installation

- Install the package using `go get -u github.com/fedemengo/repos-stat`
- Add go binary path to `$PATH` by setting `$PATH=${PATH}:${GOPATH}/bin`

## Usage

The basic usage is

```
repos-stat [options] path-to-visit [path-to-exclude]

options:
	--no-clean:		skip clean repository
	--no-broken:		skip broken repository
```

To specify which directory to exclude append a `-` to the path

## Example

In the following example all directory in `/home/user/` are inspected. In addition clean repositories are ignored.

```
repos-stat --no-clean /home/user 
```

![example](https://github.com/fedemengo/repos-stat/blob/master/res/example.png)

