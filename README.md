# Distributed Broadcast Service
This project implements a distributed broadcast system, where nodes reach eventual consistency on validated news.

Note that each node could join or leave the system at any time, make sure nodes can discover each other (with or without a centural registry).

Here lists pseudocode of behaviour offered to clients by a single node:
```go
package server
type Data struct {
    id string   
    raw string
    signature string
}
 
type BroadcastService interface {
    // Returns true if data is valid and accepted
    PostNewData(data Data)

    // Return all data received in the past 24h
    ListAllRecentData() []Data
}
``` 