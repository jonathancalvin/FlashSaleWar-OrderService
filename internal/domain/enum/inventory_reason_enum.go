package enum

type InventoryReleaseReason string

const (
	InvOrderCancelled InventoryReleaseReason = "ORDER_CANCELLED"
	InvOrderExpired   InventoryReleaseReason = "ORDER_EXPIRED"
	InvTTLExpired     InventoryReleaseReason = "TTL_EXPIRED"
)

func (r InventoryReleaseReason) IsValid() bool {
	switch r {
		case InvOrderCancelled, InvOrderExpired, InvTTLExpired:
			return true
		default:
			return false
	}
}