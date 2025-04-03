package pathx

import "os"

func UseWorkDir(workDir string, cb func()) (err error) {
	oldCwd, err := os.Getwd()
	if err != nil {
		return
	}

	err = os.Chdir(workDir)
	if err != nil {
		return
	}

	cb()

	err = os.Chdir(oldCwd)

	return
}
