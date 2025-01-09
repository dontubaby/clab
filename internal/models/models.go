package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type WorldState struct {
	SessionId int64     `db:"session_id"`
	UserId    int64     `db:"user_id"`
	AreaId    int64     `db:"area_id"`
	ActionId  int64     `db:"action_id"`
	TimeSt    time.Time `db:"time_stamp"`
}

type Resource struct {
	Id     int             `json:"id"`
	UserId int             `json:"user_id"`
	Name   string          `json:"name"`
	Value  decimal.Decimal `json:"value"`
}

type User struct {
	Id           int64           `db:"id"`
	Login        string          `db:"login"`
	Password     string          `db:"password"`
	Email        string          `db:"email"`
	Subscription bool            `db:"subscription"`
	LeagueId     int64           `db:"league_id"`
	Resources    []Resource      `db:"resources"`
	Balance      decimal.Decimal `db:"balance"`
	Level        int             `db:"level"`
}

type League struct {
	Id        int64  `db:"id"`
	Name      string `db:"name"`
	Authority int    `db:"authority"` //уровень авторитета лиги
}

type CellType int

const (
	Grass CellType = iota + 1
	Brick
	Water
	Sand
)

type Area struct {
	Id         int     `db:"id"`
	User_id    int     `db:"user_id"`
	Width      int     `db:"width"`
	Heigth     int     `db:"heigth"`
	CellTypeId int     `db:"cell_type_id"`
	ObjectsIds []int64 `db:"object_id"`
}

type Objects struct {
	Id             int64 `db:"id"`
	Area_id        int64 `db:"area_id"`
	Object_type_id int   `db:"object_type_id"`
	Object_id      int   `db:"object_id"`
}

type Building struct {
	Id               int64
	Charachteristics BuildingCharacteristics
	Level            int
}

// Возможно стоит ограничиться  типо Characteristics c полями Name и Value
// Который будет универсальным типом характеристик для Building, Unit, Enemy
type BuildingCharacteristics struct {
	Name                    string `json:"name"`
	Product                 string `json:"product"` //тип производимого продутка (популяция, еда, минералы и т.д.)
	HP                      int    `json:"hp"`
	Armor                   int    `json:"armor"`
	ProductivityCoefficient int    `json:"prod_cof"` //коэфициент производительности, который  при умножении на ед. времени дает количество продукта
	Level                   int    `json:"level"`
}

type Unit struct {
	Id               int64               `db:"id"`
	Charachteristics UnitCharacteristics `db:"charachteristics"`
	Level            int                 `db:"level"`
}

type UnitCharacteristics struct {
	Name                    string          `json:"name"`
	HP                      int             `json:"hp"`
	Armor                   int             `json:"armor"`
	Speed                   decimal.Decimal `json:"speed"`  //количество клеток(гексов) проходимых в ед.вр.
	Vision                  int             `json:"vision"` //количество клеток(гексов) на которые видит юнит
	IsRange                 bool            `json:"range"`
	Atack_range             decimal.Decimal `json:"atack_range"`
	Damage                  decimal.Decimal `json:"damage"`
	ProductivityCoefficient int             `json:"prod_cof"` //коэфициент производительности, который  при умножении на ед. времени дает количество продукта
	Level                   int             `json:"level"`
}

type Enemy struct {
	Id               int64                `db:"id"`
	Charachteristics EnemyCharacteristics `db:"charachteristics"`
	Level            int                  `db:"level"`
}

type EnemyCharacteristics struct {
	Name        string          `json:"name"`
	HP          int             `json:"hp"`
	Armor       int             `json:"armor"`
	Speed       decimal.Decimal `json:"speed"`  //количество клеток(гексов) проходимых в ед.вр.
	Vision      int             `json:"vision"` //количество клеток(гексов) на которые видит юнит
	IsRange     bool            `json:"range"`
	Atack_range decimal.Decimal `json:"atack_range"`
	Damage      decimal.Decimal `json:"damage"`
	Level       int             `json:"level"`
}

type Action struct {
	Id              int64
	UserId          int64
	AreaId          int64
	ObjectSourceId  int64
	ObjectDestId    int64
	ActionType      ActionType
	Characteristics ActionCharacteristics
	StartTime       time.Time
	Duration        time.Duration
	Status          bool //true - выполнено, false - нет

}

type ActionType int

const (
	MoveAction ActionType = iota + 1
	HArvestAction
	BuildAction
	AtackAction
)

// MoveActionCharacteristics описывает характеристики перемещения.
type MoveActionCharacteristics struct {
	From  Point           `json:"from"`
	To    Point           `json:"to"`
	Speed decimal.Decimal `json:"speed"`
}

type HarvestActionCharacteristics struct {
	Harvester    int64           `json:"harvester"`
	ResourceName string          `json:"resource"`
	Speed        decimal.Decimal `json:"speed"`
}

// BuildActionCharacteristics описывает характеристики строительства.
type BuildActionCharacteristics struct {
	Builder    int64           `json:"builder"`
	ObjectType string          `json:"object"`
	Speed      decimal.Decimal `json:"speed"`
}

// AttackActionCharacteristics описывает характеристики атаки.
type AttackActionCharacteristics struct {
	Atacker  int64           `json:"atacker"`
	Defenser int64           `json:"defenser"`
	Damage   decimal.Decimal `json:"pointx"`
}

type Point struct {
	X float64 `json:"X"`
	Y float64 `json:"Y"`
}
type ActionCharacteristics interface{} // адаптер для предоставления конструктору характеристик действия

// Этот метод будет в другом блоке, но он отражает логику адаптации характеристик действия
// в зависимости от типа действия
func NewAction(actionType ActionType, charachteristics ActionCharacteristics) (*Action, error) {
	switch actionType {
	case MoveAction:
		characteristicsPtr, ok := characteristics.(*MoveActionCharacteristics)
		if !ok {
			return nil, fmt.Errorf("неверные характеристики для типа действия %d", actionType)
		}
	case HarvestAction:
		characteristicsPtr, ok := characteristics.(*HarvestActionCharacteristics)
		if !ok {
			return nil, fmt.Errorf("неверные характеристики для типа действия %d", actionType)
		}
	case BuildAction:
		characteristicsPtr, ok := characteristics.(*BuildActionCharacteristics)
		if !ok {
			return nil, fmt.Errorf("неверные характеристики для типа действия %d", actionType)
		}
	case AttackAction:
		characteristicsPtr, ok := characteristics.(*AttackActionCharacteristics)
		if !ok {
			return nil, fmt.Errorf("неверные характеристики для типа действия %d", actionType)
		}
	default:
		return nil, errors.New("неизвестный тип действия")
	}

	action := &Action{
		ActionType:      actionType,
		Characteristics: characteristicsPtr,
	}

	return action, nil
}
