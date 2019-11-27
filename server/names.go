package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

// NameService - Struct holding API to a service managing names
type NameService struct {
	EtcdClient *clientv3.Client
}

// NewNameService - Creates a new instance of a name service
func NewNameService(endpoints []string) *NameService {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := client.Get(ctx, "ics.version")
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Debugf("ics.version: %s", resp.Kvs[0].Value)
	return &NameService{
		EtcdClient: client,
	}
}

// ReserveNickname - Reserves a nickname to be used
func (n *NameService) ReserveNickname(nickname string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	getResp, err := n.EtcdClient.Get(ctx, fmt.Sprintf("users.nicknames.%s", nickname))
	if err != nil {
		return err
	}
	if len(getResp.Kvs) > 0 {
		return fmt.Errorf("%s already in use", nickname)
	}
	if err := ValidateNickname(nickname); err != nil {
		return err
	}
	if _, err := n.EtcdClient.Put(ctx, fmt.Sprintf("users.nicknames.%s", nickname), ""); err != nil {
		return err
	}
	return nil
}

// SaveUser - Saves user's authentication information
func (n *NameService) SaveUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := n.EtcdClient.Put(ctx, fmt.Sprintf("users.%s.nickname", user.TokenStr), user.Nickname); err != nil {
		return err
	}
	if _, err := n.EtcdClient.Put(ctx, fmt.Sprintf("users.%s.addr", user.TokenStr), user.Addr); err != nil {
		return err
	}
	return nil
}

// GetUser - Gets authenticated user information
func (n *NameService) GetUser(tokenStr string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	respNickname, err := n.EtcdClient.Get(ctx, fmt.Sprintf("users.%s.nickname", tokenStr), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	respAddr, err := n.EtcdClient.Get(ctx, fmt.Sprintf("users.%s.addr", tokenStr), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	if len(respNickname.Kvs) == 0 || len(respAddr.Kvs) == 0 {
		return nil, fmt.Errorf("user not authenticated")
	}
	nickname := string(respNickname.Kvs[0].Value)
	addr := string(respAddr.Kvs[0].Value)
	return NewUser(nickname, addr, tokenStr), nil
}

// ChangeUserNickname - Updates user information
func (n *NameService) ChangeUserNickname(tokenStr, nickname string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := n.EtcdClient.Put(ctx, fmt.Sprintf("users.%s.nickname", tokenStr), nickname); err != nil {
		return err
	}
	return nil
}

// ReserveRoomName - Reserves a room name to be used
func (n *NameService) ReserveRoomName(roomID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	getResp, err := n.EtcdClient.Get(ctx, fmt.Sprintf("rooms.names.%s", roomID))
	if err != nil {
		return err
	}
	if len(getResp.Kvs) > 0 {
		return fmt.Errorf("room %s already exists", roomID)
	}
	if err := ValidateRoomName(roomID); err != nil {
		return err
	}
	if _, err := n.EtcdClient.Put(ctx, fmt.Sprintf("rooms.names.%s", roomID), ""); err != nil {
		return err
	}
	return nil
}

// GetRoom - Gets a functioning room
func (n *NameService) GetRoom(roomID string) (*Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	getResp, err := n.EtcdClient.Get(ctx, fmt.Sprintf("rooms.names.%s", roomID), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	if len(getResp.Kvs) == 0 {
		return nil, fmt.Errorf("room %s does not exist", roomID)
	}
	return NewRoom(roomID), nil
}

// GetUsersList - Gets users associated with a room
func (n *NameService) GetUsersList(roomID string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	prefix := fmt.Sprintf("rooms.%s.users.", roomID)
	getResp, err := n.EtcdClient.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		return nil, err
	}
	nicknames := []string{}
	for _, kv := range getResp.Kvs {
		nickname := strings.TrimPrefix(string(kv.Key), prefix)
		nicknames = append(nicknames, nickname)
	}
	return nicknames, nil
}

// GetRoomsList - Gets all room names
func (n *NameService) GetRoomsList() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	prefix := "rooms.names."
	getResp, err := n.EtcdClient.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		return nil, err
	}
	rooms := []string{}
	for _, kv := range getResp.Kvs {
		room := strings.TrimPrefix(string(kv.Key), prefix)
		rooms = append(rooms, room)
	}
	return rooms, nil
}
