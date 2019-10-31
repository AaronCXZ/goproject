package core

import (
	"fmt"
	"sync"
)

/*
	一个AOI地图中的个子类型
*/

type Grid struct {
	//	格子ID
	GID int
	//  格子左边的边界坐标
	MinX int
	//  格子右边的边界坐标
	MaxX int
	//  格子上边的边界坐标
	MinY int
	//  格子下边的边界坐标
	MaxY int
	//  当前格子内玩家或者物体的集合
	playerIDs map[int]bool
	//  保护集合的锁
	pIDLock sync.RWMutex
}

//初始化
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

//给格子添加一个玩家
func (g *Grid) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true
}

//从格子删除一个玩家
func (g *Grid) Remove(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

//得到当前格子中的所有玩家
func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k, _ := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}
	return playerIDs
}

//调试使用--打印出格子的基本信息
func (g *Grid) String() string {
	return fmt.Sprintf("Grid %d INFO:\n\tid=%d\n\tminX=%d\n\tmaxX=%d\n\tMinY=%d\n\tmaxY=%d\n",
		g.GID, g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY)
}
