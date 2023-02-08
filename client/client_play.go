package client

import (
	"fmt"
	"os"
)

func (c *Client) onPlayButton() {
	c.mu.RLock()
	currentPath := c.currentPath
	c.mu.RUnlock()
	c.downloadDisable()
	defer c.downloadEnable()
	c.splashImage.Hide()
	c.statusLabel.Show()
	c.progressBar.SetValue(0)
	c.logScroll.Show()

	c.logf("Opening EverQuest")

	cmd := c.createCommand(true, fmt.Sprintf("%s/eqgame.exe", currentPath), "patchme")
	cmd.Dir = currentPath
	err := cmd.Run()
	if err != nil {
		c.logf("Failed to run EverQuest: %s", err)
		return
	}

	os.Exit(0)

	/*
	   stdout, err := cmd.StdoutPipe()

	   	if err != nil {
	   		c.logf("Failed to start convert.bat: stdoutpipe: %s", err)
	   		return
	   	}

	   stderr, err := cmd.StderrPipe()

	   	if err != nil {
	   		c.logf("Failed to start convert.bat: stderrpipe: %s", err)
	   		return
	   	}

	   err = cmd.Start()

	   	if err != nil {
	   		c.logf("Failed to run convert.bat: %s", err)
	   		return
	   	}

	   reader := io.MultiReader(stdout, stderr)
	   c.progressBar.SetValue(c.addProgress(0.1))
	   err = c.processOutput(reader, currentPath, zone, "convert.log")

	   	if err != nil {
	   		c.logf("Failed stdout: %s", err)
	   		return
	   	}

	   err = cmd.Wait()

	   	if err != nil {
	   		c.logf("Failed convert.bat: %s", err)
	   		return
	   	}
	   c.logf("Created %s.eqg", zone)
	*/
}

func (c *Client) onAutoPlayCheck(value bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if value {
		c.cfg.AutoPlay = "true"
		err := c.cfg.Save()
		if err != nil {
			c.logf("failed to save autoplay: %s", err)
		}
		return
	}
	c.cfg.AutoPlay = "false"
	err := c.cfg.Save()
	if err != nil {
		c.logf("failed to save autoplay %s:", err)
	}
}
