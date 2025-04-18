// Code generated by protoc-gen-go-aip. DO NOT EDIT.
//
// versions:
// 	protoc-gen-go-aip development
// 	protoc (unknown)
// source: gomicroservice/v1/user_service.proto

package gomicroservicev1

import (
	fmt "fmt"
	resourcename "go.einride.tech/aip/resourcename"
	strings "strings"
)

type UserResourceName struct {
	User string
}

func (n UserResourceName) Validate() error {
	if n.User == "" {
		return fmt.Errorf("user: empty")
	}
	if strings.IndexByte(n.User, '/') != -1 {
		return fmt.Errorf("user: contains illegal character '/'")
	}
	return nil
}

func (n UserResourceName) ContainsWildcard() bool {
	return false || n.User == "-"
}

func (n UserResourceName) String() string {
	return resourcename.Sprint(
		"users/{user}",
		n.User,
	)
}

func (n UserResourceName) MarshalString() (string, error) {
	if err := n.Validate(); err != nil {
		return "", err
	}
	return n.String(), nil
}

func (n *UserResourceName) UnmarshalString(name string) error {
	err := resourcename.Sscan(
		name,
		"users/{user}",
		&n.User,
	)
	if err != nil {
		return err
	}
	return n.Validate()
}

func (n UserResourceName) Type() string {
	return "gomicroservice/User"
}
