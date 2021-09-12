package options

import (
	"errors"
	"flag"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Options struct {
	Name   string
	Token  string
	Prefix string
	GuildID string
	OP string
}

const (
	OpCreateOrManage string = "cm"
	OpList string = "ls"
	OpDel string = "del"
)

func Read() (*Options, error) {
	name := flag.String("name", getServerName(), "Server name")
	token := flag.String("token", "", "Bot access token")
	prefix := flag.String("prefix", "!", "Bot command prefix")
	guildID := flag.String("guild", "", "Guild id")
	op := flag.String("op", "list", "Operation ls, cm, del")

	flag.Parse()

	if *token == "" {
		return nil, errors.New("bot access Token is required")
	}

	if *op == "delete" && *guildID == "" {
		return nil, errors.New("guild id is required")
	}

	options := &Options{Name: *name, Token: *token, Prefix: *prefix, GuildID: *guildID, OP: *op}

	return options, nil
}

func GetRandStr() string {
	now := time.Now()
	rand.Seed(now.UnixNano())
	randPosFix := strings.ToUpper(strconv.FormatUint(rand.Uint64(), 16)[0:6])

	return randPosFix
}

func getServerName() string {
	return "Server - " + GetRandStr()
}