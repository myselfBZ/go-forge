package project

import "embed"

type Project struct{
    Name string
    DB        string
    Framework string

    Fs  embed.FS
}





type File struct{
    Name string
    Src string
}

type Dir struct{
    Name string
    Files []File
    Dirs  []Dir
}
