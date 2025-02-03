package game

import (
	"container/heap"
	"cyber/internal/models"
	storage "cyber/internal/storage"
	"fmt"
	"math"
)

type PathNode struct {
	Coordinate Hex
	Cost       float64
	Priority   float64
	Index      int
}

type Hex struct {
	Q float64 `json:"q"` // Координата X
	R float64 `json:"r"` // Координата Y
}

// Функция расчета стоимости перехода между 2 гексами
func (h *Hex) Cost(toPosition models.Hex) float64 {
	return math.Abs(h.Q-toPosition.Q) + math.Abs(h.R-toPosition.R) + math.Abs(h.Q+h.R-toPosition.Q-toPosition.R)
}

// Heuristic возвращает эвристическую оценку расстояния между двумя
// гексами, рассчитанную по формуле манхетенского расстояния
func (h *Hex) Heuristic(toPosition Hex) float64 {
	return (math.Abs(h.Q-toPosition.Q) + math.Abs(h.R-toPosition.R) + math.Abs(h.Q+h.R-toPosition.Q-toPosition.R)) / 2
}

// PriorityQueue реализует приоритетную очередь для PathNode.
type PriorityQueue []*PathNode

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*PathNode)
	node.Index = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	node.Index = -1
	*pq = old[0 : n-1]
	return node
}

// Функция получает из БД координаты всех объектов и маркирует их как препятсвия
// TODO: лучше разделить на 2 функции. 1 получает данные 2 - маркирует
func obstaclesMaper(areaID int64) (map[Hex]bool, error) {
	db, err := storage.New()
	if err != nil {
		return nil, fmt.Errorf("error db connection: %v\n", err)
	}
	obstacles, err := db.GetObstacles(areaID)
	if err != nil {
		return nil, err
	}
	obstaclesMap := make(map[Hex]bool, len(obstacles))
	for _, v := range obstacles {
		obstaclesMap[v] = true
	}
	return obstaclesMap, nil
}

// AStar находит кратчайший путь между двумя гексами.
func AStar(start, goal Hex, areaID int64) bool {
	obstacles, err := obstaclesMaper(areaID)
	if err != nil {
		return false
	}

	// Проверяем, не является ли цель занятой клеткой
	if obstacles[goal] {
		return false
	}

	frontier := make(PriorityQueue, 0)
	heap.Init(&frontier)

	startNode := &PathNode{
		Coordinate: start,
		Cost:       0,
		Priority:   start.Heuristic(goal),
	}
	heap.Push(&frontier, startNode)

	cameFrom := make(map[Hex]Hex)
	costSoFar := make(map[Hex]float64)
	cameFrom[start] = start
	costSoFar[start] = 0

	for frontier.Len() > 0 {
		current := heap.Pop(&frontier).(*PathNode)

		// Если достигли цели, возвращаем true
		if current.Coordinate == goal {
			return true
		}

		// Анализируем соседей
		for _, neighbor := range current.Coordinate.Neighbours() {
			if obstacles[neighbor] {
				continue
			}

			newCost := costSoFar[current.Coordinate] + current.Coordinate.Cost(neighbor)
			if cost, ok := costSoFar[neighbor]; !ok || newCost < cost {
				costSoFar[neighbor] = newCost
				priority := newCost + neighbor.Heuristic(goal)
				heap.Push(&frontier, &PathNode{
					Coordinate: neighbor,
					Cost:       newCost,
					Priority:   priority,
				})
				cameFrom[neighbor] = current.Coordinate
			}
		}
	}

	// Если путь не найден, возвращаем false
	return false
}

// Neighbours возвращает соседей гекса.
func (h *Hex) Neighbours() []Hex {
	return []Hex{
		//у каждого гекса может быть 6 соседей. Координаты которых заранее известны. Смотри файл images/neighboors 2 system
		{h.Q + 1, h.R},
		{h.Q - 1, h.R},
		{h.Q, h.R + 1},
		{h.Q, h.R - 1},
		{h.Q + 1, h.R - 1},
		{h.Q - 1, h.R + 1},
	}
}
