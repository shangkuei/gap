module github.com/shangkuei/gap/yaml

go 1.22

require (
	github.com/mitchellh/mapstructure v1.5.0
	github.com/shangkuei/gap/testhelper v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace github.com/shangkuei/gap/testhelper => ../testhelper
