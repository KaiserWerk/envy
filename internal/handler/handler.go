package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/KaiserWerk/envy/internal/configuration"
	"github.com/KaiserWerk/envy/internal/logging"
)

const envHeader = "X-Env"
const varHeader = "X-Var"
const varValueHeader = "X-Var-Value"

type Var struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func NewHandler(config *configuration.AppConfig, logger logging.Logger) *Handler {
	return &Handler{
		AppConfig: config,
		Logger:    logger,
		vars:      make(map[string]map[string]string),
		mut:       &sync.RWMutex{},
	}
}

type Handler struct {
	AppConfig *configuration.AppConfig
	Logger    logging.Logger
	vars      map[string]map[string]string
	mut       *sync.RWMutex
}

func (b *Handler) LoadVars() error {
	cont, err := os.ReadFile("vars.json")
	if err != nil {
		return err
	}

	return json.Unmarshal(cont, &b.vars)
}
func (b *Handler) StoreVars() error {
	j, err := json.Marshal(b.vars)
	if err != nil {
		return err
	}

	return os.WriteFile("vars.json", j, 0600)
}

func (b *Handler) GetVar(w http.ResponseWriter, r *http.Request) {
	env := r.Header.Get(envHeader)
	if env == "" {
		http.Error(w, `{"error": "missing X-Env header"}`, http.StatusBadRequest)
		return
	}
	varName := r.Header.Get(varHeader)
	if varName == "" {
		http.Error(w, `{"error": "missing X-Var header"}`, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	b.mut.RLock()
	defer b.mut.RUnlock()

	vars, ok := b.vars[env]
	if !ok {
		http.Error(w, `{"error": "env not found"}`, http.StatusBadRequest)
		return
	}

	varValue, ok := vars[varName]
	if !ok {
		http.Error(w, `{"error": "env var not found"}`, http.StatusBadRequest)
		return
	}

	v := Var{
		Name:  varName,
		Value: varValue,
	}

	j, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(j)
}

func (b *Handler) SetVar(w http.ResponseWriter, r *http.Request) {
	env := r.Header.Get(envHeader)
	if env == "" {
		http.Error(w, `{"error": "missing X-Env header"}`, http.StatusBadRequest)
		return
	}
	varName := r.Header.Get(varHeader)
	if varName == "" {
		http.Error(w, `{"error": "missing X-Var header"}`, http.StatusBadRequest)
		return
	}
	varValue := r.Header.Get(varValueHeader)
	if varName == "" {
		http.Error(w, `{"error": "missing X-Var-Value header"}`, http.StatusBadRequest)
		return
	}

	if len(varName) > 255 {
		http.Error(w, `{"error": "var name max length is 255"}`, http.StatusBadRequest)
		return
	}

	if len(varValue) > 1024 {
		http.Error(w, `{"error": "var value max length is 1024"}`, http.StatusBadRequest)
		return
	}

	b.mut.Lock()
	defer b.mut.Unlock()

	_, ok := b.vars[env]
	if !ok {
		b.vars[env] = map[string]string{}
	}

	b.vars[env][varName] = varValue
	_ = b.StoreVars()
}

func (b *Handler) GetAllVars(w http.ResponseWriter, r *http.Request) {
	env := r.Header.Get(envHeader)
	if env == "" {
		http.Error(w, `{"error": "missing X-Env header"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	b.mut.RLock()
	defer b.mut.RUnlock()

	vars, ok := b.vars[env]
	if !ok {
		http.Error(w, `{"error": "env not found"}`, http.StatusBadRequest)
		return
	}

	d := make([]Var, 0, len(vars))

	for n, v := range vars {
		d = append(d, Var{Name: n, Value: v})
	}

	j, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		http.Error(w, `{"error": "internal error"}`, http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(j)
}
