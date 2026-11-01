package local

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/charmbracelet/log"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/dustin/go-humanize"
	"github.com/krau/SaveAny-Bot/common/utils/fsutil"
	config "github.com/krau/SaveAny-Bot/config/storage"
	"github.com/krau/SaveAny-Bot/pkg/enums/ctxkey"
	storenum "github.com/krau/SaveAny-Bot/pkg/enums/storage"
	"github.com/krau/SaveAny-Bot/pkg/storagetypes"
)

const (
	// minFreeSpaceBeforeDownload 是開始下載前要求的最低可用空間 (1500 MiB)，
	// 不足時不下載，直接回報錯誤。
	minFreeSpaceBeforeDownload uint64 = 1500 * 1024 * 1024
	// minFreeSpace 是下載完成、寫入儲存前要求的最低可用空間 (1024 MiB)，
	// 不足時放棄該檔案。
	minFreeSpace uint64 = 1024 * 1024 * 1024
)

// availableSpace 查詢磁碟可用空間，測試時可替換。
var availableSpace = fsutil.AvailableSpace

// ErrInsufficientSpace 表示本地儲存所在磁碟可用空間不足。
var ErrInsufficientSpace = errors.New("insufficient free disk space")

// IsOutOfSpace 判斷 err 是否為磁碟空間不足：Save 的預先檢查 (ErrInsufficientSpace)，
// 或寫入時遇到的 ENOSPC（例如寫暫存檔時磁碟已滿）。這類錯誤重試也不會成功。
func IsOutOfSpace(err error) bool {
	return errors.Is(err, ErrInsufficientSpace) || errors.Is(err, syscall.ENOSPC)
}

type Local struct {
	config config.LocalStorageConfig
	logger *log.Logger
}

// checkFreeSpace 檢查儲存所在磁碟的可用空間，小於 min 時回傳 ErrInsufficientSpace。
// stage 用於錯誤訊息，說明是哪個階段的檢查。
func (l *Local) checkFreeSpace(min uint64, stage string) error {
	avail, err := availableSpace(l.config.BasePath)
	if err != nil {
		return fmt.Errorf("failed to check available disk space: %w", err)
	}
	if avail < min {
		return fmt.Errorf("%w: %s可用空間 %s，需至少 %s",
			ErrInsufficientSpace, stage, humanize.IBytes(avail), humanize.IBytes(min))
	}
	return nil
}

// CheckSpaceBeforeDownload 在開始下載前檢查可用空間，小於 1500 MiB 時回傳 ErrInsufficientSpace。
func (l *Local) CheckSpaceBeforeDownload() error {
	return l.checkFreeSpace(minFreeSpaceBeforeDownload, "下載前")
}

func (l *Local) Init(ctx context.Context, cfg config.StorageConfig) error {
	localConfig, ok := cfg.(*config.LocalStorageConfig)
	if !ok {
		return fmt.Errorf("failed to cast local config")
	}
	if err := localConfig.Validate(); err != nil {
		return err
	}
	l.config = *localConfig
	err := os.MkdirAll(localConfig.BasePath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create local storage directory: %w", err)
	}
	l.logger = log.FromContext(ctx).WithPrefix(fmt.Sprintf("local[%s]", l.config.Name))
	return nil
}

func (l *Local) Type() storenum.StorageType {
	return storenum.Local
}

func (l *Local) Name() string {
	return l.config.Name
}

func (l *Local) JoinStoragePath(path string) string {
	return filepath.Join(l.config.BasePath, path)
}

func (l *Local) Save(ctx context.Context, r io.Reader, storagePath string) error {
	if err := l.checkFreeSpace(minFreeSpace, "寫入前"); err != nil {
		return err
	}
	l.logger.Infof("Saving file to %s", storagePath)
	storagePath = l.JoinStoragePath(storagePath)

	ext := filepath.Ext(storagePath)
	base := strings.TrimSuffix(storagePath, ext)
	candidate := storagePath
	if overwrite, _ := ctx.Value(ctxkey.OverwriteExisting).(bool); !overwrite {
		for i := 1; l.existsPath(candidate); i++ {
			candidate = fmt.Sprintf("%s_%d%s", base, i, ext)
		}
	}

	absPath, err := filepath.Abs(candidate)
	if err != nil {
		return err
	}
	if err := fileutil.CreateDir(filepath.Dir(absPath)); err != nil {
		return err
	}
	file, err := os.Create(absPath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, r)
	return err
}

func (l *Local) Exists(ctx context.Context, storagePath string) bool {
	return l.existsPath(l.JoinStoragePath(storagePath))
}

func (l *Local) existsPath(storagePath string) bool {
	absPath, err := filepath.Abs(storagePath)
	if err != nil {
		return false
	}
	return fileutil.IsExist(absPath)
}

// ListFiles implements StorageListable interface
func (l *Local) ListFiles(ctx context.Context, dirPath string) ([]storagetypes.FileInfo, error) {
	absPath := l.JoinStoragePath(dirPath)

	entries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", absPath, err)
	}

	files := make([]storagetypes.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			l.logger.Warnf("Failed to get file info for %s: %v", entry.Name(), err)
			continue
		}

		filePath := filepath.Join(dirPath, entry.Name())
		files = append(files, storagetypes.FileInfo{
			Name:    entry.Name(),
			Path:    filePath,
			Size:    info.Size(),
			IsDir:   entry.IsDir(),
			ModTime: info.ModTime(),
		})
	}

	return files, nil
}

// OpenFile implements StorageReadable interface
func (l *Local) OpenFile(ctx context.Context, filePath string) (io.ReadCloser, int64, error) {
	absPath := l.JoinStoragePath(filePath)

	file, err := os.Open(absPath)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to open file %s: %w", absPath, err)
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, 0, fmt.Errorf("failed to stat file %s: %w", absPath, err)
	}

	return file, stat.Size(), nil
}
