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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/snowplow-devops/redash-client-go/redash"
	//"log"
	//"github.com/davecgh/go-spew/spew"
)

func resourceRedashQuery() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRedashQueryCreate,
		ReadContext:   resourceRedashQueryRead,
		UpdateContext: resourceRedashQueryUpdate,
		DeleteContext: resourceRedashQueryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_source_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"options": {
				Type: schema.TypeString,
				Optional: true,
			},
			// "options": {
			// 	Type: schema.TypeList,
			// 	Optional: true,
			// 	MaxItems: 1,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"parameters": {
			// 				Type: schema.TypeList,
			// 				Optional: true,
			// 				Elem: &schema.Resource{
			// 					Schema: map[string]*schema.Schema{
			// 						"name": {
			// 							Type: schema.TypeString,
			// 							Required: true,
			// 						},
			// 						"type": {
			// 							Type: schema.TypeString,
			// 							Required: true,
			// 						},
			// 						"enumOptions": {
			// 							Type: schema.TypeString,
			// 							Optional: true,
			// 						},
			// 						"value": {
			// 							Type: schema.TypeString,
			// 							Optional: true,
			// 						},
			// 					},
			// 				},
			// 			},
			// 		},
			// 	},
			// },	
			// "tags": {
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeString,
			// 	},
			// },
			// "is_draft": {
			// 	Type:     schema.TypeBool,
			// 	Optional: true,
			// 	DefaultFunc: func() (interface{}, error) {
			// 		return nil, nil
			// 	},
			// }, 
			// },
		},
	}
}

func resourceRedashQueryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	payload := redash.Query{
		Name:               d.Get("name").(string),
		DataSourceID:       d.Get("data_source_id").(int),
		Query:              d.Get("query").(string),
	}

	options := []byte(d.Get("options").(string))
	json.Unmarshal(options, &payload.Options)

	query, err := c.CreateQuery(&payload)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(query.ID))

	resourceRedashQueryRead(ctx, d, meta)

	return diags
}

func resourceRedashQueryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	query, err := c.GetQuery(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", &query.Name)
	d.Set("data_source_id", &query.DataSourceID)
	d.Set("query", &query.Query)
	//d.Set("tags", &group.Tags)

	d.SetId(fmt.Sprint(query.ID))

	return diags
}

func resourceRedashQueryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	payload := redash.Query{
		Name:               d.Get("name").(string),
		DataSourceID:       d.Get("data_source_id").(int),
		Query:              d.Get("query").(string),
	}
	options := []byte(d.Get("options").(string))
	json.Unmarshal(options, &payload.Options)

	_, err = c.UpdateQuery(id, &payload)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRedashQueryRead(ctx, d, meta)
}

func resourceRedashQueryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*redash.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteQuery(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

// func convertOptions(options *map[string]interface{}, toFormat string) map[string]interface{} {

// 	redashConversion := map[string]string{
// 		"connection_string":                  "connectionString",
// 		"db_name":                            "dbName",
// 		"json_key_file":                      "jsonKeyFile",
// 		"load_schema":                        "loadSchema",
// 		"maximum_billing_tier":               "maximumBillingTier",
// 		"project_id":                         "projectId",
// 		"replica_set_name":                   "replicaSetName",
// 		"total_mbytes_processed_limit":       "totalMBytesProcessedLimit",
// 		"use_standard_sql":                   "useStandardSql",
// 		"user_defined_function_resource_uri": "userDefinedFunctionResourceUri",
// 	}

// 	terraformConversion := map[string]string{
// 		"connectionString":               "connection_string",
// 		"dbName":                         "db_name",
// 		"jsonKeyFile":                    "json_key_file",
// 		"loadSchema":                     "load_schema",
// 		"maximumBillingTier":             "maximum_billing_tier",
// 		"projectId":                      "project_id",
// 		"replicaSetName":                 "replica_set_name",
// 		"totalMBytesProcessedLimit":      "total_mbytes_processed_limit",
// 		"useStandardSql":                 "use_standard_sql",
// 		"userDefinedFunctionResourceUri": "user_defined_function_resource_uri",
// 	}

// 	convertedOptions := map[string]interface{}{}

// 	for k, v := range *options {
// 		if toFormat == "redash" {
// 			if val, ok := redashConversion[k]; ok {
// 				convertedOptions[val] = v
// 			} else {
// 				convertedOptions[k] = v
// 			}
// 		}

// 		if toFormat == "terraform" {
// 			if val, ok := terraformConversion[k]; ok {
// 				convertedOptions[val] = v
// 			} else {
// 				convertedOptions[k] = v
// 			}
// 		}
// 	}

// 	return convertedOptions
// }
