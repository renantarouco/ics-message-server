package server

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

// NameService - Struct holding API to a service managing names
type NameService struct {
	Endpoints []string
}

// NewNameService - Creates a new instance of a name service
func NewNameService(endpoints []string) *NameService {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := cli.Get(ctx, "ics.version")
	cancel()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Debugf("ics.version: %s", resp.Kvs[0].Value)
	return &NameService{
		Endpoints: endpoints,
	}
}

// CheckNickname - Checks if a nickname can be used
func (n *NameService) CheckNickname(nickname string) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   n.Endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := cli.Get(ctx, fmt.Sprintf("users.%s", nickname), clientv3.WithPrefix())
	cancel()
	if err != nil {
		return err
	}
	if len(resp.Kvs) > 0 {
		return fmt.Errorf("%s already in use", nickname)
	}
	if err := ValidateNickname(nickname); err != nil {
		return err
	}
	return nil
}

// SaveUser - Saves user's authentication information
func (n *NameService) SaveUser(user User) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   n.Endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = cli.Put(ctx, fmt.Sprintf("users.%s.addr", user.Nickname), user.Addr)
	cancel()
	if err != nil {
		return err
	}
	_, err = cli.Put(ctx, fmt.Sprintf("users.%s.token", user.Nickname), user.TokenStr)
	cancel()
	if err != nil {
		return err
	}
	return nil
}
