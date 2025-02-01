# meganex DataStore

the intent is for this to be a fully generic DataStore that can eventually be merged into -common (like matchmaking is today)

vaguely based on [puyo_datastore](https://github.com/PretendoNetwork/puyo-puyo-tetris/blob/mm-server/datastore/puyo_datastore.go) and [smm](https://github.com/PretendoNetwork/super-mario-maker/tree/master/database/datastore)

## todo
- test the edge cases
- Ratings
- deletion and unfinished uploads (where is returning deleted data OK and not OK?)
- Under review
- Passwords need testing since idk of any games that use them
- put expiry time and referred time in the database