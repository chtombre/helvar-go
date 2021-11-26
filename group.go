package helvargo

import (
	"strconv"
	"strings"
	"sync"
)

type Group struct {
	GroupId int
	Name string
	Devices []Device
	LastScene string
}

func (g *Group) UpdateName(r *Router) {
	response, err := r.SendCommand(NewCommand(MT_COMMAND, QUERY_GROUP_DESCRIPTION, CommandParameter{GROUP, strconv.Itoa(g.GroupId)}))
	if err != nil {

	}
	g.Name = response
}

func (g *Group) UpdateDevices(r *Router) ([]Device, error){
	response, err := r.SendCommand(NewCommand(MT_COMMAND, QUERY_GROUP, CommandParameter{GROUP, strconv.Itoa(g.GroupId)}))
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	//We expect a comma separated list of addresses ids
	for _, address := range strings.Split(response, ",") {
		if err != nil {
			continue
		}
		address = strings.ReplaceAll(address, "@", "")

		//Only lookup devices with same cluster and router ID as current router
		hostParts := strings.Split(r.Host, ".")
		addressParts := strings.Split(address, ".")
		if addressParts[0] == hostParts[2] && addressParts[1] == hostParts[3] {
			d := Device{Address: address}

			wg.Add(1)

			go func(){
				defer wg.Done()
				d.UpdateName(r)
				g.Devices = append(g.Devices, d)
			}()
		}
	}
	wg.Wait()

	return g.Devices, nil
}

func (r *Router) GetRouters() ([]Router, error) {
	response, err := r.SendCommand(NewCommand(MT_COMMAND, QUERY_ROUTERS, CommandParameter{ADDRESS, "0" }))
	if err != nil {
		return nil, err
	}

	//var wg sync.WaitGroup
	var routers []Router
	for _, routerAddress := range strings.Split(response, ",") {
		if err != nil {
			continue
		}
		g := &Router{Host: routerAddress}

		routers = append(routers, *g)

	}

	return routers, nil
}

func (r *Router) GetClusters() ([]Cluster, error) {
	response, err := r.SendCommand(NewCommand(MT_COMMAND, QUERY_CLUSTERS))
	if err != nil {
		return nil, err
	}

	//var wg sync.WaitGroup
	var clusters []Cluster
	for _, clusterId := range strings.Split(response, ",") {
		id, err := strconv.Atoi(clusterId)
		if err != nil {
			continue
		}
		g := &Cluster{ClusterId: id}

		clusters = append(clusters, *g)

	}

	return clusters, nil
}

func (r *Router) GetGroups() ([]Group, error){
	response, err := r.SendCommand(NewCommand(MT_COMMAND, QUERY_GROUPS))
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	var groups []Group
	//We expect a comma separated list of group ids
	for _, groupId := range strings.Split(response, ",") {
		id, err := strconv.Atoi(groupId)
		if err != nil {
			continue
		}

		g := &Group{GroupId: id}

		wg.Add(1)

		go func(){
			defer wg.Done()

			g.UpdateName(r)
			g.UpdateDevices(r)

			groups = append(groups, *g)
		}()
	}

	wg.Wait()

	return groups, nil
}


