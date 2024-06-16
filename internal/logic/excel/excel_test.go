package excel

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/xuri/excelize/v2"
	"os"
	"testing"
)

func TestExcelStreamWrite(t *testing.T) {
	t.Run("Test Stream Write 490000", func(t *testing.T) {
		file := excelize.NewFile()
		//设置表名
		err := file.SetSheetName("Sheet1", "表1")
		if err != nil {
			g.Log().Errorf(context.Background(), err.Error())
			return
		}
		//创建流式写入
		writer, err := file.NewStreamWriter("表1")
		//修改列宽
		err = writer.SetColWidth(1, 15, 12)
		if err != nil {
			g.Log().Errorf(context.Background(), err.Error())
			return
		}
		//设置表头
		err = writer.SetRow("A1", []interface{}{"测试列名1", "测试列名2", "测试列名3", "测试列名4", "测试列名5", "测试列名6", "测试列名7", "测试列名8", "测试列名9", "测试列名10", "测试列名11", "测试列名12", "测试列名13", "测试列名14", "测试列名15"})
		if err != nil {
			g.Log().Errorf(context.Background(), err.Error())
			return
		}
		for i := 1; i <= 490000; i++ {
			//索引转单元格坐标
			cell, _ := excelize.CoordinatesToCellName(1, i+1)
			//添加的数据
			_ = writer.SetRow(cell, []interface{}{"测试数据1", "测试数据2", "测试数据3", "测试数据4", "测试数据5", "测试数据6", "测试数据7", "测试数据8", "测试数据9", "测试数据10", "测试数据11", "测试数据12", "测试数据13", "测试数据14", "测试数据15"})
		}
		//结束流式写入
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
		_ = os.Remove("test01.xlsx")
	})
}
