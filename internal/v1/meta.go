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

package v1

// Meta Common meta
type Meta struct {
	Index      int64  `json:"index"`
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	ServiceID  string `json:"service_id"`
	PluginName string `json:"plugin_name"`
}

// Equals -
func (m *Meta) Equals(c *Meta) bool {
	if m == c {
		return true
	}
	if m == nil || c == nil {
		return false
	}
	if m.Name != c.Name {
		return false
	}
	if m.Namespace != c.Namespace {
		return false
	}
	if m.PluginName != c.PluginName {
		return false
	}
	if m.ServiceID != c.ServiceID {
		return false
	}
	return true
}
