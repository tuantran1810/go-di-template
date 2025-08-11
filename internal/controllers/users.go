package controllers

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	"github.com/tuantran1810/go-di-template/internal/controllers/transformers"
	"github.com/tuantran1810/go-di-template/internal/entities"
	pb "github.com/tuantran1810/go-di-template/pkg/go_di_template/v1"
)

type UserController struct {
	pb.UnimplementedUserServiceServer
	usecase                  IUserUsecase
	userTransformer          *transformers.PbUserTransformer
	keyValuePairTransformer  *transformers.PbKeyValuePairTransformer
	userAttributeTransformer *transformers.PbUserAttributesTransformer
}

func NewUserController(usecase IUserUsecase) *UserController {
	return &UserController{
		usecase:                  usecase,
		userTransformer:          transformers.NewPbUserTransformer(),
		keyValuePairTransformer:  transformers.NewPbKeyValuePairTransformer(),
		userAttributeTransformer: transformers.NewPbUserAttributesTransformer(),
	}
}

func (c *UserController) CreateUser(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (*pb.CreateUserResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, fmt.Errorf("%w - err: %w", entities.ErrInvalid, err)
	}

	if req.User == nil {
		return nil, fmt.Errorf("%w - user is nil", entities.ErrInvalid)
	}

	user, err := c.userTransformer.ToEntity(req.User)
	if err != nil {
		return nil, fmt.Errorf("%w - cannot transform user, err: %w", entities.ErrInvalid, err)
	}

	attributes, err := c.keyValuePairTransformer.ToEntityArray_P2I(req.Attributes)
	if err != nil {
		return nil, fmt.Errorf("%w - cannot transform user attributes, err: %w", entities.ErrInvalid, err)
	}

	outUser, outAttributes, err := c.usecase.CreateUser(ctx, user, attributes)
	if err != nil {
		return nil, err
	}

	pbUser, err := c.userTransformer.FromEntity(outUser)
	if err != nil {
		return nil, fmt.Errorf("%w - cannot transform to pb user, err: %w", entities.ErrInvalid, err)
	}

	pbAttributes, err := c.userAttributeTransformer.FromEntityArray_I2P(outAttributes)
	if err != nil {
		return nil, fmt.Errorf("%w - cannot transform to pb user attributes, err: %w", entities.ErrInvalid, err)
	}

	return &pb.CreateUserResponse{
		User:       pbUser,
		Attributes: pbAttributes,
	}, nil
}

func (c *UserController) GetUserByUsername(
	ctx context.Context,
	req *pb.GetUserByUsernameRequest,
) (*pb.GetUserByUsernameResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, fmt.Errorf("%w - err: %w", entities.ErrInvalid, err)
	}

	user, atts, err := c.usecase.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("%w - cannot get user by username, err: %w", entities.ErrInvalid, err)
	}

	pbUser, err := c.userTransformer.FromEntity(user)
	if err != nil {
		return nil, fmt.Errorf("%w - cannot transform to pb user, err: %w", entities.ErrInvalid, err)
	}

	pbAttributes, err := c.userAttributeTransformer.FromEntityArray_I2P(atts)
	if err != nil {
		return nil, fmt.Errorf("%w - cannot transform to pb user attributes, err: %w", entities.ErrInvalid, err)
	}

	return &pb.GetUserByUsernameResponse{
		User:       pbUser,
		Attributes: pbAttributes,
	}, nil
}
