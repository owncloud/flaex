# Flaex
[![Build Status](https://cloud.drone.io/api/badges/owncloud/flaex/status.svg)](https://cloud.drone.io/owncloud/flaex)

Extract flags from the ocis project-components to create documentation.

## Build
Run the build target from the Makefile: `make build`.

## Usage
Running `./bin/flaex -help` shows information on how to use flaex. Flaex writes to standard output. Feel free to redirect
the content into a file.

## Example
Please consider the following command as an example:
```./bin/flaex -template=templates/CONFIGURATION.tmpl -command-path=../ocis-konnectd/pkg/command/ -flagset-path=../ocis-konnectd/pkg/flagset/ > ../ocis-konnectd/docs/configuration.md```
