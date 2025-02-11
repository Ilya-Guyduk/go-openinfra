package parser

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const sampleYAML = `
openinfra: 1.0.0
info:
  title: OpenInfra Specification
  description: A specification for describing infrastructure resources and components.
  version: 1.0.0
  contact:
    name: Your Name
    email: your.email@example.com
  license:
    name: Apache 2.0
    url: https://www.apache.org/licenses/LICENSE-2.0

providers:
  - name: local_virtualbox
    type: virtualbox
    connection:
      protocol: ssh
      host: 192.168.1.10
      port: 22
      authentication:
        method: password
        username: admin
        password: password
    capabilities:
      - name: create_vm
        description: Create a new virtual machine
        method: POST
        endpoint: /vms/create
        parameters:
          - name: name
            type: string
            required: true
          - name: cpu
            type: integer
            required: true
          - name: memory
            type: string
            required: true
      - name: delete_vm
        description: Delete an existing virtual machine
        method: DELETE
        endpoint: /vms/{vm_id}
        parameters:
          - name: vm_id
            type: string
            required: true
      - name: start_vm
        description: Start a virtual machine
        method: POST
        endpoint: /vms/{vm_id}/start
        parameters:
          - name: vm_id
            type: string
            required: true
      - name: stop_vm
        description: Stop a virtual machine
        method: POST
        endpoint: /vms/{vm_id}/stop
        parameters:
          - name: vm_id
            type: string
            required: true
      - name: restart_vm
        description: Restart a virtual machine
        method: POST
        endpoint: /vms/{vm_id}/restart
        parameters:
          - name: vm_id
            type: string
            required: true
      - name: list_vms
        description: List all virtual machines
        method: GET
        endpoint: /vms
        parameters: []

  - name: cloud_provider
    type: aws
    connection:
      protocol: https
      endpoint: https://api.cloudprovider.com
      authentication:
        method: api_key
        api_key: your_api_key_here
    capabilities:
      - name: create_instance
        description: Create a new cloud instance
        method: POST
        endpoint: /instances/create
        parameters:
          - name: instance_type
            type: string
            required: true
          - name: image_id
            type: string
            required: true
      - name: delete_instance
        description: Delete a cloud instance
        method: DELETE
        endpoint: /instances/{instance_id}
        parameters:
          - name: instance_id
            type: string
            required: true
      - name: start_instance
        description: Start a cloud instance
        method: POST
        endpoint: /instances/{instance_id}/start
        parameters:
          - name: instance_id
            type: string
            required: true
      - name: stop_instance
        description: Stop a cloud instance
        method: POST
        endpoint: /instances/{instance_id}/stop
        parameters:
          - name: instance_id
            type: string
            required: true
      - name: restart_instance
        description: Restart a cloud instance
        method: POST
        endpoint: /instances/{instance_id}/restart
        parameters:
          - name: instance_id
            type: string
            required: true
      - name: list_instances
        description: List all cloud instances
        method: GET
        endpoint: /instances
        parameters: []

components:
  - type: virtual_machine
    name: local_vm
    provider: local_virtualbox
    properties:
      cpu: 2
      memory: 4GB
      disk_size: 50GB
      os: ubuntu-22.04
      network: local_network
    actions:
      - name: start
        method: POST
        endpoint: /vms/{vm_id}/start
        parameters:
          - name: vm_id
            type: string
            required: true
      - name: stop
        method: POST
        endpoint: /vms/{vm_id}/stop
        parameters:
          - name: vm_id
            type: string
            required: true
      - name: restart
        method: POST
        endpoint: /vms/{vm_id}/restart
        parameters:
          - name: vm_id
            type: string
            required: true

  - type: network
    name: local_network
    provider: cloud_provider
    properties:
      cidr: 192.168.1.0/24
      gateway: 192.168.1.1
      dns_servers:
        - 8.8.8.8
        - 8.8.4.4
    actions:
      - name: create
        method: POST
        endpoint: /vms/create
        parameters:
          - name: cidr
            type: string
            required: true
      - name: delete
        method: DELETE
        endpoint: /vms/{network_id}
        parameters:
          - name: network_id
            type: string
            required: true

dependencies:
  - component: local_vm
    depends_on:
      - local_network
`

func TestInfo(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "openinfra-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name()) // Удаляем файл после теста

	_, err = tmpFile.WriteString(sampleYAML)
	assert.NoError(t, err)
	tmpFile.Close() // Закрываем файл, чтобы его можно было прочитать

	// Вызываем функцию ParseFile
	spec, err := ParseFile(tmpFile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, spec)

	assert.Equal(t, "OpenInfra Specification", spec.Info.Title)
	assert.Equal(t, "A specification for describing infrastructure resources and components.", spec.Info.Description)
	assert.Equal(t, "1.0.0", spec.Info.Version)
	assert.Equal(t, "Your Name", spec.Info.Contact.Name)
	assert.Equal(t, "your.email@example.com", spec.Info.Contact.Email)
	assert.Equal(t, "Apache 2.0", spec.Info.License.Name)
	assert.Equal(t, "https://www.apache.org/licenses/LICENSE-2.0", spec.Info.License.URL)

}

