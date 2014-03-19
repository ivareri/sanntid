// TODO:
//
// main.go:
// Add watchdog. Restart lift if main process dies. (As in ex6)
//
// Network:
//
// UDP-multicast not guaranteed delivery.
//
// Possible fixes:
// -- While not confirmed by any other elevator: save to file
// -- Have elevators reply as they get messages. (should still be less networkoverhead compared to unicast
//
//
// Queue Priority:
// -- Relaclucate FS and check if order should be moved to other lifts.
// -- Need new status field to transmitt over network.
// -- Need abillity to remove order from localQueue.
//(localQueue.DeleteLocalOrder deletes commands, so not that one)
//
//
// Consitency in language
// -- order/request/command might not be consistent for all packages.
package main
