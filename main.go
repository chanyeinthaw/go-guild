package main

import (
	"fmt"
	"go-guild/client"
	"go-guild/options"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	setupLogger()

	opts := readOptions()
	dg := client.New(opts.Token, opts.Prefix)
	dg.Connect()

	switch opts.OP {
	case options.OpDel:
		dg.DeleteGuild(opts.GuildID)

		fmt.Println("Guilds")
		for _, guild := range dg.Guilds() {
			fmt.Println("Guild ID: ", guild.ID)
		}
	case options.OpCreateOrManage:
		var guildIDToManage *string
		if opts.GuildID != "" {
			guildIDToManage = &opts.GuildID
		}
		guild := dg.CreateOrManageGuild(opts.Name, guildIDToManage)
		invite := dg.CreateInviteCode()
		fmt.Println("Guild   : ", guild.Name)
		fmt.Println("Guild ID: ", guild.ID)
		fmt.Println("Invite  : ", invite)
		fmt.Println("OTP     : ", dg.OTP())

		fmt.Println("Bot is now running.  Press CTRL-C to exit.")
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
	case options.OpList:
		guilds := dg.Guilds()
		fmt.Println("Guilds: ", len(guilds))
		for _, guild := range guilds {
			fmt.Println("Guild ID: ", guild.ID)
		}
	}

	dg.Close()
}

func setupLogger() {
	log.SetPrefix("go-guild: ")
	log.SetFlags(1)
}

func readOptions() options.Options {
	opts, err := options.Read()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	return *opts
}