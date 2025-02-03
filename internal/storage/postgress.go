package postgress

import (
	"context"
	models "cyber/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type PgxPool interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	// Добавьте другие методы, которые вы используете из pgxpool.Pool
}

type Storage struct {
	Db PgxPool
}

var (
	ErrNotValidAreaID    = errors.New("Invalid area ID")
	ErrNotValidUserID    = errors.New("Invalid user ID")
	ErrNotValidCoord     = errors.New("Failed to marshal coordinates")
	ErrDataBase          = errors.New("database error")
	ErrRows              = errors.New("error after iterating rows")
	ErrInvalidEmail      = errors.New("invalid email address")
	ErrNotValidChar      = errors.New("Failed to marshal characteristics")
	ErrNotValidRes       = errors.New("Failed to marshal resources")
	ErrNotValidAbilities = errors.New("Failed to marshal abilities")
)

// Storage конструктор. Пароль БД загружается из переменной окружения.
func New() (*Storage, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("cant loading .env file")
		return nil, err
	}
	pwd := os.Getenv("DBPASSWORD")

	connString := "postgres://postgres:" + pwd + "@localhost:5432/cyberball"

	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Printf("cant create new instance of DB: %v\n", err)
		return nil, err
	}
	s := Storage{
		Db: db,
	}
	return &s, nil
}

/* Дальнейшую часть кода, касающуюся работы с юзерами пока не удаляю, но ее можно не смотреть*/

// функция проверки валидности email адресса
func IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Метод добавления пользователя в БД. На вход принимает слайс объектов, возвращает ошибку, при наличии.
// TODO: добавить проверок валидности данных пользователя
func (s *Storage) AddUser(u models.User) error {
	if IsEmailValid(u.Email) {
		_, err := s.Db.Exec(context.Background(), `INSERT INTO users (login, password, email, subscription, league_id, balance, level) VALUES ($1, $2, $3, $4, $5, $6, $7);`,
			u.Login, u.Password, u.Email, u.Subscription, u.LeagueId, u.Balance, u.Level)
		if err != nil {
			//log.Fatalf("Cant add data in database! %v\n", err)
			return err
		}
	} else {
		return ErrInvalidEmail
	}

	return nil
}

func (s *Storage) GetUser(userId int64) (models.User, error) {
	if userId < 1 {

		log.Printf("Error!Invalid user id- %v", userId)
		return models.User{}, ErrNotValidUserID
	}

	rows, err := s.Db.Query(context.Background(), `SELECT id, login, password, email, subscription, league_id, balance, level FROM users WHERE id=$1`, userId)
	if err != nil {
		log.Printf("Cant read data from database: %v\n", err)
		return models.User{}, ErrDataBase
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
			&u.Balance,
			&u.Level,
		)
		if err != nil {
			log.Printf("unable scan row: %v", err)
			return models.User{}, ErrRows
		}
	}
	return u, nil
}

// GetNeutrals получает все нейтральные объекты по ID арены
func (s *Storage) GetNeutrals(areaID int64) ([]models.Neutral, error) {
	// Проверка на валидность areaID
	if areaID < 1 {
		return nil, ErrNotValidAreaID
	}
	rows, err := s.Db.Query(context.Background(), `SELECT neutrals.* FROM areas_neutrals JOIN neutrals ON areas_neutrals.neutral_id = neutrals.id WHERE areas_neutrals.area_id = $1;`, areaID)
	if err != nil {
		log.Printf("Failed to execute query GetNeutrals: %v\n", err)
		return nil, ErrDataBase
	}
	defer rows.Close()

	var neutrals []models.Neutral

	for rows.Next() {
		var n models.Neutral
		err := rows.Scan(
			&n.Id,
			&n.Name,
			&n.Product,
			&n.ProductivityCoefficient,
			&n.Capacity,
			&n.ThresholdLevel1,
			&n.ThresholdLevel2,
			&n.Coordinates,
		)
		if err != nil {
			log.Printf("unable scan row: %v", err)
			return nil, ErrRows
		}
		neutrals = append(neutrals, n)
	}
	return neutrals, nil
}

