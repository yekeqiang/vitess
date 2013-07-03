// Copyright 2013, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zktopo

import (
	"fmt"
	"sort"

	"code.google.com/p/vitess/go/vt/naming"
)

/*
This file contains the per-cell methods of ZkTopologyServer
*/

func tabletPathForAlias(alias naming.TabletAlias) string {
	return fmt.Sprintf("/zk/%v/vt/tablets/%v", alias.Cell, alias.TabletUidStr())
}

func tabletDirectoryForCell(cell string) string {
	return fmt.Sprintf("/zk/%v/vt/tablets", cell)
}

func (zkts *ZkTopologyServer) CreateOrUpdateTablet(alias naming.TabletAlias, contents string, existingVersion int) (int, error) {
	zkTabletPath := tabletPathForAlias(alias)
	stat, err := zkts.Zconn.Set(zkTabletPath, contents, existingVersion)
	if err != nil {
		return 0, err
	}
	return stat.Version(), nil
}

func (zkts *ZkTopologyServer) GetTablet(alias naming.TabletAlias) (string, int, error) {
	zkTabletPath := tabletPathForAlias(alias)
	data, stat, err := zkts.Zconn.Get(zkTabletPath)
	if err != nil {
		return "", 0, err
	}
	return data, stat.Version(), nil
}

func (zkts *ZkTopologyServer) GetTabletsByCell(cell string) ([]naming.TabletAlias, error) {
	zkTabletsPath := tabletDirectoryForCell(cell)
	children, _, err := zkts.Zconn.Children(zkTabletsPath)
	if err != nil {
		return nil, err
	}

	sort.Strings(children)
	result := make([]naming.TabletAlias, len(children))
	for i, child := range children {
		alias, err := naming.ParseTabletAliasString(child)
		if err != nil {
			return nil, err
		}
		result[i] = alias
	}
	return result, nil
}