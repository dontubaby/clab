/**/
package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// WorldState представляет состояние мира в конкретный момент времени.
type WorldState struct {
	SessionId int64     `db:"session_id"` // Идентификатор сессии
	UserId    int64     `db:"user_id"`    // Идентификатор пользователя
	AreaId    int64     `db:"area_id"`    // Идентификатор области (арены)
	ActionId  int64     `db:"action_id"`  // Идентификатор действия
	TimeSt    time.Time `db:"time_stamp"` // Временная метка состояния
}

// Resource представляет ресурс, принадлежащий пользователю.
type Resource struct {
	Id    int             `json:"id"`    // Идентификатор ресурса
	Name  string          `json:"name"`  // Название ресурса
	Value decimal.Decimal `json:"value"` // Значение ресурса (количество)
}

// User представляет пользователя системы.
type User struct {
	Id           int64           `db:"id"`           // Идентификатор пользователя
	Login        string          `db:"login"`        // Логин пользователя
	Password     string          `db:"password"`     // Пароль пользователя
	Email        string          `db:"email"`        // Электронная почта пользователя
	Subscription bool            `db:"subscription"` // Флаг подписки пользователя
	LeagueId     int64           `db:"league_id"`    // Идентификатор лиги пользователя
	Balance      decimal.Decimal `db:"balance"`      // Баланс пользователя
	Level        int             `db:"level"`        // Уровень пользователя
}

// League представляет лигу, в которой состоят пользователи.
type League struct {
	Id        int64  `db:"id"`        // Идентификатор лиги
	Name      string `db:"name"`      // Название лиги
	Authority int    `db:"authority"` // Уровень авторитета лиги
}

// CellType представляет тип клетки на арене.
type CellType int

const (
	Grass CellType = iota + 1 // Тип клетки: трава
	Brick                     // Тип клетки: кирпич
	Water                     // Тип клетки: вода
	Sand                      // Тип клетки: песок
)

// Area представляет арену игрока.
type Area struct {
	Id         int64   `db:"id"`           // Идентификатор арены
	UserId     int64   `db:"user_id"`      // Идентификатор пользователя, владеющего ареной
	Width      int     `db:"width"`        // Ширина арены
	Height     int     `db:"heigth"`       // Высота арены
	CellTypeId int     `db:"cell_type_id"` // Идентификатор типа клетки
	ObjectsIds []int64 `db:"object_id"`    // Список идентификаторов объектов на арене
}

// Neutral представляет нейтральный объект на арене.
type Neutral struct {
	Id                      int64           `db:"id"`                             // Идентификатор объекта
	Name                    string          `db:"name"`                           // Название объекта (например, золотая шахта)
	Product                 string          `db:"product"`                        // Тип производимого продукта (популяция, еда, минералы и т.д.)
	ProductivityCoefficient int             `db:"prod_cof"`                       // Коэффициент производительности
	Capacity                decimal.Decimal `db:"capacity"`                       // Емкость ресурса объекта
	ThresholdLevel1         decimal.Decimal `db:"threshold_level1"`               // Порог значения ресурса, после которого скорость добычи снижается в 5 раз
	ThresholdLevel2         decimal.Decimal `db:"threshold_level2"`               // Порог значения ресурса, после которого скорость добычи снижается в 10 раз
	Size                    int             `db:"size"`                           //Количество клеток которое занимает объект
	Coordinates             []Hex           `db:"coordinates" json:"coordinates"` // Координаты объекта на арене
}

// Building представляет здание на арене.
type Building struct {
	Id               int64                   `db:"id"`                                     // Идентификатор здания
	Name             string                  `db:"name"`                                   // Название здания
	Product          string                  `db:"product"`                                // Тип производимого продукта
	Charachteristics BuildingCharacteristics `db:"characteristics" json:"characteristics"` // Характеристики здания
	Level            int                     `db:"level"`                                  // Уровень здания
	UpgradePrice     []Resource              `db:"resources"`                              // Стоимость улучшения здания
	Coordinates      []Hex                   `db:"coordinates" json:"coordinates"`         // Координаты объекта на арене
}

// BuildingCharacteristics представляет характеристики здания.
type BuildingCharacteristics struct {
	HP                      int `json:"hp"`       // Здоровье здания
	Armor                   int `json:"armor"`    // Броня здания
	ProductivityCoefficient int `json:"prod_cof"` // Коэффициент производительности здания
	Size                    int `db:"size"`       //Количество клеток которое занимает объект
}

