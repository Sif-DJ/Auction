# Auction
A timeunit in this program is 1 second



 
## To start a servernode:

Navigate to the ```Server``` folder and run the following command.

```go run .\server.go``` 

Then you need to input of one of the following numbers {"5050","5051","5052","5053"} start with the lowest number. 
After starting the first one you cannot select a number that is lower than the one you gave it then you can start the rest in any order you want and you can't input the same number more than once. The different instances of the server nodes need to be started in their own consoles.

## To start a client:

Navigate to the ```Client``` folder and run the following command.

```go run .\client.go``` 

It is now running and you can now submit bids and get the current result.
To get the current result input "outcome".
To bid, input an Integer