// GetBuildings получает все здания по ID арены
func (s *Storage) GetBuildings(areaID int64) ([]models.Building, error) {
	// Проверка на валидность areaID
	if areaID < 1 {
		log.Printf("Invalid ara id - %v", areaID)
		return nil, ErrNotValidAreaID
	}
	rows, err := s.Db.Query(context.Background(), `SELECT buildings.* FROM areas_buildings
	JOIN buildings ON areas_buildings.buildings_id=buildings.id WHERE areas_buildings.area_id=$1`, areaID)
	if err != nil {
		log.Printf("Cant read data about buildings object from DB: %v\n", err)
		return nil, ErrDataBase
	}
	defer rows.Close()

	var buildings []models.Building

	for rows.Next() {
		var b models.Building
		err := rows.Scan(
			&b.Id,
			&b.Name,
			&b.Product,
			&b.Charachteristics,
			&b.Level,
			&b.UpgradePrice,
		)
		if err != nil {
			log.Printf("unable scan row: %v", err)
			return nil, ErrRows
		}
		buildings = append(buildings, b)
	}
	return buildings, nil
}

// GetHeroes получает всех героев по ID арены
func (s *Storage) GetHeroes(areaID int64) ([]models.Hero, error) {
	// Проверка на валидность areaID
	if areaID < 1 {
		log.Printf("Invalid ara id - %v", areaID)
		return nil, ErrNotValidAreaID
	}
	rows, err := s.Db.Query(context.Background(), `SELECT heroes.* FROM areas_heroes
	JOIN heroes ON areas_buildings.heroes_id=heroes.id WHERE areas_buildings.area_id=$1`, areaID)
	if err != nil {
		log.Printf("Cant read data about heroes object from DB: %v\n", err)
		return nil, ErrDataBase
	}
	defer rows.Close()

	var heroes []models.Hero

	for rows.Next() {
		var h models.Hero
		err := rows.Scan(
			&h.Id,
			&h.Name,
			&h.Charachteristics,
			&h.ExperienceToUp,
			&h.Level,
			&h.Abilities,
		)
		if err != nil {
			log.Printf("unable scan row: %v", err)
			return nil, ErrRows
		}
		heroes = append(heroes, h)
	}
	return heroes, nil
}

// GetObstacles производит выборку координат всех объектов конкретной арены и добавляет их в результирующий слайс []Hex
func (s *Storage) GetObstacles(areaID int64) ([]models.Hex, error) {
	// Проверка на валидность areaID
	if areaID < 1 {
		return nil, errors.New("invalid area ID")
	}

	// SQL-запрос для получения координат из всех таблиц содержих координаты объектов
	query := `
		SELECT coordinates FROM areas_neutrals
		JOIN neutrals ON areas_neutrals.neutral_id = neutrals.id
		WHERE areas_neutrals.area_id = $1
		UNION ALL
		SELECT coordinates FROM areas_buildings
		JOIN buildings ON areas_buildings.building_id = buildings.id
		WHERE areas_buildings.area_id = $1
		UNION ALL
		SELECT coordinates FROM areas_heroes
		JOIN heroes ON areas_heroes.hero_id = heroes.id
		WHERE areas_heroes.area_id = $1
		UNION ALL
		SELECT coordinates FROM areas_units
		JOIN units ON areas_units.unit_id = units.id
		WHERE areas_units.area_id = $1
		UNION ALL
		SELECT coordinates FROM areas_enemies
		JOIN enemies ON areas_enemies.enemy_id = enemies.id
		WHERE areas_enemies.area_id = $1;
	`

	// Выполняем запрос
	rows, err := s.Db.Query(context.Background(), query, areaID)
	if err != nil {
		log.Printf("Failed to execute query: %v\n", err)
		return nil, err //fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Собираем координаты в срез
	var obstacles []models.Hex
	for rows.Next() {
		var coordJSON []byte
		if err := rows.Scan(&coordJSON); err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}

		// Парсим  координаты  из JSON  в Hex
		var coords models.Hex
		if err := json.Unmarshal(coordJSON, &coords); err != nil {
			return nil, fmt.Errorf("unable to unmarshal coordinates: %w", err)
		}

		// Добавляем все координаты в общий слайс
		obstacles = append(obstacles, coords)
	}

	// Проверяем ошибки после итерации
	if err := rows.Err(); err != nil {
		log.Printf("Error after iterating rows: %v\n", err)
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}
	return obstacles, nil
}

