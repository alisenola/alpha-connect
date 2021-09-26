module gitlab.com/alphaticks/alpha-connect

go 1.16

require (
	cloud.google.com/go/storage v1.12.0
	github.com/AsynkronIT/protoactor-go v0.0.0-20210901041048-df2fc305778c
	github.com/alecthomas/participle v0.4.4 // indirect
	github.com/gogo/protobuf v1.3.2
	github.com/melaurent/gotickfile/v2 v2.0.0-20210406111104-845d7c5ec5dc
	github.com/melaurent/kafero v1.2.4-0.20210921082217-5279763aa403
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	gitlab.com/alphaticks/gorderbook v0.0.0-20210616120115-89b18aade871
	gitlab.com/alphaticks/tickobjects v0.0.0-20210921082758-753463aeda2a
	gitlab.com/alphaticks/tickstore v0.0.0-20210922125948-0d794fe0d1e0
	gitlab.com/alphaticks/tickstore-go-client v0.0.0-20210921083244-4fd9ecfc241a
	gitlab.com/alphaticks/tickstore-grpc v0.0.0-20210921083148-26dad7f5cbb0
	gitlab.com/alphaticks/xchanger v0.0.0-20210924141109-dd64bec838a3
	go.mongodb.org/mongo-driver v1.7.2
	google.golang.org/api v0.36.0
	google.golang.org/grpc v1.40.0
)

//replace gitlab.com/alphaticks/xchanger => ../xchanger
