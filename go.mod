module gitlab.com/alphaticks/alpha-connect

go 1.17

require (
	github.com/AsynkronIT/protoactor-go v0.0.0-20210901041048-df2fc305778c
	github.com/ethereum/go-ethereum v1.10.12
	github.com/gogo/protobuf v1.3.2
	github.com/melaurent/gotickfile/v2 v2.0.0-20220210143804-428af9922408
	github.com/melaurent/kafero v1.2.4-0.20210921082217-5279763aa403
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	gitlab.com/alphaticks/alpha-public-registry-grpc v0.0.0-20220228082856-580817338007
	gitlab.com/alphaticks/go-graphql-client v0.6.1-0.20211231151627-f9c2682bee23
	gitlab.com/alphaticks/gorderbook v0.0.0-20220314112354-40d564234f20
	gitlab.com/alphaticks/tickfunctors v0.0.0-20220315100643-daa5e4a18f2e
	gitlab.com/alphaticks/tickstore v0.0.0-20220221145246-aae0500ce5ae
	gitlab.com/alphaticks/tickstore-go-client v0.0.0-20220217065618-3c1fd871286a
	gitlab.com/alphaticks/tickstore-grpc v0.0.0-20211204075923-5ffd39706d88
	gitlab.com/alphaticks/tickstore-types v0.0.0-20220308115016-e21ae2e9810d
	gitlab.com/alphaticks/xchanger v0.0.0-20220322152424-70a418842535
	go.mongodb.org/mongo-driver v1.7.2
	google.golang.org/grpc v1.40.0
)

require (
	cloud.google.com/go v0.72.0 // indirect
	cloud.google.com/go/storage v1.12.0 // indirect
	github.com/StackExchange/wmi v0.0.0-20180116203802-5d049714c4a6 // indirect
	github.com/Workiva/go-datastructures v1.0.53 // indirect
	github.com/alecthomas/participle v0.4.4 // indirect
	github.com/btcsuite/btcd v0.20.1-beta // indirect
	github.com/deckarep/golang-set v0.0.0-20180603214616-504e848d77ea // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/go-errors/errors v1.0.1 // indirect
	github.com/go-ole/go-ole v1.2.1 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang-collections/collections v0.0.0-20130729185459-604e922904d3 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.1.5 // indirect
	github.com/googleapis/gax-go/v2 v2.0.5 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/jstemmer/go-junit-report v0.9.1 // indirect
	github.com/klauspost/compress v1.10.3 // indirect
	github.com/melaurent/gotickfile v0.0.0-20220126102058-08e9bfcfc230 // indirect
	github.com/orcaman/concurrent-map v0.0.0-20190107190726-7ed82d9cb717 // indirect
	github.com/rjeczalik/notify v0.9.1 // indirect
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible // indirect
	github.com/tklauser/go-sysconf v0.3.5 // indirect
	github.com/tklauser/numcpus v0.2.2 // indirect
	github.com/wangjia184/sortedset v0.0.0-20210325043434-64dd27e173e2 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.0.2 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	gitlab.com/alphaticks/gostarkware v0.0.0-20211208181336-38b492644991 // indirect
	go.opencensus.io v0.22.5 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d // indirect
	golang.org/x/oauth2 v0.0.0-20201109201403-9fd604954f58 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20210816183151-1e6c022a8912 // indirect
	golang.org/x/text v0.3.6 // indirect
	golang.org/x/tools v0.1.2 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/api v0.36.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	nhooyr.io/websocket v1.8.6 // indirect
)

//replace gitlab.com/alphaticks/xchanger => ../xchanger

//replace gitlab.com/alphaticks/gorderbook => ../gorderbook

// replace gitlab.com/alphaticks/go-graphql-client => ../go-graphql-client

//replace gitlab.com/alphaticks/tickstore => ../tickstore
