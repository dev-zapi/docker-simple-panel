package docker

import (
	"bufio"
	"os"
	"strings"
)

// ContainerInfo holds information about the current container environment
type ContainerEnvironment struct {
	IsInContainer bool
	ContainerID   string
}

// DetectContainerEnvironment detects if the application is running inside a Docker container
// and returns the container ID if available
func DetectContainerEnvironment() ContainerEnvironment {
	env := ContainerEnvironment{
		IsInContainer: false,
		ContainerID:   "",
	}

	// Method 1: Check /.dockerenv file (most reliable for Docker)
	if _, err := os.Stat("/.dockerenv"); err == nil {
		env.IsInContainer = true
	}

	// Method 2: Check cgroup file for Docker container ID
	cgroupPath := "/proc/1/cgroup"
	if file, err := os.Open(cgroupPath); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			// Look for docker container ID in cgroup entries
			// Format examples:
			// 12:memory:/docker/abc123...
			// 0::/docker/abc123...
			if strings.Contains(line, "/docker/") {
				env.IsInContainer = true
				parts := strings.Split(line, "/docker/")
				if len(parts) > 1 {
					containerID := strings.TrimSpace(parts[1])
					// Container ID can be full 64 chars or short 12 chars
					if len(containerID) >= 12 {
						env.ContainerID = containerID[:12] // Use short ID
					}
				}
				break
			}
			// Also check for containerd/podman/k8s patterns
			if strings.Contains(line, "/kubepods/") || strings.Contains(line, "/containerd/") {
				env.IsInContainer = true
			}
		}
	}

	// Method 3: Check hostname environment variable (fallback for some container runtimes)
	// Many containers set HOSTNAME to the container ID
	if env.IsInContainer && env.ContainerID == "" {
		if hostname := os.Getenv("HOSTNAME"); hostname != "" {
			// Hostname in containers is often the short container ID
			if len(hostname) == 12 || len(hostname) == 64 {
				env.ContainerID = hostname[:min(12, len(hostname))]
			}
		}
	}

	return env
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
