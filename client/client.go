package client

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"go-guild/options"
	"log"
	"os"
	"strings"
)

type Client struct {
	session *discordgo.Session
	token string
	prefix string
	guildID *string
	ownerRole *string
	otp *string
}

func New(token string, prefix string) Client {
	dg, _:= discordgo.New("Bot " + token)

	return Client{
		session: dg,
		token: token,
		prefix: prefix,
		guildID: nil,
		otp: nil,
		ownerRole: nil,
	}
}

func (client *Client) Connect() {
	err := client.session.Open()

	if err != nil {
		log.Println("Invalid bot token", client.token, err.Error())
		os.Exit(1)
	}
}

func (client *Client) Close() {
	_ = client.session.Close()
}

func (client *Client) getOwnerRoleID() *string {
	roles, err := client.session.GuildRoles(*client.guildID)
	if err != nil { return nil }

	for _, role := range roles {
		if role.Name == "owner" && role.Permissions == discordgo.PermissionAdministrator {
			return &role.ID
		}
	}

	return nil
}

func (client *Client) CreateOrManageGuild(name string, id *string) *discordgo.Guild {
	var guild *discordgo.Guild
	if id == nil {
		created, err := client.session.GuildCreate(name)
		if err != nil {
			log.Println("Unable to create guild")
			log.Println(err)
			os.Exit(1)
		}

		guild = created
	} else {
		created, err := client.session.Guild(*id)
		if err != nil {
			log.Println("Unable to find guild")
			log.Println(err)
			os.Exit(1)
		}

		guild = created
	}

	client.guildID = &guild.ID

	client.session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if !strings.HasPrefix(m.Content, client.prefix) {
			return
		}

		if strings.HasPrefix(m.Content, client.prefix + "help") {
			_, _ = client.session.ChannelMessageSend(m.ChannelID,
				"**!own <OTP>** - take **owner** role\n" +
					"**!release <OTP>** - release **owner** role\n" +
					"**!transfer <OTP>** - transfer ownership",
			)

			return
		}

		cmdParts := strings.Split(m.Content, " ")
		if len(cmdParts) < 2 {
			_,  _ = client.session.ChannelMessageSend(m.ChannelID, "Invalid OTP.")
			return
		}

		if cmdParts[1] != *client.otp {
			_,  _ = client.session.ChannelMessageSend(m.ChannelID, "Invalid OTP.")
			return
		}

		if strings.HasPrefix(m.Content, client.prefix + "release") {
			client.handleRelease(m)
		} else if strings.HasPrefix(m.Content, client.prefix + "own")  {
			client.handleOwn(m)
		} else if strings.HasPrefix(m.Content, client.prefix + "transfer") {
			client.handleTransfer(m)
		}

		fmt.Println("New OTP : ", client.OTP())
	})

	return guild
}

func (client *Client) handleTransfer(m * discordgo.MessageCreate) {
	_, err := client.session.GuildEdit(*client.guildID, discordgo.GuildParams{
		OwnerID: m.Author.ID,
	})
	if err != nil {
		_, _ = client.session.ChannelMessageSend(m.ChannelID, "Can not transfer ownership.")
		return
	}

	_, _ = client.session.ChannelMessageSend(m.ChannelID, "Transfer ownership completed.")
}

func (client *Client) handleOwn(m *discordgo.MessageCreate) {
	ownerRole := client.getOwnerRoleID()
	if ownerRole == nil {
		role, err := client.session.GuildRoleCreate(*client.guildID)
		if err != nil {
			_, _ = client.session.ChannelMessageSend(m.ChannelID, "Failed to add **owner** role.")
			return
		}

		_, _ = client.session.GuildRoleEdit(
			*client.guildID,
			role.ID,
			"owner", 0xFAFAFA, true,
			discordgo.PermissionAdministrator, true)

		ownerRole = &role.ID
	}

	err := client.session.GuildMemberRoleAdd(*client.guildID, m.Author.ID, *ownerRole)
	if err != nil {
		_, _ = client.session.ChannelMessageSend(m.ChannelID, "Failed to add **owner** role.")

		return
	}

	_, _ = client.session.ChannelMessageSend(m.ChannelID, "**owner** role added.")
}

func (client *Client) handleRelease(m *discordgo.MessageCreate) {
	ownerRole := client.getOwnerRoleID()
	if ownerRole != nil {
		err := client.session.GuildMemberRoleRemove(*client.guildID, m.Author.ID, *ownerRole)
		if err != nil {
			_, _ = client.session.ChannelMessageSend(m.ChannelID, "Failed to release **owner** role.")

			return
		}
	}
	_, _ = client.session.ChannelMessageSend(m.ChannelID, "**owner** role released.")
}

func (client *Client) CreateInviteCode() string {
	channels, err := client.session.GuildChannels(*client.guildID)
	if err != nil {
		return ""
	}

	var channelID *string
	channelID = nil
	for _, ch := range channels {
		if ch.Type == discordgo.ChannelTypeGuildText {
			channelID = &ch.ID
			break
		}
	}

	if channelID == nil {
		return ""
	}

	invite, err := client.session.ChannelInviteCreate(*channelID, discordgo.Invite{
		MaxAge: 0,
		MaxUses: 0,
		Temporary: false,
		Unique: false,
	})
	if err != nil {
		return ""
	}

	return invite.Code
}

func (client *Client) DeleteGuild(id string) {
	_, _ = client.session.GuildDelete(id)
}

func (client *Client) Guilds() []*discordgo.UserGuild {
	guilds, err := client.session.UserGuilds(10, "", "")
	if err != nil {
		return []*discordgo.UserGuild{}
	}

	var ownGuilds []*discordgo.UserGuild
	for _, guild := range guilds {
		if guild.Owner == true {
			ownGuilds = append(ownGuilds, guild)
		}
	}

	return ownGuilds
}

func (client *Client) OTP() string {
	otp := options.GetRandStr()

	client.otp = &otp

	return otp
}