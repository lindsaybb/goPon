package gopon

import "fmt"

type CpuDetail struct {
	Id	int		`json:"msanCpuDetailId"`	
	Cur	int		`json:"msanCpuDetailCurUsage"`
	Max	int		`json:"msanCpuDetailMaxUsage"`
	Min int		`json:"msanCpuDetailMinUsage"`
	Avg int		`json:"msanCpuDetailAvgUsage"`
}

func (c *CpuDetail) Tabwrite() string {
	return fmt.Sprintf("%d\t%d\t%d\t%d\t%d", c.Id, c.Cur, c.Max, c.Min, c.Avg)
}
