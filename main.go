package main

import (
	"embed"
	_ "embed"

	"github.com/myselfBZ/go-forge/internal/cmdinterface"
)

//go:embed internal/src-files/**/*
var fs embed.FS

func main() {

	app := App{
        fs: fs,
        cmdinterface: *cmdinterface.NewCmdInterface(),
        tools: tools{

            // libraries
            libs: []string{
                "Chi",
                "Gin",
                "Stdlib",
                "Echo",
                "Fiber",
            },

            // databases 
            dbs: []string{
                "Postgres",
                "Redis",
                "MongoDB",
                "MySQL",
            },

        },
	}

    app.Start()
}
