module gitlab.com/alphaticks/alpha-connect

go 1.13

require (
	cloud.google.com/go/storage v1.0.0
	github.com/AsynkronIT/protoactor-go v0.0.0-20201121081743-27a5e6684be6
	github.com/gogo/protobuf v1.3.1
	github.com/melaurent/gotickfile/v2 v2.0.0-20210105105904-87697673b134
	github.com/pkg/errors v0.8.1
	github.com/satori/go.uuid v1.2.0
	gitlab.com/alphaticks/gorderbook v0.0.0-20201219125519-1a8189af89db
	gitlab.com/alphaticks/xchanger v0.0.0-20210103120503-5dadef341bda
	gitlab.com/tachikoma.ai/tickobjects v0.0.0-20210105110910-04829fb3f735
	gitlab.com/tachikoma.ai/tickpred v0.0.0-20200722192355-60349c3e5530 // indirect
	gitlab.com/tachikoma.ai/tickstore v0.0.0-20201130132500-5507936c8c68
	gitlab.com/tachikoma.ai/tickstore-go-client v0.0.0-20200922092024-4090fea4ef71 // indirect
	gitlab.com/tachikoma.ai/tickstore-grpc v0.0.0-20201231105511-28cbcc209d89
	google.golang.org/api v0.14.0
	google.golang.org/grpc v1.28.0
)

//replace gitlab.com/tachikoma.ai/tickstore => ../../tachikoma.ai/tickstore

//replace gitlab.com/tachikoma.ai/tickobjects => ../../tachikoma.ai/tickobjects

//replace gitlab.com/alphaticks/xchanger => ../../alphaticks/xchanger
