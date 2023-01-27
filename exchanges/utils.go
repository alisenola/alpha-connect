package exchanges

import (
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	xmodels "gitlab.com/alphaticks/xchanger/models"
	"net"
)

func GetExchangesExecutor(as *actor.ActorSystem) *actor.PID {
	return actor.NewPID(as.Address(), "executor/exchanges")
}

func GetExchangeExecutor(as *actor.ActorSystem, exchange *xmodels.Exchange) *actor.PID {
	return actor.NewPID(as.Address(), "executor/exchanges/"+exchange.Name+"_executor")
}

func GetAccountListener(as *actor.ActorSystem, account string) *actor.PID {
	return actor.NewPID(as.Address(), "executor/exchanges/"+account+"_account/listener")
}

func localAddresses(name string) ([]net.Addr, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, fmt.Errorf("error getting interface: %v", err)
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, fmt.Errorf("error getting addresses: %v", err)
	}
	var res []net.Addr
	for _, a := range addrs {
		ipAddr := a.(*net.IPNet).IP.To4()
		if ipAddr != nil {
			tcpAddr := &net.TCPAddr{
				IP: ipAddr,
			}
			res = append(res, tcpAddr)
		}
	}

	return res, nil
}
