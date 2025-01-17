package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	input "github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func main() {
	termenv.ClearScreen()
	termenv.SetWindowTitle("typioca")

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	termenv.Reset()
	println("bye!")
}

func initTimerBasedTest(settings TimerBasedTestSettings) TimerBasedTest {
	return TimerBasedTest{
		settings: settings,
		timer: myTimer{
			timer:     timer.NewWithInterval(settings.timeSelections[settings.timeCursor], time.Second),
			duration:  settings.timeSelections[settings.timeCursor],
			isRunning: false,
			timedout:  false,
		},
		base: TestBase{
			wordsToEnter: NewGenerator().Generate(settings.wordListSelections[settings.wordListCursor].key),
			inputBuffer:  make([]rune, 0),
			rawInputCnt:  0,
			mistakes: mistakes{
				mistakesAt:     make(map[int]bool, 0),
				rawMistakesCnt: 0,
			},
			cursor: 0,
		},
		completed: false,
	}
}

func initWordCountBasedTest(settings WordCountBasedTestSettings) WordCountBasedTest {
	generator := NewGenerator()
	generator.Count = settings.wordCountSelections[settings.wordCountCursor]
	return WordCountBasedTest{
		settings: settings,
		stopwatch: myStopWatch{
			stopwatch: stopwatch.New(),
			isRunning: false,
		},
		base: TestBase{
			wordsToEnter: strings.TrimSpace(generator.Generate(settings.wordListSelections[settings.wordListCursor].key)),
			inputBuffer:  make([]rune, 0),
			rawInputCnt:  0,
			mistakes: mistakes{
				mistakesAt:     make(map[int]bool, 0),
				rawMistakesCnt: 0,
			},
			cursor: 0,
		},
		completed: false,
	}
}

func initSentenceCountBasedTest(settings SentenceCountBasedTestSettings) SentenceCountBasedTest {
	generator := NewGenerator()
	generator.Count = 40
	generator.Count = settings.sentenceCountSelections[settings.sentenceCountCursor]
	return SentenceCountBasedTest{
		settings: settings,
		stopwatch: myStopWatch{
			stopwatch: stopwatch.New(),
			isRunning: false,
		},
		base: TestBase{
			wordsToEnter: strings.TrimSpace(generator.Generate(settings.sentenceListSelections[settings.sentenceListCursor].key)),
			inputBuffer:  make([]rune, 0),
			rawInputCnt:  0,
			mistakes: mistakes{
				mistakesAt:     make(map[int]bool, 0),
				rawMistakesCnt: 0,
			},
			cursor: 0,
		},
		completed: false,
	}
}

func initTimerBasedTestSelection() TimerBasedTestSettings {
	return TimerBasedTestSettings{
		timeSelections: []time.Duration{time.Second * 120, time.Second * 60, time.Second * 30, time.Second * 15},
		timeCursor:     2,
		wordListSelections: []WordListSelection{
			{
				key:  "dorian-gray-words",
				show: "dorian-gray",
			},
			{
				key:  "frankenstein-words",
				show: "frankenstein",
			},
			{
				key:  "common-words",
				show: "common-words",
			},
			{
				key:  "pride-and-prejudice-words",
				show: "pride-and-prejudice",
			},
			{
				key:  "dorian-gray-sentences",
				show: "dorian-gray-sentences",
			},
			{
				key:  "frankenstein-sentences",
				show: "frankenstein-sentences",
			},
			{
				key:  "pride-and-prejudice-sentences",
				show: "pride-and-prejudice-sentences",
			},
		},
		wordListCursor: 2,
		cursor:         0,
	}
}

func initWordCountBasedTestSelection() WordCountBasedTestSettings {
	return WordCountBasedTestSettings{
		wordCountSelections: []int{100, 50, 25, 10},
		wordCountCursor:     2,
		wordListSelections: []WordListSelection{
			{
				key:  "dorian-gray-words",
				show: "dorian-gray",
			},
			{
				key:  "frankenstein-words",
				show: "frankenstein",
			},
			{
				key:  "common-words",
				show: "common-words",
			},
			{
				key:  "pride-and-prejudice-words",
				show: "pride-and-prejudice",
			},
		},
		wordListCursor: 2,
		cursor:         0,
	}
}

func initSentenceCountBasedTestSelection() SentenceCountBasedTestSettings {
	return SentenceCountBasedTestSettings{
		sentenceCountSelections: []int{30, 15, 5, 1},
		sentenceCountCursor:     2,
		sentenceListSelections: []WordListSelection{
			{
				key:  "dorian-gray-sentences",
				show: "dorian-gray",
			},
			{
				key:  "frankenstein-sentences",
				show: "frankenstein",
			},
			{
				key:  "pride-and-prejudice-sentences",
				show: "pride-and-prejudice",
			},
		},
		sentenceListCursor: 1,
		cursor:             0,
	}
}

func initMainMenu() MainMenu {
	return MainMenu{
		selections: []MainMenuSelection{
			initTimerBasedTestSelection(),
			initWordCountBasedTestSelection(),
			initSentenceCountBasedTestSelection(),
		},
		cursor: 0,
	}
}

func initialModel() model {
	profile := termenv.ColorProfile()
	fore := termenv.ForegroundColor()

	return model{
		state: initMainMenu(),
		styles: Styles{
			correct: func(str string) termenv.Style {
				return termenv.String(str).Foreground(fore)
			},
			toEnter: func(str string) termenv.Style {
				return termenv.String(str).Foreground(fore).Faint()
			},
			mistakes: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("1")).Underline()
			},
			cursor: func(str string) termenv.Style {
				return termenv.String(str).Reverse().Bold()
			},
			runningTimer: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("2"))
			},
			stoppedTimer: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("2")).Faint()
			},
			greener: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("6")).Faint()
			},
			faintGreen: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("10")).Faint()
			},
		},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		input.Blink, //we probably don't need this anymore
	)
}
