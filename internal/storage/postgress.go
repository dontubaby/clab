package postgress

import (
	"context"
	"cyberball/internal/models"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

type Storage struct {
	Db *pgxpool.Pool
}

// Storage конструктор. Пароль БД загружается из переменной окружения.
func New() (*Storage, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("cant loading .env file")
		return nil, err
	}
	pwd := os.Getenv("DBPASSWORD")

	connString := "postgres://postgres:" + pwd + "@localhost:5432/cyberball"

	db, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Println("cant create new instance of DB: %v\n", err)
		return nil, err
	}
	s := Storage{
		Db: db,
	}
	return &s, nil
}

// функция проверки валидности email адресса
func IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Метод добавления пользователя в БД. На вход принимает слайс объектов, возвращает ошибку, при наличии.
// TODO: добавить проверок валидности данных пользователя
func (s *Storage) AddUser(u models.User) error {
	if IsEmailValid(u.Email) {
		_, err := s.Db.Exec(context.Background(), `INSERT INTO users 
			(login,password,email,subscription,league_id,resources,balance,level) VALUES ($1,$2,$3,$4,$5,$6,$7,$8);`,
			u.Login, u.Password, u.Email, u.Subscription, u.LeagueId, u.Resources, u.Balance, u.Level)
		if err != nil {
			log.Fatalf("Cant add data in database! %v\n", err)
			return err
		}
	} else {
		return errors.New("invalid email adress!")
	}

	return nil
}

func (s *Storage) GetUser(userId int) (models.User, error) {
	if userId < 1 {
		err := fmt.Errorf("Error!Invalid user id- %v", userId)
		log.Println(err)
		return models.User{}, err
	}

	q := strconv.Itoa(userId)

	rows, err := s.Db.Query(context.Background(), `SELECT * FROM users WHERE id=$1`, q)
	if err != nil {
		log.Printf("Cant read data from database: %v\n", err)
		return models.User{}, err
	}
	defer rows.Close()
	var u models.User
	for rows.Next() {
		err = rows.Scan(
			&u.Id,
			&u.Login,
			&u.Password,
			&u.Email,
			&u.Subscription,
			&u.LeagueId,
			&u.Resources,
			&u.Balance,
			&u.Level,
		)
		if err != nil {
			return models.User{}, fmt.Errorf("unable scan row: %w", err)
		}
	}

	return u, nil
}
