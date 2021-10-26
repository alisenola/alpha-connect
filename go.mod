module gitlab.com/alphaticks/alpha-connect

go 1.16

require (
	cloud.google.com/go/storage v1.12.0
	github.com/AsynkronIT/protoactor-go v0.0.0-20210901041048-df2fc305778c
	github.com/alecthomas/participle v0.4.4 // indirect
	github.com/gogo/protobuf v1.3.2
	github.com/golang/snappy v0.0.3 // indirect
	github.com/melaurent/gotickfile/v2 v2.0.0-20210406111104-845d7c5ec5dc
	github.com/melaurent/kafero v1.2.4-0.20210921082217-5279763aa403
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	gitlab.com/alphaticks/gorderbook v0.0.0-20211009212430-45ffdac78c6e
	gitlab.com/alphaticks/tickobjects v0.0.0-20211009212520-18d446a7dd7f
	gitlab.com/alphaticks/tickstore v0.0.0-20211026073327-70c8d9e07de7
	gitlab.com/alphaticks/tickstore-go-client v0.0.0-20210921083244-4fd9ecfc241a
	gitlab.com/alphaticks/tickstore-grpc v0.0.0-20210921083148-26dad7f5cbb0
	gitlab.com/alphaticks/xchanger v0.0.0-20211019143513-c1a11408e148
	go.mongodb.org/mongo-driver v1.7.2
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d // indirect
	golang.org/x/sys v0.0.0-20210816183151-1e6c022a8912 // indirect
	google.golang.org/api v0.36.0
	google.golang.org/grpc v1.40.0
)

//replace gitlab.com/alphaticks/xchanger => ../xchanger
//replace gitlab.com/alphaticks/tickstore => ../tickstore
