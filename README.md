# cheatsheet

cheatsheet is a command-line reference manual providing a text-based user
interface for accessing [tldr][tldr] pages.

[tldr]: https://tldr.sh/

## Usage

Run `cheatsheet`. The TUI lists common CLI **commands** and examples by
**section**, with sections grouped into **pages**.

Navigate the TUI with these keys:

| keys               | function                              |
| ------------------ | ------------------------------------- |
| ?                  | toggle on-screen key hints            |
| q                  | quit the app                          |
| c                  | clear errors or hints                 |
| Enter              | view commands or command page         |
| Backspace          | unview commands or command page       |
| j, Down            | select the next item                  |
| k, Up              | select the previous item              |
| l, Right, PageDown | select the next page                  |
| h, Left, PageUp    | select the previous page              |

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

Please familiarize yourself with the general [contribution process][contrib].

See [#1][issue-1] for additional details and the [JSON schema][json-schema] wiki
page for a description of the dataset's JSON schema.

[issue-1]: https://github.com/atlasamerican/cheatsheet/issues/1
[json-schema]: https://github.com/atlasamerican/cheatsheet/wiki/JSON-schema
[contrib]:
  https://docs.github.com/en/get-started/quickstart/contributing-to-projects

## License

This project is licensed under the MIT License (see [LICENSE](LICENSE)).
