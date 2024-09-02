# mrpack-cli
![Current Build](https://github.com/oceanoc/mrpack-cli/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/oceanoc/mrpack-cli)](https://goreportcard.com/report/github.com/oceanoc/mrpack-cli)
###### This is a fanmade tool, It is not made or endorsed by Modrinth
This application is now in [**maintainence mode**](https://en.wikipedia.org/wiki/Maintenance_mode)

A simple command-line tool to extract .mrpack files
## Usage
```
Usage of ./mrpack-cli [mrpack] [args]:
  -download
    	Set to false to skip downloads (default true)
  -entry
    	Set to false to skip making entry in the Minecraft launcher (default true)
  -output string
    	Set where the modpack will be extracted (default "default")
```

## Getting Started
### Installing
- Get `mrpack-cli` from [releases](https://github.com/OceanOC/mrpack-cli/releases)
- Add to your PATH
### Building
#### Dependencies
[Golang](https://go.dev)
#### Steps
1. `git clone https://github.com/OceanOC/mrpack-cli.git`
2. `go build`
3. `./mrpack-cli`
## License
This project is licensed under the Apache License 2.0
## Acknowledgments
- [Modrinth](https://modrinth.com)

This project makes use of the following third-party libraries:
- **color** by Fatih Arslan, licensed under the MIT License. [Link](https://github.com/fatih/color)

For more details on the license, please refer to the `LICENSE` and `LICENSE-MIT` files in this repository.
