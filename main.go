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

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		// fmt.Println(event.Timestamp)
		// fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		// fmt.Println(event.Event)
		// fmt.Println()
	}
}

type User struct {
	Name  string
	Meals map[string]Meal
}

type Meal struct {
	Name  string
	price int
	amount int
}

var users = make(map[string]User)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	go printCommandEvents(bot.CommandEvents())

	bot.Command("{operation} {price} {meal} {number} {userName}", &slacker.CommandDefinition{
		Description: "+/- 餐  數量",
		Examples:    []string{"+ $50 排骨飯 1"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			operation := request.Param("operation")
			meal := request.Param("meal")
			price, err := strconv.Atoi(request.Param("price"))
			if err != nil {
				fmt.Println(err)
				return
			}
			number, err := strconv.Atoi(request.Param("number"))
			if err != nil {
				fmt.Println(err)
				return
			}

			var newUser User
			

			newUser.Name = request.Param("userName")
			if "" == newUser.Name {
				newUser.Name = botCtx.Event().UserName
			}

			newMeal := Meal {
				Name: meal,
				price: price,
				amount: number,
			}		
			
			newUser.Meals = map[string]Meal{newMeal.Name: newMeal}

			user, existingUser := users[newUser.Name]
			if existingUser {
				meal, existingMeal := user.Meals[newMeal.Name]
				if existingMeal {
					meal.amount += number
					user.Meals[newMeal.Name] = meal
				}else{
					user.Meals[newMeal.Name] = newMeal
				}
			}else{
				users[newUser.Name] = User{
					Name: newUser.Name,
					Meals: map[string]Meal{
						newMeal.Name: newMeal,
					},
				}
			}

			r := fmt.Sprintf("點餐成功! %s: %s %d %s", operation, meal, number, newUser.Name)

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
