package ansible

const (
	varFile           = "var.yml"
	hostFile          = "host"
	keyFile           = "key.pem"
	deploySSPlayBook  = "deploy_ss.yml"
	deployVPSPlayBook = "create-node.yml"
	deleteVPSPlayBook = "delete-node.yml"
	restartSSPlayBook = deploySSPlayBook
	attrFile          = "attr.json"

	ansibleDGOCDir               = "ansible/digitalocean/deploy"
	ansibleDGOCVarFile           = ansibleDGOCDir + "/" + varFile
	ansibleDGOCHostFile          = ansibleDGOCDir + "/" + hostFile
	ansibleDGOCSSHKeyFile        = ansibleDGOCDir + "/" + keyFile
	ansibleDGOCDeploySSPlayBook  = ansibleDGOCDir + "/" + deploySSPlayBook
	ansibleDGOCDeployVPSPlayBook = ansibleDGOCDir + "/" + deployVPSPlayBook
	ansibleDGOCDeleteVPSPlayBook = ansibleDGOCDir + "/" + deleteVPSPlayBook
	ansibleDGOCAttrFile          = ansibleDGOCDir + "/" + attrFile
	ansibleDGOCRestartSSPlayBook = ansibleDGOCDir + "/" + restartSSPlayBook
	ansibleDGOCPrivateKey        = ansibleDGOCDir + "/" + keyFile

	ansibleVulDir               = "ansible/vultr/deploy"
	ansibleVulVarFile           = ansibleVulDir + "/" + varFile
	ansibleVulHostFile          = ansibleVulDir + "/" + hostFile
	ansibleVulSSHKeyFile        = ansibleVulDir + "/" + keyFile
	ansibleVulDeploySSPlayBook  = ansibleVulDir + "/" + deploySSPlayBook
	ansibleVulDeployVPSPlayBook = ansibleVulDir + "/" + deployVPSPlayBook
	ansibleVulDeleteVPSPlayBook = ansibleVulDir + "/" + deleteVPSPlayBook
	ansibleVulRestartSSPlayBook = ansibleVulDir + "/" + restartSSPlayBook
	ansibleVulAttrFile          = ansibleVulDir + "/" + attrFile
	ansibleVulPrivateKey        = ansibleVulDir + "/" + keyFile

	ansibleUpgrade           = "ansible/update"
	ansibleUpgradeSSPlayBook = ansibleUpgrade + "/deploy_update_shadowss.yml"
	ansibleUpgradeSSHostFile = ansibleUpgrade + "/host"
	ansibleUpgradeSSHKeyFile = ansibleUpgrade + "/" + keyFile
	ansibleUpgradePrivateKey = ansibleUpgrade + "/" + keyFile
)

var privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAv/QKwWAdhxbaIDcxQA2mlKO3Xx1YKPGfe30eJvenhtLiTcwz
wY63YAzL3w1rmLdIGa50m/bncAMKms+faSc27Tey1woZGxvdmZEZgHkRiKTk9Z6+
wPvYN3X5U2n79odN8RZZVc5nn8tZJI2QkfauOqxWFwXoNamHbrSkvonRn4FxDS5J
Yvh42OY/Cowak/qGj5azu4gmS5YqYQAMmB5ZPqcz2SplS/b8gbVHy+7sNBN0fwkN
t29zxNX781xgnzLf4d/qvZr/Hcyod+V0B2z1FkhpKJa8KseYoco0Y4foh7oK/axk
Y1lL3pILOzv3cKCL3SFF2eFG/T9oPc9jFlodfwIDAQABAoIBAHvsS71IFggOosfF
mhAmP/MaNto7EZ1tUG7i+cJihE8welWLjaZaQtzJphzchyhSu0OJM1M1dXkFHaWQ
gPPcE0PWf6kApfCwbsIjwPkGMGGtQvunfrMMZCx6B3roo3gnJhSNPyN8W734BBbr
Jfh170mF1RaMA7wRNJQuH2W7iA+W3W9IA9p0ebE+OfYaS8cyzw6JylOU3kKGzkG5
raSxOQeSMmjdoAUQ9D8aVlSpm/oPXGPYwahGO7ZLrc2kKYv0QQX6mYyyxCk1290F
uCiXY76ac9CjJJAotJBcWwsrzHTK3+FQFgyl9DEWCyOqYLMjVkzSG7cJ5lsYNJ2c
D1FJh7ECgYEA31zrhzEfJbJLl0/Ym55fzVL2uactxjM4oqK5OnLtwgbTVaC8vqTf
yZXv/f7p7KU3riuYRQuyLvLT2wv45gTo/CLdLuqmS4RYplEHqX9OzD3LGIALnj76
OIcdWYAtAN1v9s7ctbyz9ss1Vwa9XONV+hxOs7MORDEMW6Er6QSsH/MCgYEA3AA3
abVNQ95lofqwxh4EyedwkkeR66+J5oQt66b7SWo4tQWK8UnYZeceyrh1w91yilIz
BcaFFU7Y+M4DPkoxOLgUOQPIOX+2wHxQrzrIM/d3n//DRWY41E4fr8AWVHaPm6d4
3DtIyij/YODUftr0z9TJd7Qk17Ylz/ZEjLulu0UCgYEAv6A3TI+++hdBtLnCypeP
91Yi59necnkFMLpMETICeoBilMbGxwQqHgbtk0o8JFMGNv2dsDa9knuvd/CIg8ZY
n9/FRHf5XTZY268O1MKstpqZABbyYLwE7bQ1YNCPS3uuj96fCaev+Z4Sz+uvT96V
p3Lbrl2CcsxlnsLiKhJhHTMCgYA+n3UmiuweeIzXicON8XeNfWrGyMaZnxMS4ecs
YBDBehIAPT6qpkmJ4Dscm1syULPM+c76QuMZCKOsVwAHWBkguw1OmWwCKf98VSam
aoYYfMW5bpVICOv+SuqsHXJ9wm3occhucBWtLfRbwEPchDkRe9GJWGbwXDHxO3mR
0cxAPQKBgQDAxigzoPKvHzfSwoDaDQ3qnhceCty0Z+KhN/8nf2Oo/YJGK6PKmaWA
fQv8BvIsPo4ZvO9nt/eu8fJDuE+Q2tuu7Zy4qLnKmcTWo3tRIPpVG5pOD+kAIlCa
I3KzXZ9RctOr9kPJczp+0kRJpI8BM9910iReRDuKPuHcJXfL1bEQ4A==
-----END RSA PRIVATE KEY-----`
