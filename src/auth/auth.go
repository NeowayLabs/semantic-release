package auth

import (
	"log"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

// GetPublicKey aims to read ssh key file and returns a go git PublicKey
// Args:
// 		sshKey ([]byte): Byte arry containing the ssh key.
// 		Note: Set up this key in the git settings and pass as a parameter to the CLI tool.
// Returns:
// 		*ssh.PublicKeys: Returns a git public key.
func GetPublicKey(sshKey []byte) (*ssh.PublicKeys, error) {
	var publicKey *ssh.PublicKeys

	publicKey, err := ssh.NewPublicKeys("git", []byte(sshKey), "")
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	return publicKey, err
}

// FormatSSHKey aims to properlly format the ssh private key.
// It is necessary once we cannot pass a multiline variable to go build. So it is not possible
// to do docker build --build-arg PRIVATE_KEY=$(PRIVATE_KEY) once this kind of file has several line breaks.
// Args:
// 		key (string): String containing sshKey with a line break holder character.
// 		lineBreaker (string): Line break holder character. I.e.: #
func FormatSSHKey(sshKey string, lineBreaker string) string {
	return strings.Replace(sshKey, lineBreaker, "\n", -1)
}
