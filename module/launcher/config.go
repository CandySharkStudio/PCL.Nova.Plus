package launcher

import (
	"NovaPlus/module/mmcll"
	"context"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/ini.v1"
)

type config struct {
	path string
}

type ExceptionHandler[T any] struct {
	Code    int    `json:"code"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type ReaderWriter struct {
	Ctx context.Context
}

func NewExceptionHandler[T any](code int, status bool, message string, data T) ExceptionHandler[T] {
	return ExceptionHandler[T]{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func NewConfig(path string) *config {
	return &config{path: path}
}

func (config config) Read(section, key string) (string, error) {
	cfg, err := ini.Load(config.path)
	if err != nil {
		return "", err
	}
	sec, err := cfg.GetSection(section)
	if err != nil {
		return "", err
	}
	k, err := sec.GetKey(key)
	if err != nil {
		return "", err
	}
	return k.String(), nil
}
func (config config) Write(section, key, value string) error {
	cfg, err := ini.Load(config.path)
	if err != nil {
		cfg = ini.Empty()
	}
	s, err := cfg.GetSection(section)
	if err != nil {
		s, _ = cfg.NewSection(section)
	}
	s.Key(key).SetValue(value)
	return cfg.SaveTo(config.path)
}

// WriteConfig 写入配置
func (rw *ReaderWriter) WriteConfig(path, section, key, value string) {
	path = mmcll.If(path == "current", rw.GetConfigIniPath(), rw.GetOtherIniPath()).(string)
	conf := NewConfig(path)
	_ = conf.Write(section, key, value)
}

// ReadConfig 读取配置
func (rw *ReaderWriter) ReadConfig(path, section, key string) string {
	path = mmcll.If(path == "current", rw.GetConfigIniPath(), rw.GetOtherIniPath()).(string)
	conf := NewConfig(path)
	if v, err := conf.Read(section, key); err == nil {
		return v
	}
	return ""
}

// GetOtherIniPath 获取隐私的配置路径（比如账号部分绑定）【ps：macos 的路径与 GetConfigIniPath 毫无两样】
func (rw *ReaderWriter) GetOtherIniPath() string {
	home, err := GetOtherDir()
	if err != nil {
		home = GetCurrentDir()
	}
	// d := mmcll.If(runtime2.GOOS == "windows", []string{"PCL.Nova"}, mmcll.If(runtime2.GOOS == "darwin", []string{"PCL.Nova"}, []string{"PCL.Nova"})).([]string)
	// d = append([]string{home}, d...)
	// d = append(d, "Other.ini")
	res := filepath.Join(home, "PCL.Nova", "Other.ini")
	err = EnsureConfigFile(res)
	return mmcll.If(err != nil, "", res).(string)
}

// GetConfigIniPath 获取元数据配置（并不隐私，只保存在当前目录下）【ps：macos 的路径与 GetOtherIniPath 毫无两样】
func (rw *ReaderWriter) GetConfigIniPath() string {
	res := filepath.Join(GetCurrentDir(), "PCL.Nova", "PCL2.Nova.ini")
	err := EnsureConfigFile(res)
	return mmcll.If(err != nil, "", res).(string)
}

// OpenDirectoryDialog 打开文件夹选择框
func (rw *ReaderWriter) OpenDirectoryDialog(title string) string {
	dialog, err := runtime.OpenDirectoryDialog(rw.Ctx, runtime.OpenDialogOptions{
		Title:                title,
		CanCreateDirectories: true,
	})
	return mmcll.If(err != nil, "", dialog).(string)
}

// OpenFileDialog 打开文件选择框
func (rw *ReaderWriter) OpenFileDialog(title string, filter []string) string {
	var fileFilter []runtime.FileFilter
	for _, f := range filter {
		fileFilter = append(fileFilter, runtime.FileFilter{
			DisplayName: f,
			Pattern:     f,
		})
	}
	if len(filter) == 0 {
		dialog, err := runtime.OpenFileDialog(rw.Ctx, runtime.OpenDialogOptions{
			Title:                title,
			CanCreateDirectories: true,
		})
		return mmcll.If(err != nil, "", dialog).(string)
	} else {
		dialog, err := runtime.OpenFileDialog(rw.Ctx, runtime.OpenDialogOptions{
			Title:                title,
			CanCreateDirectories: true,
			Filters:              fileFilter,
		})
		return mmcll.If(err != nil, "", dialog).(string)
	}
}

// 通过【参数】打开 文件资源管理器
func (rw *ReaderWriter) OpenExplorer(fpath string) bool {
	err := openExp(fpath)
	return err == nil
}

// 获取当前 可执行文件 目录（macOS 获取的是 ~/Library/Application Support 目录）
func (rw *ReaderWriter) GetCurrentExeDir() string {
	return GetCurrentDir()
}

// GetIsolationPath 通过 rootPath 和 index 准确无误的找到 MC 版本路径的 path（通常存在 ./.minecraft/versions/<版本文件名>/PCL.Nova.ini）
func (rw *ReaderWriter) GetIsolationPath(rootPath string, index int) string {
	lm := &LaunchMethod{
		Ctx: rw.Ctx,
	}
	vers := lm.GetMCAllVersion(rootPath)
	ver, _ := mmcll.SafeGet(vers, index)
	return ver
}
