// WUTONG, Application Management Platform
// Copyright (C) 2014-2019 Wutong Co., Ltd.

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

package cluster

import (
	"context"
	"net"

	"github.com/sirupsen/logrus"

	"github.com/wutong-paas/wutong-gateway/cmd/option"
	"github.com/wutong-paas/wutong-gateway/util"
)

// IPManager ip manager
// Gets all available IP addresses for synchronizing the current node
type IPManager interface {
	//Whether the IP address belongs to the current node
	IPInCurrentHost(net.IP) bool
	Start() error
	//An IP pool change triggers a forced update of the gateway policy
	NeedUpdateGatewayPolicy() <-chan util.IPEVENT
	Stop()
}

type ipManager struct {
	ctx    context.Context
	cancel context.CancelFunc
	IPPool *util.IPPool
	// lock   sync.Mutex
	config option.Config
	//An IP pool change triggers a forced update of the gateway policy
	needUpdate chan util.IPEVENT
}

// CreateIPManager create ip manage
func CreateIPManager(ctx context.Context, config option.Config) (IPManager, error) {
	newCtx, cancel := context.WithCancel(ctx)
	IPPool := util.NewIPPool(config.IgnoreInterface)
	return &ipManager{
		ctx:        newCtx,
		cancel:     cancel,
		IPPool:     IPPool,
		config:     config,
		needUpdate: make(chan util.IPEVENT, 10),
	}, nil
}

func (i *ipManager) NeedUpdateGatewayPolicy() <-chan util.IPEVENT {
	return i.needUpdate
}

// IPInCurrentHost Whether the IP address belongs to the current node
func (i *ipManager) IPInCurrentHost(in net.IP) bool {
	for _, exit := range i.IPPool.GetHostIPs() {
		if exit.Equal(in) {
			return true
		}
	}
	return false
}

func (i *ipManager) Start() error {
	logrus.Info("start ip manager.")
	go i.IPPool.LoopCheckIPs()
	i.IPPool.Ready()
	logrus.Info("ip manager is ready.")
	return nil
}

func (i *ipManager) Stop() {
	i.cancel()
}
