package system

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os/exec"
	"proxy-dev/internal/util"
	"runtime"
)

// 检查证书是否已存在于根证书存储中
func CheckExistCert(certPath string) (bool, error) {
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return false, fmt.Errorf("读取证书文件失败: %v", err)
	}

	block, _ := pem.Decode(certData)
	if block == nil {
		return false, fmt.Errorf("解析PEM格式证书失败")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false, fmt.Errorf("解析证书失败: %v", err)
	}

	switch runtime.GOOS {
	case "windows":
		output, err := util.Exec("certutil", "-user", "-verifystore", "Root")
		if err != nil {
			return false, fmt.Errorf("检查证书存储失败: %v", err)
		}

		fingerprint := util.GetSHA1(cert)
		return bytes.Contains(output, []byte(fingerprint)), nil
	case "darwin": // macOS
		// macOS - 使用 security 检查证书
		_, err := util.Exec("security", "find-certificate", "-c", cert.Subject.CommonName, "-a", "/Library/Keychains/System.keychain")
		return err == nil, nil

	case "linux":
		// Linux - 检查证书目录
		output, err := util.Exec("ls", "/usr/local/share/ca-certificates/")
		if err != nil {
			return false, fmt.Errorf("检查证书目录失败: %v", err)
		}
		return bytes.Contains(output, []byte(certPath)), nil

	default:
		return false, fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}
}

func InstallCert(certPath string) error {
	var err error

	switch runtime.GOOS {
	case "windows":
		//cmd = exec.Command("cmd", "/c", "start", "", certPath)
		_, err = util.Exec("certutil", "-user", "-addstore", "-f", "Root", certPath)
	case "darwin":
		_, err = util.Exec("open", certPath)
	case "linux":
		_, err = util.Exec("xdg-open", certPath)
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	return err
}

func InstallCertWithRoot(certPath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("certutil", "-addstore", "-f", "Root", certPath)
	case "darwin":
		cmd = exec.Command("sudo", "security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", "/Library/Keychains/System.keychain", certPath)
	case "linux":
		cmd = exec.Command("sudo", "cp", certPath, "/usr/local/share/ca-certificates/")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("复制证书失败: %v", err)
		}
		cmd = exec.Command("sudo", "update-ca-certificates")
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("证书安装失败: %v\n输出: %s", err, string(output))
	}

	return nil
}
