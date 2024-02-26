module github.com/shangkuei/gap/toml

go 1.22

require (
	github.com/mitchellh/mapstructure v1.5.0
	github.com/pelletier/go-toml/v2 v2.1.1
	github.com/shangkuei/gap/testhelper v0.0.0-00010101000000-000000000000
)

require github.com/google/go-cmp v0.6.0 // indirect

replace github.com/shangkuei/gap/testhelper => ../testhelper
