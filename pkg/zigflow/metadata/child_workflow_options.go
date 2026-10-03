/*
 * Copyright 2025 - 2026 Zigflow authors <https://github.com/zigflow/zigflow/graphs/contributors>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package metadata

import (
	"fmt"

	"github.com/open-workflow-specification/sdk-go/v4/model"
	"github.com/zigflow/zigflow/pkg/utils"
	"go.temporal.io/sdk/workflow"
)

type ChildWorkflowOptions struct {
	TaskQueue                string          `json:"taskQueue"`
	WorkflowExecutionTimeout *model.Duration `json:"workflowExecutionTimeout"`
	WorkflowRunTimeout       *model.Duration `json:"workflowRunTimeout"`
	RetryPolicy              *RetryPolicy    `json:"retryPolicy"`
}

func (c *ChildWorkflowOptions) ToTemporal(opts *workflow.ChildWorkflowOptions) {
	if c.TaskQueue != "" {
		opts.TaskQueue = c.TaskQueue
	}
	if c.WorkflowExecutionTimeout != nil {
		opts.WorkflowExecutionTimeout = utils.ToDuration(c.WorkflowExecutionTimeout)
	}
	if c.WorkflowRunTimeout != nil {
		opts.WorkflowRunTimeout = utils.ToDuration(c.WorkflowRunTimeout)
	}
	if c.RetryPolicy != nil {
		opts.RetryPolicy = c.RetryPolicy.ToTemporal(opts.RetryPolicy)
	}
}

// ApplyChildWorkflowOptions overlays the task's childWorkflowOptions
// metadata, if any, onto opts. An unset task queue keeps the parent's.
func ApplyChildWorkflowOptions(task *model.TaskBase, opts *workflow.ChildWorkflowOptions) error {
	c, ok := task.Metadata[MetadataChildWorkflowOptions]
	if !ok {
		return nil
	}

	var cwo ChildWorkflowOptions
	if err := utils.ToType(c, &cwo); err != nil {
		return fmt.Errorf("error decoding task child workflow options metadata: %w", err)
	}

	cwo.ToTemporal(opts)

	return nil
}
