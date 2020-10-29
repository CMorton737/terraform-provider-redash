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
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/snowplow-devops/redash-client-go/redash"
)

func dataSourceRedashQuery() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			// "name": {
			// 	Type:     schema.TypeString,
			// 	Required: true,
			// },
			// "data_source_id": {
			// 	Type:     schema.TypeInt,
			// 	Optional: true,
			// },
			// "description": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			// "query": {
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// },
			// "tags": {
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeString,
			// 	},
			// },
			// "is_archived": {
			// 	Type:     schema.TypeBool,
			// 	Computed: true,
			// },
			// "is_draft": {
			// 	Type:     schema.TypeBool,
			// 	Optional: true,
			// },
		},
		ReadContext: dataSourceRedashQueryRead,
	}
}

func dataSourceRedashQueryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id := d.Get("id").(int)

	query, err := c.GetQuery(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("id", query.ID)
	d.Set("name", query.Name)
	d.Set("query", query.Query)
	//d.Set("tags", query.Tags)
	d.Set("is_draft", query.IsDraft)
	d.Set("data_source_id", query.DataSourceID)

	d.SetId(fmt.Sprint(query.ID))

	return diags
}
