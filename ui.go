package main

import (
	"fmt"
	"os"

	"github.com/gizak/termui"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

func InitUI(tally *Tally) {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
		os.Exit(0)
	})

	topSrcTable := termui.NewTable()
	topSrcTable.BorderLabel = "Top Source"
	topSrcTable.FgColor = termui.ColorWhite
	topSrcTable.BgColor = termui.ColorDefault
	topSrcTable.Separator = false
	topSrcTable.Border = true
	topSrcTable.Height = 20

	topDstTable := termui.NewTable()
	topDstTable.BorderLabel = "Top Destination"
	topDstTable.FgColor = termui.ColorWhite
	topDstTable.BgColor = termui.ColorDefault
	topDstTable.Separator = false
	topDstTable.Border = true
	topDstTable.Height = 20

	sparkLineSent := termui.NewSparkline()
	sparkLineSent.LineColor = termui.ColorBlue
	sparkLineSent.Data = []int{}
	sparkLineSent.Title = "Sent: 0"
	sparkLineReceived := termui.NewSparkline()
	sparkLineReceived.LineColor = termui.ColorRed
	sparkLineReceived.Data = []int{}
	sparkLineReceived.Title = "Rcvd: 0"
	sparkLineTotalBytes := termui.NewSparklines(sparkLineSent, sparkLineReceived)
	sparkLineTotalBytes.Border = true
	sparkLineTotalBytes.BorderLabel = "Total Bytes"
	sparkLineTotalBytes.Height = 6

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(12, 0, sparkLineTotalBytes)),
		termui.NewRow(
			termui.NewCol(6, 0, topSrcTable),
			termui.NewCol(6, 0, topDstTable)))

	var lastTotalBytesSent, lastTotalBytesReceived int
	termui.Handle("/timer/1s", func(e termui.Event) {
		termHeight, _ := terminal.Height()
		topSrcTable.Height = int(termHeight) - 10
		topDstTable.Height = int(termHeight) - 10

		srcRows := [][]string{
			[]string{"Host", "Bytes"},
		}
		dstRows := [][]string{
			[]string{"Host", "Bytes"},
		}

		tally.Mux.Lock()
		hosts := tally.TopSrcHosts()
		for _, host := range hosts {
			byteCount := fmt.Sprintf("%v", tally.BytesBySrcHost[host])
			srcRows = append(srcRows, []string{host, byteCount})
		}
		hosts = tally.TopDstHosts()
		for _, host := range hosts {
			byteCount := fmt.Sprintf("%v", tally.BytesByDstHost[host])
			dstRows = append(dstRows, []string{host, byteCount})
		}
		tally.Mux.Unlock()

		topSrcTable.Rows = srcRows
		topDstTable.Rows = dstRows

		// hack to get around lib caching
		topSrcTable.FgColors = make([]termui.Attribute, 0)
		topSrcTable.BgColors = make([]termui.Attribute, 0)
		topDstTable.FgColors = make([]termui.Attribute, 0)
		topDstTable.BgColors = make([]termui.Attribute, 0)
		topSrcTable.Analysis()
		topDstTable.Analysis()
		topSrcTable.FgColors[0] = termui.ColorBlue
		topDstTable.FgColors[0] = termui.ColorBlue

		tickTotalBytesSent := tally.TotalBytesSent - lastTotalBytesSent
		tickTotalBytesReceived := tally.TotalBytesReceived - lastTotalBytesReceived
		lastTotalBytesSent = tally.TotalBytesSent
		lastTotalBytesReceived = tally.TotalBytesReceived
		sparkLineSent.Data = append([]int{tickTotalBytesSent}, sparkLineSent.Data...)
		sparkLineReceived.Data = append([]int{tickTotalBytesReceived}, sparkLineReceived.Data...)
		sparkLineTotalBytes.Lines[0].Title = fmt.Sprintf("Sent: %v", tally.TotalBytesSent)
		sparkLineTotalBytes.Lines[1].Title = fmt.Sprintf("Rcvd: %v", tally.TotalBytesReceived)
		sparkLineTotalBytes.Lines[0].Data = sparkLineSent.Data
		sparkLineTotalBytes.Lines[1].Data = sparkLineReceived.Data

		termui.Body.Align()
		termui.Render(termui.Body)
	})
	termui.Loop()
}
