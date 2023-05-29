# GoWatch

GoWatch is a CLI utility to check a program output and send a signal gain case a pattern is recognized.

## Installation

To install GoWatch, you can use the following command:

```bash
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
wget https://github.com/cheina97/gowatch/releases/download/latest/gowatch-${OS}-${ARCH}
chmod +x gowatch-${OS}-${ARCH}
mv gowatch-${OS}-${ARCH} gowatch
```

Optionally, you can move the binary to a directory in your PATH, for example:

```bash
sudo mv gowatch /usr/local/bin/gowatch
```

## Usage

GoWatch allows you to check the output of a program and send a signal when a pattern is recognized.

First of all create a **json** file containing some **regex** patterns to check:

```json
["pattern1", "pattern2", "pattern3"]
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
