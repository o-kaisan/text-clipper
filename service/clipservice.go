package service

import "github.com/o-kaisan/text-clipper/repository"

type ClipService struct {
	cs repository.ClipRepository
}

func NewClipService(cs repository.ClipRepository) ClipService {
	return ClipService{cs}
}
