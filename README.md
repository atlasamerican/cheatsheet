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
| Enter              | toggle viewing of command's tldr page |
| Backspace          | close viewed tldr page                |
| j, Down            | select the next command               |
| k, Up              | select the previous command           |
| J, n               | select the next section               |
| K, p               | select the previous section           |
| l, Right, PageDown | select the next page                  |
| h, Left, PageUp    | select the previous page              |

## Building

### Dependencies

- make [*build*]
- go >= 1.16 [*build*]
- git [*runtime*]

Ensure the above build dependencies are satisfied and build with: `make`.

### Installing

After building cheatsheet, `make install` can be used to install it.

## Contributing

cheatsheet is looking for contributors to expand its dataset of commands and
examples.

See [#1][issue-1] for details and the [JSON schema][json-schema] wiki page for a
description of the dataset's JSON schema.

[issue-1]: https://github.com/atlasamerican/cheatsheet/issues/1
[json-schema]: https://github.com/atlasamerican/cheatsheet/wiki/JSON-schema

## License

This project is licensed under the MIT License (see [LICENSE](LICENSE)).
