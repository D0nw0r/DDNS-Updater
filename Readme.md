# DDNS-updater

This project is to facilitate a VPN Domain maintenance by updating CloudFlare's DNS records periodically.

# Requirements
- Setup an .env file under the "readenv" folder with the following format:

```
EMAIL="my_email@email.com"
KEY=<your_key>
```

# How it works

DDNS-Updater fetches cloudflare with your API Key and Email on 3 endpoints:

- https://api.cloudflare.com/client/v4/zones 
- https://api.cloudflare.com/client/v4/zones/ZONE_ID/dns_records
- https://api.cloudflare.com/client/v4/zones/ZONE_ID/dns_records/RECORD_ID

The first 2 are to obtain the current IP address pointed by the record. The last one is to update it.

The program runs a loop of fetching the Public IP on ipinfo.io and comparing it to the stored records of cloudflare. 
When there is a mismatch, a request for change is made.


# TODO
- Support for the user to select which record (default just uses record[0])
- Db integration?
- better code organization and naming