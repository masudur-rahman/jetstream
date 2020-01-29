// Copyright 2019 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/nats-io/jetstream/internal/jsch"
)

type actCmd struct {
	json bool
}

func configureActCommand(app *kingpin.Application) {
	c := &actCmd{}
	act := app.Command("account", "Account management").Alias("act")
	act.Command("info", "General account information").Alias("nfo").Action(c.infoAction)
	act.Flag("json", "Produce JSON output").Short('j').BoolVar(&c.json)
}

func (c *actCmd) infoAction(pc *kingpin.ParseContext) error {
	_, err := prepareHelper(servers, natsOpts()...)
	kingpin.FatalIfError(err, "setup failed")

	info, err := jsch.JetStreamAccountInfo()
	kingpin.FatalIfError(err, "could not request account info")

	if c.json {
		printJSON(info)
		return nil
	}

	fmt.Println("JetStream Account Information:")
	fmt.Println()
	fmt.Printf("      Memory: %s of %s\n", humanize.IBytes(info.Memory), humanize.IBytes(uint64(info.Limits.MaxMemory)))
	fmt.Printf("     Storage: %s of %s\n", humanize.IBytes(info.Store), humanize.IBytes(uint64(info.Limits.MaxStore)))

	if info.Limits.MaxStreams == -1 {
		fmt.Printf("     Streams: %d of Unlimited\n", info.Streams)
	} else {
		fmt.Printf("Message Sets: %d of %d\n", info.Streams, info.Limits.MaxStreams)
	}
	fmt.Println()

	return nil
}
