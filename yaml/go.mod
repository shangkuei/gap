module github.com/shangkuei/gap/yaml

go 1.22

require (
	github.com/goccy/go-yaml v1.11.3
	github.com/mitchellh/mapstructure v1.5.0
	github.com/shangkuei/gap/testhelper v0.0.0-00010101000000-000000000000
)

require (
	github.com/fatih/color v1.10.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)

replace github.com/shangkuei/gap/testhelper => ../testhelper
