/*
Пакет полностью тестовый и на данном этапе необходим для тестирования
других функций приложения.
Пакет предназначен для создание арен и на их базе миров пользователей.
В дальнейшем это будет отдельный микросервис.
Задача пакета - создать арену с случайно расположеными на ней объектами.
*/

package game

import (
	storage "cyber/internal/storage"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/shopspring/decimal"

	"cyber/internal/models"
	//"cyber/internal/storage"
)

const (
	width  = 400
	height = 400
)

// Функция вовзращающая случайный слайс координат в зависимости от размера объекта
func generateCoordinates(size int) []models.Hex {
	rand.Seed(time.Now().UnixNano()) // Инициализация генератора случайных чисел

	// Генерация начальной точки / размещение объекта в рамках игрового поля
	startQ := rand.Intn(width - size)
	startR := rand.Intn(height - size)

	// Генерация координат для объекта
	coordinates := make([]models.Hex, 0, width*height)
	for q := startQ; q < startQ+size; q++ {
		for r := startR; r < startR+size; r++ {
			coordinates = append(coordinates, models.Hex{Q: float64(q), R: float64(r)})
		}
	}

	return coordinates
}

// HARDCODE generateNeutral создает нейтральный объект расположенный в случайном месте арены.
func generateNeutral() models.Neutral {
	return models.Neutral{
		Name:                    "Gold mine",
		Product:                 "Gold",
		ProductivityCoefficient: 4,
		Capacity:                decimal.NewFromFloat(rand.Float64() * 10000),
		ThresholdLevel1:         decimal.NewFromFloat(rand.Float64() * 5000),
		ThresholdLevel2:         decimal.NewFromFloat(rand.Float64() * 2000),
		Size:                    4,
		Coordinates:             generateCoordinates(4),
	}
}

// HARDCODE generateBuilding создает здание расположенное в случайном месте арены.
func generateBuilding() models.Building {
	return models.Building{
		Name:         "CyMan miner house",
		Product:      "Gold",
		Level:        1,
		UpgradePrice: []models.Resource{{Name: "Gold", Value: decimal.NewFromFloat(1000)}},
		Coordinates:  generateCoordinates(2),
		Charachteristics: models.BuildingCharacteristics{
			HP:                      500,
			Armor:                   10,
			ProductivityCoefficient: 3,
			Size:                    2,
		},
	}
}

// HARDCODE generateHero создает героя расположенныое в случайном месте арены.
func generateHero() models.Hero {
	return models.Hero{
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
		Coordinates: generateCoordinates(1),
	}
}

// HARDCODE generateUnit создает юнита расположенного в случайном месте арены.
func generateUnit() models.Unit {
	return models.Unit{
		Name:           "Miner",
		Experience:     decimal.NewFromFloat(50.0),
		ExperienceToUp: decimal.NewFromFloat(100.0),
		Level:          1,
		ImageId:        201,
		Charachteristics: models.UnitCharacteristics{
			HP:                      80,
			HPnow:                   80,
			Armor:                   5,
			Speed:                   decimal.NewFromFloat(4.0),
			Vision:                  8,
			IsRange:                 false,
			AtackRange:              decimal.NewFromFloat(1.0),
			Damage:                  decimal.NewFromFloat(15.0),
			ProductivityCoefficient: 4,
		},
		Coordinates: generateCoordinates(1),
	}
}

func NewArea(userId int) models.Area {
	return models.Area{
		UserId:     int64(userId),
		Width:      width,
		Height:     height,
		CellTypeId: 1, //HARDCODE - на данном этапе считаем поле покрытой одной текстурой
	}
}

// Фактически функция create world - устанавливает связи объектов и арен
func CreateWorld(userId int) error {
	db, err := storage.New()
	if err != nil {
		return fmt.Errorf("error db connection: %v\n", err)
	}
	area := NewArea(userId)
	if err != nil {
		log.Fatalf("Cant create world! %v\n", err)
		return err
	}
	areaId, err := db.AddEmptyArea(area)
	if err != nil {
		return err
	}

	neutral := generateNeutral()
	building := generateBuilding()
	hero := generateHero()
	unit := generateUnit()

	neutralId, err := db.AddNeutral(neutral)
	if err != nil {
		log.Printf("Failed when creating neutral: %v\n", err)
		return err
	}
	buildingId, err := db.AddBuilding(building)
	if err != nil {
		log.Printf("Failed when creating building: %v\n", err)
		return err
	}
	heroId, err := db.AddHero(hero)
	if err != nil {
		log.Printf("Failed when creating hero: %v\n", err)
		return err
	}
	unitId, err := db.AddUnit(unit)
	if err != nil {
		log.Printf("Failed when creating unit: %v\n", err)
		return err
	}
	db.AddNeutralAtArea(neutralId, areaId)
	db.AddBuildingAtArea(buildingId, areaId)
	db.AddHeroAtArea(heroId, areaId)
	db.AddUnitAtArea(unitId, areaId)
	return nil

}
