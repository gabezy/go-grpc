package service

import (
	"context"
	"io"
	"sync"

	"github.com/gabezy/go-grpc/internal/database"
	"github.com/gabezy/go-grpc/internal/pb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
	mu         sync.Mutex
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{
		CategoryDB: categoryDB,
	}
}

func (cs *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.Category, error) {
	categoryPersisted, err := cs.CategoryDB.Create(in.Name, in.Description)
	if err != nil {
		return &pb.Category{}, err
	}

	category := &pb.Category{
		Id:          categoryPersisted.ID,
		Name:        categoryPersisted.Name,
		Description: categoryPersisted.Description,
	}

	return category, nil
}

func (cs *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.CategoryList{}

	for {
		categoryReq, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(categories)
		}
		if err != nil {
			return err
		}

		// Lock the DB
		cs.mu.Lock()
		// persiste the current category stream
		categoryResult, err := cs.CategoryDB.Create(categoryReq.Name, categoryReq.Description)
		if err != nil {
			return err
		}
		cs.mu.Unlock()

		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryReq.Name,
			Description: categoryReq.Description,
		})
	}
}

func (cs *CategoryService) CreateCategoryStreamBidirectional(stream pb.CategoryService_CreateCategoryStreamBidirectionalServer) error {
	for {
		categoryReq, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		categoryResult, err := cs.CategoryDB.Create(categoryReq.Name, categoryReq.Description)
		if err != nil {
			return err
		}

		category := &pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryReq.Name,
			Description: categoryReq.Description,
		}

		err = stream.Send(category)
		if err != nil {
			return err
		}
	}
}

func (cs *CategoryService) ListCategory(ctx context.Context, in *pb.Blank) (*pb.CategoryList, error) {
	categories, err := cs.CategoryDB.FindAll()
	if err != nil {
		return nil, err
	}

	var categoriesResponse []*pb.Category

	for _, category := range categories {
		categoryResponse := &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		}
		categoriesResponse = append(categoriesResponse, categoryResponse)
	}

	return &pb.CategoryList{Categories: categoriesResponse}, nil
}

func (cs *CategoryService) GetCategory(ctx context.Context, in *pb.CategoryGetRequest) (*pb.Category, error) {
	categoryPersisted, err := cs.CategoryDB.FindByID(in.Id)
	if err != nil {
		return &pb.Category{}, err
	}

	category := &pb.Category{
		Id:          categoryPersisted.ID,
		Name:        categoryPersisted.Name,
		Description: categoryPersisted.Description,
	}

	return category, nil
}
