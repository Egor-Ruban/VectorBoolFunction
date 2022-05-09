package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func nonDegenerateTest() {
	name := "NonDegenerateExperiments"
	f := excelize.NewFile()
	// Set value of a cell.
	f.SetCellValue("Sheet1", "A1", "n")
	f.SetCellValue("Sheet1", "B1", "m")
	f.SetCellValue("Sheet1", "C1", "sum")
	f.SetCellValue("Sheet1", "D1", "count")
	f.SetCellValue("Sheet1", "E1", "p")
	res := make([][]Test, 40)
	for i := range res {
		res[i] = make([]Test, 40)
	}
	for i := 0; i < 10000000; i++ {
		for n := 6; n < 7; n++ {
			for m := n; m <= 6; m++ {
				bf, _ := newRandomVBF(n, m)
				if bf.isNonDegenerate() {
					res[n][m].sum++
				}
				res[n][m].count++
				//i := strconv.Itoa(1 + ((n - 1) * 30) + m)
				i := "2"
				f.SetCellValue("Sheet1", "A"+i, n)
				f.SetCellValue("Sheet1", "B"+i, m)
				f.SetCellValue("Sheet1", "C"+i, res[n][m].sum)
				f.SetCellValue("Sheet1", "D"+i, res[n][m].count)
				f.SetCellValue("Sheet1", "E"+i, float64(res[n][m].sum)/float64(res[n][m].count))
			}
		}
	}
	if err := f.SaveAs(name + ".xlsx"); err != nil {
		fmt.Println(err)
	}
}
