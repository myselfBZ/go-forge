package main

import (
	"embed"
	"os"

	"github.com/myselfBZ/go-forge/internal/cmdinterface"
	"github.com/myselfBZ/go-forge/internal/project"
)

type tools struct{
    libs []string
    dbs  []string
}

type App struct {
    tools tools
    fs embed.FS
	cmdinterface cmdinterface.CmdInterface
}


func (a *App) Start() {
    a.cmdinterface.PrintWithColor(cmdinterface.WhiteBold, "Welcome to Go Forge\n")

    var prjct project.Project
    prjct.Fs = a.fs
    prjct.Name = a.cmdinterface.Prompt("Enter the project name")


    // lib selection
    prjct.Framework = a.cmdinterface.Select(a.tools.libs, "choose your http framework")

    // db selection
    prjct.DB = a.cmdinterface.Select(a.tools.dbs, "choose your db")

    // auth
    prjct.WithAuth = a.cmdinterface.Conifrm("Do you want user authentication and authorization?")    

    if err := prjct.Build(); err != nil{
        a.cmdinterface.PrintWithColor(cmdinterface.Red, err.Error())
        os.Exit(1)
    }
}
