package helvargo

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var CommandResponseTimeout = 60 * time.Second
var KeepAliveDuration = 120 * time.Second

type Router struct {
	RouterId int
	Config *RouterConfig
	Host string
	Port int

	connection net.Conn

	connected bool

	commandReceived *sync.Cond
	commandsToSend chan Command

	commandsReceived map[string]Command

	keepAlive chan bool
}

type RouterConfig struct {
	RouterId int
	//HasValue bool //Possible workaround nullable structs
}

func (r *Router) Initialize(){
	if !r.connected {
		r.Connect()
	}

	r.commandsReceived = make(map[string]Command)
}


func (r *Router) Id() int {
	if r.Config != nil {
		return r.Config.RouterId
	}
	return r.RouterId
}

func (r *Router) Connect() {
	conn, err := net.Dial("tcp",  r.Host + ":" + strconv.Itoa(r.Port))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//setup channel for responses

	r.keepAlive = make(chan bool, 1)
	r.commandsToSend = make(chan Command)

	//r.commandReceived = make(chan Command)
	m := sync.Mutex{}
	r.commandReceived = sync.NewCond(&m)
	r.connection = conn
	r.connected = true


	go r.Keepalive()
	go r.Listener()
	go r.ExecuteCommands()
}

func (r *Router) Disconnect() {
	r.connection.Close()
	r.connected = false

	close(r.keepAlive)
	//Log
}

func (r *Router) Reconnect() {
	r.Disconnect()
	r.Connect()
}

func (r *Router) Listener(){
	c := bufio.NewReader(r.connection)
	for r.connected {
		response, err := c.ReadString('#')
		if err != nil {
			//skip for now
		}

		println("Response received", response)
		command, err := ParseCommand(response)
		if err != nil {
			//skip for now
		}

		//If someone could explain what this is supposed to do? I just followed a tutorial on internet, and it just works :-p
		r.commandReceived.L.Lock()
		r.commandsReceived[command.ToIdentifier()] = *command
		r.commandReceived.Broadcast()
		r.commandReceived.L.Unlock()
	}
}

func (r *Router) SendCommand(command Command) (string, error) {
	r.commandsToSend <- command

	//Check if command is not expecting answers
	for _, b := range COMMAND_TYPES_DONT_LISTEN_FOR_RESPONSE {
		if b == command.CommandType {
			return "", nil
		}
	}

	checkForResponse := func() Command {
		commandId := command.ToIdentifier()
		response, found := r.commandsReceived[commandId]
		if found {
			delete(r.commandsReceived, commandId)
			return response
		}
		return Command{}
	}

	responseCommand := checkForResponse()
	if responseCommand.MessageType != "" {
		return responseCommand.Result, nil
	}

	for responseCommand.MessageType == "" {

		//If someone could explain what this is supposed to do? I just followed a tutorial on internet, and it just works :-p
		r.commandReceived.L.Lock()
		r.commandReceived.Wait()
		responseCommand := checkForResponse()
		r.commandReceived.L.Unlock()
		if responseCommand.MessageType != "" {
			return responseCommand.Result, nil
		}
	}

	return "", errors.New("Probably shouldn't end up here")
}

func (r *Router) ExecuteCommands() {
	for {
		command := <- r.commandsToSend
		strCommand := command.ToString()

		println("Sendt command: %s", strCommand)
		_, err := r.connection.Write([]byte(strCommand))
		if err != nil {
			println("Write to server failed:", err.Error())
			//	os.Exit(1)
		}
		time.Sleep(10 * time.Millisecond)

		//Tell keepalive that we just sent command to router, no need for pinging
		//If it's full, continue
		select {
			case r.keepAlive <- true:
			default:
		}

		//rc := make(chan string)
		//
		//go func() {
		//	c := bufio.NewReader(r.connection)
		//	response, err := c.ReadString('#')
		//	if err != nil {
		//		//skip for now
		//	}
		//
		//	rc <- response
		//	close(rc)
		//}()
		//
		//select {
		//	case response := <- rc:
		//		responseCommand, err := ParseCommand(response)
		//		if err != nil {
		//			//skip for now
		//		}
		//		println("Received response: ", responseCommand.ToString())
		//		//return responseCommand.Result, nil
		//		r.commandReceived <- *responseCommand
		//
		//	case <- time.After(CommandResponseTimeout):
		//		//close(rc) //close the channel to indicate we're no longer interested in the response
		//		println("Response took too long")
		//		//return "", errors.New("Response took to long")
		//		r.commandReceived <-
		//}
	}
}

func (r *Router) Keepalive() {
	for {
		select {
			case <- r.keepAlive:
				continue
			case <- time.After(KeepAliveDuration): {
				if !r.connected { return }

				//send command in a goroutine, so we will go back to listen to keepalive channel to keep from blocking
				r.SendCommand(NewCommand(MT_COMMAND, QUERY_ROUTER_TIME))
				r.SendCommand(NewCommand(MT_COMMAND, QUERY_ROUTER_TIME))
			}

		}
	}
}