// AddEmptyArea добавляет пустую(без объектов на ней) арену в базу и возвращает ее ID
func (s *Storage) AddEmptyArea(a models.Area) (int64, error) {
	query := `INSERT INTO area
		(user_id, width,height,cell_type_id) 
		VALUES ($1,$2,$3,$4) RETURNING id;`

	var id int64
	err := s.Db.QueryRow(context.Background(), query,
		a.UserId,
		a.Width,
		a.Height,
		a.CellTypeId).Scan(&id)
	if err != nil {
		//log.Fatalf("Cant add data about neutral in database! %v\n", err)
		return 0, err
	}
	return id, nil
}

// AddNeutral добавляет нейтральный объект в базу и возвращает его ID
// TODO:избавиться от дополнительно маршалинга координат
func (s *Storage) AddNeutral(n models.Neutral) (int64, error) {
	query := `INSERT INTO neutral
             (name, product, prod_cof, capacity, threshold_level1, threshold_level2, size, coordinates)
             VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`

	var id int64
	coordsJSON, err := json.Marshal(n.Coordinates)
	if err != nil {
		log.Printf("Failed to marshal coordinates: %v\n", err)
		return 0, ErrNotValidCoord
	}

	err = s.Db.QueryRow(context.Background(), query,
		n.Name,
		n.Product,
		n.ProductivityCoefficient,
		n.Capacity,
		n.ThresholdLevel1,
		n.ThresholdLevel2,
		n.Size,
		coordsJSON).Scan(&id)

	if err != nil {
		log.Printf("Cant add data about neutral in database! %v\n", err)
		return 0, ErrDataBase
	}

	return id, nil
}

// AddBuilding добавляет здание в базу и возвращает его ID
// TODO:избавиться от дополнительно маршалинга координат и характеристик
func (s *Storage) AddBuilding(b models.Building) (int64, error) {
	query := `INSERT INTO building 
              (name, product, characteristics, level, resources, coordinates) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`

	var id int64
	coordsJSON, err := json.Marshal(b.Coordinates)
	if err != nil {
		log.Printf("Failed to marshal coordinates: %v\n", err)
		return 0, ErrNotValidCoord
	}

	charachteristicsJSON, err := json.Marshal(b.Charachteristics)
	if err != nil {
		log.Printf("Failed to marshal characteristics: %v\n", err)
		return 0, ErrNotValidChar
	}

	resourcesJSON, err := json.Marshal(b.UpgradePrice)
	if err != nil {
		log.Printf("Failed to marshal resources: %v\n", err)
		return 0, ErrNotValidRes
	}

	err = s.Db.QueryRow(context.Background(), query,
		b.Name,
		b.Product,
		charachteristicsJSON,
		b.Level,
		resourcesJSON,
		coordsJSON).Scan(&id)

	if err != nil {
		log.Printf("Cant add data about building in database! %v\n", err)
		return 0, err
	}

	return id, nil
}

