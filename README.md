# go-tilaa

[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)
[![Build Status](https://travis-ci.org/pascal-splotches/go-tilaa.svg?branch=master)](https://travis-ci.org/pascal-splotches/go-tilaa)

> Tilaa API Client written in Go

go-tilaa is a Go library for managing your [Tilaa](https://www.tilaa.com) Services. It currently allows you to manage your Virtual Machines, Snapshots, Metadata, SSH Keys and more through a simple interface. Currently the library implements v1 of the [Tilaa API](https://www.tilaa.com/en/api/docs).

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contribute](#contribute)
- [License](#license)

## Install

Standard `go get`:
```
$ go get github.com/pascal-splotches/go-tilaa
```

## Usage

For further documentation please see [Godoc](https://godoc.org/github.com/pascal-splotches/go-tilaa).

Basic Example:
```
package main

import (
	"splotch.es/tilaa/go-tilaa"
	"fmt"
)

func main() {
	client := go_tilaa.New("api@example.com", "***")

	presets, _ := client.Preset.List()

	templates, _ := client.Template.List()

	for _, template := range *templates {
		fmt.Printf("\nId:\t\t%v", template.Id)
		fmt.Printf("\nName:\t%v", template.Name)
		fmt.Printf("\nStorage:%v", template.Storage)
		fmt.Printf("\nRam:\t%v", template.Ram)
	}

	sites, _ := client.Site.List()

	for _, site := range *sites {
		fmt.Printf("\nId:\t\t%v", site.Id)
		fmt.Printf("\nName:\t%v", site.Name)
	}
}
```

## Maintainers

[@Pascal Scheepers](https://github.com/pascal-splotches)

## Contribute

PRs accepted.

Small note: If editing the README, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

This project is licensed under the [GNU General Public License v3.0](LICENSE)

 © Pascal Scheepers
