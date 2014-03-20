// Things that might be improved:
//
// Network:
// UDP-multicast not guaranteed delivery.
//
// Possible fixes for UDP:
// -- While not confirmed by any other elevator: save to file
// -- Have elevators reply as they get messages. (should still be less networkoverhead compared to unicast
//
// Consitency in language
// -- order/request/command might not be consistent for all packages.
//
// Cost function (figureOfSuitabillity) should take lenght of command queue into consideration
package main
