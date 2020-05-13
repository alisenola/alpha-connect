module gitlab.com/alphaticks/alphac

go 1.13

require (
	github.com/AsynkronIT/protoactor-go v0.0.0-20191119210846-07df21a705b8
	github.com/gogo/protobuf v1.3.1
	github.com/melaurent/gopactor v0.0.0-20190619082820-8be34b07aaaa // indirect
	github.com/quickfixgo/enum v0.0.0-20171007195659-2cbed3730c3e
	github.com/quickfixgo/field v0.0.0-20171007195410-74cea5ec78c7 // indirect
	github.com/quickfixgo/fix50 v0.0.0-20171007213247-d09e70735b64
	github.com/quickfixgo/fixt11 v0.0.0-20171007213433-d9788ca97f5d // indirect
	github.com/quickfixgo/quickfix v0.6.0
	github.com/quickfixgo/tag v0.0.0-20171007194743-cbb465760521 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24
	gitlab.com/alphaticks/gorderbook v0.0.0-20200513101946-928ffaaf764f
	gitlab.com/alphaticks/xchanger v0.0.0-20200513102043-880ac569064b
)

//replace gitlab.com/alphaticks/xchanger => ../../alphaticks/xchanger
