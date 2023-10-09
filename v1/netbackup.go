package v1

type Netbackup struct {
	VCloud          VCloudClient
	ProtectionLevel ProtectionLevelClient
	Machines        MachineClient
	Inventory       InventoryClient
}
