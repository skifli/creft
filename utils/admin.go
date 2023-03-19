package utils

import "math/big"

func HasAdminPerms(permissionsString *string) bool {
	permissionsBigInt, _ := new(big.Int).SetString(*permissionsString, 10)
	permissions := permissionsBigInt.Uint64()

	return (permissions & 0x8) == 0x8 // Do some bit magic to check if the user has admin perms --> https://discord.com/developers/docs/topics/permissions
}
