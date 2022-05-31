package gui

import (
	"fmt"
	"golang.org/x/term"
	"log"
	"math"
)

const DefaultWidth = 80
const DefaultPadding = 8
const FullWidth = math.MaxInt

type Line struct {
	Content string
	Indent  int    `default:0`
	Justify string `default:"left"` // "left" | "right" | "center"
}

type Jumbotron struct {
	Header    *Line
	Footer    *Line
	Lines     []Line
	Width     int
	XPadding  int
	YPadding  int
	AutoScale bool
}

func NewJumbotron(header *Line, footer *Line) Jumbotron {
	return Jumbotron{AutoScale: true, Width: DefaultWidth, Header: header, Footer: footer}
}

func (j *Jumbotron) AddLine(l Line) {
	if len(l.Content) <= j.Width {
		j.Lines = append(j.Lines, l)
	} else {
		sections := int(math.Floor(float64(len(l.Content))/float64(j.Width))) + 1

		for i := 0; i < sections; i++ {
			lowerBound := i * j.Width
			upperBound := int(math.Min(float64(len(l.Content)), float64(lowerBound+j.Width)))

			content := l.Content[lowerBound:upperBound]

			if upperBound != len(l.Content) {
				content += "-"
			}

			j.Lines = append(j.Lines, Line{Content: content, Justify: l.Justify})
		}
	}
}

func (j *Jumbotron) AddBlankLine() {
	j.Lines = append(j.Lines, Line{})
}

func (j *Jumbotron) SetHeader(header *Line) {
	j.Header = header
}

func (j *Jumbotron) SetFooter(footer *Line) {
	j.Footer = footer
}

func (j *Jumbotron) SetMaxWidth(width int) {
	termWidth := getTerminalDim().Width - (j.XPadding * 2) - DefaultPadding

	if !j.AutoScale || width <= termWidth {
		j.Width = width
	} else {
		j.Width = termWidth
	}
}

type Dimensions struct {
	Width  int
	Height int
}

func getTerminalDim() Dimensions {
	if !term.IsTerminal(0) {
		log.Fatal("Must run inside a terminal.")
	}
	width, height, err := term.GetSize(0)
	if err != nil {
		log.Fatal("Something went wrong getting terminal dimensions.")
	}

	return Dimensions{width, height}
}

func (j *Jumbotron) SetXPadding(padding int) {
	j.XPadding = padding
}

func (j *Jumbotron) SetYPadding(padding int) {
	j.YPadding = padding
}

func (j *Jumbotron) getMaxWidth() int {
	max := 0.0
	if j.Header != nil && j.Footer != nil {
		max = math.Max(float64(j.Header.GetWidth()), float64(j.Footer.GetWidth()))
	} else if j.Header != nil && j.Footer == nil {
		max = float64(j.Header.GetWidth())
	} else if j.Header == nil && j.Footer != nil {
		max = float64(j.Footer.GetWidth())
	}

	for _, l := range j.Lines {
		max = math.Max(float64(l.GetWidth()), max)
	}

	return int(max)
}

func (j *Jumbotron) printTopCap() {
	maxWidth := j.getMaxWidth()

	topBar := "╔══"

	for i := 0; i < (j.XPadding * 2); i++ {
		topBar += "═"
	}

	for i := 0; i < maxWidth; i++ {
		topBar = topBar + "═"
	}

	topBar += "╗"

	fmt.Println(topBar)
}

func (j *Jumbotron) printBottomCap() {
	maxWidth := j.getMaxWidth()

	topBar := "╚══"

	for i := 0; i < (j.XPadding * 2); i++ {
		topBar += "═"
	}

	for i := 0; i < maxWidth; i++ {
		topBar = topBar + "═"
	}

	topBar += "╝"

	fmt.Println(topBar)
}

func (j *Jumbotron) printMidCap() {
	maxWidth := j.getMaxWidth()

	topBar := "╠═"

	for i := 0; i < (j.XPadding * 2); i++ {
		topBar += "═"
	}

	for i := 0; i < maxWidth; i++ {
		topBar = topBar + "═"
	}

	topBar = topBar + "═╣"

	fmt.Println(topBar)
}

func (j *Jumbotron) printEmptyCap() {
	maxWidth := j.getMaxWidth()

	topBar := "║ "

	for i := 0; i < (j.XPadding * 2); i++ {
		topBar += " "
	}

	for i := 0; i < maxWidth; i++ {
		topBar = topBar + " "
	}

	topBar = topBar + " ║"

	fmt.Println(topBar)
}

func (j *Jumbotron) printYPadding() {
	for i := 0; i < j.YPadding; i++ {
		j.printEmptyCap()
	}
}

func (j *Jumbotron) Print() {
	w := j.getMaxWidth()

	j.printTopCap()

	j.printYPadding()

	if j.Header != nil {
		j.Header.Print(w, j.XPadding)
		j.printMidCap()
	}

	for _, l := range j.Lines {
		l.Print(w, j.XPadding)
	}

	if j.Footer != nil {
		j.printMidCap()
		j.Footer.Print(w, j.XPadding)
	}

	j.printYPadding()
	j.printBottomCap()
}

func (l *Line) getRightJustified(maxWidth int) string {
	line := ""
	for i := 0; i < (maxWidth - len(l.Content)); i++ {
		line += " "
	}
	line = line + l.Content
	return line
}

func (l *Line) getLeftJustified(maxWidth int) string {
	line := l.Content
	for i := 0; i < (maxWidth - len(l.Content)); i++ {
		line += " "
	}
	return line
}

func (l *Line) getCentered(maxWidth int) string {
	line := ""
	for i := 0; i < ((maxWidth - len(l.Content)) / 2); i++ {
		line += " "
	}
	line += l.Content
	for i := 0; i < ((maxWidth - len(l.Content)) / 2); i++ {
		line += " "
	}

	// Account for one-off error
	if (maxWidth-len(l.Content))%2 != 0 {
		line += " "
	}

	return line
}

func (l *Line) Print(maxWidth int, xpadding int) { // TODO add wrap length
	line := "║ "

	for i := 0; i < xpadding; i++ {
		line += " "
	}

	if l.Justify == "center" {
		line += l.getCentered(maxWidth)
	} else if l.Justify == "right" {
		line += l.getRightJustified(maxWidth)
	} else {
		line += l.getLeftJustified(maxWidth)
	}

	for i := 0; i < xpadding; i++ {
		line += " "
	}

	line = line + " ║"

	fmt.Println(line)
}

func (l *Line) GetWidth() int {
	return l.Indent + len(l.Content)
}
