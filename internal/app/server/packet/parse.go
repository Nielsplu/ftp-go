package packet

import (
	t "ftp/internal/app/server/types"
)

func Parse(line string, answerChan chan t.PacketOut) (t.PacketIn, error) {
	// TODO
	if (line == "END\n") {
		return t.PacketIn{ Type: t.End }, nil
	}

	if (line == "T\n") {
		return t.PacketIn{ Type: t.Terminate }, nil
	}

	return t.PacketIn{ Type: t.List, Path: "/salut man comment ça va ? \n est ce que ça marche là ? \n", AnswerChan: answerChan }, nil
}