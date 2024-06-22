package batch

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/xuri/excelize/v2"
	"testing"
)

// Document Link: https://blog.csdn.net/qq_23118345/article/details/126706626
// https://github.com/qax-os/excelize
func TestExcelStreamWrite(t *testing.T) {
	t.Run("Test Stream Write 490000", func(t *testing.T) {
		file := excelize.NewFile()
		//Set Header
		err := file.SetSheetName("Sheet1", "Table1")
		if err != nil {
			g.Log().Errorf(context.Background(), err.Error())
			return
		}
		//Create Stream Writer
		writer, err := file.NewStreamWriter("Table1")
		//Update Width Height
		err = writer.SetColWidth(1, 15, 12)
		if err != nil {
			g.Log().Errorf(context.Background(), err.Error())
			return
		}
		//Set Header
		err = writer.SetRow("A1", []interface{}{"TestRow1", "TestRow2", "TestRow3", "TestRow4", "TestRow5", "TestRow6", "TestRow7", "TestRow8", "TestRow9", "TestRow10", "TestRow11", "TestRow12", "TestRow13", "TestRow14", "TestRow15"})
		if err != nil {
			g.Log().Errorf(context.Background(), err.Error())
			return
		}
		for i := 1; i <= 490000; i++ {
			//Index to Table column
			cell, _ := excelize.CoordinatesToCellName(1, i+1)
			//Append Data
			_ = writer.SetRow(cell, []interface{}{"TestData1", "TestData2", "TestData3", "TestData4", "TestData5", "TestData6", "TestData7", "TestData8", "TestData9", "TestData10", "TestData11", "TestData12", "TestData13", "TestData14", "TestData15"})
		}
		//End Stream Writer
		err = writer.Flush()
		if err != nil {
			g.Log().Errorf(context.Background(), err.Error())
			return
		}
		err = file.SaveAs("test01.xlsx")
		if err != nil {
			g.Log().Errorf(context.Background(), err.Error())
			return
		}
		//_ = os.Remove("test01.xlsx")
	})
	t.Run("Test Non Stream Write 490000", func(t *testing.T) {
		file := excelize.NewFile()
		//Set Table Name
		_ = file.SetSheetName("Sheet1", "Table1")
		//Update Width Height
		_ = file.SetColWidth("Table1", "A", "O", 12)
		//Set Header
		_ = file.SetSheetRow("Table1", "A1", &[]interface{}{"TestRow1", "TestRow2", "TestRow3", "TestRow4", "TestRow5", "TestRow6", "TestRow7", "TestRow8", "TestRow9", "TestRow10", "TestRow11", "TestRow12", "TestRow13", "TestRow14", "TestRow15"})
		for i := 1; i <= 1000000; i++ {
			//Index to Table column
			cell, _ := excelize.CoordinatesToCellName(1, i+1)
			//Append Data
			_ = file.SetSheetRow("Table1", cell, &[]interface{}{"TestData1", "TestData2", "TestData3", "TestData4", "TestData5", "TestData6", "TestData7", "TestData8", "TestData9", "TestData10", "TestData11", "TestData12", "TestData13", "TestData14", "TestData15"})
		}
		//Save File
		_ = file.SaveAs("test01.xlsx")
	})
}
