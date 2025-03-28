package cmdinterface




const(
    Red = "red"
    WhiteBold = "whitebold"
    Reset   = "\033[0m"
)


var colors = map[string]string{
    "red": "\033[31m",
    "whitebold":"\033[1;37m",
}
