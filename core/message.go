// Copyright 2015 trivago GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"github.com/trivago/gollum/shared"
	"runtime"
	"time"
)

// MessageStreamID is the "compiled name" of a stream
type MessageStreamID uint64

// MessageState is used as a return value for the Enqueu method
type MessageState int

const (
	// LogInternalStream is the name of the internal message channel (logs)
	LogInternalStream = "_GOLLUM_"
	// WildcardStream is the name of the "all streams" channel
	WildcardStream = "*"
	// DroppedStream is the name of the stream used to store dropped messages
	DroppedStream = "_DROPPED_"
	// MessageStateOk is returned if the message could be delivered
	MessageStateOk = MessageState(iota)
	// MessageStateTimeout is returned if a message timed out
	MessageStateTimeout = MessageState(iota)
	// MessageStateDiscard is returned if a message should be discarded
	MessageStateDiscard = MessageState(iota)
)

var (
	// LogInternalStreamID is the ID of the "_GOLLUM_" stream
	LogInternalStreamID = GetStreamID(LogInternalStream)
	// WildcardStreamID is the ID of the "*" stream
	WildcardStreamID = GetStreamID(WildcardStream)
	// DroppedStreamID is the ID of the "_DROPPED_" stream
	DroppedStreamID = GetStreamID(DroppedStream)
)

// MessageSource defines methods that are common to all message sources.
// Currently this is only a placeholder.
type MessageSource interface {
}

// AsyncMessageSource extends the MessageSource interface to allow a backchannel
// that simply forwards any message coming from the producer.
type AsyncMessageSource interface {
	MessageSource

	// EnqueueResponse sends a message to the source of another message.
	EnqueueResponse(msg Message)
}

// SerialMessageSource extends the AsyncMessageSource interface to allow waiting
// for all parts of the response to be submitted.
type SerialMessageSource interface {
	AsyncMessageSource

	// Notify the end of the response stream
	ResponseDone()
}

// LinkableMessageSource extends the MessageSource interface to allow a pipe
// like behaviour between two components that communicate messages.
type LinkableMessageSource interface {
	MessageSource
	// Link the message source to the message reciever. This makes it possible
	// to create stable "pipes" between e.g. a consumer and producer.
	Link(pipe interface{})

	// IsLinked has to return true if Link executed successfull and does not
	// need to be called again.
	IsLinked() bool
}

// Message is a container used for storing the internal state of messages.
// This struct is passed between consumers and producers.
type Message struct {
	Data      []byte
	StreamID  MessageStreamID
	Source    MessageSource
	Timestamp time.Time
	Sequence  uint64
}

// NewMessage creates a new message from a given data stream
func NewMessage(source MessageSource, data []byte, sequence uint64) Message {
	return Message{
		Data:      data,
		Source:    source,
		StreamID:  WildcardStreamID,
		Timestamp: shared.LowResolutionTimeNow, //time.Now(),
		Sequence:  sequence,
	}
}

// String implements the stringer interface
func (msg Message) String() string {
	return string(msg.Data)
}

// Enqueue is a convenience function to push a message to a channel while
// waiting for a timeout instead of just blocking.
// Passing a timeout of -1 will discard the message.
// Passing a timout of 0 will always block.
// Messages that time out will be passed to the dropped queue if a Dropped
// consumer exists.
// The source parameter is used when a message is dropped, i.e. it is passed
// to the Drop function.
func (msg Message) Enqueue(channel chan<- Message, timeout time.Duration) MessageState {
	if timeout == 0 {
		channel <- msg
		return MessageStateOk // ### return, done ###
	}

	start := time.Time{}
	for {
		select {
		case channel <- msg:
			return MessageStateOk // ### return, done ###

		default:
			switch {
			// Start timeout based retries
			case start.IsZero():
				if timeout < 0 {
					return MessageStateDiscard // ### return, discard and ignore ###
				}
				start = time.Now()

			// Discard message after timeout
			case time.Since(start) > timeout:
				return MessageStateTimeout // ### return, drop and retry ###

			// Yield and try again
			default:
				runtime.Gosched()
			}
		}
	}
}

// Route enqueues this message to the given stream.
// If the stream does not exist, a default stream (broadcast) is created.
func (msg Message) Route(streamID MessageStreamID) {
	targetStream := StreamTypes.GetStreamOrFallback(streamID)
	targetStream.Enqueue(msg)
}
