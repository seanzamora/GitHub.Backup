module github.backup

go 1.21

replace cli => ./cli

replace utils => ./utils

require cli v0.0.0

require (
	github.com/spf13/cobra v1.7.0
	utils v0.0.0
)

require (
	github.com/fatih/color v1.15.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
