package game

import (
	"fmt"
)

type Monster struct {
	Character
}


func NewRat(p Pos) *Monster{
	//return &Monster{Pos:p, Rune:'R', Name:"Rat", Hitpoints:5, Strength:5, Speed:2.0, ActionPoints:0.0}
	monster := &Monster{}
	monster.Pos = p
	monster.Rune = 'R'
	monster.Name = "Rat"
	monster.Hitpoints = 5
	monster.Strength = 5
	monster.Speed = 2.0
	monster.ActionPoints = 0
	return monster
}

func NewSpider(p Pos) *Monster {
	//return &Monster{p,'S', "Spider", 10, 10, 1.0, 0.0}
	monster := &Monster{}
	monster.Pos = p
	monster.Rune = 'S'
	monster.Name = "Spider"
	monster.Hitpoints = 10
	monster.Strength = 10
	monster.Speed = 1.0
	monster.ActionPoints = 0
	return monster
}


func (m *Monster) Update(level *Level){
	m.ActionPoints += m.Speed
	playerPos := level.Player.Pos
	apInt := int(m.ActionPoints)
	positions := level.astar(m.Pos, playerPos)
	moveIndex := 1
	for i := 0; i < apInt; i++{
		if moveIndex < len(positions) {
			fmt.Println("Move")
			m.Move(positions[moveIndex], level)
			moveIndex++
			m.ActionPoints--
		}
	}
}

func (m *Monster) Move(to Pos, level *Level) {
	_, exists := level.Monsters[to]
	//TODO: check if tile to move to is valid
	if !exists && to != level.Player.Pos{
		delete(level.Monsters, m.Pos)
		level.Monsters[to] = m
		m.Pos = to
	}else{
		MonsterAttackPlayer(&level.Player, m)
		fmt.Println("Monster Attack Player")
		fmt.Println(level.Player.Hitpoints, m.Hitpoints)		
	}
}