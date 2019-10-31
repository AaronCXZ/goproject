package core

import (
	"fmt"
)

/*
	AOI区域管理模块
*/

type AOIManager struct {
	//	区域的左边界坐标
	MinX int
	//	区域的右边界坐标
	MaxX int
	//	X方向格子的数量
	CntsX int
	//	区域的上边界坐标
	MinY int
	//	区域的下边界坐标
	MaxY int
	//	Y方向格子的数量
	CntsY int
	//	当前区域中有哪些格子map， key=格子的ID，value=格子的对象
	grids map[int]*Grid
}

//初始化
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiManager := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}
	//给AOI初始化区域内所有的格子进行编号和初始化
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			//计算格子的ID，根据x，y编号
			gid := y*cntsX + x
			//	初始化gid格子
			aoiManager.grids[gid] = NewGrid(gid,
				aoiManager.MinX+x*aoiManager.gridWidth(),
				aoiManager.MinX+(x+1)*aoiManager.gridWidth(),
				aoiManager.MinY+y*aoiManager.gridWidth(),
				aoiManager.MinY+(y+1)*aoiManager.gridWidth())
		}
	}
	return aoiManager
}

//得到每个格子在X轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

//得到每个格子在Y轴方向的长度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

//打印格子的信息
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager INFO:\n\tMinX=%d,\n\tMaxX=%d\n\tCntsX=%d\n\tMinY=%d\n\tMaxY=%d\n\tCntsY=%d\n\n",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}

//根据GID得到当前GID的九宫格的GID集合
func (m *AOIManager) GetSurroundGridsByGid(gID int) (grids []*Grid) {
	//	判断当前gID时候在AOIManager中
	if _, ok := m.grids[gID]; !ok {
		return
	}
	//	初始化grids返回值切片,将gID本身加入九宫格切片
	grids = append(grids, m.grids[gID])

	//	需要判断gID的左边是否有格子，右边是否有格子
	//	需要通过gID得到当前格子的X轴坐标idx=id % nx
	idx := gID % m.CntsX

	//	判断idx编号是否左边还有格子，如果有放在gridsX集合中
	if idx > 0 {
		grids = append(grids, m.grids[gID-1])
	}

	//	判断idx编号是否右边还有格子，如果有放在gridsX集合中
	if idx < m.CntsX-1 {
		grids = append(grids, m.grids[gID+1])
	}
	//	遍历gidsX集合中每个格子的gid
	gridX := make([]int, 0, len(grids))
	for _, v := range grids {
		gridX = append(gridX, v.GID)
	}
	for _, v := range gridX {
		idy := v / m.CntsX
		//	gid上边是否还有格子
		if idy > 0 {
			grids = append(grids, m.grids[v-m.CntsX])
		}
		//	git下边是否还有格子
		if idy < m.CntsY-1 {
			grids = append(grids, m.grids[v+m.CntsX])
		}
	}
	return
}

//通过纵横坐标得到当前GID格子编号
func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) - m.MinY) / m.gridLength()
	return idy*m.CntsX + idx
}

//根据纵横坐标得到周边九宫格内全部的PalyerIDs
func (m *AOIManager) GetPidsByPos(x, y float32) (playerIDs []int) {
	//	得到当前玩家的GID格子id
	gID := m.GetGidByPos(x, y)
	//	通过GID得到周边九宫格消息
	grids := m.GetSurroundGridsByGid(gID)
	//	将九宫格的信息里的全部Player的ID放在playerIDs里
	for _, grid := range grids {
		playerIDs = append(playerIDs, grid.GetPlayerIDs()...)
	}
	return
}

//添加一个PlayerID到一个格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.grids[gID].Add(pID)
}

//移除一个格子的PlayerID
func (m *AOIManager) RemovePidFromGrid(pID, gID int) {
	m.grids[gID].Remove(pID)
}

//通过GID获取全部的PlayerID
func (m *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	playerIDs = m.grids[gID].GetPlayerIDs()
	return
}

//通过坐标将PlayerID添加到一个格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.grids[gID]
	grid.Add(pID)
}

//通过坐标将一个PlayerID从格子中移除
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.grids[gID]
	grid.Remove(gID)
}
