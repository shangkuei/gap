module github.com/shangkuei/gap/json

go 1.22

require (
	github.com/mitchellh/mapstructure v1.5.0
	github.com/shangkuei/gap/testhelper v0.0.1
)

require github.com/google/go-cmp v0.6.0 // indirect

replace github.com/shangkuei/gap/testhelper => ../testhelper
