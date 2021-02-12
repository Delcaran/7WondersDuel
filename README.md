# 7WondersDuel
Basic implementation of 7 Wonders Duel boardgame as an exercise in learning Go

## Disclaimer
Everything under copyright is property of their respective owners.

If I commit any copyright violation, please contact me and I will remove it.

## Scope

A text only implementation of the 7 Wonders Duel boardgame.

## Roadmap

Development phases will be the following:

1. [X] Definition of every entity involved in the game
2. [ ] Definition of relationships between identities
3. [ ] Implementation of interactions between entities
4. [ ] User interface
5. [ ] Local game with both players taking turns.
6. [ ] LAN game

## Further steps

These features may never be implemented:

- [ ] Internet serverless game
- [ ] Single player with adversarial A.I.

## Grid implentation 

Ages card layout is described in ages.dat as following:
- a space means "no card"
- O means visible card
- X means hidden card
- All "card" lines must have the same number of chars
- Blank line separates ages

This allows for easy layouting and maybe future customization.

