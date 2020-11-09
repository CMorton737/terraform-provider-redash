//
// Copyright (c) 2020 Snowplow Analytics Ltd. All rights reserved.
//
// This program is licensed to you under the Apache License Version 2.0,
// and you may not use this file except in compliance with the Apache License Version 2.0.
// You may obtain a copy of the Apache License Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0.
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the Apache License Version 2.0 is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the Apache License Version 2.0 for the specific language governing permissions and limitations there under.
//
package main

import (
	"encoding/json"
	"context"
	"fmt"
	"strconv"
	//"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/snowplow-devops/redash-client-go/redash"
	// "github.com/davecgh/go-spew/spew"

)

func resourceRedashVisualization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRedashVisualizationCreate,
		ReadContext:   resourceRedashVisualizationRead,
		UpdateContext: resourceRedashVisualizationUpdate,
		DeleteContext: resourceRedashVisualizationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
				Required: true,
			},	
			"type": {
				Type: schema.TypeString,
				Required: true,
			},
			"query_id": {
				Type: schema.TypeInt,
				Required: true,
			},
			"options": {
				Type: schema.TypeString,
				Required: true,
			},	
		},
	}
}

func resourceRedashVisualizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	payload := redash.Visualization{
		Name:            d.Get("name").(string),
		QueryID:         d.Get("query_id").(int),
		Type:            d.Get("type").(string),
	}
	options := []byte(d.Get("options").(string))
	json.Unmarshal(options, &payload.Options)

	visualization, err := c.CreateVisualization(&payload)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(visualization.ID))
	
	resourceRedashVisualizationRead(ctx, d, meta)

	return diags
}

func resourceRedashVisualizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	visualization, err := c.GetVisualization(d.Get("query_id").(int), id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", &visualization.Name)
	d.Set("type", &visualization.Type)
    d.Set("options", &visualization.Options)
	d.SetId(fmt.Sprint(visualization.ID))

	return diags
}

func resourceRedashVisualizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	payload := redash.Visualization{
		Name:               d.Get("name").(string),
		QueryID:            d.Get("query_id").(int),
		Type:               d.Get("type").(string),
	}
	options := []byte(d.Get("options").(string))
	json.Unmarshal(options, &payload.Options)

	_, err = c.UpdateVisualization(id, &payload)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRedashVisualizationRead(ctx, d, meta)
}

func resourceRedashVisualizationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteVisualization(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