// Heroe представляет героя.
type Hero struct {
	Id               int64               `db:"id"`                                     // Идентификатор героя
	Name             string              `db:"name"`                                   // Название героя
	Charachteristics HeroCharacteristics `db:"characteristics" json:"characteristics"` // Характеристики героя
	Experience       decimal.Decimal     `db:"experience"`                             // Текущий опыт героя
	ExperienceToUp   decimal.Decimal     `db:"experience_to_up"`                       // Опыт, необходимый для повышения уровня
	Level            int                 `db:"level"`                                  // Уровень героя
	Abilities        []Ability           `json:"abilities"`                            // Список идентификаторов способностей героя
	Coordinates      []Hex               `db:"coordinates"`                            // Координаты объекта на арене
}

// HeroCharacteristics представляет характеристики героя.
type HeroCharacteristics struct {
	HP         int             `json:"hp"`          // Общее здоровье героя
	HPnow      int             `json:"hp_now"`      // Текущее здоровье героя
	Armor      int             `json:"armor"`       // Текущая защита героя
	Speed      decimal.Decimal `json:"speed"`       // Скорость перемещения героя
	Vision     int             `json:"vision"`      // Дальность обзора героя
	IsRange    bool            `json:"range"`       // Флаг, определяющий, является ли герой дальнобойным
	AtackRange decimal.Decimal `json:"atack_range"` // Дальность атаки героя
	Damage     int             `json:"damage"`      // Текущий урон героя
}

// Ability представляет способность героя.
type Ability struct {
	Id               int64                   `db:"id"`                                     // Идентификатор способности
	Name             string                  `db:"name"`                                   // Название способности
	Charachteristics AbilitytCharacteristics `db:"characteristics" json:"characteristics"` // Характеристики способности
	Level            int                     `db:"level"`                                  // Уровень способности
	ImageId          int64                   `db:"image_id"`                               // Идентификатор изображения способности
}

// AbilitytCharacteristics представляет характеристики способности.
type AbilitytCharacteristics struct {
	IsPassive      bool            `json:"is_passive"`      // Флаг, определяющий, является ли способность пассивной
	Radius         decimal.Decimal `json:"radius"`          // Радиус действия способности
	Cooldown       time.Duration   `json:"cooldown"`        // Время перезарядки способности
	Damage         decimal.Decimal `json:"damage"`          // Урон от способности
	ProjectilSpeed decimal.Decimal `json:"projectil_speed"` // Скорость снаряда (если применимо)
}

// Unit представляет юнита.
type Unit struct {
	Id               int64               `db:"id"`                                     // Идентификатор юнита
	Name             string              `db:"name"`                                   // Название юнита
	Charachteristics UnitCharacteristics `db:"characteristics" json:"characteristics"` // Характеристики юнита
	Experience       decimal.Decimal     `db:"experience"`                             // Текущий опыт юнита
	ExperienceToUp   decimal.Decimal     `db:"experience_to_up"`                       // Опыт, необходимый для повышения уровня
	Level            int                 `db:"level"`                                  // Уровень юнита
	ImageId          int64               `db:"image_id"`                               // Идентификатор изображения юнита
	Coordinates      []Hex               `db:"coordinates"`                            // Координаты объекта на арене
}

// UnitCharacteristics представляет характеристики юнита.
type UnitCharacteristics struct {
	HP                      int             `json:"hp"`          // Здоровье юнита
	HPnow                   int             `json:"hp_now"`      // Текущее здоровье героя
	Armor                   int             `json:"armor"`       // Броня юнита
	Speed                   decimal.Decimal `json:"speed"`       // Скорость перемещения юнита в клетках в сединицу времени
	Vision                  int             `json:"vision"`      // Дальность обзора юнита в клетках
	IsRange                 bool            `json:"range"`       // Флаг, определяющий, является ли юнит дальнобойным
	AtackRange              decimal.Decimal `json:"atack_range"` // Дальность атаки юнита
	Damage                  decimal.Decimal `json:"damage"`      // Урон юнита
	ProductivityCoefficient int             `json:"prod_cof"`    // Коэффициент производительности юнита
}

