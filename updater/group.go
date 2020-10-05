package updater

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Group struct {
	// Identify the group and members:
	// Name labels unique groups
	Name string `yaml:"name"`
	// Pattern is a prefix for the dependency, or a regular expression enclosed by /'s
	Pattern string `yaml:"pattern"`

	// Parameters that apply to members:
	// Range is a comma separated list of allowed semver ranges
	Range     string    `yaml:"range"`
	Frequency Frequency `yaml:"frequency"`

	compiledPattern *regexp.Regexp
}

func NewTestGroup(name, pattern string) *Group {
	return &Group{
		Name:            name,
		Pattern:         pattern,
		compiledPattern: regexp.MustCompile(pattern),
	}
}

type Frequency string

const (
	FrequencyDaily  Frequency = "daily"
	FrequencyWeekly Frequency = "weekly"
)

// Groups is an ordered list of Group with unique names.
// Prefer a list with .Name to map[string]Group for clear iteration order.
type Groups []*Group

func ParseGroups(s string) (Groups, error) {
	ug := Groups{}
	if err := yaml.Unmarshal([]byte(s), &ug); err != nil {
		return nil, err
	}
	if err := ug.Validate(); err != nil {
		return nil, err
	}

	return ug, nil
}

func (g Groups) Validate() error {
	uniqNames := map[string]struct{}{}
	for _, group := range g {
		if group.Name == "" {
			return fmt.Errorf("groups must specify name")
		}
		if group.Pattern == "" {
			return fmt.Errorf("groups must specify pattern")
		}
		switch group.Frequency {
		case "", FrequencyDaily, FrequencyWeekly:
		default:
			return fmt.Errorf("frequency must be: [%s,%s]", FrequencyDaily, FrequencyWeekly)
		}

		if _, ok := uniqNames[group.Name]; ok {
			return fmt.Errorf("duplicate group name: %q", group.Name)
		}
		uniqNames[group.Name] = struct{}{}

		if strings.HasPrefix(group.Pattern, "/") && strings.HasSuffix(group.Pattern, "/") {
			re, err := regexp.Compile(group.Pattern[1 : len(group.Pattern)-1])
			if err != nil {
				return fmt.Errorf("compiling group %q: %w", group.Name, err)
			}
			group.compiledPattern = re
		} else {
			group.compiledPattern = regexp.MustCompile("^" + regexp.QuoteMeta(group.Pattern))
		}
	}
	return nil
}

// GroupDependencies groups dependencies according to this configuration.
func (g Groups) GroupDependencies(deps []Dependency) (byGroupName map[string][]Dependency, ungrouped []Dependency) {
	byGroupName = make(map[string][]Dependency, len(g))
	for _, dep := range deps {
		group := g.matchGroup(dep)
		if group != "" {
			byGroupName[group] = append(byGroupName[group], dep)
		} else {
			ungrouped = append(ungrouped, dep)
		}
	}
	return
}

func (g Groups) matchGroup(dep Dependency) string {
	for _, group := range g {
		if group.compiledPattern.MatchString(dep.Path) {
			return group.Name
		}
	}
	return ""
}
