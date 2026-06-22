package cli

import (
	"bufio"
	"client/client"
	"context"
	"fmt"
	"os"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

const (
	RegisterChoice = iota + 1
	LoginChoice
	IsAdminChoice
	ExitAction
)

var reader = bufio.NewReader(os.Stdin)

func EventLoop(ctx context.Context, authClient *client.Client, address string) {
	HelloMessage(address)
	for {
		choice := Menu()
		switch choice {
		case RegisterChoice:
			Register(ctx, authClient)
		case LoginChoice:
			Login(ctx, authClient)
		case IsAdminChoice:
			IsAdmin(ctx, authClient)
		case ExitAction:
			GoodbyeMessage(address)
			os.Exit(0)
		default:
			color.Yellow("Action doesn't exist")
		}
	}
}

func HelloMessage(address string) {
	color.Green("Connected to %s\n", address)
}

func GoodbyeMessage(address string) {
	color.Red("Disconnected from %s\n", address)
}

func Menu() int {
	options := []string{
		"Register",
		"Log in",
		"Check admin status",
		"Stop",
	}

	fmt.Println("╔══════════════════════════╗")
	fmt.Println("║     🔐 AUTH SERVICE      ║")
	fmt.Println("╚══════════════════════════╝")

	for i, line := range options {
		fmt.Printf("%d. %s\n", i+1, line)
	}

	fmt.Print("Choose action: ")
	var choice int
	fmt.Scan(&choice)
	reader.ReadString('\n')
	return choice
}

func Register(ctx context.Context, authClient *client.Client) {
	var email string
	fmt.Print("Enter email: ")
	fmt.Scanf("%s", &email)

	password, err := ReadPassword()
	if err != nil {
		color.Red(err.Error())
		return
	}

	userId, err := authClient.Register(ctx, email, password)
	if err != nil {
		color.Red(err.Error())
	} else {
		color.Green("Registration successful. Your id: %d\n", userId)
	}
}

func Login(ctx context.Context, authClient *client.Client) {
	var email string
	fmt.Print("Enter email: ")
	fmt.Scanf("%s", &email)

	var password string
	password, err := ReadPassword()
	if err != nil {
		color.Red(err.Error())
		return
	}
	reader.ReadString('\n')

	fmt.Println("Available apps:")
	options := []string{"Gemini", "Gmail", "Chrome"}
	for i, line := range options {
		fmt.Printf("%d - %s\n", i+1, line)
	}

	var appId int32
	fmt.Print("Choose app: ")
	fmt.Scanf("%d", &appId)

	token, err := authClient.Login(ctx, email, password, appId)
	if err != nil {
		color.Red(err.Error())
	} else {
		color.Green("Login successful.\nYour JWT token:\n%s\n", token)
	}
}

func IsAdmin(ctx context.Context, authClient *client.Client) {
	var userId int64
	fmt.Print("Enter user id: ")
	fmt.Scanf("%d", &userId)

	isAdmin, err := authClient.IsAdmin(ctx, userId)
	if err != nil {
		color.Red(err.Error())
	} else {
		if isAdmin {
			color.Green("Is admin: %t\n", isAdmin)
		} else {
			color.Red("Is admin: %t\n", isAdmin)
		}
	}
}

func ReadPassword() (string, error) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:     "Enter password: ",
		EnableMask: true,
		MaskRune:   '*',
	})
	if err != nil {
		return "", err
	}
	defer rl.Close()

	return rl.Readline()
}
