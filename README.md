# PSPSDK Example

[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen?style=flat-square)](/LICENSE)

This repository contains example code for developing applications using [PSPSDK Go](https://github.com/pspdev-go/pspsdk-go) and [pspgo](https://github.com/pspdev-go/pspgo). The examples demonstrate various features and functionalities of the PSP, including graphics, sound, input handling, and more.

## Examples

| Directory       | Description                                                |
| --------------- | ---------------------------------------------------------- |
| `hello`         | Prints text using the debug screen                         |
| `controller`    | Displays button and analog stick input                     |
| `audio-tone`    | Generates a square wave and sends it to the audio hardware |
| `clock`         | Reads and displays the local time using the RTC            |
| `power-status`  | Displays battery, power, and clock-frequency information   |
| `file-io`       | Creates, writes, and reads a file                          |
| `system-params` | Reads the PSP nickname and system preferences              |
| `wifi-info`     | Displays the WLAN switch, power state, and MAC address     |
| `triangle`      | Draws a rotating, colored triangle with GU/GUM             |
| `cube`          | Draws a rotating, colored 3D cube with GU/GUM              |
| `embed`         | Embeds a file                                              |

Each example has its own `pspgo.toml`. Build one by running `pspgo build` in
the example directory.
