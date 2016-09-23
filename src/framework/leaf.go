package framework

import (
	"framework/conf"
	"framework/console"
	"framework/log"
	"framework/module"
	"os"
	"os/signal"
)

func Run(mods ...module.Module) {
	// logger
	if conf.SvrBase.LogLevel != "" {
		logger, err := log.New(conf.SvrBase.LogLevel, conf.SvrBase.LogPath)
		if err != nil {
			panic(err)
		}
		log.Export(logger)
		defer logger.Close()
	}

	log.Release("Leaf %v starting up", version)

	// module
	for i := 0; i < len(mods); i++ {
		module.Register(mods[i])
	}
	module.Init()

	// console
	console.Init()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Release("Leaf closing down (signal: %v)", sig)
	console.Destroy()
	module.Destroy()
}
