package completion

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-ps"
)

func completionSet(ctx context.Context, ua UncAccess, params []string) ([]Element, error) {
	result := []Element{}
	base := strings.ToUpper(params[len(params)-1])
	for _, env1 := range os.Environ() {
		if strings.HasPrefix(strings.ToUpper(env1), base) {
			result = append(result, Element1(env1))
		}
	}
	return result, nil
}

func completionDir(ctx context.Context, ua UncAccess, params []string) ([]Element, error) {
	return listUpDirs(ctx, ua, params[len(params)-1])
}

func containsInList(list []Element, target Element) bool {
	for _, s := range list {
		if s.String() == target.String() {
			return true
		}
	}
	return false
}

func completionCd(ctx context.Context, ua UncAccess, params []string) ([]Element, error) {

	list, err := completionDir(ctx, ua, params)
	source := params[len(params)-1]
	if len(source) < 1 || source[0] == '.' || strings.ContainsAny(source, "/\\:") {
		return list, err
	}
	cdpath := os.Getenv("CDPATH")
	if cdpath == "" {
		return list, err
	}
	orgSlash := STD_SLASH[0]
	if UseSlash {
		orgSlash = OPT_SLASH[0]
	}
	base := strings.ToUpper(source)
	for _, cdpath1 := range filepath.SplitList(cdpath) {
		if files, err := ioutil.ReadDir(cdpath1); err == nil {
			for _, file1 := range files {
				if file1.IsDir() {
					name := strings.ToUpper(file1.Name())
					if strings.HasPrefix(name, base) {
						new1 := Element2{
							file1.Name() + string(orgSlash),
							file1.Name() + string(OPT_SLASH[0])}
						if !containsInList(list, new1) {
							list = append(list, new1)
						}
					}
				}
			}
		}
	}
	return list, nil
}

func completionEnv(ctx context.Context, ua UncAccess, param []string) ([]Element, error) {
	eq := -1
	for i := 1; i < len(param); i++ {
		if strings.Contains(param[i], "=") {
			eq = i
		}
	}
	current := len(param) - 1

	if current == eq || current == 1 {
		return completionSet(ctx, ua, param)
	} else if current == eq+1 {
		return listUpCommands(ctx, param[current])
	} else {
		return nil, nil
	}
}

func completionWhich(ctx context.Context, ua UncAccess, param []string) ([]Element, error) {
	if len(param) == 2 {
		return listUpCommands(ctx, param[len(param)-1])
	}
	return nil, nil
}

func completionProcessName(ctx context.Context, ua UncAccess, param []string) ([]Element, error) {
	processes, err := ps.Processes()
	if err != nil {
		return nil, err
	}
	uniq := map[string]struct{}{}
	base := strings.ToUpper(param[len(param)-1])
	for _, ps1 := range processes {
		name := ps1.Executable()
		if strings.HasPrefix(strings.ToUpper(name), base) {
			uniq[name] = struct{}{}
		}
	}
	result := make([]Element, 0, len(uniq))
	for name := range uniq {
		result = append(result, Element1(name))
	}
	return result, nil
}

func completionTaskKill(ctx context.Context, ua UncAccess, param []string) ([]Element, error) {
	if len(param) >= 3 && strings.EqualFold(param[len(param)-2], "/IM") {
		return completionProcessName(ctx, ua, param)
	}
	return nil, nil
}
