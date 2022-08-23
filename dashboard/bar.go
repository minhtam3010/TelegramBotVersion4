package dashboard

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

// func main() {
// 	m := pdf.NewMaroto(consts.Portrait, consts.A4)

// 	m.SetPageMargins(10, 10, 10)
// 	buildHeading(m)

// 	buildChart(m)
// 	// tableHeading := []string{"Fruit", "Description", "Price", "A", "B", "C", "D", "E"}
// 	// contents := [][]string{{"Golang", "playground", "2", "A", "B", "C", "D", "E"}, {"Python", "Easy to Learn", "2", "A", "B", "C", "D", "E"}}
// 	// buildDatabaseList(m, "Inserted", tableHeading, contents)
// 	// // tableHeading = []string{"Fruit", "Description", "Price"}
// 	// // contents = [][]string{{"a", "B", "2"}, {"A", "C", "2"}, {"A", "M", "2"}}
// 	// // buildDatabaseList(m, "Updated", tableHeading, contents)
// 	err := m.OutputFileAndClose("myFile.pdf")
// 	if err != nil {
// 		fmt.Println("Could not save PDF: ", err)
// 		os.Exit(1)
// 	}

// 	fmt.Println("PDF saved successfully")
// }

func BuildHeading(m pdf.Maroto) {
	m.RegisterHeader(func() {
		m.Row(50, func() {
			m.Col(12, func() {
				err := m.FileImage("unnamed.png", props.Rect{
					Center:  true,
					Percent: 75,
				})
				if err != nil {
					fmt.Println("Image file was not loaded: ", err)
				}
			})
		})
	})

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Prepared for you by the GiaoLangEducation", props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Center,
				Color: GetDarkPurpleColor(),
			})
		})
	})
}

func buildDatabaseList(m pdf.Maroto, tableName string, tableHeading []string, contents [][]string) {
	lightPurpleColor := getLightPurpleColor()
	m.SetBackgroundColor(getTealColor())

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(tableName, props.Text{
				Top:    2,
				Size:   13,
				Color:  color.NewWhite(),
				Family: consts.Arial,
				Style:  consts.Bold,
				Align:  consts.Center,
			})
		})
	})

	m.SetBackgroundColor(color.NewWhite())

	var a []uint
	var i int
	for i = 1; i <= 8; i++ {
		a = append(a, 12/6)
	}

	m.TableList(tableHeading, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: a,
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: a,
		},
		Align:                consts.Left,
		AlternatedBackground: &lightPurpleColor,
		HeaderContentSpace:   1,
		Line:                 false,
	})
}

func drawChart(total []int, tableName string) {
	graph := chart.BarChart{
		Title: "GiaoLang Chart " + "--- " +  strings.ToUpper(tableName) + " ---",
		TitleStyle: chart.Style{
			FontColor: drawing.ColorBlue,
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top: 30,
			},
		},
		Height:   256,
		BarWidth: 50,
		Bars: []chart.Value{
			{Value: float64(total[0]), Label: "INSERT"},
			{Value: float64(total[1]), Label: "UPDATE"},
			{Value: float64(total[2]), Label: "DELETE"},
		},
	}
	f, err := os.Create("images/result_" + tableName + ".png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	graph.Render(chart.PNG, f)
}

func BuildChart(m pdf.Maroto, total []int, tableName string) {
	drawChart(total, tableName)
	m.Row(50, func() {
		m.Col(12, func() {
			err := m.FileImage("images/result_" + tableName + ".png", props.Rect{
				Center: true,
				Percent: 100,
			})
			fmt.Println(err, tableName)
			if err != nil {
				fmt.Println("Not found any images: ", err)
			}
		})
	})
}

func getTealColor() color.Color {
	return color.Color{
		Red:   3,
		Green: 166,
		Blue:  166,
	}
}

func GetDarkPurpleColor() color.Color {
	return color.Color{
		Red:   88,
		Green: 80,
		Blue:  99,
	}
}

func getLightPurpleColor() color.Color {
	return color.Color{
		Red:   210,
		Green: 200,
		Blue:  230,
	}
}
