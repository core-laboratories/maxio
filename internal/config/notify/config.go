// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package notify

import (
	"github.com/core-laboratories/maxio/internal/event/target"
)

// Config - notification target configuration structure, holds
// information about various notification targets.
type Config struct {
	AMQP          map[string]target.AMQPArgs          `json:"amqp"`
	Elasticsearch map[string]target.ElasticsearchArgs `json:"elasticsearch"`
	Kafka         map[string]target.KafkaArgs         `json:"kafka"`
	MQTT          map[string]target.MQTTArgs          `json:"mqtt"`
	MySQL         map[string]target.MySQLArgs         `json:"mysql"`
	NATS          map[string]target.NATSArgs          `json:"nats"`
	NSQ           map[string]target.NSQArgs           `json:"nsq"`
	PostgreSQL    map[string]target.PostgreSQLArgs    `json:"postgresql"`
	Redis         map[string]target.RedisArgs         `json:"redis"`
	Webhook       map[string]target.WebhookArgs       `json:"webhook"`
}

const (
	defaultTarget = "1"
)

// NewConfig - initialize notification config.
func NewConfig() Config {
	// Make sure to initialize notification targets
	cfg := Config{
		NSQ:           make(map[string]target.NSQArgs),
		AMQP:          make(map[string]target.AMQPArgs),
		MQTT:          make(map[string]target.MQTTArgs),
		NATS:          make(map[string]target.NATSArgs),
		Redis:         make(map[string]target.RedisArgs),
		MySQL:         make(map[string]target.MySQLArgs),
		Kafka:         make(map[string]target.KafkaArgs),
		Webhook:       make(map[string]target.WebhookArgs),
		PostgreSQL:    make(map[string]target.PostgreSQLArgs),
		Elasticsearch: make(map[string]target.ElasticsearchArgs),
	}
	cfg.NSQ[defaultTarget] = target.NSQArgs{}
	cfg.AMQP[defaultTarget] = target.AMQPArgs{}
	cfg.MQTT[defaultTarget] = target.MQTTArgs{}
	cfg.NATS[defaultTarget] = target.NATSArgs{}
	cfg.Redis[defaultTarget] = target.RedisArgs{}
	cfg.MySQL[defaultTarget] = target.MySQLArgs{}
	cfg.Kafka[defaultTarget] = target.KafkaArgs{}
	cfg.Webhook[defaultTarget] = target.WebhookArgs{}
	cfg.PostgreSQL[defaultTarget] = target.PostgreSQLArgs{}
	cfg.Elasticsearch[defaultTarget] = target.ElasticsearchArgs{}
	return cfg
}
