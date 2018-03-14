package main

import (
	//"github.com/watts-kit/passwordd/passwordclib"
	"../passwordclib"

	"bufio"
	"gopkg.in/alecthomas/kingpin.v2"
	"fmt"
	"os"
	"strings"
)

const(
	version = "1.0.0"
)


var (
	versionText = "cli:"+version+", lib:"+passwordclib.Version()
	app = kingpin.New("passwordc",
		"a cli client to set/get passwords").Version(versionText)
	set_pwd     = app.Command("set", "set a secret")
	set_pwd_key = set_pwd.Arg("key", "the id of key to set").Required().String()

	get_pwd     = app.Command("get", "get a secret")
	get_pwd_key = get_pwd.Arg("key", "the id of key to get").Required().String()

    overwrite_pwd     = app.Command("overwrite", "replace a secret")
    overwrite_pwd_key = overwrite_pwd.Arg("overwrite", "the id of key to replace").Required().String()
)


func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case get_pwd.FullCommand():
		password, err := passwordclib.GetPassword(*get_pwd_key)
		if err != nil {
			fmt.Println("error getting password")
			os.Exit(3)
		}
		fmt.Printf("returned secret: '%s'\n", password)
	case set_pwd.FullCommand():
		fmt.Print("please enter the secret: ")
		reader := bufio.NewReader(os.Stdin)
		secret, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("could not read secret")
			os.Exit(2)
		}
		secret = strings.Trim(secret, "\n")
		_, err = passwordclib.SetPassword(*set_pwd_key, secret)
		if err != nil {
			fmt.Println("could not set secret")
			os.Exit(4)
		}
		fmt.Println("secret set")
	case overwrite_pwd.FullCommand():
		fmt.Print("please enter the secret: ")
		reader := bufio.NewReader(os.Stdin)
		secret, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("could not read secret")
			os.Exit(2)
		}
		secret = strings.Trim(secret, "\n")
		_, err = passwordclib.OverwritePassword(*overwrite_pwd_key, secret)
		if err != nil {
			fmt.Println("could not replace secret")
            fmt.Println(err)
			os.Exit(4)
		}
		fmt.Println("secret replaced")
	}
}
