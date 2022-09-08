// Copyright 2022 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package details

import (
	"context"
	"net"

	"github.com/cloudwego/biz-demo/bookinfo/kitex_gen/cwg/bookinfo/details"
	"github.com/cloudwego/biz-demo/bookinfo/kitex_gen/cwg/bookinfo/details/detailsservice"
	"github.com/cloudwego/biz-demo/bookinfo/pkg/constants"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"github.com/kitex-contrib/xds"
)

type Server struct {
	opts *ServerOptions
	svc  details.DetailsService
}

type ServerOptions struct {
	Addr      string `mapstructure:"addr"`
	EnableXDS bool   `mapstructure:"enableXDS"`
}

func DefaultServerOptions() *ServerOptions {
	return &ServerOptions{
		Addr:      ":8084",
		EnableXDS: false,
	}
}

func (s *Server) Run(ctx context.Context) error {
	if s.opts.EnableXDS {
		if err := xds.Init(); err != nil {
			klog.Fatal(err)
		}
	}

	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(constants.DetailsServiceName),
		provider.WithInsecure(),
	)
	defer p.Shutdown(ctx)

	addr, err := net.ResolveTCPAddr("tcp", s.opts.Addr)
	if err != nil {
		klog.Fatal(err)
	}
	svr := detailsservice.NewServer(
		s.svc,
		server.WithServiceAddr(addr),
		server.WithSuite(tracing.NewServerSuite()),
	)
	if err := svr.Run(); err != nil {
		klog.Fatal(err)
	}

	return nil
}
