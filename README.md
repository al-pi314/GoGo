## Description

GO game written in golang. Games can be played aginst geneticaly trained AI agents.

## File structure

- agents -> holds currently best trained agents
- cmd 
    - play -> main function file for playing the game
    - train -> main function file for training agents
- game -> functions required for game logic
- nn -> functions required to run NN
- population -> functions required to train, save and load populations
- player -> functions required to execute human or agent commands

## Commands

Basic <strong>game start</strong>: go run ./cmd/play/play.go <br>
Paramateres:

- white -> set to human for human player or to agent for AI player
- black -> set to human for human player or to agent for AI player
- population -> set to population.json file to loade AI players from

Basic <strong>training start</strong>: go run ./cmd/train.go <br>
Paramateres: 

- population -> set to population.json file from which to build initial population
- output -> set to file used for saving trained populations