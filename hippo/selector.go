package hippo

import (
	"errors"
	"maps"
	"strings"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/internal/gstrings"
)

var (
	errDuplicatedPlan = errors.New("duplicated plan")
	errMissingPlan    = errors.New("missing plan")
)

//goland:noinspection GoUnusedGlobalVariable
var Selector_ = &Selector{
	defaultPlan: "",
	plans:       map[string]*Plan{},
	pipelines:   map[string]*PipelineFn{},
}

type Selector struct {
	plans       map[string]*Plan
	pipelines   map[string]*PipelineFn
	defaultPlan string
}

func (p *Selector) clone() *Selector {
	s := Selector{
		defaultPlan: p.defaultPlan,
		plans:       maps.Clone(p.plans),
		pipelines:   maps.Clone(p.pipelines),
	}

	if s.plans == nil {
		s.plans = map[string]*Plan{}
		s.pipelines = map[string]*PipelineFn{}
	}

	return &s
}

func (p *Selector) String() string {
	const prefix = "Selector["
	const suffix = "]"

	value := strings.Builder{}

	if p.defaultPlan != "" {
		value.WriteString(p.defaultPlan)
		value.WriteByte('/')
	}

	value.WriteString(gstrings.JoinIt(maps.Keys(p.plans)))

	return prefix + value.String() + suffix
}

func (p *Selector) MustWithDefault(
	name string,
) *Selector {
	return M(p.WithDefault(name))
}

func (p *Selector) WithDefault(
	name string,
) (*Selector, error) {
	_, err := p.Select(name)
	if err != nil {
		return nil, err
	}

	cp := p.clone()
	cp.defaultPlan = name
	return cp, nil
}

func (p *Selector) AndPlan(
	name string,
	plan *Plan,
) (*Selector, error) {
	cp := p.clone()

	if _, ok := cp.plans[name]; ok {
		return nil, E(EF("duplicated plan: %s", name), errDuplicatedPlan)
	}

	n, err := Pipeline(plan)
	if err != nil {
		return nil, err
	}

	cp.plans[name] = plan
	cp.pipelines[name] = n

	return cp, nil
}

func (p *Selector) MustAndPlan(
	name string,
	plan *Plan,
) *Selector {
	return M(p.AndPlan(name, plan))
}

func (p *Selector) Default() (*PipelineFn, error) {
	if p.defaultPlan == "" {
		return nil, E(EF("default plan not set"), errMissingPlan)
	}

	return p.Select(p.defaultPlan)
}

func (p *Selector) MustDefault() *PipelineFn {
	return M(p.Default())
}

func (p *Selector) Select(
	name string,
) (*PipelineFn, error) {
	n, ok := p.pipelines[name]
	if ok {
		return n, nil
	}

	if p.defaultPlan != "" {
		def, defOk := p.pipelines[p.defaultPlan]
		if !defOk {
			panic(EF("default plan missing: %s", p.defaultPlan))
		}

		return def, nil
	}

	return nil, E(EF("missing plan: %s", name), errMissingPlan)
}
