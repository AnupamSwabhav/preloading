package service

import (
	"fmt"
	"preloading/test/address"
	"preloading/test/repository"
	"preloading/test/student"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type StudentService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

func NewStudentService(db *gorm.DB, repository repository.Repository) *StudentService {
	return &StudentService{
		DB:         db,
		Repository: repository,
	}
}

func (s *StudentService) Add(student *student.Student) error {
	uow := repository.NewUnitOfWork(s.DB, false)
	err := s.Repository.Add(uow, &student)
	if err != nil {
		uow.Complete()
		return err
	}

	uow.Commit()
	return nil
}

func (s *StudentService) GetAll(students *[]*student.Student) error {
	uow := repository.NewUnitOfWork(s.DB, true)
	var data []string = []string{"Address"}
	err := s.Repository.GetAll(uow, &students, data)
	if err != nil {
		return err
	}
	fmt.Println(&students)
	uow.Commit()
	return nil
}

func (s *StudentService) Get(student *student.Student, id uuid.UUID) error {
	uow := repository.NewUnitOfWork(s.DB, true)
	err := s.Repository.Get(uow, student, id, []string{"Address"})
	if err != nil {
		return err
	}
	fmt.Println(student)
	uow.Commit()
	return nil
}

func (s *StudentService) Update(student *student.Student, id uuid.UUID) error {
	uow := repository.NewUnitOfWork(s.DB, false)
	fmt.Println(student.Address[0].StudentID)
	fmt.Println("Adress : ", student.Address)
	fmt.Println("Student : ", student)
	student.ID = id
	err := s.Repository.Update(uow, &student)
	//err := s.Repository.Save(uow, &student)

	if err != nil {
		uow.Complete()
		return err
	}
	address := address.Address{}
	er := s.Repository.DeleteRemaining(uow, address, student.ID, student.Address[0].ID)
	if er != nil {
		uow.Complete()
		return er
	}
	fmt.Println(student)
	uow.Commit()
	return nil
}

func (s *StudentService) Delete(student *student.Student, id uuid.UUID) error {
	uow := repository.NewUnitOfWork(s.DB, false)
	err := s.Repository.Get(uow, student, id, []string{"Address"})
	if err != nil {
		uow.Complete()
		return err
	}
	address := address.Address{}
	er := s.Repository.DeleteByID(uow, &address, student.ID)
	if er != nil {
		uow.Complete()
		return er
	}
	err1 := s.Repository.Delete(uow, &student)
	if err1 != nil {
		uow.Complete()
		return err1
	}
	uow.Commit()
	return nil
}
