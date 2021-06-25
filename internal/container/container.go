package container

import (
	"github.com/StressTestGraphQL/test"
	"github.com/StressTestGraphQL/util"
	"net/http"
)

type components struct {
	Conf util.Config
	Transport *http.Transport
	Queries test.QueriesTest
}

type Dependency struct {
	Components components
}

func Injector() Dependency {

	conf := util.LoadConfig(".")
	caCert, _ := util.ReadCertKey(conf.PathCertification, conf.PathKey)
	transport := util.PreparTlsConfig(conf.PathCertification, conf.PathKey, caCert)
	Queries := test.MountRequests(conf.Limits)

	components := components{
		conf,
		transport,
		Queries,
	}

	return Dependency{Components: components}
}