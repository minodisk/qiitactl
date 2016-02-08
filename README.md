# qiitactl [![Circle CI](https://img.shields.io/circleci/project/minodisk/qiitactl.svg?style=flat-square)](https://circleci.com/gh/minodisk/qiitactl) [![Coverage Status](https://img.shields.io/coveralls/minodisk/qiitactl.svg?style=flat-square)](https://coveralls.io/github/minodisk/qiitactl?branch=master)

Command line interface to manage the posts in Qitta.

## Description

`qiitactl` makes it possible to CRUD posts of Qiita in your terminal.

`qiitactl` fetches posts from Qiita and write it as markdown file to current working directory. Editing the file in any editor you like, then `qiitactl update post -f path/to/file.md`, the post in Qiita will be updated. You can also generate a new file for post, create a new post and delete a post with this tool.

## Usage

### Preparing for use

1. Create a token at [https://qiita.com/settings/applications](https://qiita.com/settings/applications).
2. Set the created token to `QIITA_ACCESS_TOKEN` environment variable.

### Fetch all posts

```bash
qiitactl fetch posts
```

### Update a post

```bash
qiitactl update post -f path/to/file.md
```

### Create a new post

```bash
qiitactl generate file -t "The title of new post"
vim path/to/file.md
qiitactl create post -f path/to/file.md
```

### Others

```bash
qiitactl help
```

## Install

To install, use `go get`:

```bash
go get -d github.com/minodisk/qiitactl
```

## Contribution

1. Fork ([https://github.com/minodisk/qiitactl/fork](https://github.com/minodisk/qiitactl/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[minodisk](https://github.com/minodisk)
