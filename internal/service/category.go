package service

import (
	"context"

	"github.com/gabezy/go-grpc/internal/database"
	"github.com/gabezy/go-grpc/internal/pb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{
		CategoryDB: categoryDB,
	}
}

func (cs *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
	categoryPersisted, err := cs.CategoryDB.Create(in.Name, in.Description)
	if err != nil {
		return &pb.CategoryResponse{}, err
	}

	category := &pb.Category{
		Id:          categoryPersisted.ID,
		Name:        categoryPersisted.Name,
		Description: categoryPersisted.Description,
	}

	return &pb.CategoryResponse{Category: category}, nil
}
