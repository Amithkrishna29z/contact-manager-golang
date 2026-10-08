package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"contact-manager/repository"
	"contact-manager/service"
)

func main() {
	repository := repository.NewContactRepository()
	service := service.NewContactService(repository)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println()
		fmt.Println("===== CONTACT MANAGER =====")
		fmt.Println("1. Add Contact")
		fmt.Println("2. List Contacts")
		fmt.Println("3. Find Contact")
		fmt.Println("4. Update Contact")
		fmt.Println("5. Delete Contact")
		fmt.Println("6. Exit")
		fmt.Print("Enter choice: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		choice, err := strconv.Atoi(input)

		if err != nil {
			fmt.Println("Please enter a valid number")
			continue
		}

		switch choice {

		case 1:
			addContact(reader, service)

		case 2:
			listContacts(service)

		case 3:
			findContact(reader, service)

		case 4:
			updateContact(reader, service)

		case 5:
			deleteContact(reader, service)

		case 6:
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Invalid choice")
		}
	}
}

func addContact(
	reader *bufio.Reader,
	service *service.ContactService,
) {

	fmt.Print("Name: ")
	name, _ := reader.ReadString('\n')

	fmt.Print("Phone: ")
	phone, _ := reader.ReadString('\n')

	name = strings.TrimSpace(name)
	phone = strings.TrimSpace(phone)

	contact, err := service.AddContact(
		name,
		phone,
	)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Contact created successfully!")
	fmt.Println("ID:", contact.ID)
}

func listContacts(service *service.ContactService) {
	contacts := service.GetContacts()

	if len(contacts) == 0 {
		fmt.Println("No contacts found")
		return
	}

	fmt.Println()
	fmt.Println("===== CONTACTS =====")

	for _, contact := range contacts {
		fmt.Printf(
			"ID: %d | %s | %s\n",
			contact.ID,
			contact.Name,
			contact.Phone,
		)
	}
}

func readContactId(
	reader *bufio.Reader,
) int {
	fmt.Print("Enter contact ID: ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	id, err := strconv.Atoi(input)
	
	if err != nil {
		fmt.Println("Invalid ID")
		return 0
	}
	return id
}

func findContact(
	reader *bufio.Reader,
	service *service.ContactService,
) {
	id := readContactId(reader)
	contact, err := service.GetContact(id)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println()
	fmt.Println("===== CONTACT =====")
	fmt.Println("ID:", contact.ID)
	fmt.Println("Name:", contact.Name)
	fmt.Println("Phone:", contact.Phone)
}

func updateContact(
	reader *bufio.Reader,
	service *service.ContactService,
) {
	id := readContactId(reader)

	fmt.Println("Give detatils to update the contact")

	fmt.Print("Name: ")
	name, _ := reader.ReadString('\n')

	fmt.Print("Phone: ")
	phone, _ := reader.ReadString('\n')

	name = strings.TrimSpace(name)
	phone = strings.TrimSpace(phone)

	contact, err := service.UpdateContact(
		id,
		name,
		phone,
	)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Contact updated successfully!")
	fmt.Println(contact)
}

func deleteContact(
	reader *bufio.Reader,
	service *service.ContactService,
) {
	id := readContactId(reader)

	err := service.DeleteContact(id)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Contact deleted successfully!")
}

