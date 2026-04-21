package serviceorderhandler

const (
	msgNotFound         = "service order not found"
	msgNotCancellable   = "only orders with status 'received' can be deleted"
	msgInvalidStatus    = "invalid status transition"
	msgInvalidVehicle   = "vehicle does not belong to the specified customer"
	msgStockInsufficient = "insufficient stock for one or more parts"
)
