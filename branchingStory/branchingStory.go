package main

import (
	"os"
	"bufio"
	"fmt"
)

type storyNode struct{
	text string
	yesNode *storyNode
	noNode *storyNode
}

func (node *storyNode) play(){
	fmt.Println(node.text)
	if node.yesNode != nil && node.noNode != nil {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			scanner.Scan()
			answer := scanner.Text()
			if answer == "yes"{
				node.yesNode.play()
				break
			}else if answer == "no"{
				node.noNode.play()
				break
			}else{
				fmt.Println("Incorrect input")
			}
		}
	}	
}

func main() {

	root := storyNode{"You are at the entrance to a dark cave. Do you want to go in the cave?", nil, nil}
	winning := storyNode{"You have won!", nil, nil}
	losing := storyNode{"You have lost!", nil, nil}
	root.noNode = &winning
	root.yesNode = &losing

	root.play()
}