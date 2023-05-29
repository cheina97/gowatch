# GoWatch

GoWatch is a CLI utility to check a program output (using **regex**) and send a signal if a pattern is recognized.

## Installation

To install GoWatch, you have two options:
### Precompiled binaries

Download it from the [releases page](https://github.com/cheina97/gowatch/releases/latest)

### Build from source

```bash
git clone https://github.com/cheina97/gowatch.git
make build
chmod +x gowatch
```

Optionally, you can move the binary to a directory in your PATH, for example:

```bash
sudo mv gowatch /usr/local/bin/gowatch
```

## Usage

**GoWatch** allows you to check the output of a program using **regex** and send a **signal** when a pattern is recognized.

First of all create a **json** file containing some **regex** patterns to check:

```json
["regex1", "regex2", "regex3"]
```

Then run **gowatch** and pass the program you want to check as argument:

```bash
gowatch -f patterns.json command arg1 arg2
```

## Options

These flags can be used to customize the behavior of **gowatch**:

```bash
Usage of ./gowatch:
  -c int
        Number of concurrent workers (default 4)
  -d    Enable debug mode
  -f string
        Path to the json file which contains the patterns to check
  -q    Disable output
  -s value
        Signal to send to the command's process when its output matches a pattern (default killed)
  -v    Show version
```
