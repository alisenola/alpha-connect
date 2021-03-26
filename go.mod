module gitlab.com/alphaticks/alpha-connect

go 1.13

require (
	cloud.google.com/go/storage v1.12.0
	github.com/AsynkronIT/protoactor-go v0.0.0-20210225065513-0de6c44ed540
	github.com/alecthomas/participle v0.4.4 // indirect
	github.com/gogo/protobuf v1.3.2
	github.com/melaurent/gotickfile/v2 v2.0.0-20210306083308-17266f6bbb0c
	github.com/melaurent/kafero v1.2.4-0.20210129172623-380493ff2067
	github.com/pkg/errors v0.8.1
	github.com/satori/go.uuid v1.2.0
	gitlab.com/alphaticks/gorderbook v0.0.0-20201219125519-1a8189af89db
	gitlab.com/alphaticks/xchanger v0.0.0-20210325155741-4e88bbeb3724
	gitlab.com/tachikoma.ai/tickobjects v0.0.0-20210211101031-f6dda80a7112
	gitlab.com/tachikoma.ai/tickstore v0.0.0-20210323085801-a91ec3dad48f
	gitlab.com/tachikoma.ai/tickstore-go-client v0.0.0-20210215133608-4091e4618451
	gitlab.com/tachikoma.ai/tickstore-grpc v0.0.0-20210312094618-ca21c4db7965
	google.golang.org/api v0.36.0
	google.golang.org/grpc v1.35.0
)

//replace gitlab.com/tachikoma.ai/tickstore => ../../tachikoma.ai/tickstore

//replace gitlab.com/tachikoma.ai/tickobjects => ../../tachikoma.ai/tickobjects

//replace gitlab.com/alphaticks/xchanger => ../../alphaticks/xchanger
