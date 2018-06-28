package main

import (
	"bufio"
	"fmt"
	"os"
)

type storyNode struct {
	text    string
	yesPath *storyNode
	noPath  *storyNode
}

func (node *storyNode) printStory(depth int) {
	for i := 0; i < depth*2; i++{
		fmt.Print("  ")
	}
	fmt.Print(node.text + "\n")
	if node.yesPath != nil{
		node.yesPath.printStory(depth + 1)
	}
	if node.noPath != nil{
		node.noPath.printStory(depth + 1)
	}	
}

func (node *storyNode) play() {
	fmt.Println(node.text)
	if node.yesPath != nil && node.noPath != nil {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			scanner.Scan()
			answer := scanner.Text()
			if answer == "yes" {
				node.yesPath.play()
				break
			} else if answer == "no" {
				node.noPath.play()
				break
			} else {
				fmt.Println("Incorrect input")
			}
		}
	}
}

func main() {

	root := storyNode{"You are at the entrance to a dark cave. Do you want to go in the cave?", nil, nil}
	winning := storyNode{"You have won!", nil, nil}
	losing := storyNode{"You have lost!", nil, nil}
	root.noPath = &winning
	root.yesPath = &losing

	root.printStory(0)
}