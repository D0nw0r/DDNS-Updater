package main

import (
	sendrequest "ddns-updater/send-request"
	"fmt"
	"time"
)

func mainLoop() {

	var test bool = false

	for {

		fmt.Println("Finding Zone Details..")
		name, id, err := sendrequest.GetZoneId()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("\t[+] Zone ID:", id)
		fmt.Println("\t[+] Zone Name:", name)

		fmt.Println("Fetching Dns Records...")

		records, err := sendrequest.ListDnsRecords(id)
		if err != nil {
			fmt.Println("Error", err)
			return
		}

		for i := 0; i < len(records.Result); i++ {
			fmt.Println("[+] Record ", i+1)
			fmt.Println("\t[+] ID:", records.Result[i].ID)
			fmt.Println("\t[+] Name:", records.Result[i].Name)
			fmt.Println("\t[+] Type:", records.Result[i].Type)
			fmt.Println("\t[+] Content:", records.Result[i].Content)
		}

		fmt.Println("Fetching public IP....")
		public_ip, err := sendrequest.GetPublicIp()
		if err != nil {
			return
		}
		fmt.Println("\t [+] Public IP: ", public_ip)

		if public_ip != records.Result[0].Content {
			fmt.Println("[!] IP mismatch. Updating Cloudflare...")
			record := records.Result[0]
			new_records, err := sendrequest.OverwritteDnsrecords(id, record.ID, record.Name, public_ip)
			if err != nil {
				fmt.Println("Error", err)
				return
			}
			fmt.Println("[+] Record ")
			fmt.Println("\t[+] ID:", new_records.Result.ID)
			fmt.Println("\t[+] Name:", new_records.Result.Name)
			fmt.Println("\t[+] Type:", new_records.Result.Type)
			fmt.Println("\t[+] Content:", new_records.Result.Content)
		} else {
			fmt.Println("[+] IPs match. Sleeping for 1 minute...")
			if test {
				fmt.Println("Test activated, changing records on purpose.")
				// // Test Run
				record := records.Result[0]
				new_records, err := sendrequest.OverwritteDnsrecords(id, record.ID, record.Name, "0.0.0.0")
				test = false
				if err != nil {
					fmt.Println("Error", err)
					return
				}

				fmt.Println("[+] Record ")
				fmt.Println("\t[+] ID:", new_records.Result.ID)
				fmt.Println("\t[+] Name:", new_records.Result.Name)
				fmt.Println("\t[+] Type:", new_records.Result.Type)
				fmt.Println("\t[+] Content:", new_records.Result.Content)
			}
			time.Sleep(60 * time.Second)
		}
	}
}

func main() {

	mainLoop()
}
