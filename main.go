package main

import (
	"fmt"
	"leopard-quant/bootstrap"
	"leopard-quant/core/log"
	"leopard-quant/restful"
	"os/exec"
	"strconv"
)

func openBrowser(serv *restful.Server) {
	port, _ := strconv.Atoi(serv.Port)
	url := fmt.Sprintf("http://localhost:%d/test/welcome", port)
	cmd := exec.Command("open", url)
	if err := cmd.Start(); err != nil {
		log.Error(err)
	}
}

func main() {
	bootstrap.Init()
	restful.NewServer().ServerStartedListener(nil).Start()
}
