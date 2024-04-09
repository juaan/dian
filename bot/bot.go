package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var BotToken string

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}

func Run() {

	// create a session
	discord, err := discordgo.New("Bot " + BotToken)
	checkNilErr(err)

	// add a event handler
	discord.AddHandler(newMessage)

	// open session
	discord.Open()
	defer discord.Close() // close session, after function termination

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	ctx := context.Background()

	/* prevent bot responding to its own message
	this is achived by looking into the message author id
	if message.author.id is same as bot.author.id then just return
	*/
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// respond to user message if it contains `!help` or `!bye`
	switch {
	case strings.Contains(message.Content, "!ask"):
		question := getQuestionFromContent(message.Content)
		resp, err := askTailvy(ctx, question)
		if err != nil {
			fmt.Println("error: ", err.Error())
			discord.ChannelMessageSend(message.ChannelID, "failed processing question")
		}
		discord.ChannelMessageSend(message.ChannelID, resp.Answer)
		// add more cases if required
	}

}

func getQuestionFromContent(content string) string {
	listContentMsg := strings.Split(content, " ")
	return strings.Join(listContentMsg[2:len(listContentMsg)-1], " ")
}
