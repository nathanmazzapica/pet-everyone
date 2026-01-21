[ pet service ] -- pet_result{userID, count} --> [ router ]


## Initial Plan
1. router sends pet command {userID} to pet service
2. pet service returns {userID, count}
3. router forwards {userID, count} to leaderboard service
4. leaderboard service updates redis sorted set
5. leaderboard returns top 10 {displayName, count} IF CHANGED. Otherwise returns nothing



## Obstacles


### 1. Pet Service Shape 
The pet service does not currently have a real-time count for every actor.
It only holds `saveInterval` second deltas. 

This was designed to reduce pressure the amount of blocking I/O spent on the database.

Possible solutions include:
- 

