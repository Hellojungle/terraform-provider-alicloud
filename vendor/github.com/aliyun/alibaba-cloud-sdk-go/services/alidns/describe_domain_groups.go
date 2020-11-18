package alidns

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// DescribeDomainGroups invokes the alidns.DescribeDomainGroups API synchronously
func (client *Client) DescribeDomainGroups(request *DescribeDomainGroupsRequest) (response *DescribeDomainGroupsResponse, err error) {
	response = CreateDescribeDomainGroupsResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDomainGroupsWithChan invokes the alidns.DescribeDomainGroups API asynchronously
func (client *Client) DescribeDomainGroupsWithChan(request *DescribeDomainGroupsRequest) (<-chan *DescribeDomainGroupsResponse, <-chan error) {
	responseChan := make(chan *DescribeDomainGroupsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDomainGroups(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// DescribeDomainGroupsWithCallback invokes the alidns.DescribeDomainGroups API asynchronously
func (client *Client) DescribeDomainGroupsWithCallback(request *DescribeDomainGroupsRequest, callback func(response *DescribeDomainGroupsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDomainGroupsResponse
		var err error
		defer close(result)
		response, err = client.DescribeDomainGroups(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// DescribeDomainGroupsRequest is the request struct for api DescribeDomainGroups
type DescribeDomainGroupsRequest struct {
	*requests.RpcRequest
	PageNumber   requests.Integer `position:"Query" name:"PageNumber"`
	UserClientIp string           `position:"Query" name:"UserClientIp"`
	PageSize     requests.Integer `position:"Query" name:"PageSize"`
	Lang         string           `position:"Query" name:"Lang"`
	KeyWord      string           `position:"Query" name:"KeyWord"`
}

// DescribeDomainGroupsResponse is the response struct for api DescribeDomainGroups
type DescribeDomainGroupsResponse struct {
	*responses.BaseResponse
	RequestId    string       `json:"RequestId" xml:"RequestId"`
	TotalCount   int64        `json:"TotalCount" xml:"TotalCount"`
	PageNumber   int64        `json:"PageNumber" xml:"PageNumber"`
	PageSize     int64        `json:"PageSize" xml:"PageSize"`
	DomainGroups DomainGroups `json:"DomainGroups" xml:"DomainGroups"`
}

// CreateDescribeDomainGroupsRequest creates a request to invoke DescribeDomainGroups API
func CreateDescribeDomainGroupsRequest() (request *DescribeDomainGroupsRequest) {
	request = &DescribeDomainGroupsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Alidns", "2015-01-09", "DescribeDomainGroups", "alidns", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDescribeDomainGroupsResponse creates a response to parse from DescribeDomainGroups response
func CreateDescribeDomainGroupsResponse() (response *DescribeDomainGroupsResponse) {
	response = &DescribeDomainGroupsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