// Enemy представляет врага.
type Enemy struct {
	Id               int64                `db:"id"`                                     // Идентификатор врага
	Name             string               `db:"name"`                                   // Название врага
	Charachteristics EnemyCharacteristics `db:"characteristics" json:"characteristics"` // Характеристики врага
	Level            int                  `db:"level"`                                  // Уровень врага
	Coordinates      []Hex                `db:"coordinates"`                            // Координаты объекта на арене
}

// EnemyCharacteristics представляет характеристики врага.
type EnemyCharacteristics struct {
	HP         int             `json:"hp"`          // Здоровье врага
	Armor      int             `json:"armor"`       // Броня врага
	Speed      decimal.Decimal `json:"speed"`       // Скорость перемещения врага
	Vision     int             `json:"vision"`      // Дальность обзора врага
	IsRange    bool            `json:"range"`       // Флаг, определяющий, является ли враг дальнобойным
	AtackRange decimal.Decimal `json:"atack_range"` // Дальность атаки врага
	Damage     decimal.Decimal `json:"damage"`      // Урон врага
	Experience decimal.Decimal `json:"experience"`  // Опыт, получаемый за победу над врагом
	Level      int             `json:"level"`       // Уровень врага
}

// Action представляет действие, выполняемое пользователем.
type Action struct {
	Id              int64         `db:"id" json:"id"`                             // Идентификатор действия
	UserId          int64         `db:"user_id" json:"user_id"`                   // Идентификатор пользователя
	AreaId          int64         `db:"area_id" json:"area_id"`                   // Идентификатор арены
	ObjectSourceId  int64         `db:"object_source_id" json:"object_source_id"` // Идентификатор объекта-источника действия
	ObjectDestId    int64         `db:"object_dest_id" json:"object_dest_id"`     // Идентификатор объекта-цели действия
	ActionType      string        `db:"action_type" json:"action_type"`           // Тип действия (move, attack, build и т.д.)
	Characteristics []byte        `db:"characteristics" json:"characteristics"`   // Характеристики действия
	StartTime       time.Time     `db:"start_time" json:"start_time"`             // Время начала действия
	Duration        time.Duration `db:"duration" json:"duration"`                 // Продолжительность действия
	Status          string        `db:"status" json:"status"`                     // Статус действия (1-DONE,2-NOT_DONE,3-PROCESS)
}

// ActionType представляет тип действия.
type ActionType int

const (
	MoveAction    ActionType = iota + 1 // Действие: перемещение
	HarvestAction                       // Действие: сбор ресурсов
	BuildAction                         // Действие: строительство
	AttackAction                        // Действие: атака
)

// MoveActionCharacteristics описывает характеристики перемещения.
type MoveActionCharacteristics struct {
	From  Hex             `json:"from"`  // Начальная точка перемещения
	To    Hex             `json:"to"`    // Конечная точка перемещения
	Speed decimal.Decimal `json:"speed"` // Скорость перемещения
}

// HarvestActionCharacteristics описывает характеристики сбора ресурсов.
type HarvestActionCharacteristics struct {
	Harvester    int64           `json:"harvester"`  // Идентификатор объекта, собирающего ресурсы
	ResourceName string          `json:"resource"`   // Название ресурса
	NeutralId    int64           `json:"neutral_id"` //id нейтрального объекта из которого будет проводиться добыча ресурса
	Speed        decimal.Decimal `json:"speed"`      // Скорость сбора ресурсов в количествах ресурса в секунду
}

// BuildActionCharacteristics описывает характеристики строительства.
type BuildActionCharacteristics struct {
	Builder          int64  `json:"builder"`           // Идентификатор объекта, выполняющего строительство
	ObjectType       string `json:"object"`            // Тип объекта, который строится
	ConstructionTime int    `json:"construction_time"` // Время строительства  в секундах
	Place            Hex    `json:"place"`             // координаты на арене где будет произвоиться строителство
}

// AttackActionCharacteristics описывает характеристики атаки.
type AttackActionCharacteristics struct {
	Atacker  int64           `json:"atacker"`  // Идентификатор атакующего объекта
	Defenser int64           `json:"defenser"` // Идентификатор защищающегося объекта
	Damage   decimal.Decimal `json:"damage"`   // Урон от атаки
}

// Hex представляет точку на плоскости.
type Hex struct {
	Q float64 `json:"q"` // Координата q
	R float64 `json:"r"` // Координата r
}

// Obstacles - структура для хранения координат препятствий
type Obstacles struct {
	Coordinate []Hex
}
