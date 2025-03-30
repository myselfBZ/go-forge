package main

import (
    "github.com/myselfBZ/go-forge/internal/src-files"
	"github.com/myselfBZ/go-forge/internal/cmdinterface"
)


func main() {

	app := App{
        fs: srcfiles.FS,
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
