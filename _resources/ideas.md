
# Ideas

* Create "magic" brick that would be something good or something bad to the player's field or to all opponents' field.
  The ones that affect the player's field are usable in 1 player games.
  Ideas:
    - Good: Remove few lines at the bottom, remove few top most blocks
    - Bad: Raise by few lines at the bottom, add few blocks on top
    - Weapons: Raise by few lines at the bottom, add few blocks on top, stairs
* Create additional sweepers, like the one that removes islands of the same color that consist of 4 or more blocks
* Create simple AI for human vs computer play 
* Create a different state machine for maybe completely different type of games (for example: where blocks do not fall)
* Challenge game mode
  - Destroy the living block that slowly moves across the field and eats blocks.
  - Destroy all target blocks in a dynamic, changing field.


# Components inside the main game loop

* Player input
* Piece controllers, each with timer
* Sweepers, one per field, with timer
* Async field changer (Tracks magic blocks, periodic field changes), one per field, with timer
* Score tracker
* Render request handler
* AI controller, one for each player controller, with timer (to periodically emit actions) 
