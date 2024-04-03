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

package util

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// BackendName formats the name with weight
func BackendName(name string, ns string) string {
	name = fmt.Sprintf("%s_%s", ns, name)
	name = strings.Replace(name, ".", "_", -1)
	name = strings.Replace(name, "-", "_", -1)
	name = strings.Replace(name, ":", "_", -1)
	name = strings.Replace(name, "/", "slash", -1)
	name = strings.Replace(name, " ", "", -1)
	return name
}

// CheckAndCreateDir check and create dir
func CheckAndCreateDir(path string) error {
	if subPathExists, err := FileExists(path); err != nil {
		return fmt.Errorf("Could not determine if subPath %s exists; will not attempt to change its permissions", path)
	} else if !subPathExists {
		// Create the sub path now because if it's auto-created later when referenced, it may have an
		// incorrect ownership and mode. For example, the sub path directory must have at least g+rwx
		// when the pod specifies an fsGroup, and if the directory is not created here, Docker will
		// later auto-create it with the incorrect mode 0750
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to mkdir:%s", path)
		}

		if err := os.Chmod(path, 0755); err != nil {
			return err
		}
	}
	return nil
}

// FileExists check file exist
func FileExists(filename string) (bool, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// DirIsEmpty 验证目录是否为空
func DirIsEmpty(dir string) bool {
	infos, err := os.ReadDir(dir)
	if len(infos) == 0 || err != nil {
		return true
	}
	return false
}

// GetDirNameList get all lower level dir
func GetDirNameList(dirpath string, level int) ([]string, error) {
	var dirlist []string
	list, err := os.ReadDir(dirpath)
	if err != nil {
		return nil, err
	}
	for _, f := range list {
		if f.IsDir() {
			if level <= 1 {
				dirlist = append(dirlist, f.Name())
			} else {
				list, err := GetDirList(filepath.Join(dirpath, f.Name()), level-1)
				if err != nil {
					return nil, err
				}
				dirlist = append(dirlist, list...)
			}
		}
	}
	return dirlist, nil
}

func GetDirList(dirpath string, level int) ([]string, error) {
	var dirlist []string
	list, err := os.ReadDir(dirpath)
	if err != nil {
		return nil, err
	}
	for _, f := range list {
		if f.IsDir() {
			if level <= 1 {
				dirlist = append(dirlist, filepath.Join(dirpath, f.Name()))
			} else {
				list, err := GetDirList(filepath.Join(dirpath, f.Name()), level-1)
				if err != nil {
					return nil, err
				}
				dirlist = append(dirlist, list...)
			}
		}
	}
	return dirlist, nil
}

// LocalIP 获取本机 ip
// 获取第一个非 loopback ip
func LocalIP() (net.IP, error) {
	tables, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, t := range tables {
		addrs, err := t.Addrs()
		if err != nil {
			return nil, err
		}
		for _, a := range addrs {
			ipnet, ok := a.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() {
				continue
			}
			if v4 := ipnet.IP.To4(); v4 != nil {
				return v4, nil
			}
		}
	}
	return nil, fmt.Errorf("cannot find local IP address")
}
