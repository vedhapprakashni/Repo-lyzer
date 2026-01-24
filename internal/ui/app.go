package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/cache"
	"github.com/agnivo988/Repo-lyzer/internal/config"
	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
	stateMenu sessionState = iota
	stateInput
	stateLoading
	stateDashboard
	stateTree
	stateFileEdit
	stateSettings
	stateHelp
	stateHistory
	stateFavorites
	stateCompareInput
	stateCompareLoading
	stateCompareResult
	stateCloneInput
	stateCloning
)

// Message types for sub-models
type AnalyzeRepoMsg struct {
	repoName string
}

type BackToMenuMsg struct{}

type SwitchToInputMsg struct{}

type ErrorMsg error

type MainModel struct {
	state sessionState

	// Sub-models for different UI states
	menu             MenuModel
	input            InputModel
	loading          LoadingModel
	compareInput     CompareInputModel
	compareLoading   CompareLoadingModel
	compareResult    CompareResultModel
	settings         SettingsModel
	help             HelpModel
	history          HistoryModel
	favorites        FavoritesModel
	cloneInput       CloneInputModel
	cloning          CloningModel

	// Shared models
	dashboard DashboardModel
	tree      TreeModel
	fileEdit  FileEditModel

	// Shared state
	windowWidth  int
	windowHeight int
	cache        *cache.Cache
	appConfig    *config.AppSettings

	// Additional state fields
	historyCursor   int
	animTick        int
	compareStep     int
	compareInput1   string
	compareInput2   string
	inTokenInput    bool
	tokenInput      string
	settingsOption  string
	favoritesCursor int
	analysisType    string
	helpContent     string
	err             error
}
=======
	// Additional state fields
	historyCursor   int
	animTick        int
	compareStep     int
	compareInput1   string
	compareInput2   string
	inTokenInput    bool
	tokenInput      string
	settingsOption  string
	favoritesCursor int
	analysisType    string
	helpContent     string
	err             error
