module gitlab.com/alphaticks/alpha-connect

go 1.13

require (
	cloud.google.com/go/storage v1.12.0
	github.com/AsynkronIT/protoactor-go v0.0.0-20201121081743-27a5e6684be6
	github.com/alecthomas/participle v0.4.4 // indirect
	github.com/gogo/protobuf v1.3.2
	github.com/melaurent/gotickfile/v2 v2.0.0-20210125124203-6ee89d1c6a29
	github.com/pkg/errors v0.8.1
	github.com/satori/go.uuid v1.2.0
	gitlab.com/alphaticks/gorderbook v0.0.0-20201219125519-1a8189af89db
	gitlab.com/alphaticks/xchanger v0.0.0-20210103120503-5dadef341bda
	gitlab.com/tachikoma.ai/tickobjects v0.0.0-20210125124247-5184e45a109d
	gitlab.com/tachikoma.ai/tickstore v0.0.0-20210109101144-b8bb705c8ba9
	gitlab.com/tachikoma.ai/tickstore-grpc v0.0.0-20210125120535-c168ee502a87
	google.golang.org/api v0.36.0
	google.golang.org/grpc v1.35.0
)

//replace gitlab.com/tachikoma.ai/tickstore => ../../tachikoma.ai/tickstore

//replace gitlab.com/tachikoma.ai/tickobjects => ../../tachikoma.ai/tickobjects

//replace gitlab.com/alphaticks/xchanger => ../../alphaticks/xchanger
