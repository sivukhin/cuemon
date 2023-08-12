package lib

import (
	"bytes"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/format"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

func ReadJson[T any](input string) (data T, source string, err error) {
	var content []byte
	if input == "" {
		content, err = io.ReadAll(os.Stdin)
		source = "stdin"
	} else {
		content, err = os.ReadFile(input)
		source = fmt.Sprintf("file '%v'", input)
	}
	if err != nil {
		err = fmt.Errorf("unable to read input file '%v': %w", input, err)
		return
	}
	err = json.Unmarshal(content, &data)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal input file '%v': %w", input, err)
		return
	}
	return
}

func Bootstrap(input, module, dir string, overwrite bool) error {
	grafana, source, err := ReadJson[Grafana](input)
	if err != nil {
		return fmt.Errorf("unable to get grafana dashboard json from %v: %w", source, err)
	}

	monitoring, err := MonitoringFiles(module, dir, grafana)
	if err != nil {
		return fmt.Errorf("unable to create monitoring files: %w", err)
	}

	log.Printf("ready to write %v files for monitoring bootstrap", len(monitoring))
	for _, file := range monitoring {
		log.Printf("  %v (%v bytes)", file.Path, len(file.Content))
	}

	errs := make([]error, 0)
	for _, file := range monitoring {
		err = os.MkdirAll(path.Dir(file.Path), os.ModePerm)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		mode := os.O_CREATE | os.O_RDWR
		if !overwrite {
			mode |= os.O_EXCL
		}
		f, err := os.OpenFile(file.Path, mode, os.ModePerm)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		err = WriteFile(f, []byte(file.Content))
		if err != nil {
			errs = append(errs, err)
		}
		_ = f.Close()
	}
	return errors.Join(errs...)
}

func CutNode(content []byte, oldNode ast.Node, newNode ast.Node) ([]byte, []byte, []byte, error) {
	newPart, err := format.Node(newNode, format.Simplify())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to format new node: %v", err)
	}
	oldPart, err := format.Node(oldNode, format.Simplify())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to format old node: %v", err)
	}
	start, end := oldNode.Pos().Line()-1, oldNode.End().Line()
	lines := bytes.Split(content, []byte("\n"))
	var result []byte
	result = append(result, bytes.Join(lines[0:start], []byte("\n"))...)
	result = append(result, newPart...)
	result = append(result, bytes.Join(lines[end:], []byte("\n"))...)
	return result, oldPart, newPart, nil
}

func Update(input, dir string, overwrite bool) error {
	if dir == "" {
		return fmt.Errorf("update dir should be provided")
	}
	panel, source, err := ReadJson[JsonRaw[GrafanaPanel]](input)
	if err != nil {
		return fmt.Errorf("unable to get grafana panel json from %v: %w", source, err)
	}
	updates, err := UpdateFiles(dir, panel)
	if err != nil {
		return fmt.Errorf("unable to prepare updates: %w", err)
	}
	if len(updates) == 0 {
		return fmt.Errorf("unable to find configuration for panel with name '%v'", panel.Value.Title)
	}
	if len(updates) > 1 {
		var filenames []string
		for _, update := range updates {
			filenames = append(filenames, update.Path)
		}
		return fmt.Errorf("found more than 1 configuration for panel with name '%v': files=%v", panel.Value.Title, filenames)
	}
	log.Printf("--- OLD (%v) ---\n", updates[0].Path)
	log.Printf("%v\n", string(updates[0].BeforePart))
	log.Printf("--- NEW (%v) ---\n", updates[0].Path)
	log.Printf("%v\n", string(updates[0].AfterPart))
	if !overwrite {
		log.Printf("skipped update of file %v - pass overwrite flag to force this action", updates[0].Path)
		return nil
	}
	err = os.WriteFile(updates[0].Path, updates[0].Content, os.ModePerm)
	if err != nil {
		return err
	}
	log.Printf("updated file %v", updates[0].Path)
	return nil
}
