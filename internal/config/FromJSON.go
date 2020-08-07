package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"time"
)

type jsonFileConfig struct {
	FileName       string `json:fileName`
	EventDelimiter string `json:eventDelimiter`
	ReadInterval   string `json:readInterval`
	TimeLayout     string `json:timeLayout`
}

type jsonForwarderConfig struct {
	Enabled           *bool  `json:enabled`
	MaxBufferedEvents *int   `json:maxBufferedEvents`
	RecipientAddress  string `json:recipientAddress`
}

type jsonRecipientConfig struct {
	Enabled *bool  `json:enabled`
	Address string `json:address`
}

type jsonSqliteConfig struct {
	FileName string `json:fileName`
}

type jsonWebConfig struct {
	Enabled          *bool  `json:enabled`
	Address          string `json:address`
	UsePackagedFiles *bool  `json:usePackagedFiles`
}

type jsonConfig struct {
	Files           []jsonFileConfig `json:files`
	FieldExtractors []string         `json:fieldExtractors`

	HostName string `json:hostName`

	Forwarder *jsonForwarderConfig `json:forwarder`
	Recipient *jsonRecipientConfig `json:recipient`
	Sqlite    *jsonSqliteConfig    `json:sqlite`
	Web       *jsonWebConfig       `json:web`
}

var defaultConfig = Config{
	IndexedFiles: []IndexedFileConfig{},

	FieldExtractors: []*regexp.Regexp{
		regexp.MustCompile("(\\w+)=(\\w+)"),
		regexp.MustCompile("^(?P<_time>\\d\\d\\d\\d/\\d\\d/\\d\\d \\d\\d:\\d\\d:\\d\\d.\\d\\d\\d\\d\\d\\d)"),
	},

	Forwarder: &ForwarderConfig{
		Enabled:           false,
		MaxBufferedEvents: 1000000,
		RecipientAddress:  "http://localhost:8081",
	},

	Recipient: &RecipientConfig{
		Enabled: false,
		Address: ":8081",
	},

	SQLite: &SqliteConfig{
		DatabaseFile: "logsuck.db",
	},

	Web: &WebConfig{
		Enabled:          true,
		Address:          ":8080",
		UsePackagedFiles: true,
	},
}

var defaultEventDelimiter = regexp.MustCompile("\n")
var defaultReadInterval = 1 * time.Second
var defaultTimeLayout = "2006/01/02 15:04:05"