// AddHero добавляет героя в базу и возвращает его ID
// TODO:избавиться от дополнительно маршалинга координат и характеристик
func (s *Storage) AddHero(h models.Hero) (int64, error) {
	query := `INSERT INTO hero 
              (name, characteristics, experience, experience_to_up, level, abilities, coordinates) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
	var id int64
	coordsJSON, err := json.Marshal(h.Coordinates)
	if err != nil {
		log.Printf("Failed to marshal coordinates: %v\n", err)
		return 0, ErrNotValidCoord
	}

	charachteristicsJSON, err := json.Marshal(h.Charachteristics)
	if err != nil {
		log.Printf("Failed to marshal characteristics: %v\n", err)
		return 0, ErrNotValidChar
	}

	abilitiesJSON, err := json.Marshal(h.Abilities)
	if err != nil {
		log.Printf("Failed to marshal abilities: %v\n", err)
		return 0, ErrNotValidAbilities
	}

	err = s.Db.QueryRow(context.Background(), query,
		h.Name,
		charachteristicsJSON,
		h.Experience,
		h.ExperienceToUp,
		h.Level,
		abilitiesJSON,
		coordsJSON).Scan(&id)

	if err != nil {
		log.Printf("Cant add data about hero in database! %v\n", err)
		return 0, err
	}

	return id, nil
}

// AddHero добавляет юнита в базу и возвращает его ID
// TODO:избавиться от дополнительно маршалинга координат и характеристик
func (s *Storage) AddUnit(u models.Unit) (int64, error) {
	query := `INSERT INTO unit 
              (name, characteristics, experience, experience_to_up, level, image_id, coordinates) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`

	var id int64
	coordsJSON, err := json.Marshal(u.Coordinates)
	if err != nil {
		log.Printf("Failed to marshal coordinates: %v\n", err)
		return 0, err
	}

	charachteristicsJSON, err := json.Marshal(u.Charachteristics)
	if err != nil {
		log.Printf("Failed to marshal characteristics: %v\n", err)
		return 0, err
	}

	err = s.Db.QueryRow(context.Background(), query,
		u.Name,
		charachteristicsJSON,
		u.Experience,
		u.ExperienceToUp,
		u.Level,
		u.ImageId,
		coordsJSON).Scan(&id)

	if err != nil {
		log.Printf("Cant add data about unit in database! %v\n", err)
		return 0, err
	}

	return id, nil
}

func (s *Storage) AddNeutralAtArea(neutralId, areaId int64) error {
	query := `INSERT INTO areas_neutrals (area_id, neutral_id) VALUES $1,$2;`

	_, err := s.Db.Exec(context.Background(), query, areaId, neutralId)
	if err != nil {
		log.Printf("Cant add link between neutral ID- %v and area ID- %v database! %v\n", neutralId, areaId, err)
		return err
	}
	return nil
}

func (s *Storage) AddBuildingAtArea(buildingId, areaId int64) error {
	query := `INSERT INTO areas_buildings (area_id, building_id) VALUES $1,$2;`

	_, err := s.Db.Exec(context.Background(), query, areaId, buildingId)
	if err != nil {
		log.Printf("Cant add link between building ID- %v and area ID- %v database! %v\n", buildingId, areaId, err)
		return err
	}
	return nil
}

func (s *Storage) AddHeroAtArea(heroId, areaId int64) error {
	query := `INSERT INTO areas_heroes (area_id, hero_id) VALUES $1,$2;`

	_, err := s.Db.Exec(context.Background(), query, areaId, heroId)
	if err != nil {
		log.Printf("Cant add link between hero ID- %v and area ID- %v database! %v\n", heroId, areaId, err)
		return err
	}
	return nil
}

func (s *Storage) AddUnitAtArea(unitId, areaId int64) error {
	query := `INSERT INTO areas_heroes (area_id, hero_id) VALUES $1,$2;`

	_, err := s.Db.Exec(context.Background(), query, areaId, unitId)
	if err != nil {
		log.Printf("Cant add link between unit ID- %v and area ID- %v database! %v\n", unitId, areaId, err)
		return err
	}
	return nil
}
