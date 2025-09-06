package utils

type Role int

const (
	Manager Role = iota
	Packer
	Picker
)

func (r Role) GetRole() string {
	switch r {
		case Manager:
			return "Manager"
		case Packer:
			return "Packer"
		case Picker:
			return "Picker"
		default:
			return "Unknown"
	}
}

func StringToRole(roleStr string) Role {
	switch roleStr {
	case "Manager":
		return Manager
	case "Packer":
		return Packer
	case "Picker":
		return Picker
	default:
		return -1 // Invalid role
	}
}