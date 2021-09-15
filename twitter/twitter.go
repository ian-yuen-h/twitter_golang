package main

import(
	"os"
	"strconv"
	"proj1/server"
	"encoding/json"

)
func main() {


	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)

	var config server.Config

	if len(os.Args) != 3{

		config = server.Config{Encoder: enc, Decoder: dec, Mode: "s"}
	} else {
		consNum, _  := strconv.Atoi(os.Args[1])
		blockSize, _  := strconv.Atoi(os.Args[2])
		config = server.Config{Encoder: enc, Decoder: dec, Mode: "p", ConsumersCount: consNum, BlockSize: blockSize}
	}

	server.Run(config)

}
