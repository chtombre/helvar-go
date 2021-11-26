package helvargo

type Device struct {
	Address string
	Name string
	State int
}

func (d *Device) UpdateName(r *Router) {
	response, err := r.SendCommand(NewCommand(MT_COMMAND, QUERY_DEVICE_DESCRIPTION, CommandParameter{ADDRESS, d.Address}))
	if err == nil {
		d.Name = response
	}
}
