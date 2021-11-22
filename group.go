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

func (g *Group) UpdateDevices(r *Router){
	response, err := r.SendCommand(NewCommand(MT_COMMAND, QUERY_GROUP, CommandParameter{GROUP, strconv.Itoa(g.GroupId)}))
	if err != nil {

	}

	//We expect a comma separated list of addresses ids
	for _, address := range strings.Split(response, ",") {
		if err != nil {
			continue
		}
		d := Device{Address: strings.ReplaceAll(address, "@", "")}

		g.Devices = append(g.Devices, d)
	}
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


