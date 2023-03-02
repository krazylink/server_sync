/*
Copyright © 2023 Paul Huckabee <paul.huckabee@gmail.com>

*/
package main

import (
	"discordctl/cmd"
	_ "discordctl/cmd/channel"
	_ "discordctl/cmd/event"
)

func main() {
	cmd.Execute()
}
