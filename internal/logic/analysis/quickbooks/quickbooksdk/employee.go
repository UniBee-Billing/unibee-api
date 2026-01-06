package quickbooksdk

import (
	"errors"
	"strconv"
)

type Employee struct {
	SyncToken        string          `json:",omitempty"`
	Domain           string          `json:"domain,omitempty"`
	DisplayName      string          `json:",omitempty"`
	PrimaryPhone     TelephoneNumber `json:",omitempty"`
	PrintOnCheckName string          `json:",omitempty"`
	FamilyName       string          `json:",omitempty"`
	Active           bool            `json:",omitempty"`
	SSN              string          `json:",omitempty"`
	PrimaryAddr      PhysicalAddress `json:",omitempty"`
	BillableTime     bool            `json:",omitempty"`
	GivenName        string          `json:",omitempty"`
	Id               string          `json:",omitempty"`
	MetaData         MetaData        `json:",omitempty"`
}

// CreateEmployee creates the given employee within QuickBooks
func (c *Client) CreateEmployee(employee *Employee) (*Employee, error) {
	var resp struct {
		Employee Employee
		Time     Date
	}

	if err := c.post("employee", employee, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Employee, nil
}

// FindEmployees gets the full list of Employees in the QuickBooks account.
func (c *Client) FindEmployees() ([]Employee, error) {
	var resp struct {
		QueryResponse struct {
			Employees     []Employee `json:"Employee"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM Employee", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no employees could be found")
	}

	employees := make([]Employee, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM Employee ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.Employees == nil {
			return nil, errors.New("no employees could be found")
		}

		employees = append(employees, resp.QueryResponse.Employees...)
	}

	return employees, nil
}

// FindEmployeeById returns an employee with a given Id.
func (c *Client) FindEmployeeById(id string) (*Employee, error) {
	var resp struct {
		Employee Employee
		Time     Date
	}

	if err := c.get("employee/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Employee, nil
}

// QueryEmployees accepts an SQL query and returns all employees found using it
func (c *Client) QueryEmployees(query string) ([]Employee, error) {
	var resp struct {
		QueryResponse struct {
			Employees     []Employee `json:"Employee"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Employees == nil {
		return nil, errors.New("could not find any employees")
	}

	return resp.QueryResponse.Employees, nil
}

// UpdateEmployee updates the employee
func (c *Client) UpdateEmployee(employee *Employee) (*Employee, error) {
	if employee.Id == "" {
		return nil, errors.New("missing employee id")
	}

	existingEmployee, err := c.FindEmployeeById(employee.Id)
	if err != nil {
		return nil, err
	}

	employee.SyncToken = existingEmployee.SyncToken

	payload := struct {
		*Employee
		Sparse bool `json:"sparse"`
	}{
		Employee: employee,
		Sparse:   true,
	}

	var employeeData struct {
		Employee Employee
		Time     Date
	}

	if err = c.post("employee", payload, &employeeData, nil); err != nil {
		return nil, err
	}

	return &employeeData.Employee, err
}
