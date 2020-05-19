module gitlab.com/alphaticks/alphac

go 1.13

require (
	github.com/AsynkronIT/protoactor-go v0.0.0-20191119210846-07df21a705b8
	github.com/gogo/protobuf v1.3.1
	github.com/melaurent/gopactor v0.0.0-20190619082820-8be34b07aaaa // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24
	gitlab.com/alphaticks/gorderbook v0.0.0-20200514150550-2800aed69dfe
	gitlab.com/alphaticks/xchanger v0.0.0-20200519091333-db81f9f32919
)

//replace gitlab.com/alphaticks/xchanger => ../../alphaticks/xchanger
