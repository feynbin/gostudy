package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// 用来在centos中安装docker-ce

const DOCKER_DAEMON_PATH = "/etc/docker/daemon.json"
const DOCKER_REPO_PATH = "/etc/yum.repos.d/docker-ce.repo"

func main() {
	// 配置docker-ce.repo
	if _, err := os.Stat(DOCKER_REPO_PATH); os.IsNotExist(err) {
		fmt.Println("docker-ce.repo not exist, create it")
		if err := repoconfig(); err != nil {
			fmt.Printf("Failed to create Docker-CE repo file: %v", err)
			return
		}
		fmt.Println("Docker-CE repo file created successfully")
	} else if err == nil {
		fmt.Println("Docker-CE repo file already exist")

		if err := ReplaceDockerMirror(); err != nil {
			fmt.Printf("Failed to replace Docker mirror: %v", err)
			return
		}
		fmt.Println("Docker-CE mirror source replaced successfully")
	} else {
		fmt.Printf("Failed to access Docker-CE repo file: %v", err)
		return
	}

	// 执行yum makecache
	if err := YumMakeCache(); err != nil {
		fmt.Printf("Failed to update yum cache: %v\n", err)
		return
	}
	fmt.Println("Yum cache updated successfully.")

	// 执行yum install docker-ce
	if err := YumInstallDocker(); err != nil {
		fmt.Printf("Failed to install docker-ce: %v\n", err)
		return
	}
	fmt.Println("Docker-CE installed successfully.")

	// 配置docker-ce镜像源为腾讯云镜像源
	if err := ConfigureDockerMirror("https://mirror.ccs.tencentyun.com"); err != nil {
		fmt.Printf("Failed to configure docker-ce mirror source: %v\n", err)
		return
	}
	fmt.Println("Docker-CE mirror source configured successfully.")

	// 设置docker-ce开机自启动并立刻启动
	if err := EnableAndStartDocker(); err != nil {
		fmt.Printf("Failed to enable docker-ce service on boot and start it: %v\n", err)
		return
	}
	fmt.Println("Docker-CE service enabled on boot and started successfully.")
}

func repoconfig() error {
	const DOCKER_REPO_PATH = "/etc/yum.repos.d/docker-ce.repo"
	content := []byte("[docker-ce-stable]\nname=Docker CE Stable - $basearch\nbaseurl=https://mirrors.nju.edu.cn/docker-ce/linux/centos/9/x86_64/stable\nenabled=1\ngpgcheck=0\ngpgkey=https://mirrors.nju.edu.cn/docker-ce/linux/centos/gpg\n")
	return ioutil.WriteFile(DOCKER_REPO_PATH, content, 0644)
}

func ReplaceDockerMirror() error {
	f, err := os.Open(DOCKER_REPO_PATH)
	if err != nil {
		fmt.Printf("Failed to open Docker-CE repo file: %v", err)
		return err
	}
	defer f.Close()
	newContent := ""
	reader := bufio.NewReader(f)
	for {
		lineBytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := string(lineBytes)
		if strings.HasPrefix(line, "baseurl=") {
			newContent += "baseurl=https://mirrors.nju.edu.cn/docker-ce/linux/centos/$releasever/$basearch/stable/\n"
		} else {
			newContent += line + "\n"
		}
	}
	if err = ioutil.WriteFile(DOCKER_REPO_PATH, []byte(newContent), 0644); err != nil {
		fmt.Printf("Failed to write Docker-CE repo file: %v", err)
		return err
	}
	fmt.Println("Docker-CE mirror source replaced successfully.")
	return nil
}

func runCommandWithOutput(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command execution failed: %v", err)
	}
	return nil
}

// 更新yum缓存
func YumMakeCache() error {
	return runCommandWithOutput("sudo", "-S", "yum", "makecache")
}

// 安装Docker-CE
func YumInstallDocker() error {
	return runCommandWithOutput("sudo", "-S", "yum", "install", "docker-ce", "-y")
}

// 配置Docker-CE镜像源为腾讯云镜像源
// func ConfigureDockerMirror() error {
// 	return runCommandWithOutput("sudo", "-S", "tee", "/etc/docker/daemon.json", "& echo '{\"registry-mirrors\": [\"https://mirror.ccs.tencentyun.com\"]}'", "> /dev/null")
// }

// 设置Docker-CE开机自启动并立刻启动
func EnableAndStartDocker() error {
	if err := runCommandWithOutput("sudo", "-S", "systemctl", "enable", "docker.service"); err != nil {
		return err
	}
	if err := runCommandWithOutput("sudo", "-S", "systemctl", "start", "docker.service"); err != nil {
		return err
	}
	return nil
}

func ConfigureDockerMirror(mirrorURL string) error {
	// 检查daemon.json文件是否存在
	_, err := os.Stat(DOCKER_DAEMON_PATH)
	if err != nil {
		// 文件不存在，创建一个空文件
		if os.IsNotExist(err) {
			if err = ioutil.WriteFile(DOCKER_DAEMON_PATH, []byte("{}"), 0644); err != nil {
				return fmt.Errorf("Failed to create Docker daemon.json file: %v", err)
			}
		} else {
			return fmt.Errorf("Failed to check Docker daemon.json file: %v", err)
		}
	}

	// 读取原配置文件内容
	f, err := os.OpenFile(DOCKER_DAEMON_PATH, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("Failed to open Docker daemon.json file: %v", err)
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return fmt.Errorf("Failed to get stat of Docker daemon.json file: %v", err)
	}
	fileSize := fileInfo.Size()

	conf := make(map[string]interface{})
	if fileSize == 0 {
		if _, err = f.Seek(0, 0); err != nil {
			return fmt.Errorf("Failed to seek to the beginning of Docker daemon.json file: %v", err)
		}
		if err = f.Truncate(0); err != nil {
			return fmt.Errorf("Failed to truncate Docker daemon.json file: %v", err)
		}
		_, _ = f.Seek(0, 0)
	} else {
		if err = json.NewDecoder(f).Decode(&conf); err != nil {
			return fmt.Errorf("Failed to decode JSON from Docker daemon.json file: %v", err)
		}
	}

	// 设置registry-mirrors参数的值
	mirrors, ok := conf["registry-mirrors"]
	if !ok {
		conf["registry-mirrors"] = []string{mirrorURL}
	} else {
		mirrorList, ok := mirrors.([]interface{})
		if !ok {
			return fmt.Errorf("Unexpected value for registry-mirrors parameter in Docker daemon.json file")
		}

		var newMirrorList []string
		for _, m := range mirrorList {
			if v, ok := m.(string); ok && v != mirrorURL {
				newMirrorList = append(newMirrorList, v)
			}
		}
		newMirrorList = append(newMirrorList, mirrorURL)
		conf["registry-mirrors"] = newMirrorList
	}

	// 将新的配置写回到文件中
	f.Truncate(0)
	f.Seek(0, 0)
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "    ")
	if err = encoder.Encode(conf); err != nil {
		return fmt.Errorf("Failed to encode JSON to Docker daemon.json file: %v", err)
	}

	return nil
}
