// Copyright (C) 2014-2018 Wutong Co., Ltd.
// WUTONG, Application Management Platform

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Wutong,
// one or multiple Commercial Licenses authorized by Wutong Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"k8s.io/apiserver/pkg/server/healthz"
	"k8s.io/client-go/kubernetes"

	"github.com/wutong-paas/wutong-gateway/cmd/option"
	"github.com/wutong-paas/wutong-gateway/internal/cluster"
	"github.com/wutong-paas/wutong-gateway/internal/controller"
	"github.com/wutong-paas/wutong-gateway/internal/metric"
	"github.com/wutong-paas/wutong-gateway/util"

	k8sutil "github.com/wutong-paas/wutong-gateway/util/k8s"
)

// Run start run
func Run(s *option.GWServer) error {
	logrus.Info("start gateway...")
	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := k8sutil.NewRestConfig(s.K8SConfPath)
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	//create cluster node manage
	logrus.Debug("start creating node manager")
	node, err := cluster.CreateNodeManager(ctx, s.Config)
	if err != nil {
		return fmt.Errorf("create gateway node manage failure %s", err.Error())
	}
	if err := node.Start(); err != nil {
		return fmt.Errorf("start node manager: %v", err)
	}
	defer node.Stop()

	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewGoCollector())
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	mc := metric.NewDummyCollector()
	if s.Config.EnableMetrics {
		mc, err = metric.NewCollector(s.NodeName, reg)
		if err != nil {
			logrus.Fatalf("Error creating prometheus collector:  %v", err)
		}
	}
	mc.Start()

	gwc, err := controller.NewGWController(ctx, clientset, &s.Config, mc, node)
	if err != nil {
		return err
	}
	if gwc == nil {
		return fmt.Errorf("Fail to new GWController")
	}
	logrus.Debug("start gateway controller")
	if err := gwc.Start(errCh); err != nil {
		return fmt.Errorf("Fail to start GWController %s", err.Error())
	}
	defer gwc.Close()

	mux := chi.NewMux()
	registerHealthz(gwc, mux)
	registerMetrics(reg, mux)
	if s.Debug {
		util.ProfilerSetup(mux)
	}
	go startHTTPServer(s.ListenPorts.Health, mux)

	logrus.Info("wutong gateway start success!")

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-term:
		logrus.Warn("Received SIGTERM, exiting gracefully...")
	case err := <-errCh:
		logrus.Errorf("Received a error %s, exiting gracefully...", err.Error())
	}
	logrus.Info("See you next time!")

	return nil
}

func registerHealthz(gc *controller.GWController, mux *chi.Mux) {
	// expose health check endpoint (/healthz)
	healthz.InstallHandler(mux,
		healthz.PingHealthz,
		gc,
	)
}

func registerMetrics(reg *prometheus.Registry, mux *chi.Mux) {
	mux.Handle(
		"/metrics",
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
	)
}

func startHTTPServer(port int, mux *chi.Mux) {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%v", port),
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      300 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	logrus.Fatal(server.ListenAndServe())
}
