package practical_go

import (
	pb "charliemcelfresh/practical_go/rpc/practical_go"
	"context"
	"github.com/twitchtv/twirp"
)

type Repo interface {
	CreateItem(ctx context.Context, name string) error
	GetItem(ctx context.Context, itemID string) (Item, error)
}

// Server implements the Haberdasher service
type provider struct {
	Repository Repo
}

func NewProvider() provider {
	config := NewConfig()
	return provider{
		NewRepository(config.GetDB()),
	}
}

func (p provider) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.None, error) {
	err := p.Repository.CreateItem(ctx, req.Name)
	if err != nil {
		return &pb.None{}, twirp.NewError(twirp.FailedPrecondition, "cannot create item")
	}
	return &pb.None{}, nil
}

func (p provider) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.Item, error) {
	itemToReturn := &pb.Item{}
	item, err := p.Repository.GetItem(ctx, req.ItemId)
	if err != nil {
		return &pb.Item{}, twirp.NewError(twirp.NotFound, "item not found")
	}
	itemToReturn = &pb.Item{
		ItemId:    item.ItemID,
		Name:      item.Name,
		CreatedAt: item.CreatedAt,
	}
	return itemToReturn, nil
}
