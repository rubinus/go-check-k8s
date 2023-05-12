package check

type InventoryTask struct {
	Name  string `json:"name"`
	Shell string `json:"shell"`
}

type Inventory struct {
	Name        string          `json:"name"`
	Hosts       string          `json:"hosts"`
	GatherFacts bool            `json:"gather_facts"`
	Tasks       []InventoryTask `json:"tasks"`
}
