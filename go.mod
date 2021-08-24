module gitlab.com/alphaticks/alpha-connect

go 1.16

require (
	cloud.google.com/go/storage v1.12.0
	github.com/AsynkronIT/protoactor-go v0.0.0-20210810091324-c3c6e02d5d46
	github.com/alecthomas/participle v0.4.4 // indirect
	github.com/gogo/protobuf v1.3.2
	github.com/melaurent/gotickfile/v2 v2.0.0-20210406111104-845d7c5ec5dc
	github.com/melaurent/kafero v1.2.4-0.20210129172623-380493ff2067
	github.com/pkg/errors v0.8.1
	github.com/satori/go.uuid v1.2.0
	gitlab.com/alphaticks/gorderbook v0.0.0-20210616120115-89b18aade871
	gitlab.com/alphaticks/xchanger v0.0.0-20210824113509-0cd06c0eeab2
	gitlab.com/tachikoma.ai/tickobjects v0.0.0-20210528122836-d02ce1923e51
	gitlab.com/tachikoma.ai/tickstore v0.0.0-20210508085558-6ed77eee2a06
	gitlab.com/tachikoma.ai/tickstore-go-client v0.0.0-20210215133608-4091e4618451
	gitlab.com/tachikoma.ai/tickstore-grpc v0.0.0-20210312094618-ca21c4db7965
	google.golang.org/api v0.36.0
	google.golang.org/grpc v1.39.1
)

//replace gitlab.com/tachikoma.ai/tickstore => ../../tachikoma.ai/tickstore

//replace gitlab.com/tachikoma.ai/tickobjects => ../../tachikoma.ai/tickobjects

//replace gitlab.com/alphaticks/xchanger => ../xchanger
//replace gitlab.com/alphaticks/gorderbook => ../gorderbook
