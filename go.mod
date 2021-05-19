module gitlab.com/alphaticks/alpha-connect

go 1.16

require (
	cloud.google.com/go/storage v1.12.0
	github.com/AsynkronIT/goconsole v0.0.0-20160504192649-bfa12eebf716 // indirect
	github.com/AsynkronIT/gonet v0.0.0-20161127091928-0553637be225 // indirect
	github.com/AsynkronIT/protoactor-go v0.0.0-20210505180410-df90efd4b2b4
	github.com/alecthomas/participle v0.4.4 // indirect
	github.com/gogo/protobuf v1.3.2
	github.com/hashicorp/go.net v0.0.1 // indirect
	github.com/melaurent/gotickfile/v2 v2.0.0-20210406111104-845d7c5ec5dc
	github.com/melaurent/kafero v1.2.4-0.20210129172623-380493ff2067
	github.com/mitchellh/gox v0.4.0 // indirect
	github.com/mitchellh/iochan v1.0.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/satori/go.uuid v1.2.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.0+incompatible // indirect
	gitlab.com/alphaticks/gorderbook v0.0.0-20210414073010-806e8b077356
	gitlab.com/alphaticks/xchanger v0.0.0-20210519144648-4446edfa4ba2
	gitlab.com/tachikoma.ai/tickobjects v0.0.0-20210508085527-ab7f0d4712d8
	gitlab.com/tachikoma.ai/tickstore v0.0.0-20210508085558-6ed77eee2a06
	gitlab.com/tachikoma.ai/tickstore-go-client v0.0.0-20210215133608-4091e4618451
	gitlab.com/tachikoma.ai/tickstore-grpc v0.0.0-20210312094618-ca21c4db7965
	google.golang.org/api v0.36.0
	google.golang.org/grpc v1.37.0
)

//replace gitlab.com/tachikoma.ai/tickstore => ../../tachikoma.ai/tickstore

//replace gitlab.com/tachikoma.ai/tickobjects => ../../tachikoma.ai/tickobjects

//replace gitlab.com/alphaticks/xchanger => ../../alphaticks/xchanger
