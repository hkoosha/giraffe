package flags

import (
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hkoosha/giraffe/core/t11y/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hkoosha/giraffe/core/container/setup"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

type Flag struct {
	name        string
	def         any
	description string
	envVar      string
	sensitive   bool
	// typ         int
}

type FlagStore struct {
	lg        glog.Lg //nolint:unused
	setup     setup.Registry
	get       setup.Registry
	mu        *sync.Mutex
	vp        *viper.Viper
	name      string //nolint:unused
	envPrefix string
	flags     []*Flag
}

func (s *FlagStore) define(
	name string,
	defaultValue any,
	description string,
	sensitive bool,
) *Flag {
	s.setup.EnsureOpen("define")

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range s.flags {
		if v.name == name {
			panic(EF("flag already defined: %s", name))
		}
	}

	f := &Flag{
		name:        name,
		def:         defaultValue,
		description: description,
		envVar:      strings.ToUpper(strings.ReplaceAll(name, "-", "_")),
		sensitive:   sensitive,
	}

	s.flags = append(s.flags, f)

	return f
}

func (s *FlagStore) DefineSensitive(
	name string,
	defaultValue any,
	description string,
) *Flag {
	return s.define(name, defaultValue, description, true)
}

func (s *FlagStore) Define(
	name string,
	defaultValue any,
	description string,
) *Flag {
	return s.define(name, defaultValue, description, false)
}

func (s *FlagStore) Bind(cmd *cobra.Command) {
	s.setup.Finish("define")

	for _, v := range s.flags {
		name := v.name

		switch def := v.def.(type) {
		case bool:
			cmd.PersistentFlags().Bool(name, def, v.description)

		case string:
			cmd.PersistentFlags().String(name, def, v.description)

		case time.Duration:
			cmd.PersistentFlags().Duration(name, def, v.description)

		case int:
			cmd.PersistentFlags().Int(name, def, v.description)

		case int64:
			cmd.PersistentFlags().Int64(name, def, v.description)

		case uint:
			cmd.PersistentFlags().Uint(name, def, v.description)

		case uint16:
			cmd.PersistentFlags().Uint16(name, def, v.description)

		case uint64:
			cmd.PersistentFlags().Uint64(name, def, v.description)

		default:
			panic(EF("%s", "unsupported flag type: "+reflect.TypeOf(v.def).String()))
		}

		lu := cmd.PersistentFlags().Lookup(name)
		err := s.vp.BindPFlag(name, lu)
		OK(err)

		s.vp.MustBindEnv(name, s.envPrefix+v.envVar)
	}
}

func (s *FlagStore) Dump() []string {
	s.setup.EnsureOpen("dump")

	dump := make(map[string]string, len(s.flags))

	for _, f := range s.flags {
		name := f.name

		var v string
		switch f.def.(type) {
		case bool:
			v = strconv.FormatBool(s.vp.GetBool(name))

		case string:
			v = s.vp.GetString(name)

		case time.Duration:
			v = s.vp.GetDuration(name).String()

		case int:
			v = strconv.Itoa(s.vp.GetInt(name))

		case uint:
			v = strconv.FormatUint(uint64(s.vp.GetUint(name)), 10)

		case uint16:
			v = strconv.FormatUint(uint64(s.vp.GetUint16(name)), 10)

		default:
			panic(EF("unsupported flag type: %s", reflect.TypeOf(f.def).String()))
		}

		if f.sensitive {
			dump[name] = "***"
		} else {
			dump[name] = v
		}
	}

	keys := slices.Collect(maps.Keys(dump))
	slices.Sort(keys)

	printable := make([]string, len(keys))
	for i, k := range keys {
		printable[i] = fmt.Sprintf("%s=%s", keys[i], dump[k])
	}

	return printable
}

// ==============================================================================.

func (s *FlagStore) GetBool(flag *Flag) bool {
	s.get.EnsureOpen()

	return s.vp.GetBool(flag.name)
}

func (s *FlagStore) GetString(flag *Flag) string {
	s.get.EnsureOpen()

	v := s.vp.GetString(flag.name)

	return strings.TrimSpace(v)
}

func (s *FlagStore) GetEndpoint(flag *Flag) string {
	return s.GetString(flag)
}

func (s *FlagStore) GetPort(flag *Flag) uint16 {
	s.get.EnsureOpen()

	v := s.vp.GetUint16(flag.name)

	return v
}

func (s *FlagStore) GetUint8(flag *Flag) uint8 {
	s.get.EnsureOpen()

	v := s.vp.GetUint(flag.name)
	if v > 255 {
		panic(EF("value must be between 0 and 255: %v=%v", flag.name, v))
	}

	return uint8(v)
}

func (s *FlagStore) GetHttpTimeout(flag *Flag) time.Duration {
	return s.GetDuration(
		flag,
		10*time.Millisecond,
		1*time.Minute,
	)
}

func (s *FlagStore) GetDuration(
	flag *Flag,
	minAccepted time.Duration,
	maxAccepted time.Duration,
) time.Duration {
	s.get.EnsureOpen()

	v := s.vp.GetDuration(flag.name)

	if v < minAccepted {
		panic(EF(
			"duration not in valid range, flag=%s, accepted_range=%s~%s, adjusted: %s => %s",
			flag.name,
			minAccepted,
			maxAccepted,
			v,
			minAccepted,
		))
	}

	if v > maxAccepted {
		panic(EF("ERROR: duration not in valid range, flag=%s, "+
			"accepted_range=%s~%s, adjusted: %s => %s\n",
			flag.name,
			minAccepted,
			maxAccepted,
			v,
			maxAccepted,
		))
	}

	return v
}

// ==============================================================================.

func NewFlagStore(
	lg glog.Lg,
	name string,
	envPrefix string,
	vp *viper.Viper,
	reg setup.Registry,
) *FlagStore {
	return &FlagStore{
		lg:        lg,
		setup:     reg,
		get:       setup.New(),
		mu:        &sync.Mutex{},
		name:      name,
		vp:        vp,
		flags:     make([]*Flag, 0),
		envPrefix: envPrefix,
	}
}
