## Description

GO game written in golang. Games can be played aginst geneticaly trained AI agents.

## File structure

- cmd 
    - play -> main function file for playing and replaying games
    - train -> main function file for training agents
- game -> functions required for game logic
- nn -> functions required to run NN
- player -> functions required to execute human or agent commands
- population -> functions required to train, save and load populations
- / -> shared go types and interfaces

## Commands

Basic <strong>game start</strong>: go run ./cmd/play/play.go <br>
Paramateres:

- white -> set to human for human player or to agent for AI player
- black -> set to human for human player or to agent for AI player
- population -> set to population.json file to load AI players from
- replay -> set to game.json file to replay game <strong>(set delay flag)</strong>
- delay -> set number of miliseconds the agent should wait after the move, also <strong>effects replay speed</strong>

Basic <strong>training start</strong>: go run ./cmd/train.go <br>
Paramateres: 

- population -> set to population.json file from which to build initial population
- output -> set to file used for saving trained populations

## Encountered Problems & Solutions