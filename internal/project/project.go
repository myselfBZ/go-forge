package project

import "embed"

type Project struct{
    Name string
    DB        string
    Framework string
    WithAuth bool
    WithCache bool

    Fs  embed.FS
}





type File struct{
    Name string
    Src []byte
}

type Dir struct{
    Name string
    Files []File
    Dirs  []Dir
}
