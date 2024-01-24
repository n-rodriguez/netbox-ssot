package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// PATCH /api/ipam/ip-addresses/{id}/
func (api *NetboxAPI) PatchIPAddress(diffMap map[string]interface{}, ipId int) (*objects.IPAddress, error) {
	api.Logger.Debug("Patching IP address ", ipId, " with data: ", diffMap, " in Netbox")

	requestBody, err := json.Marshal(diffMap)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)
	response, err := api.doRequest(MethodPatch, fmt.Sprintf("/api/ipam/ip-addresses/%d/", ipId), requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var ipResponse objects.IPAddress
	err = json.Unmarshal(response.Body, &ipResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully patched IP address: ", ipResponse)
	return &ipResponse, nil
}

// POST /api/ipam/ip-addresses/
func (api *NetboxAPI) CreateIPAddress(ip *objects.IPAddress) (*objects.IPAddress, error) {
	api.Logger.Debug("Creating IP address in Netbox with data: ", ip)

	requestBody, err := utils.NetboxJsonMarshal(ip)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)

	response, err := api.doRequest(MethodPost, "/api/ipam/ip-addresses/", requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var ipResponse objects.IPAddress
	err = json.Unmarshal(response.Body, &ipResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully created IP address: ", ipResponse)

	return &ipResponse, nil
}

// PATCH /api/ipam/vlan-groups/{id}/
func (api *NetboxAPI) PatchVlanGroup(diffMap map[string]interface{}, vlanGroupId int) (*objects.VlanGroup, error) {
	api.Logger.Debug("Patching VlanGroup ", vlanGroupId, " with data: ", diffMap, " in Netbox")

	requestBody, err := json.Marshal(diffMap)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)
	response, err := api.doRequest(MethodPatch, fmt.Sprintf("/api/ipam/vlan-groups/%d/", vlanGroupId), requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var vlanGroupResponse objects.VlanGroup
	err = json.Unmarshal(response.Body, &vlanGroupResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully patched VlanGroup: ", vlanGroupResponse)
	return &vlanGroupResponse, nil
}

// POST /api/ipam/vlan-groups/
func (api *NetboxAPI) CreateVlanGroup(vlan *objects.VlanGroup) (*objects.VlanGroup, error) {
	api.Logger.Debug("Creating VlanGroup in Netbox with data: ", vlan)

	requestBody, err := utils.NetboxJsonMarshal(vlan)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)

	response, err := api.doRequest(MethodPost, "/api/ipam/vlan-groups/", requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var vlanGroupResponse objects.VlanGroup
	err = json.Unmarshal(response.Body, &vlanGroupResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully created VlanGroup: ", vlanGroupResponse)

	return &vlanGroupResponse, nil
}

// PATCH /api/ipam/vlans/{id}/
func (api *NetboxAPI) PatchVlan(diffMap map[string]interface{}, vlanId int) (*objects.Vlan, error) {
	api.Logger.Debug("Patching Vlan ", vlanId, " with data: ", diffMap, " in Netbox")

	requestBody, err := json.Marshal(diffMap)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)
	response, err := api.doRequest(MethodPatch, fmt.Sprintf("/api/ipam/vlans/%d/", vlanId), requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var vlanResponse objects.Vlan
	err = json.Unmarshal(response.Body, &vlanResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully patched Vlan: ", vlanResponse)
	return &vlanResponse, nil
}

// POST /api/ipam/vlans/
func (api *NetboxAPI) CreateVlan(vlan *objects.Vlan) (*objects.Vlan, error) {
	api.Logger.Debug("Creating Vlan in Netbox with data: ", vlan)

	requestBody, err := utils.NetboxJsonMarshal(vlan)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)

	response, err := api.doRequest(MethodPost, "/api/ipam/vlans/", requestBodyBuffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var vlanResponse objects.Vlan
	err = json.Unmarshal(response.Body, &vlanResponse)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully created Vlan: ", vlanResponse)

	return &vlanResponse, nil
}