package helvargo

import (
	"fmt"
	"regexp"
	"strings"
)

var default_helvar_termination_char = "#"

type Command struct {
	MessageType MessageType
	CommandType CommandType
	Params []CommandParameter
	Address string
	Result string
}

type CommandParameter struct {
	Type CommandParameterType
	Value string
}

func (cp *CommandParameter) ToString() string{
	if cp.Type == ADDRESS { return fmt.Sprintf("%s%s", cp.Type, cp.Value) }
	return fmt.Sprintf("%s:%s", cp.Type, cp.Value)
}

type CommandType string

const (
	QUERY_CLUSTERS CommandType = "101"
	QUERY_GROUP_DESCRIPTION CommandType = "105"
	QUERY_DEVICE_DESCRIPTION CommandType = "106"
	QUERY_DEVICE_TYPES_AND_ADDRESSES CommandType = "100"
	QUERY_DEVICE_STATE CommandType = "110"
	QUERY_WORKGROUP_NAME CommandType = "107"
	QUERY_DEVICE_LOAD_LEVEL CommandType = "152"
	QUERY_SCENE_INFO CommandType = "167"
	QUERY_ROUTER_TIME CommandType = "185"
	QUERY_LAST_SCENE_IN_GROUP CommandType = "109"
	QUERY_LAST_SCENE_IN_BLOCK CommandType = "103"
	QUERY_GROUP CommandType = "164"
	QUERY_GROUPS CommandType = "165"
	QUERY_SCENE_NAMES CommandType = "166"
	QUERY_ROUTER_VERSION CommandType = "190"
	QUERY_HELVARNET_VERSION CommandType = "191"

	//Commands
	DIRECT_LEVEL_DEVICE CommandType = "14"
	DIRECT_LEVEL_GROUP CommandType = "13"
	RECALL_SCENE CommandType = "11"
)

var COMMAND_TYPES_DONT_LISTEN_FOR_RESPONSE = []CommandType { RECALL_SCENE, DIRECT_LEVEL_DEVICE, DIRECT_LEVEL_GROUP}

type CommandParameterType string

const(
	VERSION CommandParameterType = "V"
	COMMAND CommandParameterType = "C"
	ADDRESS CommandParameterType = "@"
	GROUP CommandParameterType = "G"
	SCENE CommandParameterType = "S"
	BLOCK CommandParameterType = "B"
	FADE_TIME CommandParameterType = "F"
	LEVEL CommandParameterType = "L"
	PROPORTION CommandParameterType = "P"
	DISPLAY_SCREEN CommandParameterType = "D"
	SEQUENCE_NUMBER CommandParameterType = "Q"
	TIME CommandParameterType = "T"
	ACK CommandParameterType = "A"
	LATITUDE CommandParameterType = "L"
	LONGITUDE CommandParameterType = "E"
	TIME_ZONE_DIFFERENCE CommandParameterType = "Z"
	DAYLIGHT_SAVING_TIME CommandParameterType = "Y"
	CONSTANT_LIGHT_SCENE CommandParameterType = "K"
	FORCE_STORE_SCENE CommandParameterType = "O"
)

type MessageType string

const(
	MT_COMMAND MessageType = ">"
	MT_INTERNAL_COMMAND MessageType = "<"
	MT_REPLY MessageType = "?"
	MT_ERROR MessageType = "!"
)

func NewCommand(messageType MessageType, commandType CommandType, params ...CommandParameter) Command{
	return Command{
		MessageType: messageType,
		CommandType: commandType,
		Params: params,
	}
}

func (c *Command) buildBaseParameters() []CommandParameter {
	return []CommandParameter {
		{Type: VERSION, Value: "1"},
		{Type: COMMAND, Value: string(c.CommandType)},
	}
}

func (c *Command) ToIdentifier() string {
	var parameters []CommandParameter
	parameters = append(parameters, c.Params...)

	var stringParams []string
	for _, param := range parameters {
		stringParams = append(stringParams, param.ToString())
	}
	result := strings.Join(stringParams, ",")

	return fmt.Sprintf("%s:%s", c.CommandType, result)
}

func (c *Command) ToString() string {
	parameters := c.buildBaseParameters()

	parameters = append(parameters, c.Params...)

	if c.Address != "" {
		parameters = append(parameters, CommandParameter{ADDRESS, c.Address})
	}

	var stringParams []string
	for _, param := range parameters {
		stringParams = append(stringParams, param.ToString())
	}
	mainMessage := strings.Join(stringParams, ",")

	if c.Result != "" {
		mainMessage = fmt.Sprintf("%s=%s", mainMessage, c.Result)
	}

	return fmt.Sprintf("%s%s%s", c.MessageType, mainMessage, default_helvar_termination_char)
}


func ParseCommand(input string) (*Command, error) {
	//regexp.Compile("^(?P<type>[<>?!])V\:(?P<version>\d),C\:(?P<command>\d+),?(?P<params>[^=@#]+)?(?P<address>@[^=#]+)?(=(?P<result>[^#]*))?#?$")
	r, err := regexp.Compile("^(?P<type>[<>?!])V\\:(?P<version>\\d),C\\:(?P<command>\\d+),?(?P<params>[^=@#]+)?(?P<address>@[^=#]+)?(=(?P<result>[^#]*))?#?$")
	if err != nil {
		fmt.Println("Failed to parse command: ", input)
		return nil, err
	}

	res := r.FindStringSubmatch (input)

	messageType := res[1]
	commandType := res[3]
	params := parseParams(res[4])
	address := strings.ReplaceAll(res[5], "@", "")
	result := res[7]
	_ = result

	return &Command{
		CommandType: CommandType(commandType),
		MessageType: MessageType(messageType),
		Params: params,
		Address: address,
		Result: result,
	}, nil
}

func parseParams(input string) []CommandParameter{
	var parameters []CommandParameter

	params := strings.Split(input, ",")

	for _, p := range params {
		parts := strings.Split(p, ":")
		if len(parts) == 2 {
			parameters = append(parameters, CommandParameter{
				CommandParameterType(parts[0]),
				parts[1],
			})
		}
	}

	return parameters
}