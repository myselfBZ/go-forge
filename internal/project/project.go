package project

import (
	"embed"
	"text/template"
)

type Project struct{
    Name string
    DB        string
    Framework string

    config struct{
        ProjectName string
    }

    Fs  embed.FS
}





type File struct{
    Name string
    Src *template.Template
}

type Dir struct{
    Name string
    Files []File
    Dirs  []Dir
}
