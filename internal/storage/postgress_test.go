package postgress

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/shopspring/decimal"
	"testing"
	//"time"

	"cyber/internal/models"

	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetObstacles(t *testing.T) {
	// Создаём мок базы данных
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer mock.Close()

	// Создаём Storage с моковой базой данных
	storage := &Storage{Db: mock}

	// Тестовые данные
	areaID := int64(1)
	expectedObstacles := []models.Hex{
		{Q: 1.0, R: 2.0},
		{Q: 3.0, R: 4.0},
	}

	// Преобразуем тестовые данные в JSON
	coordJSON1, _ := json.Marshal(expectedObstacles[0])
	coordJSON2, _ := json.Marshal(expectedObstacles[1])

	// Ожидаемый SQL-запрос
	query := `
		SELECT coordinates FROM areas_neutrals
		JOIN neutrals ON areas_neutrals.neutral_id = neutrals.id
		WHERE areas_neutrals.area_id = \$1
		UNION ALL
		SELECT coordinates FROM areas_buildings
		JOIN buildings ON areas_buildings.building_id = buildings.id
		WHERE areas_buildings.area_id = \$1
		UNION ALL
		SELECT coordinates FROM areas_heroes
		JOIN heroes ON areas_heroes.hero_id = heroes.id
		WHERE areas_heroes.area_id = \$1
		UNION ALL
		SELECT coordinates FROM areas_units
		JOIN units ON areas_units.unit_id = units.id
		WHERE areas_units.area_id = \$1
		UNION ALL
		SELECT coordinates FROM areas_enemies
		JOIN enemies ON areas_enemies.enemy_id = enemies.id
		WHERE areas_enemies.area_id = \$1;
	`

	// Таблица тестовых случаев
	tests := []struct {
		name           string
		areaID         int64
		mock           func()
		expectedResult []models.Hex
		expectedError  error
	}{
		{
			name:   "Valid coordinates",
			areaID: areaID,
			mock: func() {
				rows := mock.NewRows([]string{"coordinates"}).
					AddRow(coordJSON1).
					AddRow(coordJSON2)
				mock.ExpectQuery(query).WithArgs(areaID).WillReturnRows(rows)
			},
			expectedResult: expectedObstacles,
			expectedError:  nil,
		},
		{
			name:   "InvalidAreaID",
			areaID: 0,
			mock: func() {
				// Ничего не мокаем, так как функция должна вернуть ошибку до выполнения запроса
			},
			expectedResult: nil,
			expectedError:  ErrNotValidAreaID,
		},
		{
			name:   "DBError",
			areaID: areaID,
			mock: func() {
				mock.ExpectQuery(query).WithArgs(areaID).WillReturnError(errors.New("database error"))
			},
			expectedResult: nil,
			expectedError:  ErrDataBase,
		},
		{
			name:   "ScanError",
			areaID: areaID,
			mock: func() {
				rows := mock.NewRows([]string{"coordinates"}).
					AddRow([]byte("invalid json")) // Некорректный JSON
				mock.ExpectQuery(query).WithArgs(areaID).WillReturnRows(rows)
			},
			expectedResult: nil,
			expectedError:  ErrNotValidCoord,
		},
		{
			name:   "RowsError",
			areaID: areaID,
			mock: func() {
				rows := mock.NewRows([]string{"coordinates"}).
					AddRow(string(coordJSON1)).
					RowError(0, errors.New("row error")) // Ошибка при итерации
				mock.ExpectQuery(query).WithArgs(areaID).WillReturnRows(rows)
			},
			expectedResult: nil,
			expectedError:  ErrRows,
		},
	}

	// Проходим по всем тестовым случаям
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			// Вызов тестируемой функции
			obstacles, err := storage.GetObstacles(tt.areaID)

			// Проверяем ошибки
			if tt.expectedError != nil {
				assert.Error(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err, "GetObstacles should not return an error")
			}

			// Проверяем результат
			assert.Equal(t, tt.expectedResult, obstacles, "Obstacles should match expected data")

			// Проверяем, что все ожидания выполнены
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

// Табличный тест для функции AddUser
func TestAddUser(t *testing.T) {
	tests := []struct {
		name          string
		user          models.User
		mockSetup     func(mock pgxmock.PgxPoolIface)
		expectedError error
	}{
		{
			name: "Success - User added",
			user: models.User{
				Login:        "testuser",
				Password:     "password123",
				Email:        "test@example.com",
				Subscription: true,
				LeagueId:     10,
				Balance:      decimal.NewFromFloat(100.50),
				Level:        5,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`INSERT INTO users \(login, password, email, subscription, league_id, balance, level\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\);`).
					WithArgs("testuser", "password123", "test@example.com", true, int64(10), decimal.NewFromFloat(100.50), 5).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			expectedError: nil,
		},
		{
			name: "Error - Invalid email",
			user: models.User{
				Login:        "testuser",
				Password:     "password123",
				Email:        "invalid-email", // Некорректный email
				Subscription: true,
				LeagueId:     10,
				Balance:      decimal.NewFromFloat(100.50),
				Level:        5,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				// Ничего не настраиваем, так как запрос к базе данных не должен выполняться
			},
			expectedError: ErrInvalidEmail,
		},
		{
			name: "Error - Database query failed",
			user: models.User{
				Login:        "testuser",
				Password:     "password123",
				Email:        "test@example.com",
				Subscription: true,
				LeagueId:     10,
				Balance:      decimal.NewFromFloat(100.50),
				Level:        5,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`INSERT INTO users \(login, password, email, subscription, league_id, balance, level\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\);`).
					WithArgs("testuser", "password123", "test@example.com", true, int64(10), decimal.NewFromFloat(100.50), 5).
					WillReturnError(fmt.Errorf("database error"))
			},
			expectedError: ErrDataBase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок базы данных
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			// Настраиваем мок
			tt.mockSetup(mock)

			// Создаем Storage с моком базы данных
			storage := &Storage{Db: mock}

			// Вызываем тестируемую функцию
			err = storage.AddUser(tt.user)

			// Проверяем ошибку
			if tt.expectedError != nil {
				assert.Error(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// Проверяем, что все ожидания были выполнены
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUser(t *testing.T) {
	// Создаём мок базы данных
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer mock.Close()

	// Создаём Storage с моковой базой данных
	storage := &Storage{Db: mock}

	// Тестовые данные
	expectedUser := models.User{
		Id:           int64(1),
		Login:        "test_user",
		Password:     "password",
		Email:        "test@example.com",
		Subscription: true,
		LeagueId:     2,
		Balance:      decimal.NewFromFloat(500),
		Level:        5,
	}

	// Ожидаемый SQL-запрос
	query := `SELECT id, login, password, email, subscription, league_id, balance, level FROM users WHERE id=\$1`

	// Таблица тестовых случаев
	tests := []struct {
		name           string
		userId         int64
		mock           func()
		expectedResult models.User
		expectedError  error
	}{
		{
			name:   "Valid user",
			userId: int64(1),
			mock: func() {
				rows := mock.NewRows([]string{"id", "login", "password", "email", "subscription", "league_id", "balance", "level"}).
					AddRow(
						expectedUser.Id,
						expectedUser.Login,
						expectedUser.Password,
						expectedUser.Email,
						expectedUser.Subscription,
						expectedUser.LeagueId,
						expectedUser.Balance,
						expectedUser.Level,
					)
				mock.ExpectQuery(query).WithArgs(expectedUser.Id).WillReturnRows(rows)
			},
			expectedResult: expectedUser,
			expectedError:  nil,
		},
		{
			name:   "Invalid user ID",
			userId: 0,
			mock: func() {
				// Ничего не мокаем, так как функция должна вернуть ошибку до выполнения запроса
			},
			expectedResult: models.User{},
			expectedError:  ErrNotValidUserID,
		},
		{
			name:   "DB error",
			userId: int64(1),
			mock: func() {
				mock.ExpectQuery(query).WithArgs(expectedUser.Id).WillReturnError(errors.New("database error"))
			},
			expectedResult: models.User{},
			expectedError:  ErrDataBase,
		},
		{
			name:   "Scan error",
			userId: 1,
			mock: func() {
				rows := mock.NewRows([]string{"id", "login", "password", "email", "subscription", "league_id", "balance", "level"}).
					AddRow(1, "test_user", "password", "test@example.com", true, 2, []byte("invalid json"), 5)
				mock.ExpectQuery(query).WithArgs(expectedUser.Id).WillReturnRows(rows)
			},
			expectedResult: models.User{},
			expectedError:  ErrRows,
		},
	}

	// Проходим по всем тестовым случаям
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			// Вызов тестируемой функции
			user, err := storage.GetUser(tt.userId)

			// Проверяем ошибки
			if tt.expectedError != nil {
				assert.Error(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err, "GetUser should not return an error")
			}

			// Проверяем результат
			assert.Equal(t, tt.expectedResult, user, "User should match expected data")

			// Проверяем, что все ожидания выполнены
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestGetNeutrals(t *testing.T) {
	// Создаем мок базы данных
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer mock.Close()

	// Создаем Storage с моковой базой данных
	storage := &Storage{Db: mock}

	// Тестовые данные
	areaID := int64(1)
	expected := []models.Neutral{
		{
			Id:                      int64(1),
			Name:                    "Gold mine",
			Product:                 "Gold",
			ProductivityCoefficient: 4,
			Capacity:                decimal.NewFromFloat(10000),
			ThresholdLevel1:         decimal.NewFromFloat(5000),
			ThresholdLevel2:         decimal.NewFromFloat(2000),
			Size:                    2,
			Coordinates:             []models.Hex{{Q: 1, R: 2}, {Q: 3, R: 4}},
		},
		{
			Id:                      int64(2),
			Name:                    "Wheat field",
			Product:                 "Food",
			ProductivityCoefficient: 6,
			Capacity:                decimal.NewFromFloat(1000),
			ThresholdLevel1:         decimal.NewFromFloat(500),
			ThresholdLevel2:         decimal.NewFromFloat(200),
			Size:                    2,
			Coordinates:             []models.Hex{{Q: 4, R: 5}, {Q: 6, R: 7}},
		},
	}

	// Сериализация координат в JSON
	coordJSON1, _ := json.Marshal(expected[0].Coordinates)
	coordJSON2, _ := json.Marshal(expected[1].Coordinates)

	// Ожидаемый SQL-запрос
	query := `SELECT neutrals.* FROM areas_neutrals JOIN neutrals ON areas_neutrals.neutral_id = neutrals.id WHERE areas_neutrals.area_id = \$1;`

	// Таблица тестовых случаев
	tests := []struct {
		name           string
		areaID         int64
		mock           func()
		expectedResult []models.Neutral
		expectedError  error
	}{
		{
			name:   "Valid data",
			areaID: areaID,
			mock: func() {
				rows := mock.NewRows([]string{"id", "name", "product", "productivity_coefficient", "capacity", "threshold_level1", "threshold_level2", "size", "coordinates"}).
					AddRow(expected[0].Id, expected[0].Name, expected[0].Product, expected[0].ProductivityCoefficient, expected[0].Capacity, expected[0].ThresholdLevel1, expected[0].ThresholdLevel2, expected[0].Size, string(coordJSON1)).
					AddRow(expected[1].Id, expected[1].Name, expected[1].Product, expected[1].ProductivityCoefficient, expected[1].Capacity, expected[1].ThresholdLevel1, expected[1].ThresholdLevel2, expected[1].Size, string(coordJSON2))
				mock.ExpectQuery(query).WithArgs(areaID).WillReturnRows(rows)
			},
			expectedResult: expected,
			expectedError:  nil,
		},
		{
			name:   "InvalidAreaID",
			areaID: 0,
			mock: func() {
				// Ничего не мокаем, так как функция должна вернуть ошибку до выполнения запроса
			},
			expectedResult: nil,
			expectedError:  ErrNotValidAreaID,
		},
		{
			name:   "DB error",
			areaID: areaID,
			mock: func() {
				mock.ExpectQuery(query).WithArgs(areaID).WillReturnError(errors.New("database error"))
			},
			expectedResult: nil,
			expectedError:  ErrDataBase,
		},
		{
			name:   "Scan error",
			areaID: areaID,
			mock: func() {
				rows := mock.NewRows([]string{"id", "name", "product", "productivity_coefficient", "capacity", "threshold_level1", "threshold_level2", "size", "coordinates"}).
					AddRow(1, "Gold mine", "Gold", 4, decimal.NewFromFloat(10000), decimal.NewFromFloat(5000), decimal.NewFromFloat(2000), 2, []byte("invalid json"))
				mock.ExpectQuery(query).WithArgs(areaID).WillReturnRows(rows)
			},
			expectedResult: nil,
			expectedError:  ErrRows,
		},
	}

	// Проходим по всем тестовым случаям
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			// Вызов тестируемой функции
			neutrals, err := storage.GetNeutrals(tt.areaID)

			// Проверяем ошибки
			if tt.expectedError != nil {
				assert.Error(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err, "GetNeutrals should not return an error")
			}

			// Проверяем результат
			assert.Equal(t, tt.expectedResult, neutrals, "Neutrals should match expected data")

			// Проверяем, что все ожидания выполнены
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestGetBuildings(t *testing.T) {
	// Создаем мок базы данных
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer mock.Close()

	// Создаем Storage с моковой базой данных
	storage := &Storage{Db: mock}

	// Тестовые данные
	areaID := int64(1)
	expectedBuilding1 := models.Building{
		Id:      1,
		Name:    "Barracks",
		Product: "Soldiers",
		Charachteristics: models.BuildingCharacteristics{
			HP:                      1000,
			Armor:                   50,
			ProductivityCoefficient: 4,
			Size:                    2,
		},
		Level: 3,
		UpgradePrice: []models.Resource{
			{Id: 1, Name: "Wood", Value: decimal.NewFromFloat(300)},
			{Id: 2, Name: "Stone", Value: decimal.NewFromFloat(300)},
		},
		Coordinates: []models.Hex{{Q: 1, R: 2}, {Q: 3, R: 4}},
	}
	expectedBuilding2 := models.Building{
		Id:      2,
		Name:    "Stable",
		Product: "Horses",
		Charachteristics: models.BuildingCharacteristics{
			HP:                      800,
			Armor:                   70,
			ProductivityCoefficient: 6,
			Size:                    2,
		},
		Level: 2,
		UpgradePrice: []models.Resource{
			{Id: 1, Name: "Wood", Value: decimal.NewFromFloat(300)},
			{Id: 2, Name: "Stone", Value: decimal.NewFromFloat(300)},
		},
		Coordinates: []models.Hex{{Q: 4, R: 5}, {Q: 6, R: 7}},
	}
	expectedBuildings := []models.Building{expectedBuilding1, expectedBuilding2}

	// Сериализация характеристик и цен на улучшение в JSON
	characteristicsJSON1, _ := json.Marshal(expectedBuilding1.Charachteristics)
	characteristicsJSON2, _ := json.Marshal(expectedBuilding2.Charachteristics)
	upgradePriceJSON1, _ := json.Marshal(expectedBuilding1.UpgradePrice)
	upgradePriceJSON2, _ := json.Marshal(expectedBuilding2.UpgradePrice)
	coordinatesJSON1, _ := json.Marshal(expectedBuilding1.Coordinates)
	coordinatesJSON2, _ := json.Marshal(expectedBuilding2.Coordinates)

	// Ожидаемый SQL-запрос
	query := `SELECT neutrals.* FROM areas_neutrals 
	JOIN neutrals ON areas_neutrals.neutral_id=neutrals.id WHERE areas_neutrals.area_id=\$1;`

	// Таблица тестовых случаев
	tests := []struct {
		name           string
		areaID         int64
		mock           func()
		expectedResult []models.Building
		expectedError  error
	}{
		{
			name:   "Valid data",
			areaID: areaID,
			mock: func() {
				rows := mock.NewRows([]string{"id", "name", "product", "characteristics", "level", "resources", "coordinates"}).
					AddRow(
						expectedBuilding1.Id,
						expectedBuilding1.Name,
						expectedBuilding1.Product,
						characteristicsJSON1,
						expectedBuilding1.Level,
						upgradePriceJSON1,
						coordinatesJSON1,
					).
					AddRow(
						expectedBuilding2.Id,
						expectedBuilding2.Name,
						expectedBuilding2.Product,
						characteristicsJSON2,
						expectedBuilding2.Level,
						upgradePriceJSON2,
						coordinatesJSON2,
					)
				mock.ExpectQuery(query).WithArgs(areaID).WillReturnRows(rows)
			},
			expectedResult: expectedBuildings,
			expectedError:  nil,
		},
		{
			name:   "Invalid area ID",
			areaID: 0,
			mock: func() {
				// Ничего не мокаем, так как функция должна вернуть ошибку до выполнения запроса
			},
			expectedResult: nil,
			expectedError:  ErrNotValidAreaID,
		},
		{
			name:   "DB error",
			areaID: areaID,
			mock: func() {
				mock.ExpectQuery(query).WithArgs(areaID).WillReturnError(errors.New("database error"))
			},
			expectedResult: nil,
			expectedError:  ErrDataBase,
		},
		{
			name:   "Scan error",
			areaID: areaID,
			mock: func() {
				rows := mock.NewRows([]string{"id", "name", "product", "characteristics", "level", "resources", "coordinates"}).
					AddRow(
						1,
						"Barracks",
						"Soldiers",
						[]byte("invalid json"), // Некорректный JSON для characteristics
						3,
						upgradePriceJSON1,
						coordinatesJSON1,
					)
				mock.ExpectQuery(query).WithArgs(areaID).WillReturnRows(rows)
			},
			expectedResult: nil,
			expectedError:  ErrRows,
		},
	}

	// Проходим по всем тестовым случаям
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			// Вызов тестируемой функции
			buildings, err := storage.GetBuildings(tt.areaID)

			// Проверяем ошибки
			if tt.expectedError != nil {
				assert.Error(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err, "GetBuildings should not return an error")
			}

			// Проверяем результат
			assert.Equal(t, tt.expectedResult, buildings, "Buildings should match expected data")

			// Проверяем, что все ожидания выполнены
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestGetHeroes(t *testing.T) {
	// Создаем мок базы данных
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer mock.Close()

	// Создаем Storage с моковой базой данных
	storage := &Storage{Db: mock}

	// Тестовые данные
	areaID := int64(1)
	expectedHero1 := models.Hero{
		Id:   1,
		Name: "Warrior",
		Charachteristics: models.HeroCharacteristics{
			HP:         1000,
			HPnow:      950,
			Armor:      50,
			Speed:      decimal.NewFromFloat(6.5),
			Vision:     20,
			IsRange:    false,
			AtackRange: decimal.NewFromFloat(1.5),
			Damage:     80,
		},
		Experience:     decimal.NewFromFloat(1500.0),
		ExperienceToUp: decimal.NewFromFloat(2000.0),
		Level:          5,
		Abilities: []models.Ability{
			{
				Id:   1,
				Name: "Slash",
				Charachteristics: models.AbilitytCharacteristics{
					IsPassive:      false,
					Radius:         decimal.NewFromFloat(1.0),
					Cooldown:       3 * time.Second,
					Damage:         decimal.NewFromFloat(50.0),
					ProjectilSpeed: decimal.NewFromFloat(10.0),
				},
				Level:   3,
				ImageId: 101,
			},
		},
		Coordinates: []models.Hex{{Q: 1, R: 2}, {Q: 3, R: 4}},
	}
	expectedHero2 := models.Hero{
		Id:   2,
		Name: "Archer",
		Charachteristics: models.HeroCharacteristics{
			HP:         800,
			HPnow:      750,
			Armor:      30,
			Speed:      decimal.NewFromFloat(7.0),
			Vision:     30,
			IsRange:    true,
			AtackRange: decimal.NewFromFloat(5.0),
			Damage:     50,
		},
		Experience:     decimal.NewFromFloat(1200.0),
		ExperienceToUp: decimal.NewFromFloat(1500.0),
		Level:          4,
		Abilities: []models.Ability{
			{
				Id:   2,
				Name: "Snipe",
				Charachteristics: models.AbilitytCharacteristics{
					IsPassive:      false,
					Radius:         decimal.NewFromFloat(5.0),
					Cooldown:       10 * time.Second,
					Damage:         decimal.NewFromFloat(100.0),
					ProjectilSpeed: decimal.NewFromFloat(20.0),
				},
				Level:   2,
				ImageId: 102,
			},
		},
		Coordinates: []models.Hex{{Q: 4, R: 5}, {Q: 6, R: 7}},
	}
	expectedHeroes := []models.Hero{expectedHero1, expectedHero2}

	// Сериализация характеристик, способностей и координат в JSON
	characteristicsJSON1, _ := json.Marshal(expectedHero1.Charachteristics)
	abilitiesJSON1, _ := json.Marshal(expectedHero1.Abilities)
	coordinatesJSON1, _ := json.Marshal(expectedHero1.Coordinates)
	characteristicsJSON2, _ := json.Marshal(expectedHero2.Charachteristics)
	abilitiesJSON2, _ := json.Marshal(expectedHero2.Abilities)
	coordinatesJSON2, _ := json.Marshal(expectedHero2.Coordinates)

	// Ожидаемый SQL-запрос
	query := `
        SELECT id, name, characteristics, experience, experience_to_up, level, abilities, coordinates 
        FROM areas_heroes 
        JOIN heroes ON areas_heroes.hero_id = heroes.id 
        WHERE areas_heroes.area_id = \$1;
    `

	// Таблица тестовых случаев
	tests := []struct {
		name           string
		areaID         int64
		mock           func()
		expectedResult []models.Hero
		expectedError  error
	}{
		{
			name:   "Valid data",
			areaID: areaID,
			mock: func() {
				rows := mock.NewRows([]string{"id", "name", "characteristics", "experience", "experience_to_up", "level", "abilities", "coordinates"}).
					AddRow(
						expectedHero1.Id,
						expectedHero1.Name,
						characteristicsJSON1,
						expectedHero1.Experience,
						expectedHero1.ExperienceToUp,
						expectedHero1.Level,
						abilitiesJSON1,
						coordinatesJSON1,
					).
					AddRow(
						expectedHero2.Id,
						expectedHero2.Name,
						characteristicsJSON2,
						expectedHero2.Experience,
						expectedHero2.ExperienceToUp,
						expectedHero2.Level,
						abilitiesJSON2,
						coordinatesJSON2,
					)
				mock.ExpectQuery(query).WithArgs(areaID).WillReturnRows(rows)
			},
			expectedResult: expectedHeroes,
			expectedError:  nil,
		},
		{
			name:   "Invalid area ID",
			areaID: 0,
			mock: func() {
				// Ничего не мокаем, так как функция должна вернуть ошибку до выполнения запроса
			},
			expectedResult: nil,
			expectedError:  ErrNotValidAreaID,
		},
		{
			name:   "DB error",
			areaID: areaID,
			mock: func() {
				mock.ExpectQuery(query).WithArgs(areaID).WillReturnError(errors.New("database error"))
			},
			expectedResult: nil,
			expectedError:  ErrDataBase,
		},
		{
			name:   "Scan error",
			areaID: areaID,
			mock: func() {
				rows := mock.NewRows([]string{"id", "name", "characteristics", "experience", "experience_to_up", "level", "abilities", "coordinates"}).
					AddRow(
						1,
						"Warrior",
						[]byte("invalid json"), // Некорректный JSON для characteristics
						decimal.NewFromFloat(1500.0),
						decimal.NewFromFloat(2000.0),
						5,
						abilitiesJSON1,
						coordinatesJSON1,
					)
				mock.ExpectQuery(query).WithArgs(areaID).WillReturnRows(rows)
			},
			expectedResult: nil,
			expectedError:  ErrRows,
		},
	}

	// Проходим по всем тестовым случаям
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			// Вызов тестируемой функции
			heroes, err := storage.GetHeroes(tt.areaID)

			// Проверяем ошибки
			if tt.expectedError != nil {
				assert.Error(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err, "GetHeroes should not return an error")
			}

			// Проверяем результат
			assert.Equal(t, tt.expectedResult, heroes, "Heroes should match expected data")

			// Проверяем, что все ожидания выполнены
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

// Табличный тест для функции AddEmptyArea
func TestAddEmptyArea(t *testing.T) {
	tests := []struct {
		name          string
		area          models.Area
		mockSetup     func(mock pgxmock.PgxPoolIface)
		expectedID    int64
		expectedError error
	}{
		{
			name: "Success - Area added",
			area: models.Area{
				UserId:     1,
				Width:      10,
				Height:     10,
				CellTypeId: 1,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`INSERT INTO area \(user_id, width,height,cell_type_id\) VALUES \(\$1,\$2,\$3,\$4\) RETURNING id;`).
					WithArgs(int64(1), 10, 10, 1).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int64(123)))
			},
			expectedID:    123,
			expectedError: nil,
		},
		{
			name: "Error - Database query failed",
			area: models.Area{
				UserId:     1,
				Width:      10,
				Height:     10,
				CellTypeId: 1,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`INSERT INTO area \(user_id, width,height,cell_type_id\) VALUES \(\$1,\$2,\$3,\$4\) RETURNING id;`).
					WithArgs(int64(1), 10, 10, 1).
					WillReturnError(fmt.Errorf("database error"))
			},
			expectedID:    0,
			expectedError: ErrDataBase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок базы данных
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			// Настраиваем мок
			tt.mockSetup(mock)

			// Создаем Storage с моком базы данных
			storage := &Storage{Db: mock}

			// Вызываем тестируемую функцию
			id, err := storage.AddEmptyArea(tt.area)

			// Проверяем ошибку
			if tt.expectedError != nil {
				assert.Error(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// Проверяем возвращенный ID
			assert.Equal(t, tt.expectedID, id)

			// Проверяем, что все ожидания были выполнены
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Табличный тест для функции AddNeutral
func TestAddNeutral(t *testing.T) {
	tests := []struct {
		name          string
		neutral       models.Neutral
		mockSetup     func(mock pgxmock.PgxPoolIface)
		expectedID    int64
		expectedError error
	}{
		{
			name: "Success - Neutral added",
			neutral: models.Neutral{
				Name:                    "Tree",
				Product:                 "Wood",
				ProductivityCoefficient: 1,
				Capacity:                decimal.NewFromFloat(100),
				ThresholdLevel1:         decimal.NewFromFloat(10),
				ThresholdLevel2:         decimal.NewFromFloat(5),
				Size:                    2,
				Coordinates:             []models.Hex{{Q: 1.0, R: 2.0}, {Q: 3.0, R: 4.0}},
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				coordsJSON, _ := json.Marshal([]models.Hex{{Q: 1.0, R: 2.0}, {Q: 3.0, R: 4.0}})
				mock.ExpectQuery(`INSERT INTO neutral \(name, product, prod_cof, capacity, threshold_level1, threshold_level2, size, coordinates\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\) RETURNING id;`).
					WithArgs("Tree", "Wood", 1, decimal.NewFromFloat(100), decimal.NewFromFloat(10), decimal.NewFromFloat(5), 2, coordsJSON).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int64(123)))
			},
			expectedID:    123,
			expectedError: nil,
		},
		{
			name: "Error - Failed to marshal coordinates",
			neutral: models.Neutral{
				Name:                    "Tree",
				Product:                 "Wood",
				ProductivityCoefficient: 1,
				Capacity:                decimal.NewFromFloat(100),
				ThresholdLevel1:         decimal.NewFromFloat(10),
				ThresholdLevel2:         decimal.NewFromFloat(5),
				Size:                    2,
				Coordinates:             []models.Hex{{Q: 1.0, R: 2.0}, {Q: math.Inf(1), R: 4.0}}, // Некорректные координаты
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				// Ничего не настраиваем, так как ошибка произойдет до выполнения запроса
			},
			expectedID:    0,
			expectedError: ErrNotValidCoord,
		},
		{
			name: "Error - Database query failed",
			neutral: models.Neutral{
				Name:                    "Tree",
				Product:                 "Wood",
				ProductivityCoefficient: 1,
				Capacity:                decimal.NewFromFloat(100),
				ThresholdLevel1:         decimal.NewFromFloat(10),
				ThresholdLevel2:         decimal.NewFromFloat(5),
				Size:                    2,
				Coordinates:             []models.Hex{{Q: 1.0, R: 2.0}, {Q: 3.0, R: 4.0}},
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				coordsJSON, _ := json.Marshal([]models.Hex{{Q: 1.0, R: 2.0}, {Q: 3.0, R: 4.0}})
				mock.ExpectQuery(`INSERT INTO neutral \(name, product, prod_cof, capacity, threshold_level1, threshold_level2, size, coordinates\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\) RETURNING id;`).
					WithArgs("Tree", "Wood", 1, decimal.NewFromFloat(100), decimal.NewFromFloat(10), decimal.NewFromFloat(5), 2, coordsJSON).
					WillReturnError(fmt.Errorf("database error"))
			},
			expectedID:    0,
			expectedError: ErrDataBase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок базы данных
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			// Настраиваем мок
			tt.mockSetup(mock)

			// Создаем Storage с моком базы данных
			storage := &Storage{Db: mock}

			// Вызываем тестируемую функцию
			id, err := storage.AddNeutral(tt.neutral)

			// Проверяем ошибку
			if tt.expectedError != nil {
				assert.Error(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// Проверяем возвращенный ID
			assert.Equal(t, tt.expectedID, id)

			// Проверяем, что все ожидания были выполнены
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Табличный тест для функции AddBuilding
func TestAddBuilding(t *testing.T) {
	//готовим тестовые данные
	coordsJSON, _ := json.Marshal([]models.Hex{{Q: 1.0, R: 2.0}, {Q: 3.0, R: 4.0}})
	characteristicsJSON, _ := json.Marshal(models.BuildingCharacteristics{
		HP:                      100,
		Armor:                   50,
		ProductivityCoefficient: 2,
		Size:                    4,
	})
	resourcesJSON, _ := json.Marshal([]models.Resource{
		{Id: 1, Name: "Wood"},
		{Id: 2, Name: "Stone"},
	})
	//определяем структуру тестов
	tests := []struct {
		name          string
		building      models.Building
		mockSetup     func(mock pgxmock.PgxPoolIface)
		expectedID    int64
		expectedError error
	}{
		{
			name: "Success - Building added",
			building: models.Building{
				Name:    "Farm",
				Product: "Food",
				Charachteristics: models.BuildingCharacteristics{
					HP:                      100,
					Armor:                   50,
					ProductivityCoefficient: 2,
					Size:                    4,
				},
				Level: 1,
				UpgradePrice: []models.Resource{
					{Id: 1, Name: "Wood"},
					{Id: 2, Name: "Stone"},
				},
				Coordinates: []models.Hex{{Q: 1.0, R: 2.0}, {Q: 3.0, R: 4.0}},
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`INSERT INTO building \(name, product, characteristics, level, resources, coordinates\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\) RETURNING id;`).
					WithArgs("Farm", "Food", characteristicsJSON, 1, resourcesJSON, coordsJSON).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int64(123)))
			},
			expectedID:    123,
			expectedError: nil,
		},
		{
			name: "Error - Failed to marshal coordinates",
			building: models.Building{
				Name:    "Farm",
				Product: "Food",
				Charachteristics: models.BuildingCharacteristics{
					HP:                      100,
					Armor:                   50,
					ProductivityCoefficient: 2,
					Size:                    4,
				},
				Level: 1,
				UpgradePrice: []models.Resource{
					{Id: 1, Name: "Wood"},
					{Id: 2, Name: "Stone"},
				},
				Coordinates: []models.Hex{{Q: 1.0, R: 2.0}, {Q: math.Inf(1), R: 4.0}}, // Некорректные координаты
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				// Ничего не настраиваем, так как ошибка произойдет до выполнения запроса
			},
			expectedID:    0,
			expectedError: ErrNotValidCoord,
		},
		{
			name: "Error - Failed to marshal characteristics",
			building: models.Building{
				Name:    "Farm",
				Product: "Food",
				Charachteristics: models.BuildingCharacteristics{
					HP:                      100,
					Armor:                   50,
					ProductivityCoefficient: 2,
					Size:                    4,
				},
				Level: 1,
				UpgradePrice: []models.Resource{
					{Id: 1, Name: "Wood"},
					{Id: 2, Name: "Stone"},
				},
				Coordinates: []models.Hex{{Q: 1.0, R: 2.0}, {Q: 3.0, R: 4.0}},
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				// Настраиваем мок, чтобы вернуть ошибку при маршалинге характеристик
				mock.ExpectQuery(`INSERT INTO building \(name, product, characteristics, level, resources, coordinates\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\) RETURNING id;`).
					WithArgs("Farm", "Food", characteristicsJSON, 1, resourcesJSON, coordsJSON).
					WillReturnError(fmt.Errorf("Failed to marshal characteristics"))
			},
			expectedID:    0,
			expectedError: ErrNotValidChar,
		},
		{
			name: "Error - Failed to marshal resources",
			building: models.Building{
				Name:    "Farm",
				Product: "Food",
				Charachteristics: models.BuildingCharacteristics{
					HP:                      100,
					Armor:                   50,
					ProductivityCoefficient: 2,
					Size:                    4,
				},
				Level: 1,
				UpgradePrice: []models.Resource{
					{Id: 1, Name: "Wood"},
					{Id: 2, Name: "Stone"},
				},
				Coordinates: []models.Hex{{Q: 1.0, R: 2.0}, {Q: math.Inf(1), R: 4.0}},
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				// Ничего не настраиваем, так как ошибка произойдет до выполнения запроса
			},
			expectedID:    0,
			expectedError: ErrNotValidRes,
		},
		{
			name: "Error - Database query failed",
			building: models.Building{
				Name:    "Farm",
				Product: "Food",
				Charachteristics: models.BuildingCharacteristics{
					HP:                      100,
					Armor:                   50,
					ProductivityCoefficient: 2,
					Size:                    4,
				},
				Level: 1,
				UpgradePrice: []models.Resource{
					{Id: 1, Name: "Wood"},
					{Id: 2, Name: "Stone"},
				},
				Coordinates: []models.Hex{{Q: 1.0, R: 2.0}, {Q: 3.0, R: 4.0}},
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				coordsJSON, _ := json.Marshal([]models.Hex{{Q: 1.0, R: 2.0}, {Q: 3.0, R: 4.0}})
				characteristicsJSON, _ := json.Marshal(models.BuildingCharacteristics{
					HP:                      100,
					Armor:                   50,
					ProductivityCoefficient: 2,
					Size:                    4,
				})
				resourcesJSON, _ := json.Marshal([]models.Resource{
					{Id: 1, Name: "Wood"},
					{Id: 2, Name: "Stone"},
				})
				mock.ExpectQuery(`INSERT INTO building \(name, product, characteristics, level, resources, coordinates\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\) RETURNING id;`).
					WithArgs("Farm", "Food", characteristicsJSON, 1, resourcesJSON, coordsJSON).
					WillReturnError(fmt.Errorf("database error"))
			},
			expectedID:    0,
			expectedError: ErrDataBase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок базы данных
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			// Настраиваем мок
			tt.mockSetup(mock)

			// Создаем Storage с моком базы данных
			storage := &Storage{Db: mock}

			// Вызываем тестируемую функцию
			id, err := storage.AddBuilding(tt.building)

			// Проверяем ошибку
			if tt.expectedError != nil {
				assert.Error(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// Проверяем возвращенный ID
			assert.Equal(t, tt.expectedID, id)

			// Проверяем, что все ожидания были выполнены
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAddHero(t *testing.T) {
	hero := models.Hero{
		Name:           "Ion Mash",
		Experience:     decimal.NewFromFloat(150.0),
		ExperienceToUp: decimal.NewFromFloat(200.0),
		Level:          1,
		Charachteristics: models.HeroCharacteristics{
			HP:         100,
			HPnow:      100,
			Armor:      20,
			Speed:      decimal.NewFromFloat(5.0),
			Vision:     10,
			IsRange:    false,
			AtackRange: decimal.NewFromFloat(1.0),
			Damage:     20,
		},
		Abilities: []models.Ability{
			{
				Id:      1,
				Name:    "Plasma explosion",
				Level:   1,
				ImageId: 101,
				Charachteristics: models.AbilitytCharacteristics{
					IsPassive:      false,
					Radius:         decimal.NewFromFloat(1.0),
					Cooldown:       time.Second * 7,
					Damage:         decimal.NewFromFloat(30.0),
					ProjectilSpeed: decimal.NewFromFloat(0.0),
				},
			},
		},
		Coordinates: []models.Hex{{Q: 1.0, R: 2.0}, {Q: 3.0, R: 4.0}},
	}
	coordsJSON, _ := json.Marshal(hero.Coordinates)
	characteristicsJSON, _ := json.Marshal(hero.Charachteristics)
	abilitiesJSON, _ := json.Marshal(hero.Abilities)

	tests := []struct {
		name          string
		hero          models.Hero
		mockSetup     func(mock pgxmock.PgxPoolIface)
		expectedID    int64
		expectedError error
	}{
		{
			name: "Success - Hero added",
			hero: hero,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`INSERT INTO hero \(name, characteristics, experience, experience_to_up, level, abilities, coordinates\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\) RETURNING id;`).
					WithArgs(hero.Name, characteristicsJSON, hero.Experience, hero.ExperienceToUp, hero.Level, abilitiesJSON, coordsJSON).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int64(123)))
			},
			expectedID:    123,
			expectedError: nil,
		},
		{
			name: "Error - Failed to marshal coordinates",
			hero: models.Hero{
				Name:           "Ion Mash",
				Experience:     decimal.NewFromFloat(150.0),
				ExperienceToUp: decimal.NewFromFloat(200.0),
				Level:          1,
				Charachteristics: models.HeroCharacteristics{
					HP:         100,
					HPnow:      100,
					Armor:      20,
					Speed:      decimal.NewFromFloat(5.0),
					Vision:     10,
					IsRange:    false,
					AtackRange: decimal.NewFromFloat(1.0),
					Damage:     20,
				},
				Abilities: []models.Ability{
					{
						Id:      1,
						Name:    "Plasma explosion",
						Level:   1,
						ImageId: 101,
						Charachteristics: models.AbilitytCharacteristics{
							IsPassive:      false,
							Radius:         decimal.NewFromFloat(1.0),
							Cooldown:       time.Second * 7,
							Damage:         decimal.NewFromFloat(30.0),
							ProjectilSpeed: decimal.NewFromFloat(0.0),
						},
					},
				},
				Coordinates: []models.Hex{{Q: 1.0, R: 2.0}, {Q: math.Inf(1), R: 4.0}}, // Некорректные координаты
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {

			},
			expectedID:    0,
			expectedError: ErrNotValidCoord,
		},
		{
			name: "Error - Failed to marshal characteristics",
			hero: models.Hero{
				Name:           "Ion Mash",
				Experience:     decimal.NewFromFloat(150.0),
				ExperienceToUp: decimal.NewFromFloat(200.0),
				Level:          1,
				Charachteristics: models.HeroCharacteristics{
					HP:         100,
					HPnow:      100,
					Armor:      20,
					Speed:      decimal.NewFromFloat(5.0),
					Vision:     10,
					IsRange:    false,
					AtackRange: decimal.NewFromFloat(1.0),
					Damage:     20,
				},
				Abilities: []models.Ability{
					{
						Id:      1,
						Name:    "Plasma explosion",
						Level:   1,
						ImageId: 101,
						Charachteristics: models.AbilitytCharacteristics{
							IsPassive:      false,
							Radius:         decimal.NewFromFloat(1.0),
							Cooldown:       time.Second * 7,
							Damage:         decimal.NewFromFloat(30.0),
							ProjectilSpeed: decimal.NewFromFloat(0.0),
						},
					},
				},
				Coordinates: []models.Hex{{Q: 1.0, R: 2.0}, {Q: 3.0, R: 4.0}},
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`INSERT INTO hero \(name, characteristics, experience, experience_to_up, level, abilities, coordinates\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\) RETURNING id;`).
					WithArgs("Ion Mash", characteristicsJSON, decimal.NewFromFloat(150.0), decimal.NewFromFloat(200.0), 1, abilitiesJSON, coordsJSON).
					WillReturnError(fmt.Errorf("Failed to marshal characteristics"))
			},
			expectedID:    0,
			expectedError: ErrNotValidChar,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок базы данных
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			// Настраиваем мок
			tt.mockSetup(mock)

			// Создаем Storage с моком базы данных
			storage := &Storage{Db: mock}

			// Вызываем тестируемую функцию
			id, err := storage.AddHero(tt.hero)

			// Проверяем ошибку
			if tt.expectedError != nil {
				assert.ErrorContains(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Проверяем возвращенный ID
			assert.Equal(t, tt.expectedID, id)

			// Проверяем, что все ожидания были выполнены
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
