package commands

import (
	"fmt"
	"os"

	uipack "github.com/uipack-io/cli"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func Describe(p *uipack.Package) {

	re := lipgloss.NewRenderer(os.Stdout)

	const (
		purple    = lipgloss.Color("99")
		gray      = lipgloss.Color("245")
		lightGray = lipgloss.Color("241")
	)

	var (
		// HeaderStyle is the lipgloss style used for the table headers.
		HeaderStyle = re.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
		// CellStyle is the base lipgloss style used for the table rows.
		CellStyle = re.NewStyle().Padding(0, 1).Width(14)
		// OddRowStyle is the lipgloss style used for odd-numbered table rows.
		OddRowStyle = CellStyle.Foreground(gray)
		// EvenRowStyle is the lipgloss style used for even-numbered table rows.
		EvenRowStyle = CellStyle.Foreground(lightGray)
		// BorderStyle is the lipgloss style used for the table border.
		BorderStyle = lipgloss.NewStyle().Foreground(purple)
	)

	combinations := p.Metadata.GenerateModeCombinations()

	rows := make([][]string, 1+len(p.Metadata.Variables))
	for i := range rows {
		rows[i] = make([]string, 1+len(combinations))
	}
	for i, combination := range combinations {
		vid := uipack.Variant(0)
		variant := make(map[string]string)
		for mi, mv := range combination {
			mode := p.Metadata.Modes[mi]
			variant[mode.Name] = mv.Name
			vid = vid.SetMode(mode.Identifier, uipack.Uint4(mv.Identifier))
		}
		bundle := p.GetBundle(vid)
		header := fmt.Sprintf("0x%02x", vid)
		for _, mv := range combination {
			header = header + "\n" + mv.Name
		}
		rows[0][1+i] = header
		for iv, variable := range p.Metadata.Variables {
			if i == 0 {
				header := variable.Name

				rows[iv+1][0] = header
			}
			value := bundle.Values[variable.Identifier]
			switch v := value.(type) {
			case uipack.Color:
				style := lipgloss.NewStyle().
					Background(lipgloss.Color("#" + v.ToHexString()[2:]))
				rows[iv+1][i+1] = fmt.Sprint(style.Render("  ")) + " #" + v.ToHexString()
			case uipack.TextStyle:
				rows[iv+1][i+1] = fmt.Sprint(v)
			case uipack.LinearGradient:
				rows[iv+1][i+1] = fmt.Sprint(v)
			case uipack.RadialGradient:
				rows[iv+1][i+1] = fmt.Sprint(v)
			case uipack.TypeDefinition:
				rows[iv+1][i+1] = fmt.Sprint(v)
			case uipack.Label:
				rows[iv+1][i+1] = fmt.Sprint(v)
			case uipack.Offset:
				rows[iv+1][i+1] = fmt.Sprint(v)
			case uipack.Instance:
				rows[iv+1][i+1] = fmt.Sprint(v)
			case uipack.Radius:
				rows[iv+1][i+1] = fmt.Sprint(v)
			case uipack.BorderRadius:
				rows[iv+1][i+1] = fmt.Sprint(v)
			case uint64:
				rows[iv+1][i+1] = fmt.Sprint(v)
			case bool:
				rows[iv+1][i+1] = fmt.Sprint(v)
			case float64:
				rows[iv+1][i+1] = fmt.Sprint(v)
			case string:
				rows[iv+1][i+1] = fmt.Sprint(v)

			default:
			}
		}
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(BorderStyle).
		StyleFunc(func(row, col int) lipgloss.Style {
			var style lipgloss.Style

			switch {
			case row == 0:
				return HeaderStyle
			case row%2 == 0:
				style = EvenRowStyle
			default:
				style = OddRowStyle
			}
			return style
		}).Rows(rows...)

	fmt.Println(t)
}
