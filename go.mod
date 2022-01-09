module gitlab.com/alphaticks/alpha-connect

go 1.17

require (
	cloud.google.com/go/storage v1.12.0
	github.com/AsynkronIT/protoactor-go v0.0.0-20210901041048-df2fc305778c
	github.com/gogo/protobuf v1.3.2
	github.com/melaurent/gotickfile/v2 v2.0.0-20210406111104-845d7c5ec5dc
	github.com/melaurent/kafero v1.2.4-0.20210921082217-5279763aa403
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	gitlab.com/alphaticks/gorderbook v0.0.0-20220109093018-7fc7e896f424
	gitlab.com/alphaticks/tickobjects v0.0.0-20211222104755-62bda2315570
	gitlab.com/alphaticks/tickstore v0.0.0-20211207113051-444d7194ff7a
	gitlab.com/alphaticks/tickstore-go-client v0.0.0-20210921083244-4fd9ecfc241a
	gitlab.com/alphaticks/tickstore-grpc v0.0.0-20211204075923-5ffd39706d88
	gitlab.com/alphaticks/xchanger v0.0.0-20211228094502-5203e90e1766
	go.mongodb.org/mongo-driver v1.7.2
	google.golang.org/api v0.36.0
	google.golang.org/grpc v1.40.0
)

require (
	cloud.google.com/go v0.72.0 // indirect
	github.com/Workiva/go-datastructures v1.0.53 // indirect
	github.com/alecthomas/participle v0.4.4 // indirect
	github.com/btcsuite/btcd v0.20.1-beta // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/ethereum/go-ethereum v1.10.12 // indirect
	github.com/go-errors/errors v1.0.1 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang-collections/collections v0.0.0-20130729185459-604e922904d3 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/googleapis/gax-go/v2 v2.0.5 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/jstemmer/go-junit-report v0.9.1 // indirect
	github.com/klauspost/compress v1.9.5 // indirect
	github.com/melaurent/gotickfile v0.0.0-20210921084211-01142566a54d // indirect
	github.com/orcaman/concurrent-map v0.0.0-20190107190726-7ed82d9cb717 // indirect
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
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c // indirect
	google.golang.org/protobuf v1.26.0 // indirect
)

//replace gitlab.com/alphaticks/xchanger => ../xchanger

//replace gitlab.com/alphaticks/tickstore => ../tickstore
