## v0.1.4

Fix bug

- Fix bug in multi comment
  - Parse as meta from the first of comment opener to the last of comment closer.

## v0.1.3

Remove debug code

## v0.1.2

Update CLI

- Add version number to `qiitactl --version`
- Fix description in `qiitactl --help`

## v0.1.1

Change CLI parser

- github.com/alecthomas/kingpin
- github.com/codegangsta/cli
- Don't parse response body as JSON when deleting #32

## v0.1.0

Initial release

- Add support for managing posts in Qiita and Qiita:Team
- Add support for generating file in local
