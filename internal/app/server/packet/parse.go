package packet

import (
	"errors"
	t "ftp/internal/app/server/types"
	"log/slog"
	"strings"
)

type ParseResult int

const (
	CommandeNotFound ParseResult = iota
	MissingParameter
	Ok
)

type cmdParser struct {
	cmd        string
	packetType t.PacketType
	params     bool
}

func (p cmdParser) check(
	beforeSpace, afterSpace string,
	clientId string,
) (packet t.PacketIn, result ParseResult) {

	if beforeSpace == p.cmd {
		if !p.params && afterSpace == "" {
			return packet, MissingParameter
		}
		return t.PacketIn{Type: p.packetType, Path: afterSpace, ClientId: clientId }, Ok
	}

	return
}

type cmdParserBuilder struct {
	inner cmdParser
}

func cmdFor(cmd string, packetType t.PacketType) cmdParserBuilder {
	return cmdParserBuilder{
		inner: cmdParser{
			cmd:        cmd,
			packetType: packetType,
			params:     true,
		},
	}
}

func (p cmdParserBuilder) withNoParams() cmdParserBuilder {
	p.inner.params = false
	return p
}

func (p cmdParserBuilder) build() cmdParser {
	return p.inner
}

var clientCmdParsers = []cmdParser{
	cmdFor("End", t.End).withNoParams().build(),
	cmdFor("List", t.List).build(),
	cmdFor("Get", t.Get).withNoParams().build(),
	cmdFor("Cd", t.Cd).build(),
}

var adminCmdParsers = []cmdParser{
	cmdFor("End", t.End).withNoParams().build(),
	cmdFor("Terminate", t.Terminate).withNoParams().build(),
	cmdFor("List", t.List).build(),
	cmdFor("Cd", t.Cd).withNoParams().build(),
	cmdFor("Hide", t.Hide).build(),
	cmdFor("Reveal", t.Reveal).build(),
}

func Parse(
	line string,
	admin bool,
	clientId string,
) (t.PacketIn, error) {

	slog.Debug(line)

	var cmdParsers []cmdParser
	if admin {
		cmdParsers = adminCmdParsers
	} else {
		cmdParsers = clientCmdParsers
	}

	// find first space in string
	spaceIndex := strings.Index(line, " ")

	var beforeSpace, afterSpace string
	if spaceIndex == -1 {
		beforeSpace = line
		afterSpace = ""
	} else {
		beforeSpace = line[:spaceIndex]
		// cannot be last char because
		//  last char is break line
		afterSpace = line[spaceIndex+1:len(line) - 1]
	}

	for _, cmdParser := range cmdParsers {
		packetIn, result := cmdParser.check(beforeSpace, afterSpace, clientId)
		switch result {
		case Ok:
			return packetIn, nil
		case MissingParameter:
			return packetIn, errors.New("missing Parameter")
		}
	}

	return t.PacketIn{}, errors.New("commande Not Found")

}
