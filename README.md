# yaylog

`yaylog` is a simple CLI util for arch and arch-based linux distros to list recently installed packages.

despite the name, it's not limited to `yay` and works with any package manager that logs package installations to `/var/log/pacman.log`. so it can be used with `pacman`, `yay`, `paru`, `aura`, and even `yaourt` if you're somehow still using it.

it supports optional filters for explicitly installed packages or dependencies.

## features

- view recently installed packages with timestamps.
- filter results by explicitly installed packages.
- show all installed packages with alignment for readability.

## why is it called yaylog if it works with other package managers?
because yay is my preferred aur helper and the name has a good flow.

## installation

### from AUR
install using an AUR helper like `yay`:
```bash
yay -S yaylog
```

### manual installation
clone the repo and copy the script to your bin:
```bash
git clone https://github.com/zweih/yaylog.git
cd yaylog
sudo install -m755 yaylog.sh /usr/bin/yaylog
```

## usage

```bash
yaylog [-n <number>] [-e] [-a]
```

### options
- `-n <number>`: number of recent packages to display (default: 20).
- `-e`: show only explicitly installed packages.
- `-a`: show all installed packages.
- `-h`: print help information.

### examples
1. show the last 10 installed packages:
   ```bash
   yaylog -n 10
   ```
2. show all explicitly installed packages:
   ```bash
   yaylog -ae
   ```
3. show the 15 most recent explicitly installed packages:
   ```bash
   yaylog -en 15
   ```

   **note**: the `-e` flag must be used before the `-n` flag as the n flag consumes the next argument.

## license

this project is licensed under the MIT license. see [license](license) for details.
