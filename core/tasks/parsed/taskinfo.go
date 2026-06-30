package parsed

import "path"

type TaskInfo interface {
	TaskID() string
	Site() string
	TotalResources() int64
	Downloaded() int64
	TotalBytes() int64
	DownloadedBytes() int64
	Processing() map[string]ResourceInfo
	StorageName() string
	StoragePath() string
}

// StoragePath 回傳要顯示給使用者的儲存路徑。
// 單一資源時 StorPath 通常為空（僅多資源才以 item.Title 建立子目錄），
// 實際檔案存放在 path.Join(StorPath, resource.Filename)，
// 因此單一資源時帶上檔名，避免完成訊息僅顯示空路徑。
func (t *Task) StoragePath() string {
	if len(t.item.Resources) == 1 {
		return path.Join(t.StorPath, t.item.Resources[0].Filename)
	}
	return t.StorPath
}
func (t *Task) TotalResources() int64 {
	return t.totalResources
}

func (t *Task) Downloaded() int64 {
	return t.downloaded.Load()
}

func (t *Task) StorageName() string {
	return t.Stor.Name()
}

func (t *Task) Site() string {
	return t.item.Site
}

func (t *Task) TotalBytes() int64 {
	return t.totalBytes
}

func (t *Task) DownloadedBytes() int64 {
	return t.downloadedBytes.Load()
}

func (t *Task) Processing() map[string]ResourceInfo {
	t.processingMu.RLock()
	defer t.processingMu.RUnlock()
	return t.processing
}

type ResourceInfo interface {
	FileName() string
	FileSize() int64
}
