// WUTONG, Application Management Platform
// Copyright (C) 2014-2017 Wutong Co., Ltd.

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

package metric

import "k8s.io/apimachinery/pkg/util/sets"

// NewDummyCollector returns a dummy metric collector
func NewDummyCollector() Collector {
	return &DummyCollector{}
}

// DummyCollector dummy implementation for mocks in tests
type DummyCollector struct{}

// Start ...
func (dc DummyCollector) Start() {}

// Stop ...
func (dc DummyCollector) Stop() {}

// SetServerNum -
func (dc DummyCollector) SetServerNum(httpNum, tcpNum int) {}

// SetHosts -
func (dc DummyCollector) SetHosts(hosts sets.String) {}

// RemoveHostMetric -
func (dc DummyCollector) RemoveHostMetric([]string) {}