func TestParseFile(t *testing.T) {
	// Создаем временный YAML-файл
	tmpFile, err := os.CreateTemp("", "openinfra-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name()) // Удаляем файл после теста

	_, err = tmpFile.WriteString(sampleYAML)
	assert.NoError(t, err)
	tmpFile.Close() // Закрываем файл перед чтением

	// Вызываем функцию ParseFile
	spec, err := ParseFile(tmpFile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, spec)

	// Проверяем версию спецификации
	assert.Equal(t, "1.0.0", spec.Version)

	// Проверяем информацию
	assert.Equal(t, "OpenInfra Specification", spec.Info.Title)
	assert.Equal(t, "1.0.0", spec.Info.Version)
	assert.Equal(t, "Your Name", spec.Info.Contact.Name)
	assert.Equal(t, "Apache 2.0", spec.Info.License.Name)

	// Проверяем провайдеров
	assert.Len(t, spec.Providers, 2)
	provider, exists := spec.Providers["local_virtualbox"]
	assert.True(t, exists)
	assert.Equal(t, "virtualbox", provider.Type)
	assert.Equal(t, "ssh", provider.Connection.Protocol)
	assert.Equal(t, "192.168.1.10", provider.Connection.Host)
	assert.Equal(t, 22, provider.Connection.Port)

	// Проверяем одну из возможностей провайдера
	assert.Len(t, provider.Capabilities, 6)
	capability := provider.Capabilities[0]
	assert.Equal(t, "create_vm", capability.Name)
	assert.Equal(t, "Create a new virtual machine", capability.Description)
	assert.Equal(t, "POST", capability.Method)
	assert.Equal(t, "/vms/create", capability.Endpoint)
	assert.Len(t, capability.Parameters, 3)
	assert.Equal(t, "cpu", capability.Parameters[1].Name)
	assert.Equal(t, "integer", capability.Parameters[1].Type)

	// Проверяем компоненты
	assert.Len(t, spec.Resources, 2)
	vm, exists := spec.Resources["local_vm"]
	assert.True(t, exists)
	assert.Equal(t, "virtual_machine", vm.Type)
	assert.Equal(t, "local_vm", vm.Name)
	assert.Equal(t, "local_virtualbox", vm.Provider)
	assert.Equal(t, 2, vm.Properties["cpu"])
	assert.Equal(t, "4GB", vm.Properties["memory"])

	// Проверяем действия
	assert.Len(t, vm.Actions, 3)
	action := vm.Actions[0]
	assert.Equal(t, "start", action.Name)
	assert.Equal(t, "POST", action.Method)

	// Проверяем зависимости
	assert.Len(t, spec.Dependencies, 1)
	dep := spec.Dependencies[0]
	assert.Equal(t, "local_vm", dep.Resource)
	assert.ElementsMatch(t, []string{"local_network"}, dep.DependsOn)
}

// TestParseFileErrors проверяет обработку ошибок при чтении и парсинге YAML.
func TestParseFileErrors(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		setup       func()
		expectedErr string
	}{
		{
			name:        "Файл не найден",
			filename:    "nonexistent.yaml",
			setup:       func() {}, // Ничего не создаем, файл отсутствует
			expectedErr: "ошибка: файл nonexistent.yaml не найден",
		},
		{
			name:     "Нет прав на чтение",
			filename: "no_permission.yaml",
			setup: func() {
				os.WriteFile("no_permission.yaml", []byte("openinfra: 1.0"), 0200) // Только запись
			},
			expectedErr: "ошибка: недостаточно прав для чтения файла no_permission.yaml",
		},
		{
			name:     "Пустой файл",
			filename: "empty.yaml",
			setup: func() {
				os.WriteFile("empty.yaml", []byte{}, 0644)
			},
			expectedErr: "ошибка: файл empty.yaml пуст",
		},
		{
			name:     "Некорректный YAML",
			filename: "invalid.yaml",
			setup: func() {
				os.WriteFile("invalid.yaml", []byte("invalid_yaml: [unterminated"), 0644)
			},
			expectedErr: "ошибка: некорректное форматирование YAML",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := ParseFile(tt.filename)

			if err == nil {
				t.Fatalf("Ожидалась ошибка, но её нет")
			}

			if !errors.Is(err, os.ErrPermission) && !errors.Is(err, os.ErrNotExist) && !contains(err.Error(), tt.expectedErr) {
				t.Errorf("Ожидали ошибку: %q, но получили: %q", tt.expectedErr, err)
			}

			// Удаляем тестовые файлы после проверки
			os.Remove(tt.filename)
		})
	}
}

// contains проверяет, содержит ли строка подстроку (для упрощенной проверки ошибок)
func contains(str, substr string) bool {
	return len(str) >= len(substr) && str[:len(substr)] == substr
}
