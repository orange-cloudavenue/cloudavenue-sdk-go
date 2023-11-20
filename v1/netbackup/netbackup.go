package netbackup

type Netbackup struct {
	VCloud          VcloudClient
	ProtectionLevel ProtectionLevelClient
	Machines        MachineClient
	Inventory       InventoryClient
}
