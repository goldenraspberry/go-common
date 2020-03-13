package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -config-file=env.ini
func TestConfig(t *testing.T) {
	wd, _ := exec.LookPath(os.Args[0])
	baseDir := filepath.Dir(wd)

	assert.Equal(t, GetListen(), ":9000", "get listen failure")
	assert.Equal(t, GetEnv(), "dev", "get env failure")
	assert.Equal(t, GetBaseDir(), baseDir, "get base_dir failure")
	assert.Equal(t, GetLogDir(), baseDir+"/logs3", "get log_dir failure")
	assert.Equal(t, GetCacheDir(), baseDir+"/cache2", "get cache_dir failure")
	assert.Equal(t, GetTmpDir(), "/tmp", "get tmp_dir failure")
	assert.Equal(t, GetLogPath(), baseDir+"/logs3/app.log", "get log_path failure")
	assert.Equal(t, GetAccessLogPath(), STDOUT, "get access_log_path failure")
	assert.Equal(t, GetErrorLogPath(), baseDir+"/logs3/err_log", "get error_log_path failure")
	assert.Equal(t, GetSlowLogPath(), DISABLE, "get slow_log_path failure")
	assert.Equal(t, GetLogLevel(), "debug", "get log_level failure")

	d := GetConfig("t")
	s := ""
	if _, ok := d["abc"]; ok {
		s = d["abc"]
	}
	assert.Equal(t, s, "1", "get t/abc failure")
}
