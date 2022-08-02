package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"github.com/shomali11/slacker"
)

// func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
// 	for event := range analyticsChannel {
// 		fmt.Println("Command Events")
// 		fmt.Println(event.Timestamp)
// 		fmt.Println(event.Command)
// 		fmt.Println(event.Parameters)
// 		fmt.Println(event.Event)
// 		fmt.Println()
// 	}
// }

type User struct {
	Name  string
	Meals []Meal
}

type Meal struct {
	Name  string
	price int
}

var users []User

func main() {
	
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	// go printCommandEvents(bot.CommandEvents())

	bot.Command("{operation} {price} {meal} {number} {userName}", &slacker.CommandDefinition{
		Description: "+/- 餐  數量",
		Examples:    []string{"+ 排骨飯 1"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			operation := request.Param("operation")
			meal := request.Param("meal")
			number, err := strconv.Atoi(request.Param("number"))
			if err != nil {
				fmt.Println(err)
				return
			}

			var newUser User
			var newMeal Meal

			newUser.Name = request.Param("userName")
			if "" == newUser.Name {
				newUser.Name = botCtx.Event().UserName
			}
			
			newMeal.Name = meal
			newMeal.price = number
			newUser.Meals = []Meal{newMeal}

			already := false

			for i := range users {
				if(users[i].Name == newUser.Name) {
					users[i].Meals = append(users[i].Meals, newMeal)
					already = true
					break
				}
			}

			if !already {
				users = append(users, newUser)
			}

			for i := range users {
				fmt.Println(users[i])
			}

			r := fmt.Sprintf("%s %s %d %s", operation, meal, number, newUser.Name)

			response.Reply(r)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
