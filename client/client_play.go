package client

import (
	"fmt"
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
		c.logf("Failed to run: %s", err)
	}

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