func FromJSON(r io.Reader) (*Config, error) {
	var cfg jsonConfig
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error decoding config JSON: %w", err)
	}

	indexedFiles := make([]IndexedFileConfig, len(cfg.Files))
	for i, file := range cfg.Files {
		if file.FileName == "" {
			return nil, fmt.Errorf("error reading config at files[%v]: fileName is empty", i)
		}
		indexedFiles[i].Filename = file.FileName

		if file.EventDelimiter == "" {
			log.Printf("Using default event delimiter for file=%v, defaultEventDelimiter=%v\n", file.FileName, defaultEventDelimiter)
			indexedFiles[i].EventDelimiter = defaultEventDelimiter
		} else {
			ed, err := regexp.Compile(file.EventDelimiter)
			if err != nil {
				return nil, fmt.Errorf("error reading config at files[%v]: error compiling eventDelimiter regexp: %w", i, err)
			}
			indexedFiles[i].EventDelimiter = ed
		}

		if file.ReadInterval == "" {
			log.Printf("Using default read interval for file=%v, defaultReadInterval=%v\n", file.FileName, defaultReadInterval)
			indexedFiles[i].ReadInterval = defaultReadInterval
		} else {
			ri, err := time.ParseDuration(file.ReadInterval)
			if err != nil {
				return nil, fmt.Errorf("error reading config at files[%v]: error parsing readInterval duration: %w", i, err)
			}
			indexedFiles[i].ReadInterval = ri
		}

		if file.TimeLayout == "" {
			log.Printf("Using default time layout for file=%v, defaultTimeLayout=%v\n", file.FileName, defaultTimeLayout)
			indexedFiles[i].TimeLayout = defaultTimeLayout
		} else {
			indexedFiles[i].TimeLayout = file.TimeLayout
		}
	}

	var fieldExtractors []*regexp.Regexp
	if len(cfg.FieldExtractors) == 0 {
		log.Printf("Using default field extractors. defaultFieldExtractors=%v\n", defaultConfig.FieldExtractors)
		fieldExtractors = defaultConfig.FieldExtractors
	} else {
		fieldExtractors = make([]*regexp.Regexp, len(cfg.FieldExtractors))
		for i, fe := range cfg.FieldExtractors {
			re, err := regexp.Compile(fe)
			if err != nil {
				return nil, fmt.Errorf("error reading config at fieldExtractors[%v]: error compiling regexp: %w", i, err)
			}
			fieldExtractors[i] = re
		}
	}

	var hostName string
	if cfg.HostName != "" {
		log.Printf("Using hostName=%v\n", cfg.HostName)
		hostName = cfg.HostName
	} else {
		log.Println("No hostName in configuration, will try to get host name from operating system.")
		hostName, err = os.Hostname()
		if err != nil {
			return nil, fmt.Errorf("error getting host name: %w", err)
		}
		log.Printf("Got host name from operating system. hostName=%v\n", hostName)
	}

	var forwarder *ForwarderConfig
	if cfg.Forwarder == nil {
		log.Println("Using default forwarder configuration.")
		forwarder = defaultConfig.Forwarder
	} else {
		forwarder = &ForwarderConfig{}
		if cfg.Forwarder.Enabled == nil {
			log.Println("forwarder.enabled not specified, defaulting to false")
			forwarder.Enabled = false
		} else {
			forwarder.Enabled = *cfg.Forwarder.Enabled
		}
		if cfg.Forwarder.MaxBufferedEvents == nil {
			log.Printf("Using default maxBufferedEvents for forwarder. defaultBufferedEvents=%v\n", defaultConfig.Forwarder.MaxBufferedEvents)
			forwarder.MaxBufferedEvents = defaultConfig.Forwarder.MaxBufferedEvents
		} else {
			forwarder.MaxBufferedEvents = *cfg.Forwarder.MaxBufferedEvents
		}
		if cfg.Forwarder.RecipientAddress == "" {
			log.Printf("Using default recipientAddress for forwarder. dedfaultRecipientAddress=%v\n", defaultConfig.Forwarder.RecipientAddress)
			forwarder.RecipientAddress = defaultConfig.Forwarder.RecipientAddress
		} else {
			forwarder.RecipientAddress = cfg.Forwarder.RecipientAddress
		}
	}

	var recipient *RecipientConfig
	if cfg.Recipient == nil {
		log.Println("Using default recipient configuration.")
		recipient = defaultConfig.Recipient
	} else {
		recipient = &RecipientConfig{}
		if cfg.Recipient.Enabled == nil {
			log.Println("recipient.enabled not specified, defaulting to false")
			recipient.Enabled = false
		} else {
			recipient.Enabled = *cfg.Recipient.Enabled
		}
		if cfg.Recipient.Address == "" {
			log.Printf("Using default address for recipient. defaultAddress=%v\n", defaultConfig.Recipient.Address)
			recipient.Address = defaultConfig.Recipient.Address
		} else {
			recipient.Address = cfg.Recipient.Address
		}
	}

	var sqlite *SqliteConfig
	if cfg.Sqlite == nil {
		log.Println("Using default sqlite configuration.")
		sqlite = defaultConfig.SQLite
	} else {
		sqlite = &SqliteConfig{}
		if cfg.Sqlite.FileName == "" {
			log.Printf("Using default sqlite filename. defaultFileName=%v\n", defaultConfig.SQLite.DatabaseFile)
			sqlite.DatabaseFile = defaultConfig.SQLite.DatabaseFile
		} else {
			sqlite.DatabaseFile = cfg.Sqlite.FileName
		}
	}

	var web *WebConfig
	if cfg.Web == nil {
		log.Println("Using default web configuration.")
		web = defaultConfig.Web
		if forwarder.Enabled {
			log.Println("Disabling web GUI since forwarder is enabled.")
			web.Enabled = false
		}
	} else {
		web = &WebConfig{}
		if cfg.Web.Enabled == nil {
			if forwarder.Enabled {
				log.Println("web.enabled not specified but forwarder.enabled is true. Setting web.enabled to false")
				web.Enabled = false
			} else {
				log.Println("web.enabled not specified, defaulting to true")
				web.Enabled = true
			}
		} else {
			web.Enabled = *cfg.Web.Enabled
		}
		if cfg.Web.Address == "" {
			log.Printf("Using default web address. defaultWebAddress=%v\n", defaultConfig.Web.Address)
			web.Address = defaultConfig.Web.Address
		} else {
			web.Address = cfg.Web.Address
		}
		if cfg.Web.UsePackagedFiles == nil {
			log.Println("web.usePackagedFiles not specified, defaulting to true")
			web.UsePackagedFiles = true
		} else {
			web.UsePackagedFiles = *cfg.Web.UsePackagedFiles
		}
	}

	return &Config{
		IndexedFiles:    indexedFiles,
		FieldExtractors: fieldExtractors,

		HostName: hostName,

		Forwarder: forwarder,
		Recipient: recipient,

		SQLite: sqlite,

		Web: web,
	}, nil
}
