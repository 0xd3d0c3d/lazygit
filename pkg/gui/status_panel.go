package gui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) refreshStatus() error {
	v, err := gui.g.View("status")
	if err != nil {
		panic(err)
	}
	// for some reason if this isn't wrapped in an update the clear seems to
	// be applied after the other things or something like that; the panel's
	// contents end up cleared
	gui.g.Update(func(*gocui.Gui) error {
		v.Clear()
		pushables, pullables := gui.GitCommand.UpstreamDifferenceCount()
		fmt.Fprint(v, "↑"+pushables+"↓"+pullables)
		branches := gui.State.Branches
		if err := gui.updateHasMergeConflictStatus(); err != nil {
			return err
		}
		if gui.State.HasMergeConflicts {
			fmt.Fprint(v, utils.ColoredString(" (merging)", color.FgYellow))
		}

		if len(branches) == 0 {
			return nil
		}
		branch := branches[0]
		name := utils.ColoredString(branch.Name, branch.GetColor())
		repo := utils.GetCurrentRepoName()
		fmt.Fprint(v, " "+repo+" → "+name)
		return nil
	})

	return nil
}

func (gui *Gui) handleCheckForUpdate(g *gocui.Gui, v *gocui.View) error {
	gui.Updater.CheckForNewUpdate(gui.onUserUpdateCheckFinish, true)
	return gui.createMessagePanel(v, "", gui.Tr.SLocalize("CheckingForUpdates"))
}

func (gui *Gui) handleStatusSelect(g *gocui.Gui, v *gocui.View) error {
	dashboardString := fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s\n\n%s",
		lazygitTitle,
		"Keybindings: https://github.com/jesseduffield/lazygit/blob/master/docs/Keybindings.md",
		"Config Options: https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md",
		"Tutorial: https://www.youtube.com/watch?v=VDXvbHZYeKY",
		"Raise an Issue: https://github.com/jesseduffield/lazygit/issues",
	)

	if err := gui.renderString(g, "main", dashboardString); err != nil {
		return err
	}
	return gui.renderGlobalOptions()
}

func (gui *Gui) handleOpenConfig(g *gocui.Gui, v *gocui.View) error {
	return gui.openFile(gui.Config.GetUserConfig().ConfigFileUsed())
}

func (gui *Gui) handleEditConfig(g *gocui.Gui, v *gocui.View) error {
	filename := gui.Config.GetUserConfig().ConfigFileUsed()
	return gui.editFile(filename)
}
