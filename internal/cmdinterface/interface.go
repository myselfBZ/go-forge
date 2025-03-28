package cmdinterface

import (
	"bufio"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)



type CmdInterface struct{
    scanner  bufio.Scanner   
}


func NewCmdInterface()  *CmdInterface {
    scanner := bufio.NewScanner(os.Stdin)
    return &CmdInterface{
        scanner: *scanner,
    }
}

func (i *CmdInterface) Conifrm(message string) bool {
    name := false

    prompt := &survey.Confirm{
        Message: message,
    }
    survey.AskOne(prompt, &name)

    return  name
}

func (i *CmdInterface) Prompt(prompt string) string {
    i.PrintWithColor(WhiteBold, prompt)
    fmt.Printf("> ")

    if i.scanner.Scan() {
        return i.scanner.Text()    
    }

    if err := i.scanner.Err(); err != nil{
        fmt.Println("Err occured: ", err)
    }

    return ""
}


func (i *CmdInterface) Select(opts []string, message string) string {
    var selected string

    promt := &survey.Select{
        Message: message,
        Options: opts,
    }

    survey.AskOne(promt, &selected)
    return selected
}

/* if the entered color doesn't exist, default color of is used (white)
Available options:
    whitebold
    red
*/
func (i *CmdInterface) PrintWithColor(color string, prompt string) {
    textColor := colors[color]    
    fmt.Printf("%s%s\n", textColor, prompt)
}






