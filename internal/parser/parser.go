package parser

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strings"

	"github.com/v4-nikishin/db-migrator/internal/logger"
)

type Command string

const (
	commentPrefix         = "--"
	commandPrefix         = commentPrefix + " +db-migrator "
	upCommand     Command = "Up"
	downCommand   Command = "Down"
)

type Parser struct {
	log  *logger.Logger
	path string
}

func New(logger *logger.Logger, path string) *Parser {
	return &Parser{log: logger, path: path}
}

func (p *Parser) UpMigration() ([]string, error) {
	return p.parse(upCommand)
}

func (p *Parser) DownMigration() ([]string, error) {
	return p.parse(downCommand)
}

func (p *Parser) parse(migrationCmd Command) ([]string, error) {
	f, err := os.Open(p.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	var statements []string
	currentCmd := upCommand
	var buf bytes.Buffer
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, commentPrefix) && !strings.HasPrefix(line, commandPrefix) {
			continue
		}
		p.log.Debug(line)
		if strings.HasPrefix(line, commandPrefix) {
			currentCmd, err = parseCommand(line)
			if err != nil {
				return nil, err
			}
		}
		if currentCmd == migrationCmd {
			if !strings.HasPrefix(line, commandPrefix) {
				if !endsWithSemicolon(line) {
					line = line + "\n"
				}
				if _, err := buf.WriteString(line); err != nil {
					return nil, err
				}
			}
			if endsWithSemicolon(line) {
				statements = append(statements, buf.String())
				buf.Reset()
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return statements, nil
}

func endsWithSemicolon(line string) bool {

	prev := ""
	scanner := bufio.NewScanner(strings.NewReader(line))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		if strings.HasPrefix(word, commentPrefix) {
			break
		}
		prev = word
	}

	return strings.HasSuffix(prev, ";")
}

func parseCommand(line string) (Command, error) {
	if !strings.HasPrefix(line, commandPrefix) {
		return "", errors.New("ERROR: not a sql-migrate command")
	}
	fields := strings.Fields(line[len(commandPrefix):])
	if len(fields) == 0 {
		return "", errors.New(`ERROR: incomplete migration command`)
	}
	return Command(fields[0]), nil
}
