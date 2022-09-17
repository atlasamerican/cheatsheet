# cheatsheet

[![Data CI](https://github.com/atlasamerican/cheatsheet/actions/workflows/data-ci.yml/badge.svg)](https://github.com/atlasamerican/cheatsheet/actions/workflows/data-ci.yml)

cheatsheet is a command-line reference manual providing a text-based user
interface for accessing [tldr][tldr] pages.

[tldr]: https://tldr.sh/

## Usage

Run `cheatsheet`. The TUI lists common CLI **commands** and examples by
**section**, with sections grouped into **pages**.

Navigate the TUI with these keys:

| keys               | function                              |
| ------------------ | ------------------------------------- |
| ?                  | toggle full key hints                 |
| q                  | quit the app                          |
| c                  | clear on-screen errors                |
| Enter              | view commands or command page         |
| Backspace          | go back to the previous view          |
| j, Down            | select the next item                  |
| k, Up              | select the previous item              |
| l, Right, PageDown | select the next page                  |
| h, Left, PageUp    | select the previous page              |


### Local data

cheatsheet's dataset can be supplemented with local data files in the same
format.

See [Data schema][schema] for a description of this format and the
[`data`](data) directory for working examples.

cheatsheet reads local data files from `$XDG_CONFIG_HOME/cheatsheet` or
`$HOME/.config/cheatsheet` if `$XDG_CONFIG_HOME` is not set.

[schema]: https://github.com/atlasamerican/cheatsheet/wiki/Data-schema

## Packages

- cheatsheet is packaged for the [AUR][aur].

[aur]: https://aur.archlinux.org/packages/cheatsheet-git

## Building

### Dependencies

- make [*build*]
- go >= 1.18 [*build*]
- git [*runtime*]

Ensure the above build dependencies are satisfied and build with: `make`.

### Installing

After building cheatsheet, `make install` can be used to install it.

## Contributing

cheatsheet is looking for contributors to expand its dataset of commands and
examples.

Please familiarize yourself with the general [contribution process][contrib] and
see [#1][issue-1] for additional details.

[issue-1]: https://github.com/atlasamerican/cheatsheet/issues/1
[contrib]:
  https://docs.github.com/en/get-started/quickstart/contributing-to-projects

## License

This project is licensed under the MIT License (see [LICENSE](LICENSE)).
