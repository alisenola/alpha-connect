module gitlab.com/alphaticks/alphac

go 1.13

require (
	cloud.google.com/go/storage v1.0.0
	github.com/AsynkronIT/protoactor-go v0.0.0-20200317173033-c483abfa40e2
	github.com/gogo/protobuf v1.3.1
	github.com/melaurent/gotickfile/v2 v2.0.0-20200916062636-4a9b6b157075
	github.com/melaurent/kafero v1.2.4-0.20200403144811-ea8b6e20da99
	github.com/pkg/errors v0.8.1
	github.com/satori/go.uuid v1.2.0
	gitlab.com/alphaticks/gorderbook v0.0.0-20200909075534-6bd84aeffd29
	gitlab.com/alphaticks/xchanger v0.0.0-20200922144339-69bf17c971a7
	gitlab.com/tachikoma.ai/tickobjects v0.0.0-20200923115315-24f9fcf51661
	gitlab.com/tachikoma.ai/tickpred v0.0.0-20200722192355-60349c3e5530 // indirect
	gitlab.com/tachikoma.ai/tickstore v0.0.0-20200916131501-b301eb95e765
	gitlab.com/tachikoma.ai/tickstore-go-client v0.0.0-20200916105757-9464082fef2b
	gitlab.com/tachikoma.ai/tickstore-grpc v0.0.0-20200406135554-c8e8045438c0
	google.golang.org/grpc v1.28.0
)

//replace gitlab.com/tachikoma.ai/tickobjects => ../../tachikoma.ai/tickobjects

//replace gitlab.com/alphaticks/xchanger => ../../alphaticks/xchanger
