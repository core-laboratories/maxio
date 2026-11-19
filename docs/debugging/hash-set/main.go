// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"hash/crc32"
	"log"
	"os"
	"strings"

	"github.com/dchest/siphash"
	"github.com/google/uuid"
)

// hashes the key returning an integer based on the input algorithm.
// This function currently supports
// - SIPMOD
func sipHashMod(key string, cardinality int, id [16]byte) int {
	if cardinality <= 0 {
		return -1
	}
	// use the faster version as per siphash docs
	// https://github.com/dchest/siphash#usage
	k0, k1 := binary.LittleEndian.Uint64(id[0:8]), binary.LittleEndian.Uint64(id[8:16])
	sum64 := siphash.Hash(k0, k1, []byte(key))
	return int(sum64 % uint64(cardinality))
}

// hashOrder - hashes input key to return consistent
// hashed integer slice. Returned integer order is salted
// with an input key. This results in consistent order.
// NOTE: collisions are fine, we are not looking for uniqueness
// in the slices returned.
func hashOrder(key string, cardinality int) []int {
	if cardinality <= 0 {
		// Returns an empty int slice for cardinality < 0.
		return nil
	}

	nums := make([]int, cardinality)
	keyCrc := crc32.Checksum([]byte(key), crc32.IEEETable)

	start := int(keyCrc % uint32(cardinality))
	for i := 1; i <= cardinality; i++ {
		nums[i-1] = 1 + ((start + i) % cardinality)
	}
	return nums
}

var (
	file, object, deploymentID, prefix string
	setCount, shards                   int
	verbose                            bool
)

func main() {
	flag.StringVar(&file, "file", "", "Read all objects from file, newline separated")
	flag.StringVar(&prefix, "prefix", "", "Add prefix to all objects")
	flag.StringVar(&object, "object", "", "Select an object")
	flag.StringVar(&deploymentID, "deployment-id", "", "MinIO deployment ID, obtained from 'format.json'")
	flag.IntVar(&setCount, "set-count", 0, "Total set count")
	flag.IntVar(&shards, "shards", 0, "Total shards count")
	flag.BoolVar(&verbose, "v", false, "Display all objects")

	flag.Parse()

	if deploymentID == "" {
		log.Fatalln("deployment ID is mandatory")
	}

	if setCount == 0 {
		log.Fatalln("set count cannot be zero")
	}

	id := uuid.MustParse(deploymentID)

	if file != "" {
		distrib := make([][]string, setCount)
		b, err := os.ReadFile(file)
		if err != nil {
			log.Fatalln(err)
		}
		b = bytes.ReplaceAll(b, []byte("\r"), []byte{})
		sc := bufio.NewScanner(bytes.NewBuffer(b))
		for sc.Scan() {
			object = strings.TrimSpace(sc.Text())
			set := sipHashMod(prefix+object, setCount, id)
			distrib[set] = append(distrib[set], prefix+object)
		}
		for set, files := range distrib {
			fmt.Println("Set:", set+1, "Objects:", len(files))
			if !verbose {
				continue
			}
			for _, s := range files {
				fmt.Printf("\t%s\n", s)
			}
		}
		os.Exit(0)
	}

	if object == "" {
		log.Fatalln("object name is mandatory")
	}

	if shards != 0 {
		fmt.Println("Erasure distribution for the object", hashOrder(prefix+object, shards))
	}
	fmt.Println("Erasure setNumber for the object", sipHashMod(prefix+object, setCount, id)+1)
}
