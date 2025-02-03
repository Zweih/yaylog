# yaylog

`yaylog` is a simple CLI util, written in **Go** / **Golang**, for arch and arch-based linux distros to list recently installed packages.

despite the name, it's not limited to `yay` and works with any package manager that logs package installations to `/var/log/pacman.log`. so it can be used with `pacman`, `yay`, `paru`, `aura`, and even `yaourt` if you're somehow still using it.

it supports optional filters for explicitly installed packages or dependencies.

[![Packaging status](https://repology.org/badge/vertical-allrepos/yaylog.svg)](https://repology.org/project/yaylog/versions)

## features

- view recently installed packages with timestamps
- filter results by explicitly installed packages
- filter results by packages installed as dependencies
- sort results by installation date or alphabetically
- filter results by a specific installation date
- sort results by size on disk

## why is it called yaylog if it works with other package managers?
because yay is my preferred aur helper and the name has a good flow.

## is it good?
[yes.](https://news.ycombinator.com/item?id=3067434)

## roadmap

- ~~rewrite in Golang~~ COMPLETE
- ~~additional filters~~ COMPLETE
- list possibly or confirmed stale/abandoned packages
- ~~sort by size on disk~~ COMPLETE
- dependency graph
- ~~concurrent filtering~~ COMPLETE
- ~~filter by size on disk~~ COMPLETE


## installation

### from AUR (**recommended**)
install using an AUR helper like `yay`:
```bash
yay -S yaylog
```

### building from source + manual installation
1. clone the repo:
   ```bash
   git clone https://github.com/zweih/yaylog.git
   cd yaylog
   ```
2. build the binary:
   ```bash
   go build -o yaylog ./cmd/yaylog
   ```
3. copy the binary to your system's `$PATH`:
   ```bash
   sudo install -m755 yaylog /usr/bin/yaylog
   ```
4. copy the manpage:
   ```bash
   sudo install -m644 yaylog.1 /usr/share/man/man1/yaylog.1
   ```

## usage

```bash
yaylog [options]
```

### options
- `-n <number>`: number of recent packages to display (default: 20)
- `-a`: show all installed packages (ignores `-n`)
- `-e`: show only explicitly installed packages
- `-d`: show only packages installed as dependencies
- `--date <YYYY-MM-DD>`: show packages installed on the specified date
- `--size <filter>`: filter packages by size on disk
   - size filter examples:
      - `">10MB"`: show packages larger than 10MB
      - `"<500KB"`: show packages smaller than 500KB
  - quotes are required for the filter
- `--sort <mode>`: sort results by:
  - `date` (default): sort by installation date
  - `alphabetical`: sort alphabetically by package name
  - `size:asc` / `size:desc`: sort by package size on disk; ascending or descending, respectively
- `-h`: print help info

### examples
1. show the last 10 installed packages:
   ```bash
   yaylog -n 10
   ```
2. show all explicitly installed packages:
   ```bash
   yaylog -ae
   ```
3. show only dependencies installed on a specific date:
   ```bash
   yaylog -d --date 2024-12-25
   ```
4. show all packages sorted alphabetically:
   ```bash
   yaylog -a --sort alphabetical
   ```
5. show the 15 most recent explicitly installed packages:
   ```bash
   yaylog -en 15
   ```
6. show the 20 most recently installed packages larger than 20MB:
   ```bash
   yaylog --size ">20MB"
   ```
7. show the 10 largest explicitly installed packages:
   ```bash
   yaylog -en 10 --sort size:desc
   ```
8. show all dependencies smaller than 500KB:
   ```bash
   yaylog -ad --size "<500KB"
   ```

   **note**: the `-e` flag must be used before the `-n` flag as the n flag consumes the next argument.

## license

this project is licensed under the MIT license. see [license](LICENSE) for details.
