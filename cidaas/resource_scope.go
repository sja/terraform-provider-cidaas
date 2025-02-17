package cidaas

import (
	"context"
	"fmt"
	"strings"

	"terraform-provider-cidaas/helper/cidaas"
	"terraform-provider-cidaas/helper/util"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceScope() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScopeCreate,
		ReadContext:   resourceScopeRead,
		UpdateContext: resourceScopeUpdate,
		DeleteContext: resourceScopeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locale": {
				Type:     schema.TypeString,
				Required: true,
			},
			"language": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"security_level": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scope_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_name": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"required_user_consent": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceScopeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var scope cidaas.Scope
	var scopeDescription cidaas.ScopeLocalDescription

	scopeDescription.Locale = d.Get("locale").(string)
	scopeDescription.Language = d.Get("language").(string)
	scopeDescription.Title = d.Get("title").(string)
	scopeDescription.Description = d.Get("description").(string)

	scope.LocaleWiseDescription = []cidaas.ScopeLocalDescription{scopeDescription}
	scope.SecurityLevel = strings.ToUpper(d.Get("security_level").(string))
	scope.ScopeKey = d.Get("scope_key").(string)
	scope.RequiredUserConsent = d.Get("required_user_consent").(bool)
	scope.GroupName = util.InterfaceArray2StringArray(d.Get("group_name").([]interface{}))

	cidaas_client := m.(cidaas.CidaasClient)
	response, err := cidaas_client.CreateOrUpdateScope(scope)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to create scope %+v", scope.ScopeKey),
			Detail:   err.Error(),
		})
		return diags
	}
	if err := d.Set("_id", response.Data.ID); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "error while setting _id to scope resource",
			Detail:   err.Error(),
		})
		return diags
	}
	d.SetId(response.Data.ScopeKey)
	return diags
}

func resourceScopeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	cidaas_client := m.(cidaas.CidaasClient)
	scope_key := d.Id()
	response, err := cidaas_client.GetScope(strings.ToLower(scope_key))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to read scope %+v", scope_key),
			Detail:   err.Error(),
		})
		return diags
	}
	if len(response.Data.LocaleWiseDescription) > 0 {
		if err := d.Set("locale", response.Data.LocaleWiseDescription[0].Locale); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("language", response.Data.LocaleWiseDescription[0].Language); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("description", response.Data.LocaleWiseDescription[0].Description); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("title", response.Data.LocaleWiseDescription[0].Title); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("security_level", response.Data.SecurityLevel); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("required_user_consent", response.Data.RequiredUserConsent); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("_id", response.Data.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("scope_key", response.Data.ScopeKey); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("group_name", response.Data.GroupName); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceScopeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var scope cidaas.Scope
	var scopeDescription cidaas.ScopeLocalDescription

	scopeDescription.Locale = d.Get("locale").(string)
	scopeDescription.Language = d.Get("language").(string)
	scopeDescription.Title = d.Get("title").(string)
	scopeDescription.Description = d.Get("description").(string)

	scope.LocaleWiseDescription = []cidaas.ScopeLocalDescription{scopeDescription}
	scope.SecurityLevel = d.Get("security_level").(string)
	scope.ScopeKey = d.Get("scope_key").(string)
	scope.RequiredUserConsent = d.Get("required_user_consent").(bool)
	scope.GroupName = util.InterfaceArray2StringArray(d.Get("group_name").([]interface{}))
	scope.ID = d.Get("_id").(string)

	cidaas_client := m.(cidaas.CidaasClient)
	_, err := cidaas_client.CreateOrUpdateScope(scope)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to update scope %+v", scope.ScopeKey),
			Detail:   err.Error(),
		})
	}
	return diags
}

func resourceScopeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	cidaas_client := m.(cidaas.CidaasClient)
	scope_key := d.Id()
	_, err := cidaas_client.DeleteScope(scope_key)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("failed to delete scope %+v", scope_key),
			Detail:   err.Error(),
		})
	}
	return diags
}
