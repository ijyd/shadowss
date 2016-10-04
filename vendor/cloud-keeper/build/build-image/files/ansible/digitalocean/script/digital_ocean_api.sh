#!/bin/sh
api_key="21c0ddbe1c1a15569bf81abe7f0ad1b9e316bd8b53fdf8b679a0b84395a72a22"

get_regions()
{
	curl -X GET -H "Content-Type: application/json" \
		-H "Authorization: Bearer $api_key" \
		"https://api.digitalocean.com/v2/regions"	
}

get_images()
{
	curl  -X GET -H "Content-Type: application/json" \
	-H "Authorization: Bearer $api_key" \
	"https://api.digitalocean.com/v2/images?page=1&per_page=66" -o imagelist.json
}

get_ssh_keys()
{
	curl -X GET -H "Content-Type: application/json" \
	-H "Authorization: Bearer $api_key" \
	"https://api.digitalocean.com/v2/account/keys"
}

list_droplet()
{
	curl -X GET -H "Content-Type: application/json" \
		-H "Authorization: Bearer $api_key" \
		"https://api.digitalocean.com/v2/droplets?page=1&per_page=1"
}
get_regions
get_images
# get_ssh_keys
# list_droplet