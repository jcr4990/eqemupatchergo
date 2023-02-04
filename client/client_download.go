package client

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2/widget"
)

func (c *Client) onPatchButton() {
	if c.patchButton.Text != "Patch" {
		c.patchButton.Disable()
		defer c.patchButton.Enable()
		select {
		case c.cancel <- true:
		case <-time.After(1 * time.Second):
		}
		c.patchButton.SetText("Patch")
		return
	}
	c.logScroll.Show()
	c.logf("Patching...")
	c.cancel = make(chan bool, 3)
	c.downloadDisable()
	defer c.downloadEnable()
	c.splashImage.Hide()
	c.statusLabel.Show()
	c.progressBar.SetValue(0)
	go c.asyncPatch()
}

func (c *Client) downloadDisable() {
	//c.patchButton.Disable()
	c.patchButton.SetText("Cancel")
	c.playButton.Disable()
}

func (c *Client) downloadEnable() {
	c.patchButton.SetText("Patch")
	c.playButton.Enable()
}

func (c *Client) asyncPatch() {
	err := c.patch()
	if err != nil {
		c.logf("Failed %s", err)
		return
	}
	c.mu.RLock()
	isAutoPatchPlay := c.isAutoPatchPlay
	c.mu.RUnlock()
	if isAutoPatchPlay {
		c.onPlayButton()
	}
}

func (c *Client) patch() error {
	start := time.Now()
	c.patchButton.Importance = widget.LowImportance

	c.mu.Lock()
	fileList := c.cacheFileList
	c.cacheFileList = nil
	c.mu.Unlock()
	var err error
	if fileList == nil {
		err = c.refreshFileList()
		if err != nil {
			return fmt.Errorf("filelist: %w", err)
		}
		c.mu.RLock()
		fileList = c.cacheFileList
		c.mu.RUnlock()
	}
	/*if c.cfg.LastPatchedVersion == fileList.Version {
		if len(fileList.Version) < 8 {
			c.logf("Skipping patch, we are up to date")
			return nil
		}
		c.logf("Skipping patch, we are up to date latest patch %s", fileList.Version[0:8])
		return nil
	}*/

	totalSize := int64(0)

	for _, entry := range fileList.Downloads {
		totalSize += int64(entry.Size)
	}

	progressSize := int64(1)

	totalDownloaded := int64(0)

	if len(fileList.Version) < 8 {
		c.logf("Total patch size: %s", generateSize(int(totalSize)))
	} else {
		c.logf("Total patch size: %s, version: %s", generateSize(int(totalSize)), fileList.Version[0:8])
	}
	for _, entry := range fileList.Downloads {
		select {
		case <-c.cancel:
			return fmt.Errorf("cancelled by user")
		default:
		}

		if strings.Contains(entry.Name, "/") {
			newPath := strings.TrimSuffix(entry.Name, filepath.Base(entry.Name))
			err = os.MkdirAll(newPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("mkdir %s: %w", newPath, err)
			}
		}
		_, err := os.Stat(entry.Name)
		if err != nil {
			if os.IsNotExist(err) {
				err = c.downloadPatchFile(entry)
				if err != nil {
					return fmt.Errorf("download new file: %w", err)
				}
				totalDownloaded += int64(entry.Size)
				progressSize += int64(entry.Size)
				c.progressBar.SetValue(float64(progressSize) / float64(totalSize))
				continue
			}
			return fmt.Errorf("stat %s: %w", entry.Name, err)
		}

		hash, err := md5Checksum(entry.Name)
		if err != nil {
			return fmt.Errorf("md5checksum: %w", err)
		}

		if hash == entry.Md5 {
			c.logf("%s skipped", entry.Name)
			progressSize += int64(entry.Size)
			c.progressBar.SetValue(float64(progressSize) / float64(totalSize))
			continue
		}

		err = c.downloadPatchFile(entry)
		if err != nil {
			return fmt.Errorf("download new file: %w", err)
		}
		progressSize += int64(entry.Size)
		totalDownloaded += int64(entry.Size)
		c.progressBar.SetValue(float64(progressSize) / float64(totalSize))
	}

	c.cfg.LastPatchedVersion = fileList.Version
	err = c.cfg.Save()
	if err != nil {
		c.logf("Failed to save version to eqemupatch.yml: %s", err)
	}

	if totalDownloaded == 0 {
		c.logf("Finished in %0.2f seconds", time.Since(start).Seconds())
		return nil
	}
	c.logf("Finished %s in %0.2f seconds", generateSize(int(totalDownloaded)), time.Since(start).Seconds())

	return nil
}

func (c *Client) downloadPatchFile(entry FileEntry) error {
	c.logf("%s (%s)", entry.Name, generateSize(entry.Size))
	w, err := os.Create(entry.Name)
	if err != nil {
		return fmt.Errorf("create %s: %w", entry.Name, err)
	}
	defer w.Close()
	client := http.DefaultClient

	url := fmt.Sprintf("%s/%s/%s", c.url, c.clientVersion, entry.Name)
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("download %s: %w", url, err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return fmt.Errorf("write %s: %w", entry.Name, err)
	}
	return nil
}

func md5Checksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", fmt.Errorf("new: %w", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func generateSize(in int) string {
	val := float64(in)
	if val < 1024 {
		return fmt.Sprintf("%0.2f bytes", val)
	}
	val /= 1024
	if val < 1024 {
		return fmt.Sprintf("%0.2f KB", val)
	}
	val /= 1024
	if val < 1024 {
		return fmt.Sprintf("%0.2f MB", val)
	}
	val /= 1024
	if val < 1024 {
		return fmt.Sprintf("%0.2f GB", val)
	}
	val /= 1024
	return fmt.Sprintf("%0.2f TB", val)
}
