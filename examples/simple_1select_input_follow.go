package main

import (
	"fmt"

	"gopkg.in/AlecAivazis/survey.v1"
)

// the questions to ask
var simpleQs = []*survey.Question{
	{
		Name: "name",
		Prompt: &survey.Input{
			Message: "What is your name?",
		},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name: "color",
		Prompt: &survey.Select{
			Message: "Choose a color:",
			Options: []string{"red", "blue", "green"},
		},
		Validate: survey.Required,
	},
	{
		Name: "extras",
		Prompt: &survey.Input{
			Message: "Choose your extras (type it in):",
		},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
}

func main() {
	answers := struct {
		Name   string
		Color  string
		Extras string
	}{}

	// ask the question
	err := survey.Ask(simpleQs, &answers)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// print the answers
	fmt.Printf("%s chose %s, extras=%s\n", answers.Name, answers.Color, answers.Extras)
}
