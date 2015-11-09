package main

import "github.com/kildevaeld/vault/Godeps/_workspace/src/github.com/kildevaeld/go-widgets"

func main() {
	paginated := widgets.PaginatedList{
		Message: "test",
		Paginate: func(page int) []string {
			switch page {
			case 1:
				return []string{"Page1: Item 1", "Page1: Item 2"}
			case 2:
				return []string{"Page2: Item 1", "Page2: Item 2"}
			default:
				return nil
			}

		},
	}
	paginated.Run()
	// Input
	input := widgets.Input{
		Message: "Enter name",
	}
	input.Run()

	confirm := widgets.Confirm{
		Message: "Confirm",
	}
	confirm.Run()

	password := widgets.Password{
		Message: "Password",
	}
	password.Run()

	list := widgets.List{
		Message: "List",
		Choices: []string{"Test", "Test 2"},
	}
	list.Run()

	check := widgets.Checkbox{
		Message: "Checkbox",
		Choices: []string{"Test", "Test 2"},
	}
	check.Run()
}